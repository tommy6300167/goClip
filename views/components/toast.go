package components

import (
	"image/color"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

type UndoToast struct {
	container    *fyne.Container
	background   *canvas.Rectangle
	messageLabel *widget.Label
	undoBtn      *widget.Button
	closeBtn     *widget.Button
	isVisible    bool
	hideTimer    *time.Timer
	onUndo       func()
	onClose      func()
}

func NewUndoToast() *UndoToast {
	toast := &UndoToast{
		isVisible: false,
	}

	// 創建背景
	toast.background = canvas.NewRectangle(color.NRGBA{0x33, 0x33, 0x33, 0xF0})
	toast.background.CornerRadius = 8

	// 創建消息標籤
	toast.messageLabel = widget.NewLabel("項目已刪除")
	toast.messageLabel.TextStyle = fyne.TextStyle{Bold: false}

	// 創建 Undo 按鈕
	toast.undoBtn = widget.NewButton("復原", func() {
		if toast.onUndo != nil {
			toast.onUndo()
		}
		toast.Hide()
	})
	toast.undoBtn.Importance = widget.HighImportance

	// 創建關閉按鈕
	toast.closeBtn = widget.NewButtonWithIcon("", theme.CancelIcon(), func() {
		if toast.onClose != nil {
			toast.onClose()
		}
		toast.Hide()
	})
	toast.closeBtn.Importance = widget.LowImportance

	// 創建內容容器
	content := container.NewHBox(
		toast.messageLabel,
		widget.NewSeparator(),
		toast.undoBtn,
		toast.closeBtn,
	)

	// 創建主容器 - 使用 Stack 來放置背景和內容
	toast.container = container.NewStack(
		toast.background,
		container.NewPadded(content),
	)

	// 初始時隱藏
	toast.container.Hide()

	return toast
}

func (t *UndoToast) Show(message string, onUndo func(), autoHideDuration time.Duration) {
	if t.isVisible {
		t.Hide() // 如果已經顯示，先隱藏
	}

	t.messageLabel.SetText(message)
	t.onUndo = onUndo
	t.isVisible = true
	t.container.Show()

	// 設定自動隱藏計時器
	if t.hideTimer != nil {
		t.hideTimer.Stop()
	}

	t.hideTimer = time.AfterFunc(autoHideDuration, func() {
		t.Hide()
	})
}

func (t *UndoToast) Hide() {
	if !t.isVisible {
		return
	}

	if t.hideTimer != nil {
		t.hideTimer.Stop()
		t.hideTimer = nil
	}

	t.isVisible = false
	t.container.Hide()
	t.onUndo = nil
}

func (t *UndoToast) SetOnClose(callback func()) {
	t.onClose = callback
}

func (t *UndoToast) GetWidget() *fyne.Container {
	return t.container
}

func (t *UndoToast) IsVisible() bool {
	return t.isVisible
}