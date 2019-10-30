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
	Draw(s Surface, rect Rect, extras ...interface{})
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
func (md MultiDrawer) Draw(s Surface, rect Rect, extras ...interface{}) {
	for _, d := range md.Drawers {
		d.Draw(s, rect, extras...)
	}
}

// CursorDrawer implements structure that can draw a cursor.
type CursorDrawer interface {
	Draw(s Surface, pos Vec, height float64)
}
