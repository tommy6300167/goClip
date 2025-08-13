package views

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
)

type CustomTheme struct {
	fontResource fyne.Resource
}

func NewCustomTheme(fontResource fyne.Resource) *CustomTheme {
	return &CustomTheme{
		fontResource: fontResource,
	}
}

func (t *CustomTheme) Color(n fyne.ThemeColorName, v fyne.ThemeVariant) color.Color {
	return theme.DefaultTheme().Color(n, v)
}

func (t *CustomTheme) Icon(n fyne.ThemeIconName) fyne.Resource {
	return theme.DefaultTheme().Icon(n)
}

func (t *CustomTheme) Size(n fyne.ThemeSizeName) float32 {
	return theme.DefaultTheme().Size(n)
}

func (t *CustomTheme) Font(s fyne.TextStyle) fyne.Resource {
	if t.fontResource != nil {
		return t.fontResource
	}
	return theme.DefaultTheme().Font(s)
}