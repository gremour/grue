package main

import (
	"image/color"

	"github.com/gremour/grue"
	"github.com/gremour/grue/pix"
	"github.com/gremour/grue/themes"
)

func main() {
	// This is to execute code on main thread (OpenGL requirement).
	pix.RunUI(runUI)
}

// All grue & grue/pix code should be set up in this function.
// Don't spawn goroutines which access any of grue objects:
// surfaces, window, widgets.
// Having separate goroutines for logic is OK, but
// be sure to communicate any results from them by channels.
// Good place to read results back is events handler (s.SetEvents).
// It is called after every update.
func runUI() {
	// Window configuration options.
	wcfg := grue.WindowConfig{
		Title:          "Grue example",
		WindowGeometry: grue.R(0, 0, 1000, 600),
		FPS:            60,
	}

	// Create primary surface (this includes window).
	s, err := pix.NewPrimarySurface(wcfg, grue.SurfaceConfig{PixelSize: 1})
	if err != nil {
		panic(err)
	}

	theme, err := themes.NewLight(s, "assets/caladea-bold.ttf", 20, "assets/theme-light.json")
	if err != nil {
		panic(err)
	}

	// err = s.InitImages("assets/test.json")
	// if err != nil {
	// 	panic(err)
	// }

	// Create toplevel panel.
	pn := grue.NewPanel(s.Root(), grue.Base{
		Rect:      grue.R(20, 20, 480, 280),
		BackColor: grue.RGB(0.2, 0.7, 1.0),
	})

	bt1 := grue.NewPushButton(pn, grue.Base{
		Rect:      grue.R0(120, 40),
		BackColor: grue.RGB(0.7, 0.7, 0.7),
		Text:      "Hello",
		ForeColor: color.Black,
	})
	bt1.Place(grue.V(50, 50))

	s.SetEvents(func() {
		theme.DrawButton.Draw(s, grue.R(100, 400, 300, 500))
	})

	// Run main loop.
	s.Run()
}

func runUI2() {
	// Window configuration options.
	wcfg := grue.WindowConfig{
		Title:          "Grue example",
		WindowGeometry: grue.R(0, 0, 1000, 600),
		FPS:            60,
	}

	// Create primary surface (this includes window).
	s, err := pix.NewPrimarySurface(wcfg, grue.SurfaceConfig{PixelSize: 1})
	if err != nil {
		panic(err)
	}

	s2, err := pix.NewSurface(s, grue.SurfaceConfig{PixelSize: 1})
	if err != nil {
		panic(err)
	}

	s.InitTTF("default", "themes/caladea-bold.ttf", 12, grue.Charset{})

	// Create toplevel panel.
	pn := grue.NewPanel(s.Root(), grue.Base{
		Rect:      grue.R(20, 20, 480, 280),
		BackColor: grue.RGB(0.2, 0.7, 1.0),
	})

	bt1 := grue.NewPushButton(pn, grue.Base{
		Rect:      grue.R0(80, 20),
		BackColor: grue.RGB(0.7, 0.7, 0.7),
		Text:      "Pushme",
		ForeColor: color.Black,
	})
	bt1.Place(grue.V(50, 50))

	// Note that bt2 is twice as large.
	bt2 := grue.NewPushButton(s2.Root(), grue.Base{
		Rect:      grue.R0(160, 40),
		BackColor: grue.RGB(0.7, 0.8, 0.8),
		Text:      "And me",
		ForeColor: grue.RGB(0.5, 0.0, 0.0),
	})
	// And placed at different Y offset.
	// But will appear at same height as bt1.
	// Because bt1 is at surface with pixel size = 2.
	// Also, bt1 is on panel which itself have Y offest = 20.
	bt2.Place(grue.V(400, 140))

	// Run main loop.
	s.Run()
}
