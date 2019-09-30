package grue

// Base is collection of initializable fields for widget.
type Base struct {
	Theme    *Theme
	Rect     Rect
	Text     string
	Tooltip  string
	Disabled bool

	// If true, widget is invisible (but children are).
	// OnDraw is still called and can paint.
	Phantom bool
}

// Panel is simple widget with background color and border.
// It can contain text and/or image.
type Panel struct {
	// Node implements tree-like hierarchy.
	*node

	// Virt is interface by which virtual calls can be made,
	// i. e. this panel might be embedded to button, then
	// Render calls panel's version, but virt.Render will call
	// button's one.
	// Of course, every embedding type must assign itself
	// to virt to make this work.
	virt Widget

	// Base is panel's base.
	Base

	// Interactive provides response to input events.
	Interactive

	// Custom drawing function.
	OnDraw func()

	// Graphics surface.
	Surface Surface
}

// Interactive is entity that responds to input (mouse or keyboard).
type Interactive struct {
	OnMouseIn    func()
	OnMouseOut   func()
	OnMouseMove  func()
	OnMouseDown  func(button Button)
	OnMouseUp    func(button Button)
	OnMouseClick func(button Button)
	OnMouseWheel func()
	OnKeys       func() bool
}

// NewPanel creates new panel.
func NewPanel(parent Node, b Base) *Panel {
	w := &Panel{
		node: &node{},
		Base: b,
	}
	initWidget(parent, w)
	return w
}

// initWidget initializes specific Widget behavior which must be
// repeated in all panel descendants by actual type:
// - sets up virtual;
// - attaches widget to parent;
// - derives surface.
func initWidget(parent Node, w Widget) {
	pn := w.GetPanel()
	if pn != nil {
		pn.virt = w
	}
	if parent != nil {
		parent.Foster(w)
		pw := parent.(Widget).GetPanel()
		if pw != nil && pn != nil {
			pn.Surface = pw.Surface
		}
	}
}

// Close removes panel from tree.
// It is needed over node.Close, because that operates on
// *node part of Panel, but must operate on *Panel.
// It's suitable for descendants without redefenition:
// node tree will contain Panels only, who exhibit
// virtual behaviour (via w.virt).
func (w *Panel) Close() {
	w.removeChildren()
	if w.parent != nil {
		w.parent.removeChild(w)
		w.parent = nil
	}
}

// Foster ... (same).
func (w *Panel) Foster(ch Node) {
	if ch == nil {
		return
	}
	p := ch.getParent()
	if p == w {
		return
	}
	if p != nil {
		p.removeChild(ch)
	}

	ch.setParent(w)
	w.addChild(ch)
}

// SubWidgets returns a slice of child widgets.
func (w *Panel) SubWidgets() []Widget {
	sns := w.SubNodes()
	res := make([]Widget, 0, len(sns))
	for _, v := range sns {
		if wv, ok := v.(Widget); ok {
			res = append(res, wv)
		}
	}
	return res
}

// Paint draws the widget without children.
func (w *Panel) paint() {
	if !w.Phantom {
		r := w.GlobalRect()
		theme := w.Theme
		if theme == nil {
			theme = w.Surface.GetTheme()
		}
		tdef, _ := theme.Drawers[ThemePanel]
		var tcur ThemeDrawer
		tcol := theme.TextColor
		switch {
		case w.Disabled:
			tcur, _ = theme.Drawers[ThemePanelDisabled]
			tcol = theme.DisabledTextColor
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
	if w.OnDraw != nil {
		w.OnDraw()
	}
}

// Render widget and its children on the screen.
func (w *Panel) Render() {
	w.virt.paint()
	for _, c := range w.SubWidgets() {
		c.Render()
	}
}

// ProcessMouse generates mouse events based on change in mouse coords.
func (w *Panel) ProcessMouse() {
	r := w.GlobalRect()
	lcont := r.Contains(w.Surface.PrevMousePos())
	cont := r.Contains(w.Surface.MousePos())
	if !lcont && cont {
		if w.OnMouseIn != nil {
			w.OnMouseIn()
		}
	} else if lcont && !cont {
		if w.OnMouseOut != nil {
			w.OnMouseOut()
		}
	}

	if !cont {
		return
	}

	w.Surface.SetToolTip(w.Tooltip)

	checkPress := func(bt Button) {
		if w.Surface.JustPressed(bt) {
			if w.OnMouseDown != nil {
				w.OnMouseDown(bt)
			}
		}
		if w.Surface.JustReleased(bt) {
			if w.OnMouseUp != nil {
				w.OnMouseUp(bt)
			}
			len := w.Surface.PrevMousePos().Add(
				V(-w.Surface.ClickMousePos().X, -w.Surface.ClickMousePos().Y)).Len()
			if len <= 8 && w.OnMouseClick != nil {
				w.OnMouseClick(bt)
			}
		}
	}
	checkPress(MouseButtonLeft)
	checkPress(MouseButtonRight)
	checkPress(MouseButtonMiddle)

	if w.Surface.MouseScroll() != V(0, 0) && w.OnMouseWheel != nil {
		w.OnMouseWheel()
	}

	if w.Surface.PrevMousePos() != w.Surface.MousePos() && w.OnMouseMove != nil {
		w.OnMouseMove()
	}

	for _, c := range w.SubWidgets() {
		c.ProcessMouse()
	}
}

// ProcessKeys calls keyboard handlers on the widget
// hierarchy. If any widget reports, that key is processed,
// event propagation stops.
func (w *Panel) ProcessKeys() {
	if w.OnKeys != nil && w.OnKeys() {
		return
	}
	for _, c := range w.SubWidgets() {
		c.ProcessKeys()
	}
}

// GetPanel returns panel part of the widget.
func (w *Panel) GetPanel() *Panel {
	return w
}

// GlobalRect is absolute widget rectangle (screen coords).
func (w *Panel) GlobalRect() (r Rect) {
	r = w.Rect
	if w.parent == nil {
		return
	}
	parent := w.parent.(Widget)
	r = r.Moved(parent.GetPanel().GlobalRect().Min)
	return
}

// Place moves Widget to position relative to it's parent.
// Positive numbers set position relative to left/bottom edges
// of parent. Negative -- to right/top.
func (w *Panel) Place(rel Vec) {
	parent := w.getParent()
	if parent == nil {
		return
	}
	pw := parent.(Widget)
	if pw == nil {
		return
	}
	ppn := pw.GetPanel()
	if rel.X < 0 {
		rel.X = ppn.Rect.Size().X - w.Rect.Size().X - (rel.X + 1)
	}
	if rel.Y < 0 {
		rel.Y = ppn.Rect.Size().Y - w.Rect.Size().Y - (rel.Y + 1)
	}
	sz := V(w.Rect.Size().X, w.Rect.Size().Y)
	w.Rect.Min = rel
	w.Rect.Max = w.Rect.Min.Add(sz)
}
