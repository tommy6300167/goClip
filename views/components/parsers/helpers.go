package parsers

import (
	"time"
	"unicode/utf8"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"

	"clipmini/utils"
)

// FindRowSubviews 從容器中找到列表項的 UI 組件
func FindRowSubviews(co fyne.CanvasObject) (RowSubviews, bool) {
	c, ok := co.(*fyne.Container)
	if !ok || len(c.Objects) < 3 {
		return RowSubviews{}, false
	}

	// ListItem 結構：
	// c.Objects[0] = rightContent (padded container)
	// c.Objects[1] = contentIcon
	// c.Objects[2] = deleteContainer (stack with background + button)

	rightContent, _ := c.Objects[0].(*fyne.Container)
	contentIcon, _ := c.Objects[1].(*widget.Icon)
	
	// 處理被包裝在容器中的垃圾桶按鈕
	var deleteBtn *widget.Button
	if deleteContainer, ok := c.Objects[2].(*fyne.Container); ok && len(deleteContainer.Objects) >= 2 {
		deleteBtn, _ = deleteContainer.Objects[1].(*widget.Button)
	} else {
		deleteBtn, _ = c.Objects[2].(*widget.Button)
	}

	var ct, tt *canvas.Text
	if rightContent != nil && len(rightContent.Objects) >= 1 {
		if inner, ok := rightContent.Objects[0].(*fyne.Container); ok && len(inner.Objects) >= 2 {
			// 結構：inner.Objects[0] = titlePad, inner.Objects[1] = gap, inner.Objects[2] = timePad
			if titlePad, ok := inner.Objects[0].(*fyne.Container); ok && len(titlePad.Objects) > 0 {
				ct, _ = titlePad.Objects[0].(*canvas.Text)
			}
			if len(inner.Objects) >= 3 {
				if timePad, ok := inner.Objects[2].(*fyne.Container); ok && len(timePad.Objects) > 0 {
					tt, _ = timePad.Objects[0].(*canvas.Text)
				}
			}
		}
	}

	ok = (contentIcon != nil && ct != nil && tt != nil && deleteBtn != nil)
	return RowSubviews{ContentIcon: contentIcon, ContentText: ct, TimeText: tt, DeleteBtn: deleteBtn}, ok
}

// BindDeleteByLine 為刪除按鈕綁定事件
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

// ApplyRow 將解析後的數據應用到 UI 組件
func ApplyRow(ui RowSubviews, rp RowParsed, maxLen int, loc *time.Location) {
	// 更新圖標
	if ui.ContentIcon != nil {
		if rp.Kind == RowImage {
			ui.ContentIcon.SetResource(theme.MediaPhotoIcon())
		} else {
			ui.ContentIcon.SetResource(theme.DocumentIcon())
		}
		ui.ContentIcon.Refresh()
	}

	// 更新內容文字
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

	// 更新時間文字
	if ui.TimeText != nil {
		if !rp.Time.IsZero() {
			ui.TimeText.Text = utils.FormatTimeWithTag(rp.Time, rp.Tag)
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