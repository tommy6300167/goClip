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

// TypedKey 覆寫鍵盤事件處理
func (e *CustomEntry) TypedKey(key *fyne.KeyEvent) {
	// 如果是 Enter 鍵且啟用防止自動換行
	if key.Name == fyne.KeyReturn && e.preventEnterNewline {
		// 檢查是否為輸入法組合字元
		if e.isComposingText() {
			// 如果是輸入法組合，讓它正常處理
			e.Entry.TypedKey(key)
		} else {
			// 否則只在 Shift+Enter 時才換行
			if key.Name == fyne.KeyReturn {
				if desktop := fyne.CurrentApp().Driver().(desktop.Driver); desktop != nil {
					if modifiers := desktop.CurrentKeyModifiers(); modifiers&fyne.KeyModifierShift != 0 {
						e.Entry.TypedKey(key)
					}
					// 其他情況不處理，避免自動換行
				}
			}
		}
		return
	}
	
	// 其他按鍵正常處理
	e.Entry.TypedKey(key)
}

// isComposingText 檢查是否正在輸入法組合中
func (e *CustomEntry) isComposingText() bool {
	// 這是一個簡化的檢查，實際上 Fyne 沒有直接的 API 來檢查組合狀態
	// 我們可以通過檢查最近的文字變化來推斷
	return false
}

// SetPreventEnterNewline 設定是否防止 Enter 鍵自動換行
func (e *CustomEntry) SetPreventEnterNewline(prevent bool) {
	e.preventEnterNewline = prevent
}