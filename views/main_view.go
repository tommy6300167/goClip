package views

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
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
	selectedIndex       int
}

func NewMainView(clipboardController *controllers.ClipboardController, config *models.AppConfig) *MainView {

	mv := &MainView{
		clipboardController: clipboardController,
		config:              config,
	}

	mv.statusLabel = widget.NewLabel("⌘C Auto saving")

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
	mv.listView.SetOnDelete(mv.onDeleteItem)
	mv.detailView.SetOnSave(mv.onSaveItem)
	mv.detailView.SetOnStatusUpdate(mv.updateStatus)

	mv.toolbar.SetOnCopy(mv.onCopyToClipboard)
	mv.toolbar.SetOnClear(mv.onClearHistory)
	mv.toolbar.SetOnExport(mv.onExportHistory)
	mv.toolbar.SetOnSearch(mv.onSearchChanged)
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
	mainContent := container.NewBorder(mv.toolbar.GetWidget(), nil, nil, nil, split)
	bg := canvas.NewRectangle(color.NRGBA{R: 33, G: 39, B: 44, A: 255})

	mv.content = container.NewStack(bg, mainContent)
}

func (mv *MainView) GetContent() *fyne.Container {
	return mv.content
}

func (mv *MainView) onItemSelected(index int) {
	mv.selectedIndex = index
	item := mv.clipboardController.GetHistoryItem(index)
	if item == nil {
		return
	}
	mv.currentSelectedItem = item
	mv.detailView.ShowItem(item)
}

func (mv *MainView) onDeleteItem(index int) {
	if err := mv.clipboardController.RemoveHistoryItem(index); err != nil {
		mv.updateStatus("刪除失敗: " + err.Error())
		return
	}
	mv.listView.RemoveItemAt(index)

	if mv.selectedIndex == index {
		if index >= mv.clipboardController.Count() {
			mv.selectedIndex = mv.clipboardController.Count() - 1
		} else {
			mv.selectedIndex = index
		}
	} else if mv.selectedIndex > index {
		mv.selectedIndex--
	}

	if mv.selectedIndex >= 0 {
		item := mv.clipboardController.GetHistoryItem(mv.selectedIndex)
		mv.currentSelectedItem = item
		mv.detailView.ShowItem(item)
		mv.listView.GetWidget().Select(mv.selectedIndex)
	} else {
		mv.currentSelectedItem = nil
		mv.detailView.Clear()
		mv.listView.GetWidget().UnselectAll()
	}

	mv.updateStatus("項目已刪除")
}

func (mv *MainView) onSaveItem(newContent string) {
	if mv.currentSelectedItem == nil {
		mv.updateStatus("沒有選中的項目")
		return
	}

	items := mv.clipboardController.GetHistoryItems()
	selectedIndex := -1
	for i, item := range items {
		if item == mv.currentSelectedItem {
			selectedIndex = i
			break
		}
	}

	if selectedIndex == -1 {
		mv.updateStatus("找不到選中的項目")
		return
	}

	err := mv.clipboardController.UpdateHistoryItem(selectedIndex, newContent)
	if err != nil {
		mv.updateStatus("保存失敗: " + err.Error())
		return
	}

	mv.currentSelectedItem.Content = newContent

	mv.listView.LoadFromHistory(mv.clipboardController.GetHistoryItems())

	mv.listView.SelectFirst()
	if selectedIndex < len(items) {
		mv.listView.GetWidget().Select(selectedIndex)
	}

	mv.updateStatus("已保存修改")
}

func (mv *MainView) onCopyToClipboard() error {
	if mv.currentSelectedItem == nil {
		return nil
	}

	if mv.currentSelectedItem.Type == models.ClipImage {
		return mv.clipboardController.CopyItemToClipboard(mv.currentSelectedItem)
	} else {
		// 優先使用 DetailView 的當前文字，如果為空則使用原始內容
		text := mv.detailView.GetCurrentText()
		if text == "" {
			text = mv.currentSelectedItem.Content
		}

		return mv.clipboardController.CopyItemToClipboard(&models.ClipboardItem{
			Type:    models.ClipText,
			Content: text,
		})
	}
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

func (mv *MainView) onSearchChanged(query string) {
	mv.listView.SetSearchQuery(query)
	if query == "" {
		mv.updateStatus("顯示所有項目")
	} else {
		mv.updateStatus("搜尋: " + query)
	}
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
