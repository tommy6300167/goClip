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
	parts := strings.SplitN(line, "\t", 3)
	r := RowParsed{RawLine: line}

	if len(parts) > 0 {
		if t, ok := u.ParseTimeLocal(parts[0], loc); ok {
			r.Time = t
		} else {
			r.RawTime = parts[0]
		}
	}

	if len(parts) == 3 && parts[2] == "IMAGE" {
		name := filepath.Base(parts[1])
		if name == "" {
			name = "圖片檔案"
		}
		r.Kind, r.Tag, r.Title = RowImage, "圖片", utils.Ellipsis(name, maxLen)
		return r
	}

	if len(parts) >= 2 {
		clean := u.CleanContent(parts[1])
		r.Kind, r.Tag, r.Title = RowText, "文字", utils.Ellipsis(clean, maxLen)
		return r
	}

	r.Kind, r.Tag, r.Title = RowText, "", utils.Ellipsis(line, maxLen)
	return r
}
