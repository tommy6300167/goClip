package models

import (
	"fmt"
	"strings"
	"time"
)

type History struct {
	Items       []*ClipboardItem
	MaxItems    int
	Location    *time.Location
}

func NewHistory(maxItems int) *History {
	loc, _ := time.LoadLocation("Asia/Taipei")
	return &History{
		Items:    make([]*ClipboardItem, 0),
		MaxItems: maxItems,
		Location: loc,
	}
}

func (h *History) Add(item *ClipboardItem) {
	h.Items = append([]*ClipboardItem{item}, h.Items...)
	
	if len(h.Items) > h.MaxItems {
		h.Items = h.Items[:h.MaxItems]
	}
}

func (h *History) GetItems() []*ClipboardItem {
	return h.Items
}

func (h *History) Clear() {
	h.Items = make([]*ClipboardItem, 0)
}

func (h *History) GetItem(index int) *ClipboardItem {
	if index < 0 || index >= len(h.Items) {
		return nil
	}
	return h.Items[index]
}

func (h *History) RemoveItem(index int) *ClipboardItem {
	if index < 0 || index >= len(h.Items) {
		return nil
	}
	
	removedItem := h.Items[index]
	h.Items = append(h.Items[:index], h.Items[index+1:]...)
	return removedItem
}

func (h *History) UpdateItem(index int, newContent string) bool {
	if index < 0 || index >= len(h.Items) {
		return false
	}
	
	item := h.Items[index]
	if item.Type == ClipText {
		item.Content = newContent
		return true
	}
	
	return false
}

func (h *History) ToFileFormat() []string {
	lines := make([]string, len(h.Items))
	for i := len(h.Items) - 1; i >= 0; i-- {
		item := h.Items[i]
		timestamp := item.Timestamp.In(h.Location).Format("2006-01-02 15:04:05")
		
		if item.Type == ClipImage {
			lines[len(h.Items)-1-i] = fmt.Sprintf("%s\t%s\t%s", timestamp, item.FilePath, item.Type.String())
		} else {
			lines[len(h.Items)-1-i] = fmt.Sprintf("%s\t%s", timestamp, item.Content)
		}
	}
	return lines
}

func (h *History) FromFileFormat(lines []string) {
	h.Items = make([]*ClipboardItem, 0)
	for i := len(lines) - 1; i >= 0; i-- {
		line := strings.TrimSpace(lines[i])
		if line == "" {
			continue
		}
		
		parts := strings.SplitN(line, "\t", 3)
		if len(parts) < 2 {
			continue
		}
		
		timestamp, err := time.ParseInLocation("2006-01-02 15:04:05", parts[0], h.Location)
		if err != nil {
			continue
		}
		
		item := &ClipboardItem{
			Timestamp: timestamp,
		}
		
		if len(parts) == 3 && parts[2] == "IMAGE" {
			item.Type = ClipImage
			item.FilePath = parts[1]
		} else {
			item.Type = ClipText
			item.Content = parts[1]
		}
		
		h.Items = append(h.Items, item)
	}
}