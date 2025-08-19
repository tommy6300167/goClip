package controllers

import (
	"fmt"
	"strings"
	"time"

	"clipmini/models"
	"clipmini/services"
	"clipmini/utils"
)

type ClipboardController struct {
	clipboardService *services.ClipboardService
	historyService   *services.HistoryService
	fileService      *services.FileService
	config           *models.AppConfig
	lastText         string
	lastImgHash      string
}

func NewClipboardController(config *models.AppConfig) *ClipboardController {
	return &ClipboardController{
		clipboardService: services.NewClipboardService(),
		historyService:   services.NewHistoryService(config),
		fileService:      services.NewFileService(config),
		config:           config,
	}
}

func (cc *ClipboardController) Initialize() error {
	if err := cc.historyService.LoadFromFile(); err != nil {
		return err
	}

	if cc.clipboardService.HasImageInClipboard() {
		if b, err := cc.clipboardService.ReadClipboardImage(); err == nil && len(b) > 0 {
			cc.lastImgHash = cc.clipboardService.GetImageHash(b)
		}
	} else {
		if txt, err := cc.clipboardService.ReadClipboardText(); err == nil {
			cc.lastText = txt
		}
	}

	return nil
}

func (cc *ClipboardController) PollClipboard() *models.ClipboardItem {
	loc := utils.GetTaipeiLocation()

	if cc.clipboardService.HasImageInClipboard() {
		if b, err := cc.clipboardService.ReadClipboardImage(); err == nil && len(b) > 0 {
			currentHash := cc.clipboardService.GetImageHash(b)
			if currentHash != cc.lastImgHash {
				cc.lastImgHash = currentHash
				cc.lastText = ""
				timestamp := utils.FormatTimestamp(time.Now(), loc)

				if path, err := cc.fileService.SaveImage(b, timestamp); err == nil {
					item := models.NewImageItem(path)
					item.Timestamp = time.Now()

					if err := cc.historyService.AddItem(item); err == nil {
						cc.historyService.MaintainLimit()
						return item
					}
				}
			}
		}
	} else {
		cc.lastImgHash = ""
	}

	if txt, err := cc.clipboardService.ReadClipboardText(); err == nil {
		if txt != "" && txt != cc.lastText && !cc.isSensitiveData(txt) {
			cc.lastText = txt
			cc.lastImgHash = ""
			item := models.NewTextItem(txt)

			if err := cc.historyService.AddItem(item); err == nil {
				cc.historyService.MaintainLimit()
				return item
			}
		}
	}

	return nil
}

func (cc *ClipboardController) CopyItemToClipboard(item *models.ClipboardItem) error {
	if item.Type == models.ClipImage {
		if cc.fileService.ImageExists(item.FilePath) {
			return cc.clipboardService.CopyImageToClipboard(item.FilePath)
		}
		return nil
	}
	return cc.clipboardService.CopyTextToClipboard(item.Content)
}

func (cc *ClipboardController) GetHistoryItems() []*models.ClipboardItem {
	return cc.historyService.GetItems()
}

func (cc *ClipboardController) GetHistoryItem(index int) *models.ClipboardItem {
	return cc.historyService.GetItem(index)
}

func (cc *ClipboardController) RemoveHistoryItem(index int) error {
	return cc.historyService.RemoveItem(index)
}

func (cc *ClipboardController) InsertHistoryItemAt(index int, item *models.ClipboardItem) error {
	return cc.historyService.InsertItemAt(index, item)
}

func (cc *ClipboardController) UpdateHistoryItem(index int, newContent string) error {
	return cc.historyService.UpdateItem(index, newContent)
}

func (cc *ClipboardController) ClearHistory() error {
	cc.clipboardService.ClearSystemClipboard()
	cc.lastText = ""
	cc.lastImgHash = ""
	return cc.historyService.Clear()
}

func (cc *ClipboardController) ExportHistory() string {
	items := cc.historyService.GetItems()
	lines := make([]string, 0, len(items))

	for i := len(items) - 1; i >= 0; i-- {
		item := items[i]
		timestamp := utils.FormatTimestamp(item.Timestamp, utils.GetTaipeiLocation())

		if item.Type == models.ClipImage {
			lines = append(lines, timestamp+"\t"+item.FilePath+"\t"+item.Type.String())
		} else {
			if !cc.isSensitiveData(item.Content) {
				lines = append(lines, timestamp+"\t"+item.Content)
			} else {
				lines = append(lines, timestamp+"\t[敏感資料已隱藏]")
			}
		}
	}

	return strings.Join(lines, "\n")
}

func (cc *ClipboardController) Count() int {
	return len(cc.historyService.GetItems())
}

func (cc *ClipboardController) GetHistoryService() *services.HistoryService {
	return cc.historyService
}

func (cc *ClipboardController) isSensitiveData(text string) bool {
	text = strings.TrimSpace(strings.ToLower(text))

	sensitivePatterns := []string{}

	for _, pattern := range sensitivePatterns {
		if strings.Contains(text, pattern) {
			fmt.Printf("[敏感資料檢測] 發現關鍵字 '%s'，內容: %s\n", pattern, text[:min(50, len(text))])
			return true
		}
	}

	// 檢查是否看起來像API金鑰或token（更精確的檢查）
	if len(text) > 30 && strings.Count(text, "-") >= 3 && strings.Count(text, "_") >= 2 {
		// 更嚴格的檢查：需要更長且有更多特殊字符組合
		fmt.Printf("[敏感資料檢測] 疑似API金鑰格式，長度: %d，內容: %s\n", len(text), text[:min(50, len(text))])
		return true
	}

	return false
}

// min 返回兩個整數中的較小值
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
