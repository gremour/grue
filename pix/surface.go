package pix

import (
	"fmt"
	"image"
	"image/color"
	"math"
	"os"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"github.com/faiface/pixel/pixelgl"
	"github.com/faiface/pixel/text"
	"github.com/gremour/grue"
)

// Surface implements grue.Surface.
type Surface struct {
	Config grue.SurfaceConfig
	Canvas *pixelgl.Canvas
	Window *Window
	Popups []grue.Widget

	Rect    grue.Rect
	tooltip string
	events  func()
	root    grue.Widget

	mousePos      grue.Vec
	prevMousePos  grue.Vec
	clickMousePos grue.Vec
}

// NewPrimarySurface creates new primary surface.
func NewPrimarySurface(scfg grue.SurfaceConfig) (*Surface, error) {
	pixelCfg := pixelgl.WindowConfig{
		Title:  scfg.Title,
		Bounds: PRect(scfg.WindowGeometry),
	}

	win, err := pixelgl.NewWindow(pixelCfg)
	if err != nil {
		return nil, err
	}

	window := newWindow(win, scfg.FPS)
	return createSurface(window, scfg), nil
}

// NewPrimarySurfaceWin creates new primary surface with pixel window.
func NewPrimarySurfaceWin(win *pixelgl.Window, scfg grue.SurfaceConfig, fps int) (*Surface, error) {
	if win == nil {
		return nil, fmt.Errorf("Requires non-nil window")
	}

	window := newWindow(win, fps)
	return createSurface(window, scfg), nil
}

// NewSurface creates new surface on top of given surface.
func NewSurface(surf *Surface, scfg grue.SurfaceConfig) (*Surface, error) {
	if surf == nil {
		return nil, fmt.Errorf("Requires non-nil base surface")
	}
	return createSurface(surf.Window, scfg), nil
}

func createSurface(window *Window, scfg grue.SurfaceConfig) *Surface {
	psz := float64(1)
	if scfg.PixelSize != 0 {
		psz = scfg.PixelSize
	}
	s := &Surface{
		Window: window,
		Rect: grue.R0(math.Floor(window.Bounds().W()/psz),
			math.Floor(window.Bounds().H()/psz)),
		Config: scfg,
	}
	if scfg.PixelSize != 0 {
		s.Canvas = pixelgl.NewCanvas(PRect(s.Rect))
	}
	s.Window.surfaces = append(s.Window.surfaces, s)
	s.root = grue.NewPanel(nil, grue.Base{Rect: s.Rect, Theme: &grue.Theme{}})
	s.root.GetPanel().Surface = s
	return s
}

// Target return pixel target to draw on.
func (s *Surface) Target() pixel.Target {
	if s.Canvas == nil {
		return s.Window
	}
	return s.Canvas
}

// Run the main loop for the window of surface.
func (s *Surface) Run() {
	s.Window.Run()
}

// SetEvents handler to execute for each window update.
func (s *Surface) SetEvents(handler func()) {
	s.events = handler
}

// SetToolTip ...
func (s *Surface) SetToolTip(tooltip string) {
	s.tooltip = tooltip
}

// Root returns root widget for the surface.
func (s *Surface) Root() grue.Widget {
	return s.root
}

// SetFocus ...
func (s *Surface) SetFocus(w grue.Widget) {
	s.Window.focus = w
}

// Focus ...
func (s *Surface) Focus() grue.Widget {
	return s.Window.focus
}

// PopUp ...
func (s *Surface) PopUp(w grue.Widget) {
	if w == nil {
		panic("trying to add empty popup")
	}
	for _, p := range s.Popups {
		if w == p {
			return
		}
	}
	w.GetPanel().Surface = s
	s.Popups = append(s.Popups, w.GetPanel())
}

// PopDownTo ...
func (s *Surface) PopDownTo(w grue.Widget) {
	found := false
	ind := 0
	if w == nil {
		found = true
	}
	for i, p := range s.Popups {
		if found {
			s.Popups[i] = nil
			p.Close()
		} else if w != nil && w.GetPanel() == p {
			found = true
			ind = i + 1
		}
	}
	if ind == 0 {
		s.Popups = nil
	} else {
		s.Popups = s.Popups[:ind]
	}
}

// PopUpUnder ...
func (s *Surface) PopUpUnder(pos grue.Vec) grue.Widget {
	for i := range s.Popups {
		p := s.Popups[len(s.Popups)-i-1]
		wu := p.WidgetUnder(pos)
		if wu != nil {
			return wu
		}
	}
	return nil
}

// IsPopUpMode ...
func (s *Surface) IsPopUpMode() bool {
	return len(s.Popups) > 0
}

// IsPopUp ...
func (s *Surface) IsPopUp(w grue.Widget) bool {
	if w == nil {
		return false
	}
	cnt := 100
	for w != nil && cnt > 0 {
		for _, p := range s.Popups {
			if p.Equals(w) {
				return true
			}
		}
		w = w.GetPanel().Parent
		cnt--
	}
	return false
}

