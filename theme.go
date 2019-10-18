package grue

import (
	"image/color"
)

// Theme contains info for rendering widgets
type Theme struct {
	// Font name used for titles (panels, buttons, etc)
	TitleFont   string
	TooltipFont string

	// Common text color
	TextColor color.Color

	// Text color overrides
	ButtonTextColor  color.Color
	PanelTextColor   color.Color
	EditTextColor    color.Color
	PlaceholderColor color.Color

	DisabledTextColor color.Color
	TooltipColor      color.Color

	// Pad to insert between border and text in autosized panels
	Pad float64

	// Vector to dispace test for pressed buttons
	PressDisplace Vec

	// Drawers
	Drawers      map[ThemeDrawerKey]ThemeDrawer
	CursorDrawer CursorDrawer
}

// ThemeDrawer is interface to draw rectangular panels.
type ThemeDrawer interface {
	Draw(s Surface, rect Rect)
}

// ThemeDrawerKey is used to select appropriate drawer
// based on widget and it's flags
type ThemeDrawerKey string

const (
	ThemePanel            ThemeDrawerKey = "p"
	ThemePanelDisabled    ThemeDrawerKey = "p-d"
	ThemeButton           ThemeDrawerKey = "b"
	ThemeButtonDisabled   ThemeDrawerKey = "b-d"
	ThemeButtonHL         ThemeDrawerKey = "b-h"
	ThemeButtonActive     ThemeDrawerKey = "b-a"
	ThemeLineEdit         ThemeDrawerKey = "le"
	ThemeLineEditDisabled ThemeDrawerKey = "le-d"
	ThemeLineEditHL       ThemeDrawerKey = "le-h"
	ThemeLineEditActive   ThemeDrawerKey = "le-a"
	ThemeTooltip          ThemeDrawerKey = "tip"
)

// MultiDrawer is a drawer combining several other drawers.
type MultiDrawer struct {
	Drawers []ThemeDrawer
}

// Draw ...
func (md MultiDrawer) Draw(s Surface, rect Rect) {
	for _, d := range md.Drawers {
		d.Draw(s, rect)
	}
}

// TexturedPanel contains info needed to draw texturized
// panel of arbitrary size. Panel consists of parts
// (defined my margins) that are stretched or tiled,
// depending on options:
//
//      Left   Right
//     |      |
//  [1][  2  ][3]__ Top
//  [4][  5  ][6]__ Bottom
//  [7][  8  ][9]
//
//  Parts 1, 3, 7, 9 are not stretched.
//  Other parth are either stretched or tiled to fill
//  all the remaining space in target rectangle,
//  depending on tiling options.
type TexturedPanel struct {
	Image string

	// Margins
	Left   float64
	Bottom float64
	Right  float64
	Top    float64

	// Use tiling instead of stretching (except for corner parts).
	TileHorizontal bool
	TileVertical   bool

	// Set this to some color to tint texture.
	Color color.Color
}

