package themes

import "github.com/gremour/grue"

// NewLight creates new light theme.
func NewLight(s grue.Surface, fontFile string, fontSize float64, sheetFile string) (grue.Theme, error) {
	err := s.InitTTF("light-title", fontFile, fontSize, grue.Charset{})
	if err != nil {
		return grue.Theme{}, err
	}
	err = s.InitImages(sheetFile)
	if err != nil {
		return grue.Theme{}, err
	}

	theme := grue.Theme{
		TitleFont:         "light-title",
		TooltipFont:       "light-title",
		TextColor:         grue.RGB(0, 0, 0),
		DisabledTextColor: grue.RGB(0.3, 0.3, 0.3),
		PlaceholderColor:  grue.RGB(0.8, 0.8, 0.8),
		TooltipColor:      grue.RGB(0, 0, 0),
		Pad:               8,
		Drawers: map[grue.ThemeDrawerKey]grue.ThemeDrawer{
			grue.ThemePanel: grue.PlainRect{
				BackColor:   grue.RGB(0.7, 0.7, 0.7),
				BorderColor: grue.RGB(0.4, 0.4, 0.4),
				BorderSize:  1,
			},
			grue.ThemePanelDisabled: grue.PlainRect{
				BackColor:   grue.RGB(0.5, 0.5, 0.5),
				BorderColor: grue.RGB(0.3, 0.3, 0.3),
				BorderSize:  1,
			},
			grue.ThemeButton: grue.TexturedPanel{
				Image: "light-bt",
				Left:  4, Right: 4, Top: 4, Bottom: 4,
			},
			grue.ThemeButtonDisabled: grue.TexturedPanel{
				Image: "light-bt",
				Color: grue.RGB(0.8, 0.8, 0.8),
				Left:  4, Right: 4, Top: 4, Bottom: 4,
			},
			grue.ThemeButtonHL: grue.TexturedPanel{
				Image: "light-bt",
				Color: grue.RGB(0.8, 1, 1),
				Left:  4, Right: 4, Top: 4, Bottom: 4,
			},
			grue.ThemeButtonActive: grue.TexturedPanel{
				Image: "light-bt-act",
				Color: grue.RGB(0.8, 1, 1),
				Left:  4, Right: 4, Top: 4, Bottom: 4,
			},
			grue.ThemeLineEdit: grue.TexturedPanel{
				Image: "light-le",
				Left:  4, Right: 4, Top: 4, Bottom: 4,
			},
			grue.ThemeTooltip: grue.PlainRect{
				BackColor:   grue.RGB(1, 0.95, 0.8),
				BorderColor: grue.RGB(0, 0, 0),
				BorderSize:  1,
			},
		},
		CursorDrawer: grue.RectCursorDrawer{
			Color1:        grue.RGB(0, 0, 0),
			Color2:        grue.RGBA(0, 0, 0, 0),
			Width:         2,
			PulseInterval: 1,
		},
	}

	if s.GetTheme() == nil {
		s.SetTheme(&theme)
	}
	return theme, nil
}
