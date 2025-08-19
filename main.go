package main

import (
	_ "embed"
	"log"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"

	"clipmini/controllers"
	"clipmini/models"
	"clipmini/services"
	"clipmini/views"
)

var notoTC []byte

func main() {
	config := models.NewAppConfig()
	authService := services.NewAuthService(config)

	fyneApp := app.NewWithID("superClip")

	if len(notoTC) > 0 {
		fontResource := fyne.NewStaticResource("NotoSansTC-VariableFont_wght.ttf", notoTC)
		customTheme := views.NewCustomTheme(fontResource)
		fyneApp.Settings().SetTheme(customTheme)
	}

	window := fyneApp.NewWindow("超吉貼")
	window.Resize(fyne.NewSize(820, 520))
	window.CenterOnScreen()

	// 檢查認證狀態
	needsAuth := authService.IsFirstLaunch() || !authService.IsAuthValid()
	
	if needsAuth {
		authDialog := views.NewAuthDialog(window)
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
		// 已有有效認證，直接進入主程式
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

	stopChannel := make(chan struct{})
	authService := services.NewAuthService(config)
	
	go func() {
		clipboardTicker := time.NewTicker(time.Duration(config.PollingInterval) * time.Millisecond)
		authTicker := time.NewTicker(1 * time.Hour) // 每小時檢查一次認證
		defer clipboardTicker.Stop()
		defer authTicker.Stop()

		for {
			select {
			case <-clipboardTicker.C:
				if newItem := clipboardController.PollClipboard(); newItem != nil {
					mainView.OnNewClipboardItem(newItem)
				}
			case <-authTicker.C:
				// 定期檢查認證是否過期
				if !authService.IsAuthValid() {
					// 認證過期，需要重新驗證
					go func() {
						window.Close()
						// 重新啟動應用程式的認證流程
					}()
					return
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
}
