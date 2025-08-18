package views

import (
	"image/color"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"

	"clipmini/controllers"
	"clipmini/views/components"
)

type Toolbar struct {
	container           *fyne.Container
	copyBtn             *components.PillButton
	clearBtn            *components.PillButton
	exportBtn           *components.PillButton
	searchEntry         *widget.Entry
	clipboardController *controllers.ClipboardController
	onCopy              func() error
	onClear             func() error
	onExport            func() string
	onSearch            func(string)
	onStatusUpdate      func(string)
	window              fyne.Window
}

func NewToolbar(window fyne.Window, clipboardController *controllers.ClipboardController) *Toolbar {
	tb := &Toolbar{
		clipboardController: clipboardController,
		window:              window,
	}

	// 按鈕統一高度，和搜尋同高
	neutral := components.DefaultPillStyle
	neutral.Size = fyne.NewSize(88, 36)

	accent := neutral
	accent.Bg = color.NRGBA{86, 212, 194, 255}
	accent.BgHover = color.NRGBA{78, 190, 174, 255}
	accent.Fg = color.NRGBA{15, 23, 28, 255}

	tb.copyBtn = components.NewPillButtonWithStyle("Copy", theme.ContentCopyIcon(), tb.handleCopy, neutral)
	tb.clearBtn = components.NewPillButtonWithStyle("Clear", theme.DeleteIcon(), tb.handleClear, neutral)
	tb.exportBtn = components.NewPillButtonWithStyle("Export", theme.DocumentSaveIcon(), tb.handleExport, accent)

	tb.searchEntry = widget.NewEntry()
	tb.searchEntry.SetPlaceHolder("搜尋剪貼簿內容...")

	searchIcon := widget.NewIcon(theme.SearchIcon())
	clearBtn := widget.NewButtonWithIcon("", theme.ContentClearIcon(), func() {
		tb.searchEntry.SetText("")
		if tb.onSearch != nil {
			tb.onSearch("")
		}
	})
	clearBtn.Importance = widget.LowImportance
	clearBtn.Hide()

	tb.searchEntry.OnChanged = func(text string) {
		t := strings.TrimSpace(text)
		if t == "" {
			clearBtn.Hide()
		} else {
			clearBtn.Show()
		}
		if tb.onSearch != nil {
			tb.onSearch(t)
		}
	}

	searchBg := canvas.NewRectangle(color.NRGBA{45, 55, 65, 255})
	searchBg.CornerRadius = 8

	iconBox := container.NewCenter(searchIcon)
	clearBox := container.NewCenter(clearBtn)
	entryBox := container.NewStack(tb.searchEntry)

	searchInner := container.NewBorder(nil, nil, iconBox, clearBox, entryBox)
	searchWrapper := container.NewStack(searchBg, searchInner)

	fixedSearch := container.New(layout.NewGridWrapLayout(fyne.NewSize(200, 36)), searchWrapper)
	searchWithMargin := container.NewVBox(
		layout.NewSpacer(),
		fixedSearch,
	)

	tb.container = container.NewHBox(
		container.NewPadded(tb.clearBtn),
		container.NewPadded(tb.copyBtn),
		container.NewPadded(tb.exportBtn),
		layout.NewSpacer(),
		searchWithMargin,
	)

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

func (tb *Toolbar) SetOnSearch(callback func(string)) {
	tb.onSearch = callback
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
