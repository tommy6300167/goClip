package views

import (
	"image/color"

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
			contentIcon := components.NewContentIcon(false)
			contentLabel := widget.NewLabelWithStyle("", fyne.TextAlignLeading, fyne.TextStyle{Bold: true})
			contentLabel.Wrapping = fyne.TextWrapOff
			contentLabel.Truncation = fyne.TextTruncateEllipsis
			titlePad := container.NewPadded(contentLabel)
			timeText := canvas.NewText("", color.NRGBA{0xB0, 0xB0, 0xB0, 0xFF})
			timeText.TextSize = 12
			timeText.Alignment = fyne.TextAlignLeading
			timePad := container.NewPadded(timeText)
			gap := canvas.NewRectangle(color.Transparent)
			gap.SetMinSize(fyne.NewSize(1, 2))

			rightInner := container.NewVBox(titlePad, gap, timePad)
			rightContent := container.NewPadded(rightInner)
			deleteBtn := components.NewDeleteButton()
			return container.NewBorder(nil, nil, contentIcon, deleteBtn, rightContent)
		},
		func(di binding.DataItem, co fyne.CanvasObject) {
			line, _ := di.(binding.String).Get()
			ui, ok := FindRowSubviews(co)
			if !ok {
				return
			}
			loc := utils.GetTaipeiLocation()

			rp := ParseHistoryLine(line, lv.config.MaxDisplayLength, loc)
			ApplyRow(ui, rp, loc)
			BindDeleteByLine(ui.DeleteBtn, lv.historyData, line, func(idx int) {
				if lv.onDelete != nil {
					lv.onDelete(idx)
				}
			})
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
