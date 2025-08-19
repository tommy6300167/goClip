package views

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"

	"clipmini/models"
)

type DetailView struct {
	container      *fyne.Container
	textEntry      *CustomEntry
	richText       *widget.RichText
	imageCard      *widget.Card
	saveButton     *widget.Button
	editButton     *widget.Button
	currentItem    *models.ClipboardItem
	onSave         func(string)
	onStatusUpdate func(string)
	originalText   string
	isEditing      bool
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
		if dv.onStatusUpdate != nil {
			dv.onStatusUpdate("Press ⌘C Auto saving")
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
		if dv.onStatusUpdate != nil {
			dv.onStatusUpdate("Use Shift + Enter to add a line break.")
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
	if dv.onStatusUpdate != nil {
		dv.onStatusUpdate("Press ⌘C Auto saving")
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
