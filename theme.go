package grue

import "image/color"

// Theme contains info for rendering widgets
type Theme struct {
	// Font name used for titles (panels, buttons, etc)
	TitleFont         string
	TooltipFont       string
	TextColor         color.Color
	DisabledTextColor color.Color
	TooltipColor      color.Color

	// Pad to insert between border and text in autosized panels
	Pad float64

	// Drawers
	Drawers map[ThemeDrawerKey]ThemeDrawer
}

// ThemeDrawer is interface to draw rectangular panels.
type ThemeDrawer interface {
	Draw(s Surface, rect Rect)
}

// ThemeDrawerKey is used to select appropriate drawer
// based on widget and it's flags
type ThemeDrawerKey string

const (
	ThemePanel          ThemeDrawerKey = "p"
	ThemePanelDisabled  ThemeDrawerKey = "p-d"
	ThemeButton         ThemeDrawerKey = "b"
	ThemeButtonDisabled ThemeDrawerKey = "b-d"
	ThemeButtonHL       ThemeDrawerKey = "b-h"
	ThemeButtonActive   ThemeDrawerKey = "b-a"
	ThemeTooltip        ThemeDrawerKey = "tip"
)

// TexturedPanel contains info needed to draw texturized
// panel of arbitrary size.
type TexturedPanel struct {
	// Prefix of image name for the parts of panel.
	// By adding -tl, -t, -tr,
	//           -l,  -c, -r,
	//           -bl, -b, -br
	// to this, image name of the part is formed.
	ImagePrefix string

	// Use tiling instead of stretching for these parts:
	// l, r, c, t, b.
	TileHorizontal bool
	TileVertical   bool

	// Set this to some color to tint texture.
	Color color.Color
}

// Draw ...
func (tp TexturedPanel) Draw(s Surface, rect Rect) {
	imls, _ := s.GetImageSize(tp.ImagePrefix + "-l")
	imrs, _ := s.GetImageSize(tp.ImagePrefix + "-r")
	imts, _ := s.GetImageSize(tp.ImagePrefix + "-t")
	imbs, _ := s.GetImageSize(tp.ImagePrefix + "-b")

	// TODO: tile options
	s.DrawImageStretched(tp.ImagePrefix+"-c",
		R(rect.Min.X+imls.X, rect.Min.Y+imbs.Y, rect.Max.X-imrs.X, rect.Max.Y-imts.Y),
		tp.Color)
	s.DrawImageStretched(tp.ImagePrefix+"-l",
		R(rect.Min.X, rect.Min.Y+imbs.Y, rect.Min.X+imls.X, rect.Max.Y-imts.Y),
		tp.Color)
	s.DrawImageStretched(tp.ImagePrefix+"-r",
		R(rect.Max.X-imrs.X, rect.Min.Y+imbs.Y, rect.Max.X, rect.Max.Y-imts.Y),
		tp.Color)
	s.DrawImageStretched(tp.ImagePrefix+"-b",
		R(rect.Min.X+imls.X, rect.Min.Y, rect.Max.X-imrs.X, rect.Min.Y+imbs.Y),
		tp.Color)
	s.DrawImageStretched(tp.ImagePrefix+"-t",
		R(rect.Min.X+imls.X, rect.Max.Y-imts.Y, rect.Max.X-imrs.X, rect.Max.Y),
		tp.Color)

	s.DrawImageAligned(tp.ImagePrefix+"-tl",
		V(rect.Min.X, rect.Max.Y).ZR(),
		AlignTopLeft,
		tp.Color)
	s.DrawImageAligned(tp.ImagePrefix+"-tr",
		V(rect.Max.X, rect.Max.Y).ZR(),
		AlignTopRight,
		tp.Color)
	s.DrawImageAligned(tp.ImagePrefix+"-bl",
		V(rect.Min.X, rect.Min.Y).ZR(),
		AlignBottomLeft,
		tp.Color)
	s.DrawImageAligned(tp.ImagePrefix+"-br",
		V(rect.Max.X, rect.Min.Y).ZR(),
		AlignBottomRight,
		tp.Color)
}

// PlainRect draws plain rectangle with optional color.
type PlainRect struct {
	BackColor   color.Color
	BorderSize  float64
	BorderInset float64
	BorderColor color.Color
}

// Draw ...
func (pr PlainRect) Draw(s Surface, rect Rect) {
	if pr.BackColor != nil {
		s.DrawFillRect(rect, pr.BackColor)
	}
	if pr.BorderColor != nil && pr.BorderSize > 0 {
		s.DrawRect(rect.Expanded(-pr.BorderInset), pr.BorderColor, pr.BorderSize)
	}
}
