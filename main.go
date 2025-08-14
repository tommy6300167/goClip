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
	config := models.NewAppConfig()

	clipboardController := controllers.NewClipboardController(config)
	if err := clipboardController.Initialize(); err != nil {
		log.Fatal("Failed to initialize clipboard controller:", err)
	}

	fyneApp := app.NewWithID("superClip")

	if len(notoTC) > 0 {
		fontResource := fyne.NewStaticResource("NotoSansTC-VariableFont_wght.ttf", notoTC)
		customTheme := views.NewCustomTheme(fontResource)
		fyneApp.Settings().SetTheme(customTheme)
	}

	window := fyneApp.NewWindow("超吉貼")
	window.Resize(fyne.NewSize(820, 520))

	mainView := views.NewMainView(clipboardController, config)
	mainView.Initialize(window)

	window.SetContent(mainView.GetContent())

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

	window.SetCloseIntercept(func() {
		close(stopChannel)
		window.Close()
	})

	window.ShowAndRun()
}
