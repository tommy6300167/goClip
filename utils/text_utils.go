package utils

import (
	"strings"
)

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

func ReverseSlice(ss []string) []string {
	n := len(ss)
	out := make([]string, n)
	for i := range ss {
		out[n-1-i] = ss[i]
	}
	return out
}