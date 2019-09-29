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

	drwPn := PlainRect{
		BackColor:   grue.RGB(0.7, 0.7, 0.7),
		BorderColor: grue.RGB(0.4, 0.4, 0.4),
		BorderSize:  1,
	}
	drwBt := TexturedPanel{
		ImagePrefix: "light-bt",
	}

	theme := grue.Theme{
		TitleFont:  "light-title",
		DrawPanel:  drwPn,
		DrawButton: drwBt,
	}

	s.SetTheme(theme)
	return theme, nil
}