// Draw ...
func (tp TexturedPanel) Draw(s Surface, rect Rect) {
	// Points of interest:
	//  x1 x2        x3 x4
	//   _______________   y4
	//  |   _________   |  y3
	//  |  |         |  |
	//  |  |_________|  |  y2
	//  |_______________|  y1

	imsz := s.GetImageRect(tp.Image).Size()

	x1i := 0.0
	x2i := tp.Left
	x3i := imsz.X - tp.Right
	x4i := imsz.X
	y1i := 0.0
	y2i := tp.Bottom
	y3i := imsz.Y - tp.Top
	y4i := imsz.Y

	x1r := rect.Min.X
	x2r := rect.Min.X + tp.Left
	x3r := rect.Max.X - tp.Right
	x4r := rect.Max.X
	y1r := rect.Min.Y
	y2r := rect.Min.Y + tp.Bottom
	y3r := rect.Max.Y - tp.Top
	y4r := rect.Max.Y

	// Bottomleft
	s.DrawImagePart(tp.Image,
		R(x1i, y1i, x2i, y2i),
		R(x1r, y1r, x2r, y2r), tp.Color)

	// Topleft
	s.DrawImagePart(tp.Image,
		R(x1i, y3i, x2i, y4i),
		R(x1r, y3r, x2r, y4r), tp.Color)

	// Bottomright
	s.DrawImagePart(tp.Image,
		R(x3i, y1i, x4i, y2i),
		R(x3r, y1r, x4r, y2r), tp.Color)

	// TopRight
	s.DrawImagePart(tp.Image,
		R(x3i, y3i, x4i, y4i),
		R(x3r, y3r, x4r, y4r), tp.Color)

	if tp.TileHorizontal {
		wi := x3i - x2i
		wr := x3r - x2r
		nw := wr / wi
		for i, x, w := 0, 0.0, wi; i <= int(nw); i++ {
			if x+w > wr {
				w = wr - x
			}
			// Top tiled
			s.DrawImagePart(tp.Image,
				R(x2i, y3i, x2i+w, y4i),
				R(x2r+x, y3r, x2r+x+w, y4r), tp.Color)

			// Bottom tiled
			s.DrawImagePart(tp.Image,
				R(x2i, y1i, x2i+w, y2i),
				R(x2r+x, y1r, x2r+x+w, y2r), tp.Color)

			if !tp.TileVertical {
				// Center tiled horizontally, stretched vertically
				s.DrawImagePart(tp.Image,
					R(x2i, y2i, x2i+w, y3i),
					R(x2r+x, y2r, x2r+x+w, y3r), tp.Color)
			}
			x += w
		}
	} else {
		// Top stretched
		s.DrawImagePart(tp.Image,
			R(x2i, y3i, x3i, y4i),
			R(x2r, y3r, x3r, y4r), tp.Color)

		// Bottom stretched
		s.DrawImagePart(tp.Image,
			R(x2i, y1i, x3i, y2i),
			R(x2r, y1r, x3r, y2r), tp.Color)
	}

	if tp.TileVertical {
		hi := y3i - y2i
		hr := y3r - y2r
		nh := hr / hi
		for j, y, h := 0, 0.0, hi; j <= int(nh); j++ {
			if y+h > hr {
				h = hr - y
			}
			// Left tiled
			s.DrawImagePart(tp.Image,
				R(x1i, y2i, x2i, y2i+h),
				R(x1r, y2r+y, x2r, y2r+y+h), tp.Color)

			// Right tiled
			s.DrawImagePart(tp.Image,
				R(x3i, y2i, x4i, y2i+h),
				R(x3r, y2r+y, x4r, y2r+y+h), tp.Color)

			if !tp.TileHorizontal {
				// Center tiled vertically, stretched horizontally
				s.DrawImagePart(tp.Image,
					R(x2i, y2i, x3i, y2i+h),
					R(x2r, y2r+y, x3r, y2r+y+h), tp.Color)
			}
			y += h
		}
	} else {
		// Left stretched
		s.DrawImagePart(tp.Image,
			R(x1i, y2i, x2i, y3i),
			R(x1r, y2r, x2r, y3r), tp.Color)

		// Right stretched
		s.DrawImagePart(tp.Image,
			R(x3i, y2i, x4i, y3i),
			R(x3r, y2r, x4r, y3r), tp.Color)
	}

	if !tp.TileHorizontal && !tp.TileVertical {
		// Center stretched
		s.DrawImagePart(tp.Image,
			R(x2i, y2i, x3i, y3i),
			R(x2r, y2r, x3r, y3r), tp.Color)
	} else if tp.TileHorizontal && tp.TileVertical {
		// Center tiled
		hi := y3i - y2i
		hr := y3r - y2r
		nh := hr / hi
		for j, y, h := 0, 0.0, hi; j <= int(nh); j++ {
			if y+h > hr {
				h = hr - y
			}
			wi := x3i - x2i
			wr := x3r - x2r
			nw := wr / wi
			for i, x, w := 0, 0.0, wi; i <= int(nw); i++ {
				if x+w > wr {
					w = wr - x
				}
				s.DrawImagePart(tp.Image,
					R(x2i, y2i, x2i+w, y2i+h),
					R(x2r+x, y2r+y, x2r+x+w, y2r+y+h), tp.Color)
				x += w
			}
			y += h
		}
	}
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

// RectCursorDrawer ...
type RectCursorDrawer struct {
	Color1        color.Color
	Color2        color.Color
	Width         float64
	PulseInterval float64
}

// CursorDrawer implements structure that can draw a cursor.
type CursorDrawer interface {
	Draw(s Surface, pos Vec, height float64)
}

// Draw ...
func (cd RectCursorDrawer) Draw(s Surface, pos Vec, height float64) {
	if cd.PulseInterval == 0 {
		cd.PulseInterval = 1
	}
	if cd.Width == 0 {
		cd.Width = 2
	}
	col := ColorInterpolate(cd.Color1, cd.Color2, s.Pulse(cd.PulseInterval))
	s.DrawFillRect(R(pos.X, pos.Y, pos.X+cd.Width, pos.Y+height), col)
}
