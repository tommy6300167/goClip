package parsers

import (
	"time"

	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/widget"
)

type RowKind int

const (
	RowText RowKind = iota
	RowImage
)

type RowParsed struct {
	RawLine string
	Kind    RowKind
	Title   string
	Tag     string
	Time    time.Time
	RawTime string
}

type RowSubviews struct {
	ContentIcon *widget.Icon
	ContentText *canvas.Text
	TimeText    *canvas.Text
	DeleteBtn   *widget.Button
}
