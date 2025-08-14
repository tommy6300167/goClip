package views

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"

	"clipmini/controllers"
	"clipmini/models"
)

type MainView struct {
	content             *fyne.Container
	listView            *ListView
	detailView          *DetailView
	toolbar             *Toolbar
	statusLabel         *widget.Label
	clipboardController *controllers.ClipboardController
	config              *models.AppConfig
	currentSelectedItem *models.ClipboardItem
}

func NewMainView(clipboardController *controllers.ClipboardController, config *models.AppConfig) *MainView {
	mv := &MainView{
		clipboardController: clipboardController,
		config:              config,
	}
	
	mv.statusLabel = widget.NewLabel("監看中… ⌘C 內容會自動記錄（文字＋圖片）")
	
	return mv
}

func (mv *MainView) Initialize(window fyne.Window) {
	mv.listView = NewListView(mv.config)
	mv.detailView = NewDetailView()
	mv.toolbar = NewToolbar(window, mv.clipboardController)
	
	mv.setupEventHandlers()
	mv.loadInitialData()
	mv.buildLayout()
}

func (mv *MainView) setupEventHandlers() {
	mv.listView.SetOnSelected(mv.onItemSelected)
	
	mv.toolbar.SetOnCopy(mv.onCopyToClipboard)
	mv.toolbar.SetOnClear(mv.onClearHistory)
	mv.toolbar.SetOnExport(mv.onExportHistory)
	mv.toolbar.SetOnStatusUpdate(mv.updateStatus)
}

func (mv *MainView) loadInitialData() {
	items := mv.clipboardController.GetHistoryItems()
	mv.listView.LoadFromHistory(items)
}

func (mv *MainView) buildLayout() {
	left := container.NewBorder(nil, nil, nil, nil, mv.listView.GetWidget())
	right := container.NewBorder(nil, mv.statusLabel, nil, nil, mv.detailView.GetWidget())
	
	split := container.NewHSplit(left, right)
	split.SetOffset(0.35)
	
	mv.content = container.NewBorder(mv.toolbar.GetWidget(), nil, nil, nil, split)
}

func (mv *MainView) GetContent() *fyne.Container {
	return mv.content
}

func (mv *MainView) onItemSelected(index int) {
	item := mv.clipboardController.GetHistoryItem(index)
	if item == nil {
		return
	}
	
	mv.currentSelectedItem = item
	mv.detailView.ShowItem(item)
}

func (mv *MainView) onCopyToClipboard() error {
	if mv.currentSelectedItem == nil {
		return nil
	}
	
	if mv.currentSelectedItem.Type == models.ClipImage {
		return mv.clipboardController.CopyItemToClipboard(mv.currentSelectedItem)
	} else {
		text := mv.detailView.GetCurrentText()
		if text != "" {
			return mv.clipboardController.CopyItemToClipboard(&models.ClipboardItem{
				Type:    models.ClipText,
				Content: text,
			})
		}
	}
	return nil
}

func (mv *MainView) onClearHistory() error {
	err := mv.clipboardController.ClearHistory()
	if err == nil {
		mv.listView.Clear()
		mv.detailView.Clear()
		mv.currentSelectedItem = nil
	}
	return err
}

func (mv *MainView) onExportHistory() string {
	return mv.clipboardController.ExportHistory()
}

func (mv *MainView) updateStatus(message string) {
	mv.statusLabel.SetText(message)
}

func (mv *MainView) OnNewClipboardItem(item *models.ClipboardItem) {
	fyne.Do(func() {
		mv.listView.PrependItem(item)
		
		if item.Type == models.ClipImage {
			mv.updateStatus("圖片已記錄")
		} else {
			mv.updateStatus("文字已記錄")
		}
	})
}