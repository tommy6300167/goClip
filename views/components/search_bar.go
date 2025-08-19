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

func (sb *SearchBar) Tapped(ev *fyne.PointEvent) {
	sb.Entry.Tapped(ev)
}

func (sb *SearchBar) FocusGained() {
	sb.Entry.FocusGained()
}

func (sb *SearchBar) FocusLost() {
	sb.Entry.FocusLost()
}

func (sb *SearchBar) Focusable() bool {
	return true
}

func (sb *SearchBar) KeyDown(key *fyne.KeyEvent) {
	sb.Entry.KeyDown(key)
}

func (sb *SearchBar) TypedRune(r rune) {
	sb.Entry.TypedRune(r)
}

func (sb *SearchBar) TypedKey(key *fyne.KeyEvent) {
	sb.Entry.TypedKey(key)
}

func NewSearchBar() *SearchBar {
	entry := widget.NewEntry()
	entry.SetPlaceHolder(" Search...")

	sb := &SearchBar{
		Entry: entry,
	}

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
