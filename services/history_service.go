package services

import (
	"clipmini/models"
)

type HistoryService struct {
	history     *models.History
	fileService *FileService
}

func NewHistoryService(config *models.AppConfig) *HistoryService {
	return &HistoryService{
		history:     models.NewHistory(config.MaxHistoryItems),
		fileService: NewFileService(config),
	}
}

func (hs *HistoryService) LoadFromFile() error {
	lines, err := hs.fileService.ReadHistoryLines()
	if err != nil {
		return nil // File might not exist yet, that's ok
	}
	
	hs.history.FromFileFormat(lines)
	return nil
}

func (hs *HistoryService) SaveToFile() error {
	lines := hs.history.ToFileFormat()
	return hs.fileService.WriteHistoryLines(lines)
}

func (hs *HistoryService) AddItem(item *models.ClipboardItem) error {
	hs.history.Add(item)
	return hs.SaveToFile()
}

func (hs *HistoryService) GetItems() []*models.ClipboardItem {
	return hs.history.GetItems()
}

func (hs *HistoryService) GetItem(index int) *models.ClipboardItem {
	return hs.history.GetItem(index)
}

func (hs *HistoryService) Clear() error {
	// Clean up image files first
	imagePaths := make([]string, 0)
	for _, item := range hs.history.GetItems() {
		if item.Type == models.ClipImage && item.FilePath != "" {
			imagePaths = append(imagePaths, item.FilePath)
		}
	}
	hs.fileService.CleanupImageFiles(imagePaths)
	
	// Clear history
	hs.history.Clear()
	
	// Delete files
	hs.fileService.DeleteHistoryFile()
	hs.fileService.DeleteImageDirectory()
	
	return nil
}

func (hs *HistoryService) MaintainLimit() {
	items := hs.history.GetItems()
	if len(items) > hs.history.MaxItems {
		// Clean up excess images
		excessItems := items[hs.history.MaxItems:]
		imagePaths := make([]string, 0)
		for _, item := range excessItems {
			if item.Type == models.ClipImage && item.FilePath != "" {
				imagePaths = append(imagePaths, item.FilePath)
			}
		}
		hs.fileService.CleanupImageFiles(imagePaths)
		
		// Trim history
		hs.history.Items = items[:hs.history.MaxItems]
		hs.SaveToFile()
	}
}