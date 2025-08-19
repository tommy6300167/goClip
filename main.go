package main

import (
	"context"
	_ "embed"
	"log"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"

	"clipmini/controllers"
	"clipmini/models"
	"clipmini/services"
	"clipmini/views"
	"clipmini/views/dialogs"
	"clipmini/views/theme"
)

var notoTC []byte

func main() {
	config := models.NewAppConfig()
	authService := services.NewAuthService(config)

	fyneApp := app.NewWithID("superClip")

	if len(notoTC) > 0 {
		fontResource := fyne.NewStaticResource("NotoSansTC-VariableFont_wght.ttf", notoTC)
		customTheme := theme.NewCustomTheme(fontResource)
		fyneApp.Settings().SetTheme(customTheme)
	}

	window := fyneApp.NewWindow("超吉貼")
	window.Resize(fyne.NewSize(float32(config.WindowWidth), float32(config.WindowHeight)))
	window.CenterOnScreen()

	needsAuth := authService.IsFirstLaunch() || !authService.IsAuthValid()

	if needsAuth {
		authDialog := dialogs.NewAuthDialog(window)
		authenticated := make(chan bool)

		authDialog.SetOnAuth(func(key string) bool {
			if authService.ValidateKey(key) {
				if err := authService.SaveAuthConfig(); err != nil {
					log.Printf("Failed to save auth config")
					return false
				}
				authenticated <- true
				return true
			}
			return false
		})

		authDialog.SetOnCancel(func() {
			authenticated <- false
		})

		window.Show()
		authDialog.Show()

		go func() {
			if !<-authenticated {
				fyneApp.Quit()
				return
			}

			initializeMainApp(config, window)
		}()

		fyneApp.Run()
	} else {
		initializeMainApp(config, window)
		window.ShowAndRun()
	}
}

func initializeMainApp(config *models.AppConfig, window fyne.Window) {
	clipboardController := controllers.NewClipboardController(config)
	if err := clipboardController.Initialize(); err != nil {
		log.Fatal("Failed to initialize clipboard controller:", err)
	}

	mainView := views.NewMainView(clipboardController, config)
	mainView.Initialize(window)

	window.SetContent(mainView.GetContent())

	ctx, cancel := context.WithCancel(context.Background())
	authService := services.NewAuthService(config)

	clipboardService := services.NewClipboardService()
	clipboardService.StartClipboardWatcher(ctx, func() {
		if newItem := clipboardController.PollClipboard(); newItem != nil {
			mainView.OnNewClipboardItem(newItem)
		}
	})

	go func() {
		authTicker := time.NewTicker(1 * time.Hour)
		defer authTicker.Stop()

		for {
			select {
			case <-authTicker.C:
				if !authService.IsAuthValid() {
					go func() {
						cancel()
						window.Close()
					}()
					return
				}
			case <-ctx.Done():
				return
			}
		}
	}()

	window.SetCloseIntercept(func() {
		cancel()
		if historyService := clipboardController.GetHistoryService(); historyService != nil {
			historyService.ForceFlush()
		}
		window.Close()
	})
}
