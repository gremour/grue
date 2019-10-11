package grue

import (
	"image"
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

	// Show widget as popup. Popups are closed if mouse is
	// clicked outside of any popup widget, or
	// if escape is pressed. Mouse click that closed popup
	// is not otherwise processed.
	PopUp(w Widget)

	// Close popups added after this widget.
	// If passed nil, closes all popups.
	PopDownTo(w Widget)

	// Returns true, if there are any active popups.
	IsPopUpMode() bool

	// Returns true, if widget or any of its parents is popup.
	IsPopUp(w Widget) bool

	// PopUpUnder returns widget from popup list (if any),
	// if it's under pointer coords
	PopUpUnder(pos Vec) Widget

	// Draw functions
	DrawFillRect(r Rect, col color.Color)
	DrawRect(r Rect, col color.Color, thick float64)

	// Font is a font name that was previously initialized.
	// alh, alv -- horizontal and vertical text alignment (see Aling type)
	DrawText(msg, font string, r Rect, col color.Color, al Align)

	GetTextRect(msg, font string) Rect

	// Draws image centered around pos.
	DrawImage(name string, pos Vec, col color.Color)

	// Draws image aligned relative to rect.
	DrawImageAligned(name string, rect Rect, al Align, col color.Color)

	// Draw image into target rectangle.
	DrawImageStretched(name string, rect Rect, col color.Color)

	// Mouse & keyboard
	MousePos() Vec
	PrevMousePos() Vec
	ClickMousePos() Vec
	JustPressed(button Button) bool
	JustReleased(button Button) bool
	MouseScroll() Vec

	// Load and init TTF font that will be known under given name
	InitTTF(fontName, fileName string, size float64, charset Charset) error

	// Load and init images from sheet described by JSON file
	InitImages(configFileName string) error

	// Get size of the previously loaded image
	GetImageSize(name string) (Vec, error)

	// Init images from sheet described by config structure.
	// Images will be named as described in configuration.
	// DrawImage can be used with corresponding image names.
	// If JSON sheets configuration file exists,
	// you may want to use InitImages method.
	InitImageSheets(config ImageSheetConfig) error

	SetTheme(theme Theme)
	GetTheme() *Theme
}

// Widget is an interface for widgets functionality.
// This is implemented by Panel type.
// Derived types must embed *Panel type to inherit this
// interface.
type Widget interface {
	// This should be redefined in derived types
	Paint()

	// Panel implementation of these should be sufficient
	GetPanel() *Panel
	Equals(w Widget) bool
	Close()
	Foster(ch Widget)
	Render()
	ProcessMouse(wu Widget)
	ProcessKeys()
	WidgetUnder(pos Vec) Widget
	GlobalRect() Rect
	Place(rel Vec)

	// Private methods -- manipulating tree structure
	addChild(ch Widget)
	removeChild(ch Widget)
	removeChildren()
}

// Charset represents character set for use in fonts
type Charset struct {
	// TODO: implement charset.
	// For now, ASCII is hardcoded.
}

// ImageSheetConfig contains configuration for sheets containing subimages (sprites).
// OpenGL allows limited number of textures to be loaded into videocard memory.
// Because of that, images are loaded as atlases -- a big texture containing
// a lot of subimages. Sheet is a continous line of images of same width and height
// within atlas. Atlas can hold multiple sheets.
type ImageSheetConfig struct {
	// Either atlas or file should be set.
	// If atlas is nil, file is used as filename to load image.
	Atlas image.Image
	// File is the name of the image file to load.
	File string `json:"file"`
	// Sheets containing continuous line of uniform sized subimages.
	Sheets []SheetConfig `json:"sheets"`
}

// SheetConfig is configuration for oine sheet
type SheetConfig struct {
	XOffset float64  `json:"x_offset"`
	YOffset float64  `json:"y_offset"`
	W       float64  `json:"width"`
	H       float64  `json:"height"`
	Names   []string `json:"names"`
}