// DrawFillRect draws filled rectangle.
func (s *Surface) DrawFillRect(r grue.Rect, col color.Color) {
	imd := imdraw.New(nil)
	imd.Color = col
	imd.Push(PVec(r.Min))
	imd.Push(PVec(r.Max))
	imd.Rectangle(0)
	imd.Draw(s.Target())
}

// DrawRect draws rectlangle with given line thickness.
func (s *Surface) DrawRect(r grue.Rect, col color.Color, thick float64) {
	imd := imdraw.New(nil)
	imd.Color = col
	imd.Push(PVec(r.Min))
	imd.Push(pixel.V(r.Max.X, r.Min.Y+thick))
	imd.Rectangle(0)
	imd.Push(pixel.V(r.Min.X, r.Max.Y-thick))
	imd.Push(PVec(r.Max))
	imd.Rectangle(0)
	imd.Push(pixel.V(r.Min.X, r.Min.Y+thick))
	imd.Push(pixel.V(r.Min.X+thick, r.Max.Y-thick))
	imd.Rectangle(0)
	imd.Push(pixel.V(r.Max.X-thick, r.Min.Y+thick))
	imd.Push(pixel.V(r.Max.X, r.Max.Y-thick))
	imd.Rectangle(0)
	imd.Draw(s.Target())
}

// DrawText draws text with given color, font and alignment.
func (s *Surface) DrawText(msg, font string, r grue.Rect, col color.Color, al grue.Align) {
	if len(msg) == 0 {
		return
	}
	atl, ok := s.Window.fonts[font]
	if !ok {
		atl = text.Atlas7x13
	}
	txt := text.New(pixel.ZV, atl)
	tsz := txt.BoundsOf(msg)
	tsz.Max.Y -= atl.LineHeight() / 2
	txt.Color = col
	fmt.Fprintf(txt, msg)
	pos := GRect(tsz).AlignToRect(r, al)
	pos = pos.Sub(grue.V(tsz.W()/2, tsz.H()/2))
	txt.Draw(s.Target(), pixel.IM.Moved(PVec(pos)))
}

// GetTextRect ...
func (s *Surface) GetTextRect(msg, font string) grue.Rect {
	if len(msg) == 0 {
		return grue.Rect{}
	}
	atl, ok := s.Window.fonts[font]
	if !ok {
		atl = text.Atlas7x13
	}
	txt := text.New(pixel.ZV, atl)
	tsz := txt.BoundsOf(msg)
	//tsz.Max.Y -= atl.LineHeight() / 2
	return GRect(tsz)
}

// FitText ...
func (s *Surface) FitText(msg, font string, width float64) string {
	if s.GetTextRect(msg, font).W() <= width {
		return msg
	}
	avgSym := s.GetTextRect("a", font).W()
	avgLen := int(width / avgSym)
	for s.GetTextRect(msg[:avgLen], font).W() <= width {
		avgLen++
	}
	for s.GetTextRect(msg[:avgLen], font).W() > width {
		avgLen--
	}
	return msg[:avgLen]
}

// DrawImage ...
func (s *Surface) DrawImage(name string, pos grue.Vec, col color.Color) {
	im, err := s.GetImage(name)
	if err != nil {
		return
	}
	im.DrawColorMask(s.Target(), pixel.IM.Moved(PVec(pos)), col)
}

// DrawImageStretched ...
func (s *Surface) DrawImageStretched(name string, rect grue.Rect, col color.Color) {
	im, err := s.GetImage(name)
	if err != nil {
		return
	}
	imsz, _ := s.GetImageSize(name)
	if imsz.X == 0 ||
		imsz.Y == 0 {
		return
	}
	rctr := PVec(rect.Center())
	scv := pixel.V(rect.W()/imsz.X, rect.H()/imsz.Y)
	im.DrawColorMask(s.Target(), pixel.IM.Moved(rctr).ScaledXY(rctr, scv), col)
}

// DrawImageAligned ...
func (s *Surface) DrawImageAligned(name string, rect grue.Rect, al grue.Align, col color.Color) {
	im, err := s.GetImage(name)
	if err != nil {
		return
	}
	imsz, err := s.GetImageSize(name)
	if err != nil {
		return
	}
	pos := grue.Rect{Max: imsz}.AlignToRect(rect, al)
	im.DrawColorMask(s.Target(), pixel.IM.Moved(PVec(pos)), col)
}

