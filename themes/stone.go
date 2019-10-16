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
		PressDisplace:     grue.V(1, -1),
		Drawers: map[grue.ThemeDrawerKey]grue.ThemeDrawer{
			grue.ThemePanel: grue.TexturedPanel{
				Image: "stone-pn",
				Left:  2, Right: 2, Top: 2, Bottom: 2,
				TileHorizontal: true, TileVertical: true,
			},
			grue.ThemePanelDisabled: grue.TexturedPanel{
				Image: "stone-pn",
				Left:  2, Right: 2, Top: 2, Bottom: 2,
				TileHorizontal: true, TileVertical: true,
			},
			grue.ThemeButton: grue.TexturedPanel{
				Image: "stone-bt",
				Color: grue.RGB(0.9, 0.9, 0.9),
				Left:  2, Right: 2, Top: 2, Bottom: 2,
				TileHorizontal: true, TileVertical: true,
			},
			grue.ThemeButtonDisabled: grue.TexturedPanel{
				Image: "stone-bt",
				Color: grue.RGB(0.9, 0.9, 0.9),
				Left:  2, Right: 2, Top: 2, Bottom: 2,
				TileHorizontal: true, TileVertical: true,
			},
			grue.ThemeButtonHL: grue.TexturedPanel{
				Image: "stone-bt",
				Color: grue.RGB(0.9, 0.9, 0.9),
				Left:  2, Right: 2, Top: 2, Bottom: 2,
				TileHorizontal: true, TileVertical: true,
			},
			grue.ThemeButtonActive: grue.TexturedPanel{
				Image: "stone-bt-act",
				Color: grue.RGB(0.8, 1, 1),
				Left:  2, Right: 2, Top: 2, Bottom: 2,
				TileHorizontal: true, TileVertical: true,
			},
			grue.ThemeLineEdit: grue.TexturedPanel{
				Image: "stone-bt",
				Left:  2, Right: 2, Top: 2, Bottom: 2,
				Color: grue.RGB(0.8, 0.8, 1),
			},
			grue.ThemeTooltip: grue.TexturedPanel{
				Image: "stone-pn",
				Left:  2, Right: 2, Top: 2, Bottom: 2,
				TileHorizontal: true, TileVertical: true,
				Color: grue.RGB(1, 0.95, 0.8),
			},
		},
		CursorDrawer: grue.RectCursorDrawer{
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
