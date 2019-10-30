package grue

// PushButton is pressable and optionally, checkable (TODO), button.
type PushButton struct {
	*Panel
	Pressed bool

	OnPress func()
}

// NewPushButton creates new button.
func NewPushButton(parent Widget, b Base) *PushButton {
	pb := &PushButton{
		Panel: NewPanel(nil, b),
	}
	InitWidget(parent, pb)

	pb.OnMouseOut = func() {
		pb.Pressed = false
	}
	pb.OnMouseDown = func(bt Button) {
		if bt != MouseButtonLeft {
			return
		}
		pb.Pressed = true
	}
	pb.OnMouseUp = func(bt Button) {
		if bt != MouseButtonLeft {
			return
		}
		pb.Pressed = false
		if pb.OnPress != nil {
			pb.OnPress()
		}
	}
	return pb
}

// Paint draws the widget without children.
func (pb *PushButton) Paint() {
	r := pb.GlobalRect()
	theme := pb.MyTheme()
	tdef, _ := theme.Drawers[ThemeButton]
	var tcur ThemeDrawer
	tcol := theme.ButtonTextColor
	if tcol == nil {
		tcol = theme.TextColor
	}
	var disp Vec
	switch {
	case pb.Disabled:
		tcur, _ = theme.Drawers[ThemeButtonDisabled]
		tcol = theme.DisabledTextColor
	case pb.Pressed:
		tcur, _ = theme.Drawers[ThemeButtonActive]
		disp = theme.PressDisplace
	case pb.PointerInside:
		tcur, _ = theme.Drawers[ThemeButtonHL]
	}
	if tcur != nil {
		tdef = tcur
	}
	if tdef != nil {
		tdef.Draw(pb.Surface, r, pb.Extras...)
	}
	pb.DrawImageAndText(pb.Image, pb.Text, tcol, pb.ImageAlign, pb.TextAlign, disp)
}
