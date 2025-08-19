package views

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"

	"clipmini/models"
)

type DetailView struct {
	container        *fyne.Container
	textEntry        *CustomEntry
	richText         *widget.RichText
	imageCard        *widget.Card
	saveButton       *widget.Button
	editButton       *widget.Button
	currentItem      *models.ClipboardItem
	onSave           func(string)
	onStatusUpdate   func(string)
	onEditingMode    func()
	onNormalMode     func()
	originalText     string
	isEditing        bool
	saveShortcut     *fyne.Shortcut
}

func NewDetailView() *DetailView {
	dv := &DetailView{
		textEntry: NewCustomEntry(),
		richText:  widget.NewRichText(),
		isEditing: false,
	}

	dv.textEntry.SetPlaceHolder("No item selected")
	dv.textEntry.Hide()

	dv.richText.Wrapping = fyne.TextWrapOff
	dv.richText.Scroll = container.ScrollBoth

	dv.saveButton = widget.NewButtonWithIcon("Save", theme.DocumentSaveIcon(), func() {
		if dv.onSave != nil {
			dv.onSave(dv.textEntry.Text)
		}
		dv.isEditing = false
		dv.showAsRichText()
		if dv.onNormalMode != nil {
			dv.onNormalMode()
		}
	})
	dv.saveButton.Importance = widget.HighImportance
	dv.saveButton.Hide()

	dv.editButton = widget.NewButtonWithIcon("Edit", theme.DocumentCreateIcon(), func() {
		if dv.currentItem == nil || dv.currentItem.Type != models.ClipText {
			return
		}
		dv.isEditing = true
		dv.showAsTextEntry()
		dv.saveButton.Show()
		if dv.onEditingMode != nil {
			dv.onEditingMode()
		}
	})

	dv.textEntry.OnChanged = func(text string) {
		if dv.isEditing {
			dv.saveButton.Show()
		}
	}

	dv.textEntry.SetPreventEnterNewline(true)

	dv.editButton.Disable()

	buttonContainer := container.NewHBox(
		container.NewPadded(dv.editButton),
		container.NewPadded(dv.saveButton),
	)
	dv.container = container.NewBorder(nil, buttonContainer, nil, nil, dv.richText)

	return dv
}

func (dv *DetailView) GetWidget() *fyne.Container {
	return dv.container
}

func (dv *DetailView) ShowItem(item *models.ClipboardItem) {
	dv.currentItem = item

	if item.Type == models.ClipImage {
		dv.showImage(item)
	} else {
		dv.showText(item)
	}
}

func (dv *DetailView) showText(item *models.ClipboardItem) {
	dv.originalText = item.Content
	dv.textEntry.SetText(item.Content)
	dv.richText.ParseMarkdown("```\n" + item.Content + "\n```")
	dv.saveButton.Hide()
	dv.editButton.Show()
	dv.editButton.Enable()
	dv.imageCard = nil
	dv.isEditing = false
	dv.showAsRichText()
}

func (dv *DetailView) showAsRichText() {
	dv.textEntry.Hide()
	dv.richText.Show()
	dv.container.Objects[0] = dv.richText
	dv.saveButton.Hide()
	dv.editButton.Show()
	dv.container.Refresh()
	if dv.onNormalMode != nil {
		dv.onNormalMode()
	}
}

func (dv *DetailView) showAsTextEntry() {
	dv.richText.Hide()
	dv.textEntry.Show()
	dv.container.Objects[0] = dv.textEntry
	dv.editButton.Hide()
	dv.container.Refresh()
}

func (dv *DetailView) showImage(item *models.ClipboardItem) {
	if res, err := fyne.LoadResourceFromPath(item.FilePath); err == nil {
		img := canvas.NewImageFromResource(res)
		img.FillMode = canvas.ImageFillContain
		img.SetMinSize(fyne.NewSize(360, 260))

		timestamp := item.Timestamp.Format("2006-01-02 15:04:05")
		dv.imageCard = widget.NewCard("Image", timestamp, img)

		dv.textEntry.Hide()
		dv.richText.Hide()
		dv.editButton.Hide()
		dv.container.Objects[0] = dv.imageCard
		dv.container.Refresh()
	} else {
		dv.textEntry.SetText("[IMAGE ERROR] " + item.FilePath)
		dv.richText.Hide()
		dv.editButton.Hide()
		dv.textEntry.Show()
		dv.imageCard = nil
		dv.container.Objects[0] = dv.textEntry
		dv.container.Refresh()
	}
}

func (dv *DetailView) ShowError(message string) {
	dv.textEntry.SetText("[ERROR] " + message)
	dv.richText.Hide()
	dv.editButton.Hide()
	dv.textEntry.Show()
	dv.imageCard = nil
	dv.container.Objects[0] = dv.textEntry
	dv.container.Refresh()
}

func (dv *DetailView) Clear() {
	dv.textEntry.SetText("")
	dv.richText.ParseMarkdown("")
	dv.richText.Hide()
	dv.textEntry.Show()
	dv.saveButton.Hide()
	dv.editButton.Hide()
	dv.editButton.Disable()
	dv.imageCard = nil
	dv.currentItem = nil
	dv.originalText = ""
	dv.isEditing = false
	dv.container.Objects[0] = dv.textEntry
	dv.container.Refresh()
}

func (dv *DetailView) GetCurrentText() string {
	if dv.isEditing {
		return dv.textEntry.Text
	}
	return dv.originalText
}

func (dv *DetailView) GetCurrentItem() *models.ClipboardItem {
	return dv.currentItem
}

func (dv *DetailView) SetOnSave(callback func(string)) {
	dv.onSave = callback
}

func (dv *DetailView) SetOnStatusUpdate(callback func(string)) {
	dv.onStatusUpdate = callback
}

func (dv *DetailView) SetOnEditingMode(callback func()) {
	dv.onEditingMode = callback
}

func (dv *DetailView) SetOnNormalMode(callback func()) {
	dv.onNormalMode = callback
}

func (dv *DetailView) SetupShortcuts(window fyne.Window) {
	println("DEBUG: Setting up shortcuts using menu approach...")
	
	// 嘗試通過 Menu 系統設定快速鍵
	// 這是唯一可靠的方式，因為許多快速鍵被作業系統攔截
	
	// 建立一個不可見的選單項目來綁定快速鍵
	saveMenuItem := fyne.NewMenuItem("Save", func() {
		println("DEBUG: Menu save triggered!")
		if dv.isEditing {
			dv.saveButton.OnTapped()
		}
	})
	
	// 設定快速鍵
	saveMenuItem.Shortcut = &desktop.CustomShortcut{
		KeyName:  fyne.KeyS,
		Modifier: desktop.SuperModifier, // Cmd 鍵
	}
	
	// 建立選單但不顯示，只是為了註冊快速鍵
	menu := fyne.NewMenu("Edit", saveMenuItem)
	mainMenu := fyne.NewMainMenu(menu)
	window.SetMainMenu(mainMenu)
	
	println("DEBUG: Shortcuts setup via menu complete")
}
