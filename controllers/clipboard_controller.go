package controllers

import (
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
			cc.lastText = strings.TrimSpace(txt)
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
		normalized := strings.TrimSpace(txt)
		if normalized != "" && normalized != cc.lastText && !cc.isSensitiveData(normalized) {
			cc.lastText = normalized
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
	var lines []string

	for i := len(items) - 1; i >= 0; i-- {
		item := items[i]
		timestamp := utils.FormatTimestamp(item.Timestamp, utils.GetTaipeiLocation())

		if item.Type == models.ClipImage {
			lines = append(lines, timestamp+"\t"+item.FilePath+"\t"+item.Type.String())
		} else {
			// 過濾敏感資料
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

func (cc *ClipboardController) isSensitiveData(text string) bool {
	text = strings.TrimSpace(strings.ToLower(text))
	
	// 檢查常見敏感資料模式
	sensitivePatterns := []string{
		"password",
		"passwd",
		"token",
		"secret",
		"api_key",
		"apikey",
		"private_key",
		"privatekey",
		"credential",
		"密碼",
		"金鑰",
		"令牌",
		"憑證",
	}
	
	for _, pattern := range sensitivePatterns {
		if strings.Contains(text, pattern) {
			return true
		}
	}
	
	// 檢查是否看起來像API金鑰或token（含有特定格式）
	if len(text) > 20 && (strings.Contains(text, "-") || strings.Contains(text, "_")) {
		// 簡單的啟發式檢查：長度超過20且包含特殊字符可能是金鑰
		return true
	}
	
	return false
}
