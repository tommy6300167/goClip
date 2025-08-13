package views

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"

	"clipmini/controllers"
)

type Toolbar struct {
	container            *fyne.Container
	copyBtn              *widget.Button
	clearBtn             *widget.Button
	exportBtn            *widget.Button
	clipboardController  *controllers.ClipboardController
	onCopy               func() error
	onClear              func() error
	onExport             func() string
	onStatusUpdate       func(string)
	window               fyne.Window
}

func NewToolbar(window fyne.Window, clipboardController *controllers.ClipboardController) *Toolbar {
	tb := &Toolbar{
		clipboardController: clipboardController,
		window:              window,
	}
	
	tb.copyBtn = widget.NewButton("複製回剪貼簿", tb.handleCopy)
	tb.clearBtn = widget.NewButton("清空", tb.handleClear)
	tb.exportBtn = widget.NewButton("匯出到檔案", tb.handleExport)
	
	tb.container = container.NewHBox(tb.copyBtn, tb.clearBtn, tb.exportBtn)
	
	return tb
}

func (tb *Toolbar) GetWidget() *fyne.Container {
	return tb.container
}

func (tb *Toolbar) SetOnCopy(callback func() error) {
	tb.onCopy = callback
}

func (tb *Toolbar) SetOnClear(callback func() error) {
	tb.onClear = callback
}

func (tb *Toolbar) SetOnExport(callback func() string) {
	tb.onExport = callback
}

func (tb *Toolbar) SetOnStatusUpdate(callback func(string)) {
	tb.onStatusUpdate = callback
}

func (tb *Toolbar) handleCopy() {
	if tb.onCopy != nil {
		if err := tb.onCopy(); err != nil {
			tb.updateStatus("複製失敗")
		} else {
			tb.updateStatus("已複製到剪貼簿")
		}
	}
}

func (tb *Toolbar) handleClear() {
	dialog.ShowConfirm("清空", "確定要清空歷史？（會刪除已存圖片）", func(ok bool) {
		if !ok {
			return
		}
		
		if tb.onClear != nil {
			if err := tb.onClear(); err != nil {
				tb.updateStatus("清空失敗")
			} else {
				tb.updateStatus("已清空")
			}
		}
	}, tb.window)
}

func (tb *Toolbar) handleExport() {
	fd := dialog.NewFileSave(func(uc fyne.URIWriteCloser, err error) {
		if err != nil || uc == nil {
			return
		}
		defer uc.Close()
		
		if tb.onExport != nil {
			content := tb.onExport()
			uc.Write([]byte(content))
			tb.updateStatus("已匯出")
		}
	}, tb.window)
	
	fd.SetFileName("clipmini_history.txt")
	fd.Show()
}

func (tb *Toolbar) updateStatus(message string) {
	if tb.onStatusUpdate != nil {
		tb.onStatusUpdate(message)
	}
}