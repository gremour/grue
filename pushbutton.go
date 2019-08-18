package grue

// PushButton is pressable and optionally, checkable (TODO), button.
type PushButton struct {
	*Panel
	hilited bool
	pressed bool

	OnPress func()
}

// DefaultPushButton is default Base for pushbutton.
var DefaultPushButton = Base{
	BackColor:   RGB(0.7, 0.7, 0.7),
	Border:      1,
	BorderInset: 2,
}

// NewPushButton creates new button.
func NewPushButton(parent Node, b Base) *PushButton {
	w := &PushButton{
		Panel: NewPanel(nil, initDefaultBase(b, DefaultPushButton)),
	}
	initWidget(parent, w)

	w.OnMouseIn = func() {
		w.hilited = true
	}
	w.OnMouseOut = func() {
		w.hilited = false
		w.pressed = false
	}
	w.OnMouseDown = func(bt Button) {
		if bt != MouseButtonLeft {
			return
		}
		w.pressed = true
	}
	w.OnMouseUp = func(bt Button) {
		if bt != MouseButtonLeft {
			return
		}
		w.pressed = false
		if w.OnPress != nil {
			w.OnPress()
		}
	}
	return w
}

// Paint draws the widget without children.
func (w *PushButton) paint() {
	bc := w.BackColor
	if w.hilited {
		w.BackColor = RGB(0.6, 1, 1)
	}
	bw := w.Border
	if w.pressed {
		w.Border = bw * 2
	}
	w.Panel.paint()
	w.Border = bw
	w.BackColor = bc
}
