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
	
	// 清理不可見字符和控制字符
	s = CleanInvisibleChars(s)
	
	if utf8.RuneCountInString(s) <= max {
		return s
	}

	runes := []rune(s)
	if max == 1 {
		return "…"
	}
	return string(runes[:max-1]) + "…"
}

// CleanInvisibleChars 清理不可見字符和控制字符
func CleanInvisibleChars(s string) string {
	// 替換常見的不可見字符和控制字符
	s = strings.ReplaceAll(s, "\n", " ")
	s = strings.ReplaceAll(s, "\r", " ")
	s = strings.ReplaceAll(s, "\t", " ")
	s = strings.ReplaceAll(s, "\v", " ") // 垂直制表符
	s = strings.ReplaceAll(s, "\f", " ") // 換頁符
	s = strings.ReplaceAll(s, "\u00A0", " ") // 不間斷空格
	s = strings.ReplaceAll(s, "\u2000", " ") // En Quad
	s = strings.ReplaceAll(s, "\u2001", " ") // Em Quad
	s = strings.ReplaceAll(s, "\u2002", " ") // En Space
	s = strings.ReplaceAll(s, "\u2003", " ") // Em Space
	s = strings.ReplaceAll(s, "\u2004", " ") // Three-Per-Em Space
	s = strings.ReplaceAll(s, "\u2005", " ") // Four-Per-Em Space
	s = strings.ReplaceAll(s, "\u2006", " ") // Six-Per-Em Space
	s = strings.ReplaceAll(s, "\u2007", " ") // Figure Space
	s = strings.ReplaceAll(s, "\u2008", " ") // Punctuation Space
	s = strings.ReplaceAll(s, "\u2009", " ") // Thin Space
	s = strings.ReplaceAll(s, "\u200A", " ") // Hair Space
	s = strings.ReplaceAll(s, "\u200B", "")  // 零寬空格
	s = strings.ReplaceAll(s, "\u200C", "")  // 零寬非連字符
	s = strings.ReplaceAll(s, "\u200D", "")  // 零寬連字符
	s = strings.ReplaceAll(s, "\u2060", "")  // 單詞連接符
	s = strings.ReplaceAll(s, "\uFEFF", "")  // 零寬非斷空格(BOM)
	
	// 清理多餘的空格
	s = strings.TrimSpace(s)
	
	return s
}
