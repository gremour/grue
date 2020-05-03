package grue

import "fmt"

// PopupMenu is menu to use as popup.
// Menu contains a number of options
// which are represented as buttons
// and handlers for each option activation.
// To create hierarchical menu,
// create new popup menu from option handler.
type PopupMenu struct {
	*Panel

	opts []optWidget
}

type optWidget struct {
	opt MenuOption
	w   Widget
}

// MenuOption is one option for popup menu.
type MenuOption struct {
	// Identifier string to search options by it.
	ID       string
	Text     string
	Image    string
	Disabled bool
	// Handler is called when option is activated
	Handler func(pm *PopupMenu) bool
}

// NewPopupMenu creates new popup menu.
// Menu is located relative to topleft of
// provided Rect (in Base). Height is used
// as distance between options vertically
// (including height of button).
func NewPopupMenu(parent Widget, b Base, mo ...MenuOption) *PopupMenu {
	pm := &PopupMenu{
		Panel: NewPanel(nil, b),
	}
	InitWidget(parent, pm)
	pad := pm.MyTheme().Pad
	btH := b.Rect.H() - pad
	btW := b.Rect.W() - pad*2
	pm.Rect = R(pm.Rect.Min.X, pm.Rect.Max.Y-pm.Rect.H()*float64(len(mo))-pad,
		pm.Rect.Max.X, pm.Rect.Max.Y)
	y := pm.Rect.H() - pad
	for _, o := range mo {
		bt := NewPushButton(pm, Base{
			Rect:     R(pad, y-btH, pad+btW, y),
			Text:     o.Text,
			Image:    o.Image,
			Disabled: o.Disabled,
		})
		y -= btH + pad
		o := o
		bt.OnPress = func() {
			fmt.Printf("popup menu onpress before popdown o=%+v\n", o)
			bt.Surface.PopDownTo(pm)
			fmt.Printf("popup menu onpress after popdown o=%+v\n", o)
			if o.Handler != nil {
				fmt.Printf("popup menu onpress before handler\n")
				close := o.Handler(pm)
				fmt.Printf("popup menu onpress after handler =%v\n", close)
				if close {
					pm.Surface.PopDownTo(nil)
				}
			}
		}
		ow := optWidget{
			opt: o,
			w:   bt,
		}
		pm.opts = append(pm.opts, ow)
	}
	pm.Surface.PopUp(pm)
	return pm
}

// GetButton returns option's button by ID.
func (pm *PopupMenu) GetButton(id string) *PushButton {
	for _, ow := range pm.opts {
		if ow.opt.ID == id {
			pb, ok := ow.w.(*PushButton)
			if !ok {
				return nil
			}
			return pb
		}
	}
	return nil
}
