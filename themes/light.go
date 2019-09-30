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
		TextColor:         grue.RGB(0, 0, 0),
		DisabledTextColor: grue.RGB(0.3, 0.3, 0.3),
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
				ImagePrefix: "light-bt",
			},
			grue.ThemeButtonDisabled: grue.TexturedPanel{
				ImagePrefix: "light-bt",
				Color:       grue.RGB(0.8, 0.8, 0.8),
			},
			grue.ThemeButtonHL: grue.TexturedPanel{
				ImagePrefix: "light-bt",
				Color:       grue.RGB(0.8, 1, 1),
			},
			grue.ThemeButtonActive: grue.TexturedPanel{
				ImagePrefix: "light-bt-act",
				Color:       grue.RGB(0.8, 1, 1),
			},
		},
	}

	s.SetTheme(theme)
	return theme, nil
}
