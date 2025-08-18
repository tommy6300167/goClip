package utils

import (
	"time"
)

func FormatTimestamp(t time.Time, location *time.Location) string {
	return t.In(location).Format("2006-01-02 15:04:05")
}

func FormatTimestampForDisplay(t time.Time, location *time.Location) string {
	localTime := t.In(location)
	now := time.Now().In(location)
	
	// 檢查是否為今天
	if localTime.Year() == now.Year() && localTime.YearDay() == now.YearDay() {
		return "今天"
	}
	
	// 其他日期顯示年/月/日格式
	return localTime.Format("2006/01/02")
}

func GetTaipeiLocation() *time.Location {
	loc, _ := time.LoadLocation("Asia/Taipei")
	return loc
}