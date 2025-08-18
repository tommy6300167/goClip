package views

import (
	"image/color"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"

	"clipmini/models"
	"clipmini/utils"
)

type ListView struct {
	list          *widget.List
	historyData   binding.StringList
	filteredData  binding.StringList
	allItems      []string
	config        *models.AppConfig
	onSelected    func(int)
	onDelete      func(int)
	selectedIndex int
	searchQuery   string
}

func NewListView(config *models.AppConfig) *ListView {
	lv := &ListView{
		historyData:   binding.NewStringList(),
		filteredData:  binding.NewStringList(),
		allItems:      []string{},
		config:        config,
		selectedIndex: -1,
	}

	lv.list = widget.NewListWithData(
		lv.historyData,
		func() fyne.CanvasObject {
			contentIcon := widget.NewIcon(theme.FileTextIcon())

			contentText := canvas.NewText("", color.NRGBA{0xEE, 0xEE, 0xEE, 0xFF})
			contentText.TextStyle = fyne.TextStyle{Bold: true}
			contentText.Alignment = fyne.TextAlignLeading
			titlePad := container.NewPadded(contentText)

			timeText := canvas.NewText("", color.NRGBA{0xB0, 0xB0, 0xB0, 0xFF})
			timeText.TextSize = 12
			timeText.Alignment = fyne.TextAlignLeading
			timePad := container.NewPadded(timeText)

			gap := canvas.NewRectangle(color.Transparent)
			gap.SetMinSize(fyne.NewSize(1, 2))

			rightInner := container.NewVBox(titlePad, gap, timePad)
			rightContent := container.NewPadded(rightInner)

			deleteBtn := widget.NewButtonWithIcon("", theme.DeleteIcon(), nil)
			deleteBtn.Importance = widget.LowImportance
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

			ApplyRow(ui, rp, lv.config.MaxDisplayLength, loc)

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
	lv.allItems = lines
	lv.applySearch()
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
	lv.allItems = append([]string{line}, lv.allItems...)
	lv.applySearch()
	lv.SelectFirst()
}

func (lv *ListView) Clear() {
	lv.allItems = []string{}
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
	displayedItems, _ := lv.historyData.Get()
	if index < 0 || index >= len(displayedItems) {
		return
	}
	
	// 找到要刪除的項目
	itemToRemove := displayedItems[index]
	
	// 從所有項目中移除
	for i, item := range lv.allItems {
		if item == itemToRemove {
			lv.allItems = append(lv.allItems[:i], lv.allItems[i+1:]...)
			break
		}
	}
	
	// 重新應用搜尋
	lv.applySearch()
}

// ---- 搜尋功能 ----

// SetSearchQuery 設置搜尋查詢並過濾結果
func (lv *ListView) SetSearchQuery(query string) {
	lv.searchQuery = strings.TrimSpace(strings.ToLower(query))
	lv.applySearch()
}

// GetSearchQuery 獲取當前搜尋查詢
func (lv *ListView) GetSearchQuery() string {
	return lv.searchQuery
}

// applySearch 應用搜尋過濾
func (lv *ListView) applySearch() {
	if lv.searchQuery == "" {
		// 沒有搜尋查詢，顯示所有項目
		lv.historyData.Set(lv.allItems)
	} else {
		// 過濾項目
		filtered := []string{}
		for _, item := range lv.allItems {
			if lv.matchesSearch(item) {
				filtered = append(filtered, item)
			}
		}
		lv.historyData.Set(filtered)
	}
	
	// 重置選擇
	lv.selectedIndex = -1
	lv.list.UnselectAll()
}

// matchesSearch 檢查項目是否匹配搜尋查詢
func (lv *ListView) matchesSearch(item string) bool {
	if lv.searchQuery == "" {
		return true
	}
	
	// 解析項目內容
	parts := strings.SplitN(item, "\t", 3)
	if len(parts) < 2 {
		return false
	}
	
	searchLower := lv.searchQuery
	
	// 檢查內容
	if len(parts) == 3 && parts[2] == "IMAGE" {
		// 圖片項目：搜尋檔案名稱
		filename := strings.ToLower(parts[1])
		return strings.Contains(filename, searchLower)
	} else {
		// 文字項目：搜尋內容
		content := strings.ToLower(parts[1])
		return strings.Contains(content, searchLower)
	}
}
