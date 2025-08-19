package views

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/widget"
)

type CustomEntry struct {
	widget.Entry
	preventEnterNewline bool
}

func NewCustomEntry() *CustomEntry {
	entry := &CustomEntry{
		preventEnterNewline: true,
	}
	entry.ExtendBaseWidget(entry)
	entry.MultiLine = true
	entry.Wrapping = fyne.TextWrapOff
	return entry
}

func (e *CustomEntry) TypedKey(key *fyne.KeyEvent) {
	if key.Name != "LeftSuper" && key.Name != "RightSuper" &&
		key.Name != "LeftControl" && key.Name != "RightControl" {
	}

	if key.Name == fyne.KeyReturn && e.preventEnterNewline {
		if desktop := fyne.CurrentApp().Driver().(desktop.Driver); desktop != nil {
			if modifiers := desktop.CurrentKeyModifiers(); modifiers&fyne.KeyModifierShift != 0 {
				e.Entry.TypedKey(key)
			}
		}
		return
	}
	e.Entry.TypedKey(key)
}

func (e *CustomEntry) SetPreventEnterNewline(prevent bool) {
	e.preventEnterNewline = prevent
}
