package components

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
)

// CustomMultiLineEntry 是一個自定義的多行輸入框，可以控制 Enter 鍵的行為
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
	// 不需要 ExtendBaseWidget，因為我們只是包裝現有的 Entry
	return entry
}

// SetPreventEnterNewline 設置是否阻止 Enter 鍵換行
func (c *CustomMultiLineEntry) SetPreventEnterNewline(prevent bool) {
	c.preventEnterNewline = prevent
}

// TypedKey 處理鍵盤輸入事件
func (c *CustomMultiLineEntry) TypedKey(key *fyne.KeyEvent) {
	// 如果設置了阻止 Enter 換行，且按下的是 Enter 鍵，則不處理
	if c.preventEnterNewline && key.Name == fyne.KeyReturn {
		// 什麼都不做，忽略 Enter 鍵
		return
	}
	
	// 對於其他所有鍵，都正常處理
	c.Entry.TypedKey(key)
}

// TypedRune 處理字符輸入 - 確保正常文字輸入功能
func (c *CustomMultiLineEntry) TypedRune(r rune) {
	c.Entry.TypedRune(r)
}

// AllowNewlineMode 允許換行模式（正常的多行輸入行為）
func (c *CustomMultiLineEntry) AllowNewlineMode() {
	c.SetPreventEnterNewline(false)
}

// PreventNewlineMode 阻止換行模式（Enter 鍵不會換行）
func (c *CustomMultiLineEntry) PreventNewlineMode() {
	c.SetPreventEnterNewline(true)
}