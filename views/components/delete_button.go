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
	Size:            fyne.NewSize(24, 24),
	IconColor:       color.NRGBA{255, 255, 255, 255},
	BackgroundColor: color.NRGBA{0, 0, 0, 0},        // 完全透明背景
	HoverColor:      color.NRGBA{255, 255, 255, 20},  // 輕微的白色透明 hover 效果
	Shape:           "square",
	CornerRadius:    6,
	Padding:         2,
	IconResource:    loadTrashcanIcon(),
	IconText:        "🗑️",
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
	icon       fyne.CanvasObject
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
		img := canvas.NewImageFromResource(button.config.IconResource)
		img.FillMode = canvas.ImageFillContain
		img.SetMinSize(fyne.NewSize(float32(button.config.Size.Width)*0.6, float32(button.config.Size.Height)*0.6))
		icon = img
	} else {

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

func (r *deleteButtonRenderer) Layout(alloc fyne.Size) {
	target := r.button.config.Size

	px := (alloc.Width - target.Width) / 2
	py := (alloc.Height - target.Height) / 2

	switch bg := r.background.(type) {
	case *canvas.Circle:
		bg.Resize(target)
		bg.Move(fyne.NewPos(px, py))
	case *canvas.Rectangle:
		bg.Resize(target)
		bg.Move(fyne.NewPos(px, py))
		bg.CornerRadius = r.button.config.CornerRadius
	}

	p := r.button.config.Padding
	iconW := target.Width - p*2
	iconH := target.Height - p*2
	if iconW < 1 {
		iconW = 1
	}
	if iconH < 1 {
		iconH = 1
	}
	iconSize := fyne.NewSize(iconW, iconH)
	iconPos := fyne.NewPos(px+p, py+p)

	if img, ok := r.icon.(*canvas.Image); ok {
		img.FillMode = canvas.ImageFillContain
		img.Resize(iconSize)
		img.Move(iconPos)
	} else if txt, ok := r.icon.(*canvas.Text); ok {
		txt.TextSize = iconH * 0.8
		txt.Alignment = fyne.TextAlignCenter
		txt.Resize(target)
		txt.Move(fyne.NewPos(px, py))
	} else {
		r.icon.Resize(iconSize)
		r.icon.Move(iconPos)
	}
}
func (r *deleteButtonRenderer) MinSize() fyne.Size {
	return r.button.config.Size
}

func (r *deleteButtonRenderer) Refresh() {
	bgCol := r.button.config.BackgroundColor
	if r.button.hovered {
		bgCol = r.button.config.HoverColor
	}

	if c, ok := r.background.(*canvas.Circle); ok {
		c.FillColor = bgCol
		c.Refresh()
	}
	if rect, ok := r.background.(*canvas.Rectangle); ok {
		rect.FillColor = bgCol
		rect.CornerRadius = r.button.config.CornerRadius
		rect.Refresh()
	}

	if txt, ok := r.icon.(*canvas.Text); ok {
		txt.Color = r.button.config.IconColor
		txt.Refresh()
	}
	if img, ok := r.icon.(*canvas.Image); ok {
		img.Refresh()
	}
}
