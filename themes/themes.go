package themes

import (
	"image/color"

	"github.com/gremour/grue"
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
func (tp TexturedPanel) Draw(s grue.Surface, rect grue.Rect) {
	imls, _ := s.GetImageSize(tp.ImagePrefix + "-l")
	imrs, _ := s.GetImageSize(tp.ImagePrefix + "-r")
	imts, _ := s.GetImageSize(tp.ImagePrefix + "-t")
	imbs, _ := s.GetImageSize(tp.ImagePrefix + "-b")

	// TODO: tile options
	s.DrawImageStretched(tp.ImagePrefix+"-c",
		grue.R(rect.Min.X+imls.X, rect.Min.Y+imbs.Y, rect.Max.X-imrs.X, rect.Max.Y-imts.Y),
		tp.Color)
	s.DrawImageStretched(tp.ImagePrefix+"-l",
		grue.R(rect.Min.X, rect.Min.Y+imbs.Y, rect.Min.X+imls.X, rect.Max.Y-imts.Y),
		tp.Color)
	s.DrawImageStretched(tp.ImagePrefix+"-r",
		grue.R(rect.Max.X-imrs.X, rect.Min.Y+imbs.Y, rect.Max.X, rect.Max.Y-imts.Y),
		tp.Color)
	s.DrawImageStretched(tp.ImagePrefix+"-b",
		grue.R(rect.Min.X+imls.X, rect.Min.Y, rect.Max.X-imrs.X, rect.Min.Y+imbs.Y),
		tp.Color)
	s.DrawImageStretched(tp.ImagePrefix+"-t",
		grue.R(rect.Min.X+imls.X, rect.Max.Y-imts.Y, rect.Max.X-imrs.X, rect.Max.Y),
		tp.Color)

	s.DrawImageAligned(tp.ImagePrefix+"-tl",
		grue.V(rect.Min.X, rect.Max.Y),
		grue.AlignLeft, grue.AlignTop,
		tp.Color)
	s.DrawImageAligned(tp.ImagePrefix+"-tr",
		grue.V(rect.Max.X, rect.Max.Y),
		grue.AlignRight, grue.AlignTop,
		tp.Color)
	s.DrawImageAligned(tp.ImagePrefix+"-bl",
		grue.V(rect.Min.X, rect.Min.Y),
		grue.AlignLeft, grue.AlignBottom,
		tp.Color)
	s.DrawImageAligned(tp.ImagePrefix+"-br",
		grue.V(rect.Max.X, rect.Min.Y),
		grue.AlignRight, grue.AlignBottom,
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
func (pr PlainRect) Draw(s grue.Surface, rect grue.Rect) {
	if pr.BackColor != nil {
		s.DrawFillRect(rect, pr.BackColor)
	}
	if pr.BorderColor != nil && pr.BorderSize > 0 {
		s.DrawRect(rect.Expanded(-pr.BorderInset), pr.BorderColor, pr.BorderSize)
	}
}
