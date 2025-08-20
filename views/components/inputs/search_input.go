package inputs

import (
	"image/color"
	"strings"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

type SearchInput struct {
	container *fyne.Container
	entry     *widget.Entry
	clearBtn  *widget.Button
	onChanged func(string)
}

func NewSearchInput() *SearchInput {
	si := &SearchInput{}

	si.entry = widget.NewEntry()
	si.entry.SetPlaceHolder(" Search...")

	searchIcon := widget.NewIcon(theme.SearchIcon())

	si.clearBtn = widget.NewButtonWithIcon("", theme.ContentClearIcon(), func() {
		si.entry.SetText("")
		if si.onChanged != nil {
			si.onChanged("")
		}
	})
	si.clearBtn.Importance = widget.LowImportance
	si.clearBtn.Hide()

	si.entry.OnChanged = func(text string) {
		t := strings.TrimSpace(text)
		if t == "" {
			si.clearBtn.Hide()
		} else {
			si.clearBtn.Show()
		}
		if si.onChanged != nil {
			si.onChanged(t)
		}
	}

	searchBg := canvas.NewRectangle(color.NRGBA{45, 55, 65, 255})
	searchBg.CornerRadius = 8

	iconBox := container.NewCenter(searchIcon)
	clearBox := container.NewCenter(si.clearBtn)
	entryBox := container.NewStack(si.entry)

	searchInner := container.NewBorder(nil, nil, iconBox, clearBox, entryBox)
	searchWrapper := container.NewStack(searchBg, searchInner)

	fixedSearch := container.New(layout.NewGridWrapLayout(fyne.NewSize(200, 36)), searchWrapper)
	si.container = container.NewVBox(
		layout.NewSpacer(),
		fixedSearch,
	)

	return si
}

func (si *SearchInput) GetWidget() *fyne.Container {
	return si.container
}

func (si *SearchInput) SetOnChanged(callback func(string)) {
	si.onChanged = callback
}

func (si *SearchInput) GetText() string {
	return strings.TrimSpace(si.entry.Text)
}

func (si *SearchInput) SetText(text string) {
	si.entry.SetText(text)
}

func (si *SearchInput) Clear() {
	si.entry.SetText("")
}

func (si *SearchInput) Focus() {

	go func() {
		time.Sleep(10 * time.Millisecond)
		fyne.Do(func() {

			si.entry.CursorRow = 0
			si.entry.CursorColumn = len(si.entry.Text)

			if canvas := fyne.CurrentApp().Driver().CanvasForObject(si.entry); canvas != nil {
				canvas.Focus(si.entry)
			}
		})
	}()
}
