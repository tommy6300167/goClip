package views

import (
	"image/color"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"

	"clipmini/controllers"
	"clipmini/models"
	"clipmini/views/components/displays"
	"clipmini/views/components/feedback"
	"clipmini/views/dialogs"
)

type MainView struct {
	content             *fyne.Container
	listView            *ListView
	detailView          *DetailView
	toolbar             *Toolbar
	statusBar           *displays.StatusBar
	undoToast           *feedback.UndoToast
	settingsBtn         *widget.Button
	settingsDialog      *dialogs.SettingsDialog
	clipboardController *controllers.ClipboardController
	config              *models.AppConfig
	currentSelectedItem *models.ClipboardItem
	selectedIndex       int
	pendingDeleteItem   *models.ClipboardItem
	pendingDeleteIndex  int
	window              fyne.Window
}

func NewMainView(clipboardController *controllers.ClipboardController, config *models.AppConfig) *MainView {

	mv := &MainView{
		clipboardController: clipboardController,
		config:              config,
		pendingDeleteIndex:  -1,
	}

	mv.statusBar = displays.NewStatusBar("Press ⌘C Auto saving")
	mv.undoToast = feedback.NewUndoToast()

	return mv
}

func (mv *MainView) Initialize(window fyne.Window) {
	mv.window = window
	mv.listView = NewListView(mv.config)
	mv.detailView = NewDetailView()
	mv.toolbar = NewToolbar(window, mv.clipboardController)
	
	// 創建設置按鈕
	mv.settingsBtn = widget.NewButtonWithIcon("", theme.SettingsIcon(), func() {
		mv.showSettings()
	})
	mv.settingsBtn.Importance = widget.LowImportance
	mv.settingsBtn.Resize(fyne.NewSize(24, 24))
	
	// 創建設置對話框
	mv.settingsDialog = dialogs.NewSettingsDialog(window, mv.config)
	mv.settingsDialog.SetOnSave(mv.onSettingsSave)

	mv.setupEventHandlers()
	mv.detailView.SetupShortcuts(window)
	mv.loadInitialData()
	mv.buildLayout()
}

func (mv *MainView) setupEventHandlers() {
	mv.listView.SetOnSelected(mv.onItemSelected)
	mv.listView.SetOnDelete(mv.onDeleteItem)
	mv.detailView.SetOnSave(mv.onSaveItem)
	mv.detailView.SetOnStatusUpdate(mv.updateStatus)
	mv.detailView.SetOnEditingMode(mv.setEditingMode)
	mv.detailView.SetOnNormalMode(mv.setNormalMode)
	mv.detailView.SetOnFocusSearch(mv.onFocusSearch)

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
	
	// 創建右側面板，在右下角放置設置按鈕
	rightBottom := container.NewBorder(nil, nil, nil, mv.settingsBtn, mv.statusBar.GetWidget())
	right := container.NewBorder(nil, rightBottom, nil, nil, mv.detailView.GetWidget())

	split := container.NewHSplit(left, right)
	split.SetOffset(0.35)
	mainContent := container.NewBorder(mv.toolbar.GetWidget(), nil, nil, nil, split)
	bg := canvas.NewRectangle(color.NRGBA{R: 33, G: 39, B: 44, A: 255})

	toastContainer := container.NewBorder(nil, mv.undoToast.GetWidget(), nil, nil, nil)

	mv.content = container.NewStack(bg, mainContent, toastContainer)
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
	// 先獲取要刪除的項目，用於 undo 功能
	itemToDelete := mv.clipboardController.GetHistoryItem(index)
	if itemToDelete == nil {
		mv.updateStatus("找不到要刪除的項目")
		return
	}

	// 保存刪除項目的信息，用於 undo
	mv.pendingDeleteItem = &models.ClipboardItem{
		Type:      itemToDelete.Type,
		Content:   itemToDelete.Content,
		FilePath:  itemToDelete.FilePath,
		Timestamp: itemToDelete.Timestamp,
	}
	mv.pendingDeleteIndex = index

	// 從控制器中移除項目（但保存在 pending 中以便復原）
	if err := mv.clipboardController.RemoveHistoryItem(index); err != nil {
		mv.updateStatus("刪除失敗: " + err.Error())
		return
	}

	// 從 UI 中移除
	mv.listView.RemoveItemAt(index)

	// 調整選中項目索引
	if mv.selectedIndex == index {
		if index >= mv.clipboardController.Count() {
			mv.selectedIndex = mv.clipboardController.Count() - 1
		} else {
			mv.selectedIndex = index
		}
	} else if mv.selectedIndex > index {
		mv.selectedIndex--
	}

	if mv.selectedIndex >= 0 && mv.selectedIndex < mv.clipboardController.Count() {
		item := mv.clipboardController.GetHistoryItem(mv.selectedIndex)
		mv.currentSelectedItem = item
		mv.detailView.ShowItem(item)
		mv.listView.GetWidget().Select(mv.selectedIndex)
	} else {
		mv.currentSelectedItem = nil
		mv.detailView.Clear()
		mv.listView.GetWidget().UnselectAll()
	}

	mv.undoToast.Show("item deleted", func() {
		mv.undoDelete()
	}, time.Duration(mv.config.UndoCountdownSecs)*time.Second)

	go func() {
		time.Sleep(time.Duration(mv.config.UndoCountdownSecs) * time.Second)
		if mv.pendingDeleteItem != nil && !mv.undoToast.IsVisible() {
			mv.pendingDeleteItem = nil
			mv.pendingDeleteIndex = -1
		}
	}()
}

