package views

import (
	"path/filepath"
	"strings"
	"time"

	"clipmini/utils"
	u "clipmini/utils"
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

func ParseHistoryLine(line string, maxLen int, loc *time.Location) RowParsed {
	r := RowParsed{RawLine: line}

	first := strings.SplitN(line, "\t", 2)
	if len(first) > 0 {
		if t, ok := u.ParseTimeLocal(first[0], loc); ok {
			r.Time = t
		} else {
			r.RawTime = first[0]
		}
	}
	if len(first) < 2 {
		r.Kind, r.Tag, r.Title = RowText, "", utils.Ellipsis(line, maxLen)
		return r
	}

	rest := first[1]

	if strings.HasSuffix(rest, "\tIMAGE") {
		path := rest[:len(rest)-len("\tIMAGE")]
		name := filepath.Base(path)
		if name == "" {
			name = "圖片檔案"
		}
		r.Kind, r.Tag, r.Title = RowImage, "Image", utils.Ellipsis(name, maxLen)
		return r
	}

	r.Kind, r.Tag, r.Title = RowText, "Text", utils.Ellipsis(rest, maxLen)
	return r
}