// DrawTooltip ...
func (s *Surface) DrawTooltip() {
	if s.tooltip == "" {
		return
	}
	theme := s.GetTheme()
	drw, _ := theme.Drawers[grue.ThemeTooltip]
	if drw == nil {
		return
	}
	r := s.GetTextRect(s.tooltip, theme.TooltipFont)
	r = r.Moved(s.MousePos()).Expanded(theme.Pad).Moved(grue.V(theme.Pad, theme.Pad))
	drw.Draw(s, r)
	s.DrawText(s.tooltip, theme.TooltipFont, r, theme.TooltipColor, grue.AlignCenter)
}

func (s *Surface) updateMousePos(pos grue.Vec, click bool) {
	psz := float64(1)
	if s.Config.PixelSize != 0 {
		psz = s.Config.PixelSize
	}
	if psz != 1 {
		pos = grue.V(math.Floor(pos.X/psz), math.Floor(pos.Y/psz))
	}
	s.prevMousePos = s.mousePos
	s.mousePos = pos
	if click {
		s.clickMousePos = pos
	}
}

// MousePos getter.
func (s *Surface) MousePos() grue.Vec {
	return s.mousePos
}

// PrevMousePos getter.
func (s *Surface) PrevMousePos() grue.Vec {
	return s.prevMousePos
}

// ClickMousePos getter.
func (s *Surface) ClickMousePos() grue.Vec {
	return s.clickMousePos
}

// JustPressed getter.
func (s *Surface) JustPressed(button grue.Button) bool {
	return s.Window.JustPressed(pixelgl.Button(button))
}

// JustReleased getter.
func (s *Surface) JustReleased(button grue.Button) bool {
	return s.Window.JustReleased(pixelgl.Button(button))
}

// KeysInput ...
func (s *Surface) KeysInput() string {
	return s.Window.Typed()
}

// Repeated ...
func (s *Surface) Repeated(button grue.Button) bool {
	return s.Window.Repeated(pixelgl.Button(button))
}

// MouseScroll getter.
func (s *Surface) MouseScroll() grue.Vec {
	return GVec(s.Window.MouseScroll())
}

// InitTTF ...
func (s *Surface) InitTTF(fontName, fileName string, size float64, charset grue.Charset) error {
	face, err := grue.LoadTTF(fileName, size)
	if err != nil {
		return err
	}
	atlas := text.NewAtlas(face, text.ASCII)
	s.Window.fonts[fontName] = atlas
	return nil
}

// InitImageSheets ...
func (s *Surface) InitImageSheets(config grue.ImageSheetConfig) error {
	if config.Atlas == nil {
		imageFile, err := os.Open(config.File)
		if err != nil {
			return err
		}
		defer imageFile.Close()

		config.Atlas, _, err = image.Decode(imageFile)
		if err != nil {
			return err
		}
	}
	pic := pixel.PictureDataFromImage(config.Atlas)
	b := pic.Bounds()
	for _, sh := range config.Sheets {
		pos := b.Min
		pos.X += sh.XOffset
		pos.Y += sh.YOffset
		size := pixel.V(sh.W, sh.H)
		if pos.X+size.X > b.Max.X ||
			pos.Y+size.Y > b.Max.Y {
			return fmt.Errorf("offest exceeds image size: offset=%v,%v, image size=%v,%v",
				sh.XOffset, sh.YOffset, b.Max.X, b.Max.Y)
		}
		for _, n := range sh.Names {
			if len(n) > 0 {
				r := pixel.R(
					pos.X,
					b.Max.Y-pos.Y-size.Y,
					pos.X+size.X,
					b.Max.Y-pos.Y)
				s.Window.sprites[n] = pixel.NewSprite(pic, r)
				pos.X += size.X
			}
			if len(n) == 0 || pos.X+size.X > b.Max.X {
				pos.X = b.Min.X
				pos.Y += size.Y
			}
			if pos.Y+size.Y > b.Max.Y {
				break
			}
		}
	}
	return nil
}

// InitImages ...
func (s *Surface) InitImages(configFileName string) error {
	sheets, err := grue.LoadImages(configFileName)
	if err != nil {
		return err
	}
	err = s.InitImageSheets(sheets)
	if err != nil {
		return err
	}
	return err
}

// GetImageSize ...
func (s *Surface) GetImageSize(name string) (grue.Vec, error) {
	im, err := s.GetImage(name)
	if err != nil {
		return grue.Vec{}, err
	}
	return GRect(im.Frame()).Size(), nil
}

// GetImage ...
func (s *Surface) GetImage(name string) (*pixel.Sprite, error) {
	spr, ok := s.Window.sprites[name]
	if !ok {
		return nil, fmt.Errorf(`image "%v" not found`, name)
	}
	return spr, nil
}

// SetTheme ...
func (s *Surface) SetTheme(theme grue.Theme) {
	s.Window.theme = theme
}

// GetTheme ...
func (s *Surface) GetTheme() *grue.Theme {
	return &s.Window.theme
}

// Pulse ...
func (s *Surface) Pulse(dur float64) float64 {
	return math.Abs(math.Sin(s.Window.totalTime * math.Pi / dur))
}
