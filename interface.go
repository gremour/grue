package grue

import (
	"image/color"
)

// Surface is an interface representing surface (graphical layer) to draw on.
type Surface interface {
	Run()
	SetEvents(handler func())

	Root() Widget
	SetToolTip(tooltip string)

	// Draw functions
	DrawFillRect(r Rect, col color.Color)
	DrawRect(r Rect, col color.Color, thick float64)
	DrawText(r Rect, col color.Color, font, msg string, alh, alv int)

	// Mouse postion
	MousePos() Vec
	PrevMousePos() Vec
	ClickMousePos() Vec

	JustPressed(button Button) bool
	JustReleased(button Button) bool
	MouseScroll() Vec
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
