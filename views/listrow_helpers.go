package views

import (
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/widget"

	u "clipmini/utils"
	"clipmini/views/components"
)

type RowSubviews struct {
	ContentLabel *widget.Label
	TimeText     *canvas.Text
	DeleteBtn    *components.DeleteButton
}

func FindRowSubviews(co fyne.CanvasObject) (RowSubviews, bool) {
	c, ok := co.(*fyne.Container)
	if !ok || len(c.Objects) < 3 {
		return RowSubviews{}, false
	}

	rightContent, _ := c.Objects[0].(*fyne.Container)
	deleteBtn, _ := c.Objects[2].(*components.DeleteButton)

	var lbl *widget.Label
	var ttxt *canvas.Text

	if rightContent != nil && len(rightContent.Objects) >= 1 {
		if inner, ok := rightContent.Objects[0].(*fyne.Container); ok && len(inner.Objects) >= 3 {
			if titlePad, ok := inner.Objects[0].(*fyne.Container); ok && len(titlePad.Objects) > 0 {
				lbl, _ = titlePad.Objects[0].(*widget.Label)
			}
			if timePad, ok := inner.Objects[2].(*fyne.Container); ok && len(timePad.Objects) > 0 {
				ttxt, _ = timePad.Objects[0].(*canvas.Text)
			}
		}
	}

	ok = (lbl != nil && ttxt != nil && deleteBtn != nil)
	return RowSubviews{ContentLabel: lbl, TimeText: ttxt, DeleteBtn: deleteBtn}, ok
}

func BindDeleteByLine(btn *components.DeleteButton, list binding.StringList, line string, onDelete func(int)) {
	if btn == nil || onDelete == nil {
		return
	}
	btn.SetOnTapped(func() {
		items, _ := list.Get()
		for i, it := range items {
			if it == line {
				onDelete(i)
				return
			}
		}
	})
}

func ApplyRow(ui RowSubviews, rp RowParsed, loc *time.Location) {
	if ui.ContentLabel != nil {
		ui.ContentLabel.SetText(rp.Title)
	}
	if ui.TimeText != nil {
		if !rp.Time.IsZero() {
			ui.TimeText.Text = u.FormatTimeWithTag(rp.Time, rp.Tag)
		} else {
			suffix := ""
			if rp.Tag != "" {
				suffix = " · " + rp.Tag
			}
			ui.TimeText.Text = rp.RawTime + suffix
		}
		ui.TimeText.Refresh()
	}
}
