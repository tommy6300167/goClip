// Package parsers 提供列表項的解析、渲染和事件處理功能
//
// 主要功能：
// - 解析歷史記錄行數據 (parser.go)
// - 管理列表項 UI 組件 (helpers.go)
// - 定義數據類型和結構 (types.go)
//
// 使用範例：
//
//	// 解析歷史記錄
//	rp := parsers.ParseHistoryLine(line, maxLen, loc)
//	
//	// 找到 UI 組件
//	ui, ok := parsers.FindRowSubviews(container)
//	
//	// 應用數據到 UI
//	parsers.ApplyRow(ui, rp, maxLen, loc)
//
package parsers