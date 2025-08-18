package views

import (
	"image/color"
	"path/filepath"
	"strings"
	"time"
	"unicode/utf8"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/widget"

	"clipmini/models"
	"clipmini/utils"
	"clipmini/views/components"
)

type ListView struct {
	list          *widget.List
	historyData   binding.StringList
	config        *models.AppConfig
	onSelected    func(int)
	onDelete      func(int)
	selectedIndex int
}

func NewListView(config *models.AppConfig) *ListView {
	lv := &ListView{
		historyData:   binding.NewStringList(),
		config:        config,
		selectedIndex: -1,
	}

	lv.list = widget.NewListWithData(
		lv.historyData,
		func() fyne.CanvasObject {
			iconLabel := widget.NewLabelWithStyle("📄", fyne.TextAlignCenter, fyne.TextStyle{Bold: true})
			iconLabel.Resize(fyne.NewSize(40, 40))

			contentLabel := widget.NewLabelWithStyle("", fyne.TextAlignLeading, fyne.TextStyle{Bold: true})
			contentLabel.Wrapping = fyne.TextWrapOff

			timeText := canvas.NewText("", color.NRGBA{R: 0xB0, G: 0xB0, B: 0xB0, A: 0xFF})
			timeText.TextSize = 12
			timeText.Alignment = fyne.TextAlignLeading

			gap := canvas.NewRectangle(color.Transparent)
			gap.SetMinSize(fyne.NewSize(1, 2))

			rightInner := container.NewVBox(contentLabel, gap, timeText)
			rightContent := container.NewPadded(rightInner)

			deleteBtn := components.NewDeleteButton()

			return container.NewBorder(nil, nil, iconLabel, deleteBtn, rightContent)
		},
		func(di binding.DataItem, co fyne.CanvasObject) {
			str, _ := di.(binding.String).Get()
			parts := strings.SplitN(str, "\t", 3)

			containerObj := co.(*fyne.Container)

			var (
				iconLabel    *widget.Label
				deleteBtn    *components.DeleteButton
				rightContent *fyne.Container
				contentLabel *widget.Label
				timeText     *canvas.Text
			)

			if len(containerObj.Objects) >= 3 {
				rightContent, _ = containerObj.Objects[0].(*fyne.Container) // center (Padded)
				iconLabel, _ = containerObj.Objects[1].(*widget.Label)      // left
				deleteBtn, _ = containerObj.Objects[2].(*components.DeleteButton) // right

				if rightContent != nil && len(rightContent.Objects) >= 1 {
					if inner, ok := rightContent.Objects[0].(*fyne.Container); ok && len(inner.Objects) >= 3 {
						contentLabel, _ = inner.Objects[0].(*widget.Label)
						timeText, _ = inner.Objects[2].(*canvas.Text)
					}
				}
			}

			if deleteBtn != nil {
				items, _ := lv.historyData.Get()
				currentIndex := -1
				for i, item := range items {
					if item == str {
						currentIndex = i
						break
					}
				}
				deleteBtn.SetOnTapped(func() {
					if lv.onDelete != nil && currentIndex >= 0 {
						lv.onDelete(currentIndex)
					}
				})
			}

			if iconLabel != nil {
				if len(parts) == 3 && parts[2] == "IMAGE" {
					iconLabel.SetText("🖼️")
				} else {
					iconLabel.SetText("📄")
				}
			}

			if contentLabel != nil && timeText != nil {
				if len(parts) == 3 && parts[2] == "IMAGE" {

					filePath := parts[1]
					fileName := filepath.Base(filePath)
					if fileName == "" {
						fileName = "圖片檔案"
					}
					contentLabel.SetText(ellipsis(fileName, lv.config.MaxDisplayLength))

					if originalTime, err := time.Parse("2006-01-02 15:04:05", parts[0]); err == nil {
						timeText.Text = originalTime.In(utils.GetTaipeiLocation()).Format("2006-01-02 15:04:05") + " · 圖片"
					} else {
						timeText.Text = parts[0] + " · 圖片"
					}
					timeText.Refresh()

				} else if len(parts) >= 2 {
					content := parts[1]
					clean := strings.ReplaceAll(content, "\n", " ")
					clean = strings.ReplaceAll(clean, "\r", " ")
					clean = strings.ReplaceAll(clean, "\t", " ")
					for strings.Contains(clean, "  ") {
						clean = strings.ReplaceAll(clean, "  ", " ")
					}
					clean = strings.TrimSpace(clean)

					contentLabel.SetText(ellipsis(clean, lv.config.MaxDisplayLength))

					if originalTime, err := time.Parse("2006-01-02 15:04:05", parts[0]); err == nil {
						timeText.Text = originalTime.In(utils.GetTaipeiLocation()).Format("2006-01-02 15:04:05") + " · 文字"
					} else {
						timeText.Text = parts[0] + " · 文字"
					}
					timeText.Refresh()

				} else {
					contentLabel.SetText(ellipsis(str, lv.config.MaxDisplayLength))
					timeText.Text = ""
					timeText.Refresh()
				}
			}
		},
	)

	lv.list.OnSelected = func(id widget.ListItemID) {
		lv.selectedIndex = id
		if lv.onSelected != nil {
			lv.onSelected(id)
		}
	}

	return lv
}

