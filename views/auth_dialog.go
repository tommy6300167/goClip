package views

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

type AuthDialog struct {
	window       fyne.Window
	keyEntry     *widget.Entry
	onAuth       func(string) bool
	onCancel     func()
	dialogWidget dialog.Dialog
}

func NewAuthDialog(window fyne.Window) *AuthDialog {
	return &AuthDialog{
		window: window,
	}
}

func (ad *AuthDialog) SetOnAuth(callback func(string) bool) {
	ad.onAuth = callback
}

func (ad *AuthDialog) SetOnCancel(callback func()) {
	ad.onCancel = callback
}

func (ad *AuthDialog) Show() {
	ad.keyEntry = widget.NewPasswordEntry()
	ad.keyEntry.SetPlaceHolder("請輸入金鑰...")

	title := widget.NewLabel("首次使用需要驗證")
	title.Alignment = fyne.TextAlignCenter

	description := widget.NewLabel("請輸入您的金鑰以開始使用超吉貼")
	description.Alignment = fyne.TextAlignCenter

	confirmButton := widget.NewButton("確認", func() {
		key := ad.keyEntry.Text
		if ad.onAuth != nil && ad.onAuth(key) {
			ad.dialogWidget.Hide()
		} else {
			// 顯示錯誤提示
			ad.showError("金鑰錯誤，請重新輸入")
		}
	})

	cancelButton := widget.NewButton("取消", func() {
		if ad.onCancel != nil {
			ad.onCancel()
		}
		ad.dialogWidget.Hide()
	})

	// 支持Enter鍵提交
	ad.keyEntry.OnSubmitted = func(text string) {
		confirmButton.OnTapped()
	}

	buttons := container.NewHBox(cancelButton, confirmButton)
	content := container.NewVBox(
		title,
		widget.NewSeparator(),
		description,
		ad.keyEntry,
		widget.NewSeparator(),
		buttons,
	)

	ad.dialogWidget = dialog.NewCustom("金鑰驗證", "", content, ad.window)
	ad.dialogWidget.Resize(fyne.NewSize(350, 200))
	ad.dialogWidget.Show()

	// 自動聚焦到輸入框
	ad.window.Canvas().Focus(ad.keyEntry)
}

func (ad *AuthDialog) showError(message string) {
	errorDialog := dialog.NewInformation("錯誤", message, ad.window)
	errorDialog.Show()
}

func (ad *AuthDialog) Hide() {
	if ad.dialogWidget != nil {
		ad.dialogWidget.Hide()
	}
}