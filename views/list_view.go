package views

import (
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/widget"

	"clipmini/models"
	"clipmini/utils"
	"clipmini/views/components/displays"
	"clipmini/views/components/parsers"
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
	itemCache     map[string]int
	lastUpdate    int64
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

			listItem := displays.NewListItem(displays.ListItemConfig{
				MaxDisplayLength: lv.config.MaxDisplayLength,
				OnDelete:         lv.onDelete,
			})
			return listItem.GetWidget()
		},
		func(di binding.DataItem, co fyne.CanvasObject) {
			line, _ := di.(binding.String).Get()

			ui, ok := parsers.FindRowSubviews(co)
			if !ok {
				return
			}

			loc := utils.GetTaipeiLocation()
			rp := parsers.ParseHistoryLine(line, lv.config.MaxDisplayLength, loc)

			parsers.ApplyRow(ui, rp, lv.config.MaxDisplayLength, loc)

			parsers.BindDeleteByLine(ui.DeleteBtn, lv.historyData, line, func(idx int) {
				if lv.onDelete != nil {
					lv.onDelete(idx)
				}
			})
		},
	)

	lv.list.OnSelected = func(id widget.ListItemID) {
		lv.selectedIndex = id
		if lv.onSelected != nil {
			originalIndex := lv.getOriginalIndex(id)
			lv.onSelected(originalIndex)
		}
	}

	return lv
}

// ---- 事件綁定/資料操作 ----

func (lv *ListView) GetWidget() *widget.List { return lv.list }

func (lv *ListView) SetOnSelected(callback func(int)) { lv.onSelected = callback }

func (lv *ListView) SetOnDelete(callback func(int)) { lv.onDelete = callback }

func (lv *ListView) LoadFromHistory(items []*models.ClipboardItem) {
	if lv.shouldFullRebuild(items) {
		lv.fullRebuild(items)
	} else {
		lv.incrementalUpdate(items)
	}
}

func (lv *ListView) shouldFullRebuild(items []*models.ClipboardItem) bool {
	if abs(len(items)-len(lv.allItems)) > len(lv.allItems)/4 {
		return true
	}
	if len(lv.allItems) == 0 {
		return true
	}
	return false
}

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

	if lv.hasChanges(newLines) {
		lv.allItems = newLines
		lv.itemCache = newCache
		lv.applySearch()
	}
}

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

	lv.allItems = append([]string{line}, lv.allItems...)

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
			originalIndex := lv.getOriginalIndex(0)
			lv.onSelected(originalIndex)
		}
	}
}

func (lv *ListView) RemoveItemAt(index int) {
	displayedItems, _ := lv.historyData.Get()
	if index < 0 || index >= len(displayedItems) {
		return
	}

	itemToRemove := displayedItems[index]

	if cachedIndex, exists := lv.itemCache[itemToRemove]; exists {
		lv.allItems = append(lv.allItems[:cachedIndex], lv.allItems[cachedIndex+1:]...)

		delete(lv.itemCache, itemToRemove)
		for key, idx := range lv.itemCache {
			if idx > cachedIndex {
				lv.itemCache[key] = idx - 1
			}
		}
	} else {
		for i, item := range lv.allItems {
			if item == itemToRemove {
				lv.allItems = append(lv.allItems[:i], lv.allItems[i+1:]...)
				break
			}
		}
		lv.rebuildCache()
	}

	lv.applySearch()
}

func (lv *ListView) rebuildCache() {
	lv.itemCache = make(map[string]int, len(lv.allItems))
	for i, item := range lv.allItems {
		lv.itemCache[item] = i
	}
}

func (lv *ListView) SetSearchQuery(query string) {
	newQuery := strings.TrimSpace(strings.ToLower(query))

	if lv.searchQuery == newQuery {
		return
	}

	wasEmpty := lv.searchQuery == ""
	lv.searchQuery = newQuery
	lv.applySearch()

	if (wasEmpty && newQuery != "") || (!wasEmpty && newQuery == "") {
		if items, _ := lv.historyData.Get(); len(items) > 0 {
			fyne.Do(func() {
				lv.SelectFirst()
			})
		}
	}
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

func (lv *ListView) getOriginalIndex(filteredIndex int) int {
	displayedItems, _ := lv.historyData.Get()
	if filteredIndex < 0 || filteredIndex >= len(displayedItems) {
		return -1
	}

	if lv.searchQuery == "" {
		return filteredIndex
	}

	selectedItem := displayedItems[filteredIndex]

	if originalIndex, exists := lv.itemCache[selectedItem]; exists {
		return originalIndex
	}

	for i, item := range lv.allItems {
		if item == selectedItem {
			return i
		}
	}

	return -1
}
