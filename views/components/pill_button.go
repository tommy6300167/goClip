package components

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

type PillButtonStyle struct {
	Size         fyne.Size
	Bg           color.Color
	BgHover      color.Color
	Fg           color.Color
	CornerRadius float32
	Padding      float32
}

var DefaultPillStyle = PillButtonStyle{
	Size:         fyne.NewSize(80, 36),
	Bg:           color.NRGBA{255, 255, 255, 28},  // ~rgba(255,255,255,0.11)
	BgHover:      color.NRGBA{255, 255, 255, 46},  // ~rgba(255,255,255,0.18)
	Fg:           color.NRGBA{230, 233, 236, 255}, // 淺灰字
	CornerRadius: 18,
	Padding:      12,
}

type PillButton struct {
	widget.BaseWidget
	Icon     fyne.Resource
	Text     string
	OnTapped func()

	style PillButtonStyle
	hover bool
	press bool

	bg   *canvas.Rectangle
	img  *canvas.Image
	lbl  *canvas.Text
	box  *fyne.Container
	root *fyne.Container
}

func NewPillButton(text string, icon fyne.Resource, onTap func()) *PillButton {
	return NewPillButtonWithStyle(text, icon, onTap, DefaultPillStyle)
}

func NewPillButtonWithStyle(text string, icon fyne.Resource, onTap func(), st PillButtonStyle) *PillButton {
	p := &PillButton{
		Text:     text,
		Icon:     icon,
		OnTapped: onTap,
		style:    st,
	}
	p.ExtendBaseWidget(p)

	p.bg = canvas.NewRectangle(st.Bg)
	p.bg.CornerRadius = st.CornerRadius

	if icon != nil {
		p.img = canvas.NewImageFromResource(icon)
		p.img.FillMode = canvas.ImageFillContain
		p.img.SetMinSize(fyne.NewSize(16, 16))
	}
	p.lbl = canvas.NewText(text, st.Fg)
	p.lbl.TextSize = 14

	inner := []fyne.CanvasObject{}
	if p.img != nil {
		inner = append(inner, p.img)
	}
	inner = append(inner, p.lbl)

	p.box = container.NewHBox(inner...)
	content := container.NewCenter(p.box)

	p.root = container.NewStack(p.bg, content)
	p.root.Resize(st.Size)
	return p
}

func (p *PillButton) CreateRenderer() fyne.WidgetRenderer {
	return widget.NewSimpleRenderer(p.root)
}
func (p *PillButton) MinSize() fyne.Size { return p.style.Size }

func (p *PillButton) Tapped(*fyne.PointEvent) {
	if p.OnTapped != nil {
		p.OnTapped()
	}
}
func (p *PillButton) MouseIn(*fyne.PointEvent)   { p.hover = true; p.refreshBg() }
func (p *PillButton) MouseOut()                  { p.hover = false; p.refreshBg() }
func (p *PillButton) MouseDown(*fyne.PointEvent) { p.press = true; p.refreshBg() }
func (p *PillButton) MouseUp(*fyne.PointEvent)   { p.press = false; p.refreshBg() }

func (p *PillButton) refreshBg() {
	col := p.style.Bg
	if p.hover {
		col = p.style.BgHover
	}
	p.bg.FillColor = col
	p.bg.Refresh()
	p.lbl.Color = p.style.Fg
	p.lbl.Refresh()
}
