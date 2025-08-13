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
	
	// 首先檢查圖片
	if cc.clipboardService.HasImageInClipboard() {
		if b, err := cc.clipboardService.ReadClipboardImage(); err == nil && len(b) > 0 {
			currentHash := cc.clipboardService.GetImageHash(b)
			if currentHash != cc.lastImgHash {
				cc.lastImgHash = currentHash
				// 重置文字追蹤，因為現在是圖片
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
		// 沒有圖片時，重置圖片哈希
		cc.lastImgHash = ""
	}

	// 然後檢查文字（無論是否有圖片都要檢查）
	if txt, err := cc.clipboardService.ReadClipboardText(); err == nil {
		normalized := strings.TrimSpace(txt)
		if normalized != "" && normalized != cc.lastText {
			cc.lastText = normalized
			// 重置圖片追蹤，因為現在是文字
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

func (cc *ClipboardController) ClearHistory() error {
	cc.clipboardService.ClearSystemClipboard()
	cc.lastText = ""
	cc.lastImgHash = ""
	return cc.historyService.Clear()
}

func (cc *ClipboardController) ExportHistory() string {
	items := cc.historyService.GetItems()
	lines := make([]string, len(items))
	
	for i := len(items) - 1; i >= 0; i-- {
		item := items[i]
		timestamp := utils.FormatTimestamp(item.Timestamp, utils.GetTaipeiLocation())
		
		if item.Type == models.ClipImage {
			lines[len(items)-1-i] = timestamp + "\t" + item.FilePath + "\t" + item.Type.String()
		} else {
			lines[len(items)-1-i] = timestamp + "\t" + item.Content
		}
	}
	
	return strings.Join(lines, "\n")
}