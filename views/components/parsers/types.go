package parsers

import (
	"time"

	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/widget"
)

// RowKind 定義列表項類型
type RowKind int

const (
	RowText RowKind = iota
	RowImage
)

// RowParsed 解析後的列表項數據
type RowParsed struct {
	RawLine string
	Kind    RowKind
	Title   string
	Tag     string
	Time    time.Time
	RawTime string
}

// RowSubviews 列表項的 UI 組件
type RowSubviews struct {
	ContentIcon *widget.Icon
	ContentText *canvas.Text
	TimeText    *canvas.Text
	DeleteBtn   *widget.Button
}