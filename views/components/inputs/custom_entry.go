package inputs

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/widget"
)

// CustomEntry 自定義輸入框，支援更靈活的鍵盤處理
type CustomEntry struct {
	widget.Entry
	preventEnterNewline bool
}

// NewCustomEntry 創建新的自定義輸入框
func NewCustomEntry() *CustomEntry {
	entry := &CustomEntry{
		preventEnterNewline: true,
	}
	entry.ExtendBaseWidget(entry)
	entry.MultiLine = true
	entry.Wrapping = fyne.TextWrapOff
	return entry
}

// TypedKey 處理鍵盤輸入事件
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

// SetPreventEnterNewline 設定是否阻止 Enter 鍵換行
func (e *CustomEntry) SetPreventEnterNewline(prevent bool) {
	e.preventEnterNewline = prevent
}

// AllowNewlineMode 允許換行模式
func (e *CustomEntry) AllowNewlineMode() {
	e.SetPreventEnterNewline(false)
}

// PreventNewlineMode 阻止換行模式
func (e *CustomEntry) PreventNewlineMode() {
	e.SetPreventEnterNewline(true)
}

// CustomMultiLineEntry 自定義多行輸入框（向後相容性）
type CustomMultiLineEntry struct {
	*widget.Entry
	preventEnterNewline bool
}

// NewCustomMultiLineEntry 創建一個新的自定義多行輸入框
// 注意：建議使用 NewCustomEntry() 代替
func NewCustomMultiLineEntry() *CustomMultiLineEntry {
	baseEntry := widget.NewMultiLineEntry()
	entry := &CustomMultiLineEntry{
		Entry: baseEntry,
	}
	return entry
}

// SetPreventEnterNewline 設定是否阻止 Enter 鍵換行
func (c *CustomMultiLineEntry) SetPreventEnterNewline(prevent bool) {
	c.preventEnterNewline = prevent
}

// TypedKey 處理鍵盤輸入事件
func (c *CustomMultiLineEntry) TypedKey(key *fyne.KeyEvent) {
	if c.preventEnterNewline && key.Name == fyne.KeyReturn {
		return
	}
	c.Entry.TypedKey(key)
}

// TypedRune 處理字符輸入事件
func (c *CustomMultiLineEntry) TypedRune(r rune) {
	c.Entry.TypedRune(r)
}

// AllowNewlineMode 允許換行模式
func (c *CustomMultiLineEntry) AllowNewlineMode() {
	c.SetPreventEnterNewline(false)
}

// PreventNewlineMode 阻止換行模式
func (c *CustomMultiLineEntry) PreventNewlineMode() {
	c.SetPreventEnterNewline(true)
}