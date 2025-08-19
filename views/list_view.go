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
	// 增量更新相關
	itemCache  map[string]int // 項目內容到索引的映射
	lastUpdate int64          // 最後更新時間戳
}

func NewListView(config *models.AppConfig) *ListView {
	lv := &ListView{
		historyData:   binding.NewStringList(),
		filteredData:  binding.NewStringList(),
		allItems:      make([]string, 0, config.MaxHistoryItems),
		config:        config,
		selectedIndex: -1,
		itemCache:     make(map[string]int, config.MaxHistoryItems),
		lastUpdate:    0,
	}

	lv.list = widget.NewListWithData(
		lv.historyData,
		func() fyne.CanvasObject {
			contentIcon := widget.NewIcon(theme.DocumentIcon())

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
	// 增量更新：只在必要時重建
	if lv.shouldFullRebuild(items) {
		lv.fullRebuild(items)
	} else {
		lv.incrementalUpdate(items)
	}
}

// shouldFullRebuild 判斷是否需要完整重建
func (lv *ListView) shouldFullRebuild(items []*models.ClipboardItem) bool {
	// 如果項目數量差異很大，進行完整重建
	if abs(len(items)-len(lv.allItems)) > len(lv.allItems)/4 {
		return true
	}
	// 如果是初次載入，進行完整重建
	if len(lv.allItems) == 0 {
		return true
	}
	return false
}

// fullRebuild 完整重建列表
func (lv *ListView) fullRebuild(items []*models.ClipboardItem) {
	lines := make([]string, len(items))
	lv.itemCache = make(map[string]int, len(items))

	for i, item := range items {
		timestamp := utils.FormatTimestamp(item.Timestamp, utils.GetTaipeiLocation())
		var line string
		if item.Type == models.ClipImage {
			line = timestamp + "\t" + item.FilePath + "\t" + item.Type.String()
		} else {
			line = timestamp + "\t" + item.Content
		}
		lines[i] = line
		lv.itemCache[line] = i
	}

	lv.allItems = lines
	lv.applySearch()
	if len(lines) > 0 {
		lv.SelectFirst()
	}
}

// incrementalUpdate 增量更新列表
func (lv *ListView) incrementalUpdate(items []*models.ClipboardItem) {
	newLines := make([]string, len(items))
	newCache := make(map[string]int, len(items))

	for i, item := range items {
		timestamp := utils.FormatTimestamp(item.Timestamp, utils.GetTaipeiLocation())
		var line string
		if item.Type == models.ClipImage {
			line = timestamp + "\t" + item.FilePath + "\t" + item.Type.String()
		} else {
			line = timestamp + "\t" + item.Content
		}
		newLines[i] = line
		newCache[line] = i
	}

	// 檢查是否有實質變化
	if lv.hasChanges(newLines) {
		lv.allItems = newLines
		lv.itemCache = newCache
		lv.applySearch()
	}
}

// hasChanges 檢查列表是否有變化
func (lv *ListView) hasChanges(newLines []string) bool {
	if len(newLines) != len(lv.allItems) {
		return true
	}

	for i, line := range newLines {
		if i >= len(lv.allItems) || lv.allItems[i] != line {
			return true
		}
	}
	return false
}

// abs 計算絕對值
func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func (lv *ListView) PrependItem(item *models.ClipboardItem) {
	timestamp := utils.FormatTimestamp(item.Timestamp, utils.GetTaipeiLocation())
	var line string
	if item.Type == models.ClipImage {
		line = timestamp + "\t" + item.FilePath + "\t" + item.Type.String()
	} else {
		line = timestamp + "\t" + item.Content
	}

	// 增量更新：只添加新項目
	lv.allItems = append([]string{line}, lv.allItems...)

	// 更新快取映射
	for key, index := range lv.itemCache {
		lv.itemCache[key] = index + 1
	}
	lv.itemCache[line] = 0

	lv.applySearch()
	lv.SelectFirst()
}

func (lv *ListView) Clear() {
	lv.allItems = make([]string, 0, lv.config.MaxHistoryItems)
	lv.itemCache = make(map[string]int, lv.config.MaxHistoryItems)
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

	itemToRemove := displayedItems[index]

	// 增量更新：從快取中查找並移除
	if cachedIndex, exists := lv.itemCache[itemToRemove]; exists {
		// 從 allItems 中移除
		lv.allItems = append(lv.allItems[:cachedIndex], lv.allItems[cachedIndex+1:]...)

		// 更新快取映射
		delete(lv.itemCache, itemToRemove)
		for key, idx := range lv.itemCache {
			if idx > cachedIndex {
				lv.itemCache[key] = idx - 1
			}
		}
	} else {
		// 回退到原方法（應該很少觸發）
		for i, item := range lv.allItems {
			if item == itemToRemove {
				lv.allItems = append(lv.allItems[:i], lv.allItems[i+1:]...)
				break
			}
		}
		// 重建快取
		lv.rebuildCache()
	}

	lv.applySearch()
}

// rebuildCache 重建項目快取
func (lv *ListView) rebuildCache() {
	lv.itemCache = make(map[string]int, len(lv.allItems))
	for i, item := range lv.allItems {
		lv.itemCache[item] = i
	}
}

func (lv *ListView) SetSearchQuery(query string) {
	newQuery := strings.TrimSpace(strings.ToLower(query))

	// 如果搜尋查詢沒有變化，跳過處理
	if lv.searchQuery == newQuery {
		return
	}

	lv.searchQuery = newQuery
	lv.applySearch()
}

func (lv *ListView) GetSearchQuery() string {
	return lv.searchQuery
}

func (lv *ListView) applySearch() {
	if lv.searchQuery == "" {
		lv.historyData.Set(lv.allItems)
	} else {
		filtered := make([]string, 0, len(lv.allItems))
		for _, item := range lv.allItems {
			if lv.matchesSearch(item) {
				filtered = append(filtered, item)
			}
		}
		lv.historyData.Set(filtered)
	}
	lv.selectedIndex = -1
	lv.list.UnselectAll()
}

func (lv *ListView) matchesSearch(item string) bool {
	if lv.searchQuery == "" {
		return true
	}

	// 快速檢查：如果查詢字串較短，直接在原字串中搜尋
	if len(lv.searchQuery) <= 3 {
		return strings.Contains(strings.ToLower(item), lv.searchQuery)
	}

	parts := strings.SplitN(item, "\t", 3)
	if len(parts) < 2 {
		return false
	}

	searchLower := lv.searchQuery

	if len(parts) == 3 && parts[2] == "IMAGE" {
		filename := strings.ToLower(parts[1])
		return strings.Contains(filename, searchLower)
	} else {
		content := strings.ToLower(parts[1])
		return strings.Contains(content, searchLower)
	}
}
