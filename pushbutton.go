package grue

// PushButton is pressable and optionally, checkable (TODO), button.
type PushButton struct {
	*Panel
	Hilited bool
	Pressed bool

	OnPress func()
}

// NewPushButton creates new button.
func NewPushButton(parent Node, b Base) *PushButton {
	w := &PushButton{
		Panel: NewPanel(nil, b),
	}
	initWidget(parent, w)

	w.OnMouseIn = func() {
		w.Hilited = true
	}
	w.OnMouseOut = func() {
		w.Hilited = false
		w.Pressed = false
	}
	w.OnMouseDown = func(bt Button) {
		if bt != MouseButtonLeft {
			return
		}
		w.Pressed = true
	}
	w.OnMouseUp = func(bt Button) {
		if bt != MouseButtonLeft {
			return
		}
		w.Pressed = false
		if w.OnPress != nil {
			w.OnPress()
		}
	}
	return w
}

// Paint draws the widget without children.
func (w *PushButton) paint() {
	r := w.GlobalRect()
	theme := w.Theme
	if theme == nil {
		theme = w.Surface.GetTheme()
	}
	tdef, _ := theme.Drawers[ThemeButton]
	var tcur ThemeDrawer
	tcol := theme.TextColor
	switch {
	case w.Disabled:
		tcur, _ = theme.Drawers[ThemeButtonDisabled]
		tcol = theme.DisabledTextColor
	case w.Pressed:
		tcur, _ = theme.Drawers[ThemeButtonActive]
	case w.Hilited:
		tcur, _ = theme.Drawers[ThemeButtonHL]
	}
	if tcur != nil {
		tdef = tcur
	}
	if tdef != nil {
		tdef.Draw(w.Surface, r)
	}
	if len(w.Text) > 0 {
		w.Surface.DrawText(w.Text, theme.TitleFont, r, tcol, AlignCenter, AlignCenter)
	}
}
