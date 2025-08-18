package components

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/widget"
)

func loadTrashcanIcon() fyne.Resource {
	if resource, err := fyne.LoadResourceFromPath("assets/trashcan.png"); err == nil {
		return resource
	}
	return nil
}

type DeleteButtonConfig struct {
	Size            fyne.Size
	IconColor       color.Color
	BackgroundColor color.Color
	HoverColor      color.Color
	Shape           string
	CornerRadius    float32
	Padding         float32
	IconResource    fyne.Resource
	IconText        string
}

var DefaultDeleteButtonConfig = DeleteButtonConfig{
	Size:            fyne.NewSize(40, 40),
	IconColor:       color.NRGBA{R: 255, G: 255, B: 255, A: 255}, // 白色圖示
	BackgroundColor: color.NRGBA{R: 200, G: 40, B: 30, A: 180},   // 深色紅色背景
	HoverColor:      color.NRGBA{R: 220, G: 50, B: 40, A: 220},   // hover 時更深
	Shape:           "square",
	CornerRadius:    8.0,
	Padding:         4.0,
	IconResource:    loadTrashcanIcon(), // 載入垃圾桶圖示
	IconText:        "🗑️",               // 備用文字圖示
}

type DeleteButton struct {
	widget.BaseWidget

	config   DeleteButtonConfig
	onTapped func()

	hovered bool
	pressed bool
}

func NewDeleteButton() *DeleteButton {
	return NewDeleteButtonWithConfig(DefaultDeleteButtonConfig)
}

func NewDeleteButtonWithConfig(config DeleteButtonConfig) *DeleteButton {
	btn := &DeleteButton{
		config: config,
	}
	btn.ExtendBaseWidget(btn)
	return btn
}

func (d *DeleteButton) SetOnTapped(callback func()) {
	d.onTapped = callback
}

func (d *DeleteButton) SetSize(size fyne.Size) {
	d.config.Size = size
	d.Resize(size)
}

func (d *DeleteButton) SetIconColor(color color.Color) {
	d.config.IconColor = color
	d.Refresh()
}

func (d *DeleteButton) SetHoverColor(color color.Color) {
	d.config.HoverColor = color
	d.Refresh()
}

func (d *DeleteButton) SetBackgroundColor(color color.Color) {
	d.config.BackgroundColor = color
	d.Refresh()
}

func (d *DeleteButton) SetIconResource(resource fyne.Resource) {
	d.config.IconResource = resource
	d.Refresh()
}

func (d *DeleteButton) SetIconText(text string) {
	d.config.IconText = text
	d.Refresh()
}

func (d *DeleteButton) GetConfig() DeleteButtonConfig {
	return d.config
}

func (d *DeleteButton) IsHovered() bool {
	return d.hovered
}
func (d *DeleteButton) IsPressed() bool {
	return d.pressed
}

func (d *DeleteButton) CreateRenderer() fyne.WidgetRenderer {
	return newDeleteButtonRenderer(d)
}

func (d *DeleteButton) MinSize() fyne.Size {
	return d.config.Size
}

func (d *DeleteButton) Tapped(_ *fyne.PointEvent) {
	if d.onTapped != nil {
		d.onTapped()
	}
}

func (d *DeleteButton) MouseIn(_ *fyne.PointEvent) {
	d.hovered = true
	d.Refresh()
}

func (d *DeleteButton) MouseOut() {
	d.hovered = false
	d.Refresh()
}

func (d *DeleteButton) MouseDown(_ *fyne.PointEvent) {
	d.pressed = true
	d.Refresh()
}

func (d *DeleteButton) MouseUp(_ *fyne.PointEvent) {
	d.pressed = false
	d.Refresh()
}

type deleteButtonRenderer struct {
	button     *DeleteButton
	background fyne.CanvasObject
	icon       fyne.CanvasObject // 可以是 *canvas.Text 或 *canvas.Image
	objects    []fyne.CanvasObject
}

func newDeleteButtonRenderer(button *DeleteButton) *deleteButtonRenderer {
	var bg fyne.CanvasObject
	if button.config.Shape == "circle" {
		bg = canvas.NewCircle(color.Transparent)
	} else {
		rect := canvas.NewRectangle(color.Transparent)
		if button.config.CornerRadius > 0 {
			rect.CornerRadius = button.config.CornerRadius
		}
		bg = rect
	}

	var icon fyne.CanvasObject
	if button.config.IconResource != nil {
		// 使用圖片資源
		img := canvas.NewImageFromResource(button.config.IconResource)
		img.FillMode = canvas.ImageFillContain
		img.SetMinSize(fyne.NewSize(float32(button.config.Size.Width)*0.6, float32(button.config.Size.Height)*0.6))
		icon = img
	} else {
		// 使用文字圖示
		textIcon := canvas.NewText(button.config.IconText, button.config.IconColor)
		textIcon.Alignment = fyne.TextAlignCenter
		textIcon.TextSize = float32(button.config.Size.Height) * 0.6
		icon = textIcon
	}

	renderer := &deleteButtonRenderer{
		button:     button,
		background: bg,
		icon:       icon,
		objects:    []fyne.CanvasObject{bg, icon},
	}

	renderer.Refresh()
	return renderer
}

func (r *deleteButtonRenderer) Objects() []fyne.CanvasObject {
	return r.objects
}

func (r *deleteButtonRenderer) Destroy() {}

func (r *deleteButtonRenderer) Layout(size fyne.Size) {
	padding := r.button.config.Padding
	bgSize := fyne.NewSize(size.Width-padding*2, size.Height-padding*2)
	bgPos := fyne.NewPos(padding, padding)

	r.background.Resize(bgSize)
	r.background.Move(bgPos)

	r.icon.Resize(size)
	r.icon.Move(fyne.NewPos(0, 0))
}

func (r *deleteButtonRenderer) MinSize() fyne.Size {
	return r.button.config.Size
}

func (r *deleteButtonRenderer) Refresh() {
	if r.button.hovered {
		if circle, ok := r.background.(*canvas.Circle); ok {
			circle.FillColor = r.button.config.HoverColor
		} else if rect, ok := r.background.(*canvas.Rectangle); ok {
			rect.FillColor = r.button.config.HoverColor
		}
	} else {
		if circle, ok := r.background.(*canvas.Circle); ok {
			circle.FillColor = r.button.config.BackgroundColor
		} else if rect, ok := r.background.(*canvas.Rectangle); ok {
			rect.FillColor = r.button.config.BackgroundColor
		}
	}

	// 更新圖示
	if textIcon, ok := r.icon.(*canvas.Text); ok {
		textIcon.Color = r.button.config.IconColor
		textIcon.TextSize = float32(r.button.config.Size.Height) * 0.6
	}

	r.background.Refresh()
	r.icon.Refresh()
}
