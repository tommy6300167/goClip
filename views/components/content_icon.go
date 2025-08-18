package components

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/widget"
)

func loadTextIcon() fyne.Resource {
	if resource, err := fyne.LoadResourceFromPath("assets/text.png"); err == nil {
		return resource
	}
	return nil
}

func loadImageIcon() fyne.Resource {
	if resource, err := fyne.LoadResourceFromPath("assets/image.png"); err == nil {
		return resource
	}
	return nil
}

type ContentIconConfig struct {
	Size              fyne.Size
	IconColor         color.Color
	BackgroundColor   color.Color
	Shape             string
	CornerRadius      float32
	Padding           float32
	TextIconResource  fyne.Resource
	ImageIconResource fyne.Resource
	TextIconText      string
	ImageIconText     string
}

var DefaultContentIconConfig = ContentIconConfig{
	Size:              fyne.NewSize(20, 20),
	IconColor:         color.NRGBA{R: 0, G: 0, B: 0, A: 255},
	BackgroundColor:   color.NRGBA{R: 240, G: 240, B: 240, A: 255},
	Shape:             "square",
	CornerRadius:      6.0,
	Padding:           2.0,
	TextIconResource:  loadTextIcon(),
	ImageIconResource: loadImageIcon(),
	TextIconText:      "📄",
	ImageIconText:     "🖼️",
}

type ContentIcon struct {
	widget.BaseWidget

	config  ContentIconConfig
	isImage bool
}

func NewContentIcon(isImage bool) *ContentIcon {
	return NewContentIconWithConfig(DefaultContentIconConfig, isImage)
}

func NewContentIconWithConfig(config ContentIconConfig, isImage bool) *ContentIcon {
	icon := &ContentIcon{
		config:  config,
		isImage: isImage,
	}
	icon.ExtendBaseWidget(icon)
	return icon
}

func (c *ContentIcon) SetIsImage(isImage bool) {
	c.isImage = isImage
	c.Refresh()
}

func (c *ContentIcon) SetSize(size fyne.Size) {
	c.config.Size = size
	c.Resize(size)
}

func (c *ContentIcon) SetIconColor(color color.Color) {
	c.config.IconColor = color
	c.Refresh()
}

func (c *ContentIcon) SetBackgroundColor(color color.Color) {
	c.config.BackgroundColor = color
	c.Refresh()
}

func (c *ContentIcon) GetConfig() ContentIconConfig {
	return c.config
}

func (c *ContentIcon) IsImage() bool {
	return c.isImage
}

func (c *ContentIcon) CreateRenderer() fyne.WidgetRenderer {
	return newContentIconRenderer(c)
}

func (c *ContentIcon) MinSize() fyne.Size {
	return c.config.Size
}

type contentIconRenderer struct {
	icon       *ContentIcon
	background fyne.CanvasObject
	iconObj    fyne.CanvasObject
	objects    []fyne.CanvasObject
}

func newContentIconRenderer(icon *ContentIcon) *contentIconRenderer {
	var bg fyne.CanvasObject
	if icon.config.Shape == "circle" {
		bg = canvas.NewCircle(color.Transparent)
	} else {
		rect := canvas.NewRectangle(color.Transparent)
		if icon.config.CornerRadius > 0 {
			rect.CornerRadius = icon.config.CornerRadius
		}
		bg = rect
	}

	var iconObj fyne.CanvasObject
	if icon.isImage {
		if icon.config.ImageIconResource != nil {
			img := canvas.NewImageFromResource(icon.config.ImageIconResource)
			img.FillMode = canvas.ImageFillContain
			img.SetMinSize(fyne.NewSize(float32(icon.config.Size.Width)*0.6, float32(icon.config.Size.Height)*0.6))
			iconObj = img
		} else {
			textIcon := canvas.NewText(icon.config.ImageIconText, icon.config.IconColor)
			textIcon.Alignment = fyne.TextAlignCenter
			textIcon.TextSize = float32(icon.config.Size.Height) * 0.6
			iconObj = textIcon
		}
	} else {
		if icon.config.TextIconResource != nil {
			img := canvas.NewImageFromResource(icon.config.TextIconResource)
			img.FillMode = canvas.ImageFillContain
			img.SetMinSize(fyne.NewSize(float32(icon.config.Size.Width)*0.6, float32(icon.config.Size.Height)*0.6))
			iconObj = img
		} else {
			textIcon := canvas.NewText(icon.config.TextIconText, icon.config.IconColor)
			textIcon.Alignment = fyne.TextAlignCenter
			textIcon.TextSize = float32(icon.config.Size.Height) * 0.6
			iconObj = textIcon
		}
	}

	renderer := &contentIconRenderer{
		icon:       icon,
		background: bg,
		iconObj:    iconObj,
		objects:    []fyne.CanvasObject{bg, iconObj},
	}

	renderer.Refresh()
	return renderer
}

func (r *contentIconRenderer) Objects() []fyne.CanvasObject {
	return r.objects
}

func (r *contentIconRenderer) Destroy() {}

func (r *contentIconRenderer) Layout(alloc fyne.Size) {
	target := r.icon.config.Size
	px := (alloc.Width - target.Width) / 2
	py := (alloc.Height - target.Height) / 2

	switch bg := r.background.(type) {
	case *canvas.Circle:
		bg.Resize(target)
		bg.Move(fyne.NewPos(px, py))
	case *canvas.Rectangle:
		bg.Resize(target)
		bg.Move(fyne.NewPos(px, py))
		bg.CornerRadius = r.icon.config.CornerRadius
	}

	p := r.icon.config.Padding
	w := target.Width - p*2
	h := target.Height - p*2
	if w < 1 {
		w = 1
	}
	if h < 1 {
		h = 1
	}
	iconSize := fyne.NewSize(w, h)
	iconPos := fyne.NewPos(px+p, py+p)

	if img, ok := r.iconObj.(*canvas.Image); ok {
		img.FillMode = canvas.ImageFillContain
		img.Resize(iconSize)
		img.Move(iconPos)
	} else if txt, ok := r.iconObj.(*canvas.Text); ok {
		txt.TextSize = h * 0.8
		txt.Alignment = fyne.TextAlignCenter
		txt.Resize(target)
		txt.Move(fyne.NewPos(px, py))
	} else {
		r.iconObj.Resize(iconSize)
		r.iconObj.Move(iconPos)
	}
}

func (r *contentIconRenderer) MinSize() fyne.Size {
	return r.icon.config.Size
}

func (r *contentIconRenderer) Refresh() {
	if c, ok := r.background.(*canvas.Circle); ok {
		c.FillColor = r.icon.config.BackgroundColor
		c.Refresh()
	}
	if rect, ok := r.background.(*canvas.Rectangle); ok {
		rect.FillColor = r.icon.config.BackgroundColor
		rect.CornerRadius = r.icon.config.CornerRadius
		rect.Refresh()
	}

	if img, ok := r.iconObj.(*canvas.Image); ok {
		img.Resource = nil
		if r.icon.isImage && r.icon.config.ImageIconResource != nil {
			img.Resource = r.icon.config.ImageIconResource
		}
		if !r.icon.isImage && r.icon.config.TextIconResource != nil {
			img.Resource = r.icon.config.TextIconResource
		}
		img.Refresh()
	}
	if txt, ok := r.iconObj.(*canvas.Text); ok {
		if r.icon.isImage {
			txt.Text = r.icon.config.ImageIconText
		} else {
			txt.Text = r.icon.config.TextIconText
		}
		txt.Color = r.icon.config.IconColor
		txt.Refresh()
	}
}
