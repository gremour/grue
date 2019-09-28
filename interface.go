package grue

import (
	"image/color"
)

// Surface is an interface representing surface (graphical layer) to draw on.
type Surface interface {
	// Run main loop (blocking)
	Run()

	// Set function to be called at every screen update.
	SetEvents(handler func())

	// Get root widget to use as parent to other UI elements
	Root() Widget

	// Set current tooltip
	SetToolTip(tooltip string)

	// Draw functions
	DrawFillRect(r Rect, col color.Color)
	DrawRect(r Rect, col color.Color, thick float64)
	DrawText(r Rect, col color.Color, font, msg string, alh, alv int)

	// Mouse & keyboard
	MousePos() Vec
	PrevMousePos() Vec
	ClickMousePos() Vec
	JustPressed(button Button) bool
	JustReleased(button Button) bool
	MouseScroll() Vec

	// Load and init TTF font that will be known under given name
	InitTTF(fontName, fileName string, size float64, charset Charset) error
}

// Widget is an interface for widgets functionality.
type Widget interface {
	Node

	SubWidgets() []Widget
	GetPanel() *Panel
	Render()
	ProcessMouse()
	ProcessKeys()
	GlobalRect() Rect

	paint()
}

// Charset represents character set for use in fonts
type Charset struct {
}
