package main

import (
	"math"

	"github.com/gremour/grue"
	"github.com/gremour/grue/particles"
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
		WindowGeometry: grue.R(0, 0, 500, 400),
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

	// _, err = themes.NewLight(s, "assets/caladea-bold.ttf", 20, "assets/theme-light.json")
	// if err != nil {
	// 	panic(err)
	// }
	_, err = themes.NewStone(s, "assets/caladea-bold.ttf", 20, "assets/theme-stone.json")
	if err != nil {
		panic(err)
	}

	err = s.InitImages("assets/images.json")
	if err != nil {
		panic(err)
	}

	// Create toplevel panel.
	pn := grue.NewPanel(s.Root(), grue.Base{
		Rect: grue.R(20, 20, 480, 380),
	})

	pn1 := grue.NewPanel(pn, grue.Base{
		Rect: grue.R0(250, 200),
		Text: ":)",
	})
	pn1.Place(grue.V(50, 120))

	le := grue.NewLineEdit(pn1, grue.Base{
		Rect:            grue.R0(230, 40),
		PlaceholderText: "placeholder",
	})
	le.OnTextChanged = func() {
		pn1.Text = le.Text
	}
	le.Place(grue.V(10, 10))
	polish(le.Panel)

	bt1 := grue.NewPushButton(pn, grue.Base{
		Rect: grue.R0(120, 40),
		Text: "Hello",
	})
	bt1.Place(grue.V(50, 50))
	polish(bt1.Panel)

	bt2 := grue.NewPushButton(pn, grue.Base{
		Rect: grue.R0(200, 40),
		Text: "Grue",
		// TextAlign:  grue.AlignRight,
		Tooltip:  "Graphical UI lib",
		Disabled: true,
		Image:    "grue-logo20",
		// ImageAlign: grue.AlignRight,
	})
	bt2.Place(grue.V(180, 50))
	polish(bt2.Panel)

	bt1.OnPress = func() {
		btPop := grue.NewPushButton(s.Root(), grue.Base{
			Rect: grue.R0(120, 40).Moved(s.MousePos()),
			Text: "Pop",
		})
		polish(btPop.Panel)
		s.PopUp(btPop)
	}

	s.SetEvents(func() {
	})

	// grue.PrintWidgets(s.Root(), "")

	// Run main loop.
	s.Run()
}

// btPolish initializes particles generator for button.
// Particles generator is used in theme if
// theme has ParticleDrawer.
func polish(pn *grue.Panel) {
	pg := &particles.Group{
		Generator: &particles.GlitterEdge{
			Rect:           pn.GlobalRect().Expanded(-2),
			Placer:         particles.BorderPlacer{},
			Image:          "ptc-star",
			MaxSize:        4,
			LifeTime:       1,
			SpawnTempo:     2,
			MinParticles:   4,
			MaxParticles:   16,
			SizeFunc:       math.Sin,
			SizeFuncMaxArg: math.Pi,
			Color:          grue.RGBA(1, 0.85, 0.5, 0.5),
		},
	}
	pn.Extras = []interface{}{pg}
}
