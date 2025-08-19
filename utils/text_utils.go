package utils

import (
	"strings"
	"time"
	"unicode/utf8"
)

const TimeLayout = "2006-01-02 15:04:05"

func TruncateText(s string, max int) string {
	s = strings.ReplaceAll(s, "\n", " ")
	s = strings.ReplaceAll(s, "\t", " ")
	s = strings.TrimSpace(s)
	r := []rune(s)
	if len(r) > max {
		return string(r[:max]) + "..."
	}
	return s
}


func ParseTimeLocal(ts string, loc *time.Location) (time.Time, bool) {
	t, err := time.Parse(TimeLayout, ts)
	if err != nil {
		return time.Time{}, false
	}
	return t.In(loc), true
}

func FormatTimeWithTag(t time.Time, tag string) string {
	base := t.Format(TimeLayout)
	if tag == "" {
		return base
	}
	return base + " · " + tag
}

func Ellipsis(s string, max int) string {
	if max <= 0 {
		return ""
	}
	if utf8.RuneCountInString(s) <= max {
		return s
	}

	runes := []rune(s)
	if max == 1 {
		return "…"
	}
	return string(runes[:max-1]) + "…"
}
