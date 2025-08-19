package views

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"

	"clipmini/controllers"
	"clipmini/views/components"
)

type Toolbar struct {
	container           *fyne.Container
	copyBtn             *components.PillButton
	clearBtn            *components.PillButton
	exportBtn           *components.PillButton
	searchInput         *components.SearchInput
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

	neutral := components.DefaultPillStyle
	neutral.Size = fyne.NewSize(88, 36)

	transparent := neutral
	transparent.Bg = color.NRGBA{0, 0, 0, 0}
	transparent.BgHover = color.NRGBA{255, 255, 255, 20}

	accent := neutral
	accent.Bg = color.NRGBA{45, 150, 135, 255}      // 更深的藍綠色
	accent.BgHover = color.NRGBA{35, 130, 118, 255} // 更深的 hover 效果
	accent.Fg = color.NRGBA{255, 255, 255, 255}     // 白色文字確保對比度

	tb.copyBtn = components.NewPillButtonWithStyle("Copy", theme.ContentCopyIcon(), tb.handleCopy, transparent)
	tb.clearBtn = components.NewPillButtonWithStyle("Clear", theme.DeleteIcon(), tb.handleClear, transparent)
	tb.exportBtn = components.NewPillButtonWithStyle("Export", theme.DocumentSaveIcon(), tb.handleExport, accent)

	tb.searchInput = components.NewSearchInput()
	tb.searchInput.SetOnChanged(func(text string) {
		if tb.onSearch != nil {
			tb.onSearch(text)
		}
	})

	tb.container = container.NewHBox(
		container.NewPadded(tb.clearBtn),
		container.NewPadded(tb.copyBtn),
		container.NewPadded(tb.exportBtn),
		layout.NewSpacer(),
		tb.searchInput.GetWidget(),
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
			tb.updateStatus("Failed to copy")
		} else {
			tb.updateStatus("Copyed")
		}
	}
}

func (tb *Toolbar) handleClear() {
	dialog.ShowConfirm("Clear", "Confirm clearing history? ", func(ok bool) {
		if !ok {
			return
		}

		if tb.onClear != nil {
			if err := tb.onClear(); err != nil {
				tb.updateStatus("Failed to clear")
			} else {
				tb.updateStatus("Cleared")
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
			tb.updateStatus("Export successful")
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
