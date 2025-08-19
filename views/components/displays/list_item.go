package displays

import (
	"image/color"
	"path/filepath"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"

	"clipmini/models"
	"clipmini/utils"
)

type ListItem struct {
	container       *fyne.Container
	contentIcon     *widget.Icon
	contentText     *canvas.Text
	timeText        *canvas.Text
	deleteBtn       *widget.Button
	deleteBackground *canvas.Rectangle
	item            *models.ClipboardItem
	index           int
	onDelete        func(int)
	maxDisplayLength int
}

type ListItemConfig struct {
	MaxDisplayLength int
	OnDelete         func(int)
}

func NewListItem(config ListItemConfig) *ListItem {
	li := &ListItem{
		maxDisplayLength: config.MaxDisplayLength,
		onDelete:         config.OnDelete,
	}

	// 創建圖標
	li.contentIcon = widget.NewIcon(theme.DocumentIcon())

	// 創建內容文字
	li.contentText = canvas.NewText("", color.NRGBA{0xEE, 0xEE, 0xEE, 0xFF})
	li.contentText.TextStyle = fyne.TextStyle{Bold: true}
	li.contentText.Alignment = fyne.TextAlignLeading

	// 創建時間文字
	li.timeText = canvas.NewText("", color.NRGBA{0xB0, 0xB0, 0xB0, 0xFF})
	li.timeText.TextSize = 12
	li.timeText.Alignment = fyne.TextAlignLeading

	// 創建刪除按鈕
	li.deleteBtn = widget.NewButtonWithIcon("", theme.DeleteIcon(), func() {
		if li.onDelete != nil && li.index >= 0 {
			li.onDelete(li.index)
		}
	})
	li.deleteBtn.Importance = widget.LowImportance

	// 為刪除按鈕添加背景
	li.deleteBackground = canvas.NewRectangle(color.NRGBA{0x00, 0x00, 0x00, 0x66})
	deleteContainer := container.NewStack(li.deleteBackground, li.deleteBtn)

	// 創建文字容器
	titlePad := container.NewWithoutLayout(li.contentText)
	timePad := container.NewWithoutLayout(li.timeText)

	gap := canvas.NewRectangle(color.Transparent)
	rightInner := container.NewVBox(titlePad, gap, timePad)
	rightContent := container.NewPadded(rightInner)

	// 創建主容器
	li.container = container.NewBorder(nil, nil, li.contentIcon, deleteContainer, rightContent)

	return li
}

func (li *ListItem) UpdateItem(item *models.ClipboardItem, index int) {
	li.item = item
	li.index = index

	// 更新圖標
	if item.Type == models.ClipImage {
		li.contentIcon.SetResource(theme.MediaPhotoIcon())
	} else {
		li.contentIcon.SetResource(theme.DocumentIcon())
	}
	li.contentIcon.Refresh()

	// 更新內容文字
	var displayText string
	var timeTag string

	if item.Type == models.ClipImage {
		// 處理圖片項目
		filename := filepath.Base(item.FilePath)
		if filename == "" {
			filename = "圖片檔案"
		}
		displayText = filename
		timeTag = "Image"
	} else {
		// 處理文字項目
		displayText = item.Content
		timeTag = "Text"
	}

	// 清理和截斷文字
	displayText = utils.CleanInvisibleChars(displayText)
	if li.maxDisplayLength > 0 {
		displayText = utils.Ellipsis(displayText, li.maxDisplayLength)
	}

	li.contentText.Text = displayText
	li.contentText.Refresh()

	// 更新時間文字
	loc := utils.GetTaipeiLocation()
	timeStr := utils.FormatTimestamp(item.Timestamp, loc)
	if timeTag != "" {
		timeStr += " · " + timeTag
	}
	
	li.timeText.Text = timeStr
	li.timeText.Refresh()
}

func (li *ListItem) SetOnDelete(callback func(int)) {
	li.onDelete = callback
	li.deleteBtn.OnTapped = func() {
		if li.onDelete != nil && li.index >= 0 {
			li.onDelete(li.index)
		}
	}
}

func (li *ListItem) GetWidget() *fyne.Container {
	return li.container
}

func (li *ListItem) GetContentText() *canvas.Text {
	return li.contentText
}

func (li *ListItem) GetTimeText() *canvas.Text {
	return li.timeText
}

func (li *ListItem) GetContentIcon() *widget.Icon {
	return li.contentIcon
}

func (li *ListItem) GetDeleteButton() *widget.Button {
	return li.deleteBtn
}

func (li *ListItem) GetIndex() int {
	return li.index
}

func (li *ListItem) GetItem() *models.ClipboardItem {
	return li.item
}