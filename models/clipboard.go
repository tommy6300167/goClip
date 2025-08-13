package models

import "time"

type ClipboardItem struct {
	ID        string
	Timestamp time.Time
	Content   string
	Type      ClipType
	FilePath  string // for images
}

type ClipType int

const (
	ClipText ClipType = iota
	ClipImage
)

func (t ClipType) String() string {
	switch t {
	case ClipText:
		return "TEXT"
	case ClipImage:
		return "IMAGE"
	default:
		return "UNKNOWN"
	}
}

func NewTextItem(content string) *ClipboardItem {
	return &ClipboardItem{
		Timestamp: time.Now(),
		Content:   content,
		Type:      ClipText,
	}
}

func NewImageItem(filePath string) *ClipboardItem {
	return &ClipboardItem{
		Timestamp: time.Now(),
		Type:      ClipImage,
		FilePath:  filePath,
	}
}