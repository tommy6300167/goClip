package components

import (
	"fmt"
	"image/color"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

type UndoToast struct {
	container        *fyne.Container
	background       *canvas.Rectangle
	messageLabel     *widget.Label
	undoBtn          *widget.Button
	closeBtn         *widget.Button
	isVisible        bool
	hideTimer        *time.Timer
	countdownTimer   *time.Ticker
	remainingSeconds int
	onUndo           func()
	onClose          func()
}

func NewUndoToast() *UndoToast {
	toast := &UndoToast{
		isVisible: false,
	}

	toast.background = canvas.NewRectangle(color.NRGBA{0x33, 0x33, 0x33, 0xF0})
	toast.background.CornerRadius = 8

	toast.messageLabel = widget.NewLabel("Item deleted")
	toast.messageLabel.TextStyle = fyne.TextStyle{Bold: false}

	toast.undoBtn = widget.NewButton("Undo", func() {
		if toast.onUndo != nil {
			toast.onUndo()
		}
		toast.Hide()
	})
	toast.undoBtn.Importance = widget.HighImportance

	toast.closeBtn = widget.NewButtonWithIcon("", theme.CancelIcon(), func() {
		if toast.onClose != nil {
			toast.onClose()
		}
		toast.Hide()
	})
	toast.closeBtn.Importance = widget.LowImportance

	content := container.NewHBox(
		toast.messageLabel,
		widget.NewSeparator(),
		toast.undoBtn,
		toast.closeBtn,
	)

	toast.container = container.NewStack(
		toast.background,
		container.NewPadded(content),
	)

	toast.container.Hide()

	return toast
}

func (t *UndoToast) Show(message string, onUndo func(), autoHideDuration time.Duration) {
	if t.isVisible {
		t.Hide() // 如果已經顯示，先隱藏
	}

	t.onUndo = onUndo
	t.isVisible = true
	t.remainingSeconds = int(autoHideDuration.Seconds())

	// 初始顯示包含倒數的消息
	t.updateMessageWithCountdown(message)
	t.container.Show()

	// 設定倒數計時器
	if t.countdownTimer != nil {
		t.countdownTimer.Stop()
	}

	t.countdownTimer = time.NewTicker(1 * time.Second)
	go func() {
		for range t.countdownTimer.C {
			if !t.isVisible {
				return
			}
			t.remainingSeconds--
			if t.remainingSeconds <= 0 {
				fyne.Do(func() {
					t.Hide()
				})
				return
			}
			fyne.Do(func() {
				t.updateMessageWithCountdown(message)
			})
		}
	}()

	if t.hideTimer != nil {
		t.hideTimer.Stop()
	}

	t.hideTimer = time.AfterFunc(autoHideDuration, func() {
		t.Hide()
	})
}

func (t *UndoToast) updateMessageWithCountdown(baseMessage string) {
	countdownMessage := fmt.Sprintf("%s (%d Second)", baseMessage, t.remainingSeconds)
	t.messageLabel.SetText(countdownMessage)
}

func (t *UndoToast) Hide() {
	if !t.isVisible {
		return
	}

	if t.hideTimer != nil {
		t.hideTimer.Stop()
		t.hideTimer = nil
	}

	if t.countdownTimer != nil {
		t.countdownTimer.Stop()
		t.countdownTimer = nil
	}

	t.isVisible = false
	t.container.Hide()
	t.onUndo = nil
	t.remainingSeconds = 0
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
