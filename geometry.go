package grue

import "math"

// Types here are abstracted from specific graphics library.
// Although pixel has same types, we need not to create dependency
// on pixel.

// Vec decribes point or vector on surface.
type Vec struct {
	X, Y float64
}

// Add returns sum of vectors
func (v Vec) Add(u Vec) Vec {
	return Vec{v.X + u.X, v.Y + u.Y}
}

// Sub return result of subtraction of u from v
func (v Vec) Sub(u Vec) Vec {
	return Vec{v.X - u.X, v.Y - u.Y}
}

// Len returns the length of the vector v.
func (v Vec) Len() float64 {
	return math.Hypot(v.X, v.Y)
}

// Half of the rectangle
func (v Vec) Half() Vec {
	return Vec{
		X: v.X / 2,
		Y: v.X / 2,
	}
}

// ZR return zero size rectangle centered at the point.
func (v Vec) ZR() Rect {
	return Rect{
		Min: v,
		Max: v,
	}
}

// Rect describes rectangular area on surface.
type Rect struct {
	Min, Max Vec
}

// Size of the rectangle
func (r Rect) Size() Vec {
	return Vec{
		X: r.Max.X - r.Min.X,
		Y: r.Max.Y - r.Min.Y,
	}
}

// Center of the rectangle
func (r Rect) Center() Vec {
	return Vec{
		X: (r.Max.X + r.Min.X) / 2,
		Y: (r.Max.Y + r.Min.Y) / 2,
	}
}

// W calculates the width of the rectangle
func (r Rect) W() float64 {
	return r.Max.X - r.Min.X
}

// H calculates the height of the rectangle
func (r Rect) H() float64 {
	return r.Max.Y - r.Min.Y
}

// Expanded returns expanded rectangle by given distance in pixels.
// If d is negative, rectangle is shrunk instead.
func (r Rect) Expanded(d float64) Rect {
	r.Min.X -= d
	r.Min.Y -= d
	r.Max.X += d
	r.Max.Y += d
	return r
}

// Extended returns expanded rectangle by given distances in pixels.
// If distance is negative, rectangle is shrunk instead.
func (r Rect) Extended(left, bottom, right, top float64) Rect {
	r.Min.X -= left
	r.Min.Y -= bottom
	r.Max.X += right
	r.Max.Y += top
	return r
}

// Moved returns the Rect moved (both Min and Max) by the given vector delta.
func (r Rect) Moved(delta Vec) Rect {
	return Rect{
		Min: r.Min.Add(delta),
		Max: r.Max.Add(delta),
	}
}

// Contains checks whether a vector u is contained within this Rect (including it's borders).
func (r Rect) Contains(u Vec) bool {
	return r.Min.X <= u.X && u.X <= r.Max.X && r.Min.Y <= u.Y && u.Y <= r.Max.Y
}

// V returns initialized Vec.
func V(x, y float64) Vec {
	return Vec{x, y}
}

// R returns initialized Rect.
func R(x1, y1, x2, y2 float64) Rect {
	return Rect{
		Min: Vec{x1, y1},
		Max: Vec{x2, y2},
	}
}

// R0 creates rect with Min in (0, 0).
func R0(w, h float64) Rect {
	return Rect{Max: V(w, h)}
}

// Align defines object alignment relative to parent.
type Align int

const (
	// AlignDefault is default align (center in most cases)
	AlignDefault Align = 0
	// AlignLeft ...
	AlignLeft Align = 1
	// AlignRight ...
	AlignRight Align = 2
	// AlignTop ...
	AlignTop Align = 3
	// AlignBottom ...
	AlignBottom Align = 4

	// AlignTopLeft ...
	AlignTopLeft Align = 5
	// AlignTopRight ...
	AlignTopRight Align = 6
	// AlignBottomLeft ...
	AlignBottomLeft Align = 7
	// AlignBottomRight ...
	AlignBottomRight Align = 8

	// AlignCenter is explicit center alignment
	AlignCenter Align = 10
)

// AlignToRect returns a Vec that src Rect have to be
// moved by in order to align relative to dst Rect given
// alignments alh, alv.
func (r Rect) AlignToRect(dst Rect, al Align) Vec {
	delta := dst.Center()
	switch al {
	case AlignDefault:
		fallthrough
	case AlignCenter:
	case AlignTopRight:
		fallthrough
	case AlignBottomRight:
		fallthrough
	case AlignRight:
		delta.X += (dst.W() - r.W()) / 2
	case AlignTopLeft:
		fallthrough
	case AlignBottomLeft:
		fallthrough
	case AlignLeft:
		delta.X -= (dst.W() - r.W()) / 2
	default:
	}
	switch al {
	case AlignDefault:
		fallthrough
	case AlignCenter:
	case AlignTopLeft:
		fallthrough
	case AlignTopRight:
		fallthrough
	case AlignTop:
		delta.Y += (dst.H() - r.H()) / 2
	case AlignBottomLeft:
		fallthrough
	case AlignBottomRight:
		fallthrough
	case AlignBottom:
		delta.Y -= (dst.H() - r.H()) / 2
	default:
	}
	return delta
}

// AlignToPoint returns a Vec that src Rect have to be
// moved by in order to align relative to dst Point given
// alignments alh, alv.
func (r Rect) AlignToPoint(dst Vec, al Align) Vec {
	return r.AlignToRect(Rect{dst, dst}, al)
}
