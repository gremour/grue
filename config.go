package grue

import "image/color"

// SurfaceConfig contains configuration for surface.
type SurfaceConfig struct {
	// Window specific options -- only applicable
	// to primary surface.
	Title          string
	WindowGeometry Rect
	IconFile       string
	FPS            int

	// Surface options
	PixelSize float64
	BackColor color.Color
	// PixelPrecise bool -- make all drawing converted to integer coordinates
}
