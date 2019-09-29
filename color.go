package grue

import "image/color"

// RGB constructs color with float R, G, B components in range 0..1.
func RGB(r, g, b float64) color.Color {
	return color.RGBA{
		uint8(0xff * r),
		uint8(0xff * g),
		uint8(0xff * b),
		0xff,
	}
}

// RGBA constructs color with float R, G, B, A components in range 0..1.
func RGBA(r, g, b, a float64) color.Color {
	return color.RGBA{
		uint8(0xff * r),
		uint8(0xff * g),
		uint8(0xff * b),
		uint8(0xff * a),
	}
}
