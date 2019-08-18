package pix

import (
	"fmt"
	"image/color"
	"math"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"github.com/faiface/pixel/pixelgl"
	"github.com/gremour/grue"
)

// Surface implements grue.Surface.
type Surface struct {
	Config grue.SurfaceConfig
	Canvas *pixelgl.Canvas
	Window *Window

	Rect    grue.Rect
	tooltip string
	events  func()
	root    grue.Widget

	mousePos      grue.Vec
	prevMousePos  grue.Vec
	clickMousePos grue.Vec
}

// NewPrimarySurface creates new primary surface.
func NewPrimarySurface(wcfg grue.WindowConfig, scfg grue.SurfaceConfig) (*Surface, error) {
	pixelCfg := pixelgl.WindowConfig{
		Title:  wcfg.Title,
		Bounds: PRect(wcfg.WindowGeometry),
	}

	win, err := pixelgl.NewWindow(pixelCfg)
	if err != nil {
		return nil, err
	}

	window := &Window{
		Window: win,
		fps:    wcfg.FPS,
	}
	return createSurface(window, scfg), nil
}

// NewPrimarySurfaceWin creates new primary surface with pixel window.
func NewPrimarySurfaceWin(win *pixelgl.Window, scfg grue.SurfaceConfig, fps int) (*Surface, error) {
	if win == nil {
		return nil, fmt.Errorf("Requires non-nil window")
	}

	window := &Window{
		Window: win,
		fps:    fps,
	}

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
	s.root = grue.NewPanel(nil, grue.Base{Rect: s.Rect, Phantom: true})
	s.root.GetPanel().Surface = s
	return s
}

func (s *Surface) target() pixel.Target {
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

// DrawFillRect draws filled rectangle.
func (s *Surface) DrawFillRect(r grue.Rect, col color.Color) {
	imd := imdraw.New(nil)
	imd.Color = col
	imd.Push(PVec(r.Min))
	imd.Push(PVec(r.Max))
	imd.Rectangle(0)
	imd.Draw(s.target())
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
	imd.Draw(s.target())
}

// DrawText draws text with given color, font and alignment.
func (s *Surface) DrawText(r grue.Rect, col color.Color, font, msg string, alh, alv int) {
	if len(msg) == 0 {
		return
	}
	// atl := l.Gcon.GetAtlas(font)
	// txt := text.New(pixel.ZV, atl)
	// tsz := txt.BoundsOf(msg)
	// tsz.Max.Y -= atl.LineHeight() / 2
	// txt.Color = col
	// delta := r.Min
	// switch alh {
	// case 0:
	// 	delta.X += (r.W() - tsz.W()) / 2
	// case 1:
	// 	delta.X += r.W() - tsz.W()
	// default:
	// }
	// switch alv {
	// case 0:
	// 	delta.Y += (r.H() - tsz.H()) / 2
	// case 1:
	// 	delta.Y += r.H() - tsz.H()
	// default:
	// }
	// fmt.Fprintf(txt, msg)
	// txt.Draw(l, pixel.IM.Moved(delta))
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

// MouseScroll getter.
func (s *Surface) MouseScroll() grue.Vec {
	return GVec(s.Window.MouseScroll())
}