func (mv *MainView) undoDelete() {
	if mv.pendingDeleteItem == nil {
		return
	}

	err := mv.clipboardController.InsertHistoryItemAt(mv.pendingDeleteIndex, mv.pendingDeleteItem)
	if err != nil {
		mv.updateStatus("復原失敗: " + err.Error())
		return
	}

	items := mv.clipboardController.GetHistoryItems()
	mv.listView.LoadFromHistory(items)

	mv.selectedIndex = mv.pendingDeleteIndex
	mv.currentSelectedItem = mv.pendingDeleteItem
	mv.detailView.ShowItem(mv.pendingDeleteItem)
	mv.listView.GetWidget().Select(mv.pendingDeleteIndex)

	mv.pendingDeleteItem = nil
	mv.pendingDeleteIndex = -1

	mv.updateStatus("Undo Item")
}

func (mv *MainView) onSaveItem(newContent string) {
	if mv.currentSelectedItem == nil {
		mv.updateStatus("No Selected")
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

	mv.updateStatus("Saved")
}

func (mv *MainView) onCopyToClipboard() error {
	if mv.currentSelectedItem == nil {
		return nil
	}

	if mv.currentSelectedItem.Type == models.ClipImage {
		return mv.clipboardController.CopyItemToClipboard(mv.currentSelectedItem)
	} else {
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
		mv.updateStatus("Search: " + query)
	}
}

func (mv *MainView) updateStatus(message string) {
	mv.statusBar.UpdateText(message)
}

func (mv *MainView) setEditingMode() {
	mv.statusBar.SetEditingMode()
}

func (mv *MainView) setNormalMode() {
	mv.statusBar.SetNormalMode()
}

func (mv *MainView) setErrorStatus(message string) {
	mv.statusBar.SetErrorStatus(message)
}

func (mv *MainView) setSuccessStatus(message string) {
	mv.statusBar.SetSuccessStatus(message)
}

func (mv *MainView) onFocusSearch() {
	mv.toolbar.FocusSearch()
	mv.updateStatus("Search focused")
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

func (mv *MainView) showSettings() {
	mv.settingsDialog.Show()
}

func (mv *MainView) onSettingsSave(newConfig *models.AppConfig) {
	// 更新當前配置
	mv.config.UndoCountdownSecs = newConfig.UndoCountdownSecs
	mv.config.WindowWidth = newConfig.WindowWidth
	mv.config.WindowHeight = newConfig.WindowHeight
	mv.config.MaxHistoryItems = newConfig.MaxHistoryItems
	mv.config.MaxDisplayLength = newConfig.MaxDisplayLength
	mv.config.PollingInterval = newConfig.PollingInterval
	
	// 保存配置到文件
	if err := mv.config.SaveConfig(); err != nil {
		mv.updateStatus("保存配置失敗: " + err.Error())
		return
	}
	
	// 更新視窗大小
	mv.window.Resize(fyne.NewSize(float32(newConfig.WindowWidth), float32(newConfig.WindowHeight)))
	
	mv.updateStatus("設置已保存")
}
