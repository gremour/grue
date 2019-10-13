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

// ColorInterpolate returns color in between of two colors defined by dist from 0 to 1.
func ColorInterpolate(c1, c2 color.Color, dist float64) color.Color {
	r1, g1, b1, a1 := c1.RGBA()
	r2, g2, b2, a2 := c2.RGBA()

	return RGBA(
		(float64(r2)/0xffff-float64(r1)/0xffff)*dist+float64(r1)/0xffff,
		(float64(g2)/0xffff-float64(g1)/0xffff)*dist+float64(g1)/0xffff,
		(float64(b2)/0xffff-float64(b1)/0xffff)*dist+float64(b1)/0xffff,
		(float64(a2)/0xffff-float64(a1)/0xffff)*dist+float64(a1)/0xffff,
	)
}
