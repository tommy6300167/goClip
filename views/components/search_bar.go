package components

import (
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
)

type SearchBar struct {
	*widget.Entry

	onChanged func(string)
}

func (sb *SearchBar) MinSize() fyne.Size {
	return fyne.NewSize(280, 40)
}

// Tapped handles tap events and ensures proper focus
func (sb *SearchBar) Tapped(ev *fyne.PointEvent) {
	sb.Entry.Tapped(ev)
}

// FocusGained handles when the widget gains focus
func (sb *SearchBar) FocusGained() {
	sb.Entry.FocusGained()
}

// FocusLost handles when the widget loses focus
func (sb *SearchBar) FocusLost() {
	sb.Entry.FocusLost()
}

// Focusable returns whether this widget can receive focus
func (sb *SearchBar) Focusable() bool {
	return true
}

// KeyDown handles key press events
func (sb *SearchBar) KeyDown(key *fyne.KeyEvent) {
	sb.Entry.KeyDown(key)
}

// TypedRune handles typed character events
func (sb *SearchBar) TypedRune(r rune) {
	sb.Entry.TypedRune(r)
}

// TypedKey handles typed key events
func (sb *SearchBar) TypedKey(key *fyne.KeyEvent) {
	sb.Entry.TypedKey(key)
}

func NewSearchBar() *SearchBar {
	entry := widget.NewEntry()
	entry.SetPlaceHolder("搜尋剪貼簿內容...")
	
	sb := &SearchBar{
		Entry: entry,
	}

	// 直接使用Entry的OnChanged事件
	sb.Entry.OnChanged = func(text string) {
		text = strings.TrimSpace(text)
		if sb.onChanged != nil {
			sb.onChanged(text)
		}
	}

	return sb
}

func (sb *SearchBar) SetOnChanged(callback func(string)) {
	sb.onChanged = callback
}

func (sb *SearchBar) GetText() string {
	return strings.TrimSpace(sb.Entry.Text)
}

func (sb *SearchBar) SetText(text string) {
	sb.Entry.SetText(text)
}

func (sb *SearchBar) Clear() {
	sb.Entry.SetText("")
}
