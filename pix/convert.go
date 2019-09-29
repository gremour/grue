package pix

import (
	"github.com/faiface/pixel"
	"github.com/gremour/grue"
)

// PVec converts grue.Vec to pixel.Vec.
func PVec(v grue.Vec) pixel.Vec {
	return pixel.Vec{
		X: v.X,
		Y: v.Y,
	}
}

// PRect converts grue.Rect to pixel.Rect.
func PRect(r grue.Rect) pixel.Rect {
	return pixel.Rect{
		Min: PVec(r.Min),
		Max: PVec(r.Max),
	}
}

// GVec converts pixel.Vec to grue.Vec.
func GVec(v pixel.Vec) grue.Vec {
	return grue.Vec{
		X: v.X,
		Y: v.Y,
	}
}

// GRect converts pixel.Rect to grue.Rect.
func GRect(r pixel.Rect) grue.Rect {
	return grue.Rect{
		Min: GVec(r.Min),
		Max: GVec(r.Max),
	}
}
