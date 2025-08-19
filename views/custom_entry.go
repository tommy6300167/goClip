package views

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/widget"
)

// CustomEntry 自訂文字輸入元件，防止中文輸入時的自動換行
type CustomEntry struct {
	widget.Entry
	preventEnterNewline bool
}

// NewCustomEntry 建立新的自訂文字輸入元件
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
	if key.Name == fyne.KeyReturn && e.preventEnterNewline {
		if e.isComposingText() {
			e.Entry.TypedKey(key)
		} else {
			if key.Name == fyne.KeyReturn {
				if desktop := fyne.CurrentApp().Driver().(desktop.Driver); desktop != nil {
					if modifiers := desktop.CurrentKeyModifiers(); modifiers&fyne.KeyModifierShift != 0 {
						e.Entry.TypedKey(key)
					}
				}
			}
		}
		return
	}

	e.Entry.TypedKey(key)
}

func (e *CustomEntry) isComposingText() bool {
	return false
}

func (e *CustomEntry) SetPreventEnterNewline(prevent bool) {
	e.preventEnterNewline = prevent
}
