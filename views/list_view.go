package views

import (
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/widget"

	"clipmini/models"
	"clipmini/utils"
)

type ListView struct {
	list        *widget.List
	historyData binding.StringList
	config      *models.AppConfig
	onSelected  func(int)
	selectedIndex int
}

func NewListView(config *models.AppConfig) *ListView {
	lv := &ListView{
		historyData: binding.NewStringList(),
		config:      config,
		selectedIndex: -1,
	}
	
	lv.list = widget.NewListWithData(
		lv.historyData,
		func() fyne.CanvasObject { 
			label := widget.NewLabel("item")
			label.Wrapping = fyne.TextWrapOff
			return label
		},
		func(di binding.DataItem, co fyne.CanvasObject) {
			str, _ := di.(binding.String).Get()
			parts := strings.SplitN(str, "\t", 3)
			lbl := co.(*widget.Label)
			
			if len(parts) == 3 && parts[2] == "IMAGE" {
				displayText := parts[0] + " [IMAGE]"
				lbl.SetText(displayText)
			} else if len(parts) >= 2 {
				timestamp := parts[0]
				content := parts[1]
				truncatedContent := utils.TruncateText(content, config.MaxDisplayLength)
				displayText := timestamp + " " + truncatedContent
				lbl.SetText(displayText)
			} else {
				displayText := utils.TruncateText(str, config.MaxDisplayLength)
				lbl.SetText(displayText)
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

func (lv *ListView) GetWidget() *widget.List {
	return lv.list
}

func (lv *ListView) SetOnSelected(callback func(int)) {
	lv.onSelected = callback
}

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
	
	// 自動選取第一筆記錄
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
	
	// 自動選取新添加的第一筆項目
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
		// 觸發選擇回調以確保 UI 狀態同步
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