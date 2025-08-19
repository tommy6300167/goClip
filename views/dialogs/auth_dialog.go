package dialogs

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
	ad.keyEntry.SetPlaceHolder("Press")

	title := widget.NewLabel("Verification is required on first use")
	title.Alignment = fyne.TextAlignCenter

	description := widget.NewLabel("Please enter your key to start using SuperClip")
	description.Alignment = fyne.TextAlignCenter

	confirmButton := widget.NewButton("Confirm", func() {
		key := ad.keyEntry.Text
		if ad.onAuth != nil && ad.onAuth(key) {
			ad.dialogWidget.Hide()
		} else {
			ad.showError("nvalid key. Please try again.")
		}
	})

	cancelButton := widget.NewButton("Cancel", func() {
		if ad.onCancel != nil {
			ad.onCancel()
		}
		ad.dialogWidget.Hide()
	})

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

	ad.dialogWidget = dialog.NewCustomWithoutButtons("Key Verification", content, ad.window)
	ad.dialogWidget.Resize(fyne.NewSize(350, 200))
	ad.dialogWidget.Show()

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
