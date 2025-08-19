package services

import (
	"clipmini/models"
	"context"
	"sync"
	"time"
)

type HistoryService struct {
	history        *models.History
	fileService    *FileService
	mu             sync.RWMutex
	pendingChanges bool
	batchTimer     *time.Timer
	batchInterval  time.Duration
	ctx            context.Context
	cancel         context.CancelFunc
}

func NewHistoryService(config *models.AppConfig) *HistoryService {
	ctx, cancel := context.WithCancel(context.Background())
	hs := &HistoryService{
		history:       models.NewHistory(config.MaxHistoryItems),
		fileService:   NewFileService(config),
		batchInterval: 3 * time.Second, // 3秒批量寫入
		ctx:           ctx,
		cancel:        cancel,
	}
	return hs
}

func (hs *HistoryService) LoadFromFile() error {
	lines, err := hs.fileService.ReadHistoryLines()
	if err != nil {
		return nil
	}

	hs.history.FromFileFormat(lines)
	return nil
}

func (hs *HistoryService) SaveToFile() error {
	lines := hs.history.ToFileFormat()
	return hs.fileService.WriteHistoryLines(lines)
}

func (hs *HistoryService) AddItem(item *models.ClipboardItem) error {
	hs.mu.Lock()
	hs.history.Add(item)
	hs.scheduleBatchSave()
	hs.mu.Unlock()
	return nil
}

func (hs *HistoryService) GetItems() []*models.ClipboardItem {
	return hs.history.GetItems()
}

func (hs *HistoryService) GetItem(index int) *models.ClipboardItem {
	return hs.history.GetItem(index)
}

func (hs *HistoryService) RemoveItem(index int) error {
	hs.mu.Lock()
	defer hs.mu.Unlock()

	removedItem := hs.history.RemoveItem(index)
	if removedItem == nil {
		return nil
	}

	if removedItem.Type == models.ClipImage && removedItem.FilePath != "" {
		hs.fileService.CleanupImageFiles([]string{removedItem.FilePath})
	}

	hs.scheduleBatchSave()
	return nil
}

func (hs *HistoryService) UpdateItem(index int, newContent string) error {
	hs.mu.Lock()
	defer hs.mu.Unlock()

	if hs.history.UpdateItem(index, newContent) {
		hs.scheduleBatchSave()
	}
	return nil
}

func (hs *HistoryService) Clear() error {
	items := hs.history.GetItems()
	imagePaths := make([]string, 0, len(items))
	for _, item := range items {
		if item.Type == models.ClipImage && item.FilePath != "" {
			imagePaths = append(imagePaths, item.FilePath)
		}
	}
	hs.fileService.CleanupImageFiles(imagePaths)

	hs.history.Clear()

	hs.fileService.DeleteHistoryFile()
	hs.fileService.DeleteImageDirectory()

	return nil
}

func (hs *HistoryService) MaintainLimit() {
	items := hs.history.GetItems()
	if len(items) > hs.history.MaxItems {
		excessItems := items[hs.history.MaxItems:]
		imagePaths := make([]string, 0, len(excessItems))
		for _, item := range excessItems {
			if item.Type == models.ClipImage && item.FilePath != "" {
				imagePaths = append(imagePaths, item.FilePath)
			}
		}
		hs.fileService.CleanupImageFiles(imagePaths)

		hs.history.Items = items[:hs.history.MaxItems]
		hs.scheduleBatchSave()
	}
}

// scheduleBatchSave 安排批量保存
func (hs *HistoryService) scheduleBatchSave() {
	if hs.pendingChanges {
		return
	}

	hs.pendingChanges = true
	if hs.batchTimer != nil {
		hs.batchTimer.Stop()
	}

	hs.batchTimer = time.AfterFunc(hs.batchInterval, func() {
		hs.flushPendingChanges()
	})
}

// flushPendingChanges 立即寫入所有待保存的變更
func (hs *HistoryService) flushPendingChanges() {
	hs.mu.Lock()
	defer hs.mu.Unlock()

	if !hs.pendingChanges {
		return
	}

	hs.pendingChanges = false
	if err := hs.SaveToFile(); err != nil {
		// 記錄錯誤但不阻塞
		// log.Printf("Failed to save history: %v", err)
	}
}

// ForceFlush 強制立即寫入（用於應用程式關閉前）
func (hs *HistoryService) ForceFlush() error {
	hs.mu.Lock()
	defer hs.mu.Unlock()

	if hs.batchTimer != nil {
		hs.batchTimer.Stop()
		hs.batchTimer = nil
	}

	if hs.pendingChanges {
		hs.pendingChanges = false
		return hs.SaveToFile()
	}
	return nil
}

// Close 關閉服務並保存未寫入的變更
func (hs *HistoryService) Close() error {
	hs.cancel()
	return hs.ForceFlush()
}
