package utils

import (
	"time"
)

func FormatTimestamp(t time.Time, location *time.Location) string {
	return t.In(location).Format("2006-01-02 15:04:05")
}

func GetTaipeiLocation() *time.Location {
	loc, _ := time.LoadLocation("Asia/Taipei")
	return loc
}