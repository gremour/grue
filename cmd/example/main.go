package main

import (
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
	// Primary surface configuration options.
	scfg := grue.SurfaceConfig{
		// Window options
		Title:          "Grue example",
		WindowGeometry: grue.R(0, 0, 1000, 600),
		FPS:            60,
		// Surface options
		PixelSize: 1,
		BackColor: grue.RGB(0.1, 0, 0),
	}

	// Create primary surface (this includes window).
	s, err := pix.NewPrimarySurface(scfg)
	if err != nil {
		panic(err)
	}

	_, err = themes.NewLight(s, "assets/caladea-bold.ttf", 20, "assets/theme-light.json")
	if err != nil {
		panic(err)
	}

	// Create toplevel panel.
	pn := grue.NewPanel(s.Root(), grue.Base{
		Rect: grue.R(20, 20, 480, 280),
	})

	pn1 := grue.NewPanel(pn, grue.Base{
		Rect: grue.R0(250, 100),
		Text: ":)",
	})
	pn1.Place(grue.V(50, 120))

	bt1 := grue.NewPushButton(pn, grue.Base{
		Rect: grue.R0(120, 40),
		Text: "Hello",
	})
	bt1.Place(grue.V(50, 50))

	bt2 := grue.NewPushButton(pn, grue.Base{
		Rect:     grue.R0(120, 40),
		Text:     "Grue",
		Tooltip:  "Graphical UI lib",
		Disabled: true,
	})
	bt2.Place(grue.V(180, 50))

	bt1.OnPress = func() {
		btPop := grue.NewPushButton(s.Root(), grue.Base{
			Rect: grue.R0(120, 40).Moved(s.MousePos()),
			Text: "Pop",
		})
		s.PopUp(btPop)
	}

	s.SetEvents(func() {
	})

	// Run main loop.
	s.Run()
}
