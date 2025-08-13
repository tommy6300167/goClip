package main

import (
	_ "embed"
	"log"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"

	"clipmini/controllers"
	"clipmini/models"
	"clipmini/views"
)

var notoTC []byte

func main() {
	// 初始化配置
	config := models.NewAppConfig()
	
	// 初始化控制器
	clipboardController := controllers.NewClipboardController(config)
	if err := clipboardController.Initialize(); err != nil {
		log.Fatal("Failed to initialize clipboard controller:", err)
	}
	
	// 初始化 Fyne 應用
	fyneApp := app.NewWithID("superClip")
	
	// 設置字型主題
	if len(notoTC) > 0 {
		fontResource := fyne.NewStaticResource("NotoSansTC-VariableFont_wght.ttf", notoTC)
		customTheme := views.NewCustomTheme(fontResource)
		fyneApp.Settings().SetTheme(customTheme)
	}
	
	// 創建主窗口
	window := fyneApp.NewWindow("超吉貼")
	window.Resize(fyne.NewSize(820, 520))
	
	// 初始化 UI 組件
	mainView := views.NewMainView(clipboardController, config)
	mainView.Initialize(window)
	
	// 設置窗口內容
	window.SetContent(mainView.GetContent())
	
	// 開始剪貼簿監控
	stopChannel := make(chan struct{})
	go func() {
		ticker := time.NewTicker(time.Duration(config.PollingInterval) * time.Millisecond)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				if newItem := clipboardController.PollClipboard(); newItem != nil {
					mainView.OnNewClipboardItem(newItem)
				}
			case <-stopChannel:
				return
			}
		}
	}()
	
	// 設置窗口關閉處理
	window.SetCloseIntercept(func() {
		close(stopChannel)
		window.Close()
	})
	
	// 運行應用
	window.ShowAndRun()
}