package grue

// WindowConfig contains configuration for window.
type WindowConfig struct {
	Title          string
	WindowGeometry Rect
	IconFile       string
	FPS            int
}

// SurfaceConfig contains configuration for surface.
type SurfaceConfig struct {
	PixelSize float64
	// PixelPrecise bool -- make all drawing converted to integer coordinates
}
