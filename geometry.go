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

// Len returns the length of the vector v.
func (v Vec) Len() float64 {
	return math.Hypot(v.X, v.Y)
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

// Expanded returns expanded rectangle by given distance in pixels.
// If d is negative, rectangle is shrunk instead.
func (r Rect) Expanded(d float64) Rect {
	r.Min.X -= d
	r.Min.Y -= d
	r.Max.X += d
	r.Max.Y += d
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
