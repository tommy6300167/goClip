package displays

import (
	"fyne.io/fyne/v2/widget"
)

type StatusType int

const (
	StatusNormal StatusType = iota
	StatusEditing
	StatusError
	StatusSuccess
)

type StatusBar struct {
	label         *widget.Label
	currentStatus StatusType
	defaultText   string
}

func NewStatusBar(defaultText string) *StatusBar {
	sb := &StatusBar{
		label:         widget.NewLabel(defaultText),
		currentStatus: StatusNormal,
		defaultText:   defaultText,
	}
	return sb
}

func (sb *StatusBar) GetWidget() *widget.Label {
	return sb.label
}

func (sb *StatusBar) SetStatus(message string, statusType StatusType) {
	sb.label.SetText(message)
	sb.currentStatus = statusType
}

func (sb *StatusBar) SetNormalMode() {
	sb.SetStatus(sb.defaultText, StatusNormal)
}

func (sb *StatusBar) SetEditingMode() {
	sb.SetStatus("Use Shift + Enter to add a line break.", StatusEditing)
}

func (sb *StatusBar) SetErrorStatus(message string) {
	sb.SetStatus(message, StatusError)
}

func (sb *StatusBar) SetSuccessStatus(message string) {
	sb.SetStatus(message, StatusSuccess)
}

func (sb *StatusBar) GetCurrentStatus() StatusType {
	return sb.currentStatus
}

func (sb *StatusBar) UpdateText(message string) {
	sb.label.SetText(message)
}