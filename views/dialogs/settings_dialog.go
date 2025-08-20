package dialogs

import (
	"errors"
	"strconv"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"

	"clipmini/models"
)

type SettingsDialog struct {
	dialog       dialog.Dialog
	window       fyne.Window
	config       *models.AppConfig
	onSave       func(*models.AppConfig)
	undoEntry    *widget.Entry
	widthEntry   *widget.Entry
	heightEntry  *widget.Entry
	historyEntry *widget.Entry
	displayEntry *widget.Entry
	pollingEntry *widget.Entry
}

func NewSettingsDialog(window fyne.Window, config *models.AppConfig) *SettingsDialog {
	sd := &SettingsDialog{
		window: window,
		config: config,
	}

	sd.buildDialog(window)
	return sd
}

func (sd *SettingsDialog) buildDialog(window fyne.Window) {
	sd.undoEntry = widget.NewEntry()
	sd.undoEntry.SetText(strconv.Itoa(sd.config.UndoCountdownSecs))

	sd.widthEntry = widget.NewEntry()
	sd.widthEntry.SetText(strconv.Itoa(sd.config.WindowWidth))

	sd.heightEntry = widget.NewEntry()
	sd.heightEntry.SetText(strconv.Itoa(sd.config.WindowHeight))

	sd.historyEntry = widget.NewEntry()
	sd.historyEntry.SetText(strconv.Itoa(sd.config.MaxHistoryItems))

	sd.displayEntry = widget.NewEntry()
	sd.displayEntry.SetText(strconv.Itoa(sd.config.MaxDisplayLength))

	sd.pollingEntry = widget.NewEntry()
	sd.pollingEntry.SetText(strconv.Itoa(sd.config.PollingInterval))

	form := &widget.Form{
		Items: []*widget.FormItem{
			{Text: "刪除倒數 (秒)", Widget: sd.undoEntry},
			{Text: "視窗寬度", Widget: sd.widthEntry},
			{Text: "視窗高度", Widget: sd.heightEntry},
			{Text: "最大歷史項目", Widget: sd.historyEntry},
			{Text: "最大顯示長度", Widget: sd.displayEntry},
			{Text: "輪詢間隔 (毫秒)", Widget: sd.pollingEntry},
		},
		OnSubmit: func() {
			if sd.validateAndSave() {
				sd.dialog.Hide()
			}
		},
		OnCancel: func() {
			sd.dialog.Hide()
		},
	}

	form.SubmitText = "保存"
	form.CancelText = "取消"

	resetBtn := widget.NewButtonWithIcon("重置為默認值", theme.ViewRefreshIcon(), func() {
		sd.resetToDefaults()
	})
	resetBtn.Importance = widget.LowImportance

	content := container.NewVBox(
		widget.NewLabel("設置"),
		widget.NewSeparator(),
		form,
		widget.NewSeparator(),
		container.NewCenter(resetBtn),
	)

	sd.dialog = dialog.NewCustom("設置", "關閉", content, window)
	sd.dialog.Resize(fyne.NewSize(400, 500))
}

func (sd *SettingsDialog) validateAndSave() bool {
	undoSecs, err := strconv.Atoi(strings.TrimSpace(sd.undoEntry.Text))
	if err != nil || undoSecs < 1 || undoSecs > 60 {
		dialog.ShowError(errors.New("刪除倒數必須是 1-60 之間的整數"), sd.window)
		return false
	}

	width, err := strconv.Atoi(strings.TrimSpace(sd.widthEntry.Text))
	if err != nil || width < 400 || width > 2000 {
		dialog.ShowError(errors.New("視窗寬度必須是 400-2000 之間的整數"), sd.window)
		return false
	}

	height, err := strconv.Atoi(strings.TrimSpace(sd.heightEntry.Text))
	if err != nil || height < 300 || height > 1500 {
		dialog.ShowError(errors.New("視窗高度必須是 300-1500 之間的整數"), sd.window)
		return false
	}

	maxHistory, err := strconv.Atoi(strings.TrimSpace(sd.historyEntry.Text))
	if err != nil || maxHistory < 10 || maxHistory > 1000 {
		dialog.ShowError(errors.New("最大歷史項目必須是 10-1000 之間的整數"), sd.window)
		return false
	}

	maxDisplay, err := strconv.Atoi(strings.TrimSpace(sd.displayEntry.Text))
	if err != nil || maxDisplay < 10 || maxDisplay > 200 {
		dialog.ShowError(errors.New("最大顯示長度必須是 10-200 之間的整數"), sd.window)
		return false
	}

	polling, err := strconv.Atoi(strings.TrimSpace(sd.pollingEntry.Text))
	if err != nil || polling < 100 || polling > 5000 {
		dialog.ShowError(errors.New("輪詢間隔必須是 100-5000 之間的整數"), sd.window)
		return false
	}

	newConfig := &models.AppConfig{
		UndoCountdownSecs: undoSecs,
		WindowWidth:       width,
		WindowHeight:      height,
		MaxHistoryItems:   maxHistory,
		MaxDisplayLength:  maxDisplay,
		PollingInterval:   polling,
		LogDirPath:        sd.config.LogDirPath,
		LogFilePath:       sd.config.LogFilePath,
		ImageDirPath:      sd.config.ImageDirPath,
		ConfigFilePath:    sd.config.ConfigFilePath,
		AuthFilePath:      sd.config.AuthFilePath,
		AuthApiEndpoint:   sd.config.AuthApiEndpoint,
	}

	if sd.onSave != nil {
		sd.onSave(newConfig)
	}

	return true
}

func (sd *SettingsDialog) resetToDefaults() {
	sd.undoEntry.SetText(strconv.Itoa(models.DefaultUndoCountdownSecs))
	sd.widthEntry.SetText(strconv.Itoa(models.DefaultWindowWidth))
	sd.heightEntry.SetText(strconv.Itoa(models.DefaultWindowHeight))
	sd.historyEntry.SetText(strconv.Itoa(models.DefaultMaxHistoryItems))
	sd.displayEntry.SetText(strconv.Itoa(models.DefaultMaxDisplayLength))
	sd.pollingEntry.SetText(strconv.Itoa(models.DefaultPollingInterval))
}

func (sd *SettingsDialog) SetOnSave(callback func(*models.AppConfig)) {
	sd.onSave = callback
}

func (sd *SettingsDialog) Show() {
	sd.dialog.Show()
}
