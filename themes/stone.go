package themes

import "github.com/gremour/grue"

// NewStone creates new stone theme.
func NewStone(s grue.Surface, fontFile string, fontSize float64, sheetFile string) (grue.Theme, error) {
	err := s.InitTTF("stone-title", fontFile, fontSize, grue.Charset{})
	if err != nil {
		return grue.Theme{}, err
	}
	err = s.InitImages(sheetFile)
	if err != nil {
		return grue.Theme{}, err
	}

	pnmd := grue.MultiDrawer{Drawers: []grue.ThemeDrawer{
		TexturedPanel{
			Image:          "stone-pn",
			TileHorizontal: true, TileVertical: true,
		},
		TexturedPanel{
			Image: "stone-orn3",
			Left:  10, Right: 10, Top: 10, Bottom: 10,
		},
	}}
	btmd := grue.MultiDrawer{Drawers: []grue.ThemeDrawer{
		TexturedPanel{
			Image:          "stone-bt",
			TileHorizontal: true, TileVertical: true,
		},
		TexturedPanel{
			Image: "stone-orn2",
			Left:  4, Right: 4, Top: 4, Bottom: 4,
		},
	}}
	btmda := grue.MultiDrawer{Drawers: []grue.ThemeDrawer{
		TexturedPanel{
			Image:          "stone-bt",
			TileHorizontal: true, TileVertical: true,
			Color: grue.RGB(0.7, 0.7, 0.7),
		},
		TexturedPanel{
			Image: "stone-orn2",
			Left:  4, Right: 4, Top: 4, Bottom: 4,
			Color: grue.RGB(0.7, 0.7, 0.7),
		},
	}}
	btmdhl := btmd
	btmdhl.Drawers = append(btmdhl.Drawers, ParticleDrawer{})
	lemd := grue.MultiDrawer{Drawers: []grue.ThemeDrawer{
		TexturedPanel{
			Image:          "stone-le",
			TileHorizontal: true, TileVertical: true,
			Left: 6, Right: 6, Top: 6, Bottom: 6,
		},
		TexturedPanel{
			Image: "stone-orn2",
			Left:  6, Right: 6, Top: 6, Bottom: 6,
		},
	}}
	lemdhl := lemd
	lemdhl.Drawers = append(lemdhl.Drawers, ParticleDrawer{})

	theme := grue.Theme{
		TitleFont:         "stone-title",
		TooltipFont:       "stone-title",
		TextColor:         grue.RGB(0.9, 0.7, 0.55),
		PanelTextColor:    grue.RGB(0, 0, 0),
		EditTextColor:     grue.RGB(1, 1, 1),
		DisabledTextColor: grue.RGB(0.8, 0.5, 0.5),
		PlaceholderColor:  grue.RGB(0.7, 0.7, 0.7),
		TooltipColor:      grue.RGB(0, 0, 0),
		Pad:               8,
		//		PressDisplace:     grue.V(1, -1),
		Drawers: map[grue.ThemeDrawerKey]grue.ThemeDrawer{
			grue.ThemePanel:        pnmd,
			grue.ThemeButton:       btmd,
			grue.ThemeButtonActive: btmda,
			grue.ThemeButtonHL:     btmdhl,
			grue.ThemeLineEdit:     lemd,
			grue.ThemeLineEditHL:   lemdhl,
			grue.ThemeTooltip: PlainRect{
				BackColor:   grue.RGB(1, 0.95, 0.8),
				BorderColor: grue.RGB(0, 0, 0),
				BorderSize:  1,
			},
		},
		CursorDrawer: RectCursorDrawer{
			Color1:        grue.RGB(1, 1, 1),
			Color2:        grue.RGBA(0, 0, 0, 0),
			Width:         3,
			PulseInterval: 1,
		},
	}

	if s.GetTheme() == nil {
		s.SetTheme(&theme)
	}
	return theme, nil
}
