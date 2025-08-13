package views

import (
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"

	"clipmini/models"
)

type DetailView struct {
	container   *fyne.Container
	textEntry   *widget.Entry
	imageCard   *widget.Card
	currentItem *models.ClipboardItem
}

func NewDetailView() *DetailView {
	dv := &DetailView{
		textEntry: widget.NewMultiLineEntry(),
	}
	
	dv.textEntry.SetPlaceHolder("左側選一筆來預覽內容")
	dv.textEntry.Wrapping = fyne.TextWrapWord
	
	dv.container = container.NewBorder(nil, nil, nil, nil, dv.textEntry)
	
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
	dv.textEntry.SetText(item.Content)
	dv.textEntry.Show()
	dv.imageCard = nil
	dv.container.Objects[0] = dv.textEntry
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
		dv.container.Objects[0] = dv.imageCard
		dv.container.Refresh()
	} else {
		dv.textEntry.SetText("[IMAGE ERROR] " + item.FilePath)
		dv.textEntry.Show()
		dv.imageCard = nil
		dv.container.Objects[0] = dv.textEntry
		dv.container.Refresh()
	}
}

func (dv *DetailView) ShowError(message string) {
	dv.textEntry.SetText("[ERROR] " + message)
	dv.textEntry.Show()
	dv.imageCard = nil
	dv.container.Objects[0] = dv.textEntry
	dv.container.Refresh()
}

func (dv *DetailView) Clear() {
	dv.textEntry.SetText("")
	dv.textEntry.Show()
	dv.imageCard = nil
	dv.currentItem = nil
	dv.container.Objects[0] = dv.textEntry
	dv.container.Refresh()
}

func (dv *DetailView) GetCurrentText() string {
	return strings.TrimSpace(dv.textEntry.Text)
}

func (dv *DetailView) GetCurrentItem() *models.ClipboardItem {
	return dv.currentItem
}