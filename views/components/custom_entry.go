package components

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
)

type CustomMultiLineEntry struct {
	*widget.Entry
	preventEnterNewline bool
}

// NewCustomMultiLineEntry 創建一個新的自定義多行輸入框
func NewCustomMultiLineEntry() *CustomMultiLineEntry {
	baseEntry := widget.NewMultiLineEntry()
	entry := &CustomMultiLineEntry{
		Entry: baseEntry,
	}

	return entry
}

func (c *CustomMultiLineEntry) SetPreventEnterNewline(prevent bool) {
	c.preventEnterNewline = prevent
}

func (c *CustomMultiLineEntry) TypedKey(key *fyne.KeyEvent) {
	if c.preventEnterNewline && key.Name == fyne.KeyReturn {
		return
	}

	c.Entry.TypedKey(key)
}

func (c *CustomMultiLineEntry) TypedRune(r rune) {
	c.Entry.TypedRune(r)
}

func (c *CustomMultiLineEntry) AllowNewlineMode() {
	c.SetPreventEnterNewline(false)
}

func (c *CustomMultiLineEntry) PreventNewlineMode() {
	c.SetPreventEnterNewline(true)
}