// ---- 事件綁定/資料操作 ----

func (lv *ListView) GetWidget() *widget.List { return lv.list }

func (lv *ListView) SetOnSelected(callback func(int)) { lv.onSelected = callback }

func (lv *ListView) SetOnDelete(callback func(int)) { lv.onDelete = callback }

func (lv *ListView) LoadFromHistory(items []*models.ClipboardItem) {
	lines := make([]string, len(items))
	for i, item := range items {
		timestamp := utils.FormatTimestamp(item.Timestamp, utils.GetTaipeiLocation())
		if item.Type == models.ClipImage {
			lines[i] = timestamp + "\t" + item.FilePath + "\t" + item.Type.String()
		} else {
			lines[i] = timestamp + "\t" + item.Content
		}
	}
	lv.historyData.Set(lines)
	if len(lines) > 0 {
		lv.SelectFirst()
	}
}

func (lv *ListView) PrependItem(item *models.ClipboardItem) {
	timestamp := utils.FormatTimestamp(item.Timestamp, utils.GetTaipeiLocation())
	var line string
	if item.Type == models.ClipImage {
		line = timestamp + "\t" + item.FilePath + "\t" + item.Type.String()
	} else {
		line = timestamp + "\t" + item.Content
	}
	lv.historyData.Prepend(line)
	lv.SelectFirst()
}

func (lv *ListView) Clear() {
	lv.historyData.Set([]string{})
	lv.selectedIndex = -1
}

func (lv *ListView) SelectFirst() {
	if items, _ := lv.historyData.Get(); len(items) > 0 {
		lv.selectedIndex = 0
		lv.list.Select(0)
		if lv.onSelected != nil {
			lv.onSelected(0)
		}
	}
}

func (lv *ListView) GetSelectedItem() (string, bool) {
	items, _ := lv.historyData.Get()
	if lv.selectedIndex < 0 || lv.selectedIndex >= len(items) {
		return "", false
	}
	return items[lv.selectedIndex], true
}

func (lv *ListView) RemoveItemAt(index int) {
	items, _ := lv.historyData.Get()
	if index < 0 || index >= len(items) {
		return
	}
	newItems := append(items[:index], items[index+1:]...)
	lv.historyData.Set(newItems)

	if lv.selectedIndex >= index {
		if lv.selectedIndex > 0 {
			lv.selectedIndex--
		} else if len(newItems) == 0 {
			lv.selectedIndex = -1
		}
	}
	if lv.selectedIndex >= 0 && lv.selectedIndex < len(newItems) {
		lv.list.Select(lv.selectedIndex)
	} else {
		lv.list.UnselectAll()
	}
}

// ---- 工具：手動省略號（rune 安全） ----

func ellipsis(s string, max int) string {
	if max <= 0 {
		return ""
	}
	if utf8.RuneCountInString(s) <= max {
		return s
	}
	// 預留 1 格給「…」
	runes := []rune(s)
	if max == 1 {
		return "…"
	}
	return string(runes[:max-1]) + "…"
}
