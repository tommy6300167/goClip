package views

import (
	"time"
	"unicode/utf8"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"

	u "clipmini/utils"
)

type RowSubviews struct {
	ContentIcon *widget.Icon
	ContentText *canvas.Text
	TimeText    *canvas.Text
	DeleteBtn   *widget.Button
}

func FindRowSubviews(co fyne.CanvasObject) (RowSubviews, bool) {
	c, ok := co.(*fyne.Container)
	if !ok || len(c.Objects) < 3 {
		return RowSubviews{}, false
	}

	rightContent, _ := c.Objects[0].(*fyne.Container)
	contentIcon, _ := c.Objects[1].(*widget.Icon)
	deleteBtn, _ := c.Objects[2].(*widget.Button)

	var ct, tt *canvas.Text
	if rightContent != nil && len(rightContent.Objects) >= 1 {
		if inner, ok := rightContent.Objects[0].(*fyne.Container); ok && len(inner.Objects) >= 3 {
			if titlePad, ok := inner.Objects[0].(*fyne.Container); ok && len(titlePad.Objects) > 0 {
				ct, _ = titlePad.Objects[0].(*canvas.Text)
			}
			if timePad, ok := inner.Objects[2].(*fyne.Container); ok && len(timePad.Objects) > 0 {
				tt, _ = timePad.Objects[0].(*canvas.Text)
			}
		}
	}

	ok = (contentIcon != nil && ct != nil && tt != nil && deleteBtn != nil)
	return RowSubviews{ContentIcon: contentIcon, ContentText: ct, TimeText: tt, DeleteBtn: deleteBtn}, ok
}
func BindDeleteByLine(btn *widget.Button, list binding.StringList, line string, onDelete func(int)) {
	if btn == nil || onDelete == nil {
		return
	}
	btn.OnTapped = func() {
		items, _ := list.Get()
		for i, it := range items {
			if it == line {
				onDelete(i)
				return
			}
		}
	}
}
func ApplyRow(ui RowSubviews, rp RowParsed, maxLen int, loc *time.Location) {
	if ui.ContentIcon != nil {
		if rp.Kind == RowImage {
			ui.ContentIcon.SetResource(theme.MediaPhotoIcon())
		} else {
			ui.ContentIcon.SetResource(theme.DocumentIcon())
		}
		ui.ContentIcon.Refresh()
	}

	if ui.ContentText != nil {
		if maxLen > 0 && utf8.RuneCountInString(rp.Title) > maxLen {
			runes := []rune(rp.Title)
			if maxLen == 1 {
				ui.ContentText.Text = "…"
			} else {
				ui.ContentText.Text = string(runes[:maxLen-1]) + "…"
			}
		} else {
			ui.ContentText.Text = rp.Title
		}
		ui.ContentText.Refresh()
	}

	if ui.TimeText != nil {
		if !rp.Time.IsZero() {
			ui.TimeText.Text = u.FormatTimeWithTag(rp.Time, rp.Tag)
		} else {
			txt := rp.RawTime
			if rp.Tag != "" {
				txt += " · " + rp.Tag
			}
			ui.TimeText.Text = txt
		}
		ui.TimeText.Refresh()
	}
}
