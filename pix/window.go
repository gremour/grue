// Package pix is implementation of grue based on faiface/pixel
package pix

import (
	"time"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"github.com/gremour/grue"
)

// Window represents underlying graphical window.
type Window struct {
	*pixelgl.Window
	surfaces []*Surface

	frameTime float64
	totalTime float64
	fps       int
}

// Run the main loop.
func (w *Window) Run() {
	var fps <-chan time.Time
	if w.fps != 0 {
		fps = time.Tick(time.Second / time.Duration(w.fps))
	}
	lastTime := time.Now()
	for _, s := range w.surfaces {
		s.updateMousePos(GVec(w.MousePosition()), false)
	}
	for !w.Closed() {
		click := w.JustPressed(pixelgl.MouseButtonLeft) ||
			w.JustPressed(pixelgl.MouseButtonRight) ||
			w.JustPressed(pixelgl.MouseButtonMiddle)

		for _, s := range w.surfaces {
			s.updateMousePos(GVec(w.MousePosition()), click)
			if s.root != nil {
				s.root.ProcessMouse()
				s.root.ProcessKeys()
				s.root.Render()
			}
		}

		// w.ShowTooltip()

		for _, s := range w.surfaces {
			if s.Canvas != nil {
				if s.Config.PixelSize == 1 {
					s.Canvas.Draw(w, pixel.IM.Moved(w.Bounds().Center()))
				} else {
					s.Canvas.Draw(w, pixel.IM.Scaled(pixel.ZV, s.Config.PixelSize).Moved(w.Bounds().Center()))
				}
				s.Canvas.Clear(grue.RGBA(0.15, 0.15, 0.15, 0))
			}
		}
		w.Update()
		for _, s := range w.surfaces {
			if s.events != nil {
				s.events()
			}
		}

		if fps != nil {
			<-fps
		}
		w.frameTime = time.Since(lastTime).Seconds()
		w.totalTime += w.frameTime
		lastTime = time.Now()
	}
}

// RunUI is used to run code on main thread.
// Put any code that creates grue surfaces in the closure and pass
// it to this function.
func RunUI(f func()) {
	pixelgl.Run(f)
}
