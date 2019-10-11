package grue

import "fmt"

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
	Parent   Widget
	Children []Widget

	// Virt is interface by which virtual calls can be made,
	// i. e. this panel might be embedded to button, then
	// Render calls panel's version, but virt.Render will call
	// button's one.
	// Of course, every embedding type must assign itself
	// to virt to make this work.
	// WARNING: Virt is required to be used in place of
	// Widget interface values in any of Panel and
	// it's derivative types methods. Otherwise "slicing" will
	// happen and Widget will be in fact Panel (or the type that
	// used pointer to self when creating value of Widget type)
	// instead of derivative type. This may have some side-effects,
	// although comparision is already handled by Equals method.
	Virt Widget

	// Base is panel's base.
	Base

	// True if pointer is inside the widget
	PointerInside bool

	// Interactive provides response to input events.
	Interactive

	// Custom drawing function. Called from Paint.
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

// InitWidget initializes specific Widget behavior which must be
// repeated in all panel descendants by actual type:
// - sets up virtual;
// - attaches widget to parent;
// - derives surface.
func InitWidget(parent Widget, w Widget) {
	w.GetPanel().Virt = w
	if parent == nil {
		return
	}
	parent.Foster(w)
	w.GetPanel().Surface = parent.GetPanel().Surface
}

// NewPanel creates new panel.
func NewPanel(parent Widget, b Base) *Panel {
	p := &Panel{
		Base: b,
	}
	InitWidget(parent, p)
	return p
}

// Paint draws panel without children.
func (p *Panel) Paint() {
	if !p.Phantom {
		r := p.GlobalRect()
		theme := p.Theme
		if theme == nil {
			theme = p.Surface.GetTheme()
		}
		tdef, _ := theme.Drawers[ThemePanel]
		var tcur ThemeDrawer
		tcol := theme.TextColor
		switch {
		case p.Disabled:
			tcur, _ = theme.Drawers[ThemePanelDisabled]
			tcol = theme.DisabledTextColor
		}
		if tcur != nil {
			tdef = tcur
		}
		if tdef != nil {
			tdef.Draw(p.Surface, r)
		}
		if len(p.Text) > 0 {
			p.Surface.DrawText(p.Text, theme.TitleFont, r, tcol, AlignCenter, AlignCenter)
		}
	}
	if p.OnDraw != nil {
		p.OnDraw()
	}
}

// GetPanel returns widget's panel
func (p *Panel) GetPanel() *Panel {
	return p
}

// Equals checks if panel and passed
// widget are both the same object
func (p *Panel) Equals(w Widget) bool {
	if p == nil && w == nil {
		return true
	}
	if p == nil || w == nil {
		return false
	}
	return p.Virt == w.GetPanel().Virt
}

// Render widget and its children on the screen.
func (p *Panel) Render() {
	p.Virt.Paint()
	for _, c := range p.Children {
		c.Render()
	}
}

// ProcessMouse generates mouse events based on change in mouse coords.
// wu holds top widget under mouse.
func (p *Panel) ProcessMouse(wu Widget) {
	r := p.GlobalRect()
	cont := r.Contains(p.Surface.MousePos())

	if cont && !p.PointerInside && p.OnMouseIn != nil {
		p.OnMouseIn()
	}

	if !cont && p.PointerInside && p.OnMouseOut != nil {
		p.OnMouseOut()
	}

	p.PointerInside = cont

	if !cont {
		return
	}

	p.Surface.SetToolTip(p.Tooltip)

	checkPress := func(bt Button) {
		if p.Surface.JustPressed(bt) {
			if p.OnMouseDown != nil {
				p.OnMouseDown(bt)
			}
		}
		if p.Surface.JustReleased(bt) {
			if p.OnMouseUp != nil {
				p.OnMouseUp(bt)
			}
			len := p.Surface.PrevMousePos().Add(
				V(-p.Surface.ClickMousePos().X, -p.Surface.ClickMousePos().Y)).Len()
			if len <= 8 && p.OnMouseClick != nil {
				p.OnMouseClick(bt)
			}
		}
	}
	if p.Equals(wu) {
		checkPress(MouseButtonLeft)
		checkPress(MouseButtonRight)
		checkPress(MouseButtonMiddle)
	}

	if p.Surface.MouseScroll() != V(0, 0) && p.OnMouseWheel != nil {
		p.OnMouseWheel()
	}

	if p.Surface.PrevMousePos() != p.Surface.MousePos() && p.OnMouseMove != nil {
		p.OnMouseMove()
	}

	for _, c := range p.Children {
		c.ProcessMouse(wu)
	}
}

// ProcessKeys calls keyboard handlers on the widget
// hierarchy. If any widget reports, that key is processed,
// event propagation stops.
func (p *Panel) ProcessKeys() {
	if p.OnKeys != nil && p.OnKeys() {
		return
	}
	for _, c := range p.Children {
		c.ProcessKeys()
	}
}

// WidgetUnder finds widget that is under given pointer coordinates
func (p *Panel) WidgetUnder(pos Vec) Widget {
	r := p.GlobalRect()
	if !r.Contains(pos) {
		return nil
	}

	for _, c := range p.Children {
		wu := c.WidgetUnder(pos)
		if wu != nil {
			return wu
		}
	}
	return p.Virt
}

// GlobalRect is absolute widget rectangle (screen coords).
func (p *Panel) GlobalRect() (r Rect) {
	r = p.Rect
	if p.Parent == nil {
		return
	}
	parent := p.Parent
	r = r.Moved(parent.GlobalRect().Min)
	return
}

// Place moves Widget to position relative to it's parent.
// Positive numbers set position relative to left/bottom edges
// of parent. Negative -- to right/top.
func (p *Panel) Place(rel Vec) {
	parent := p.Parent
	if parent == nil {
		return
	}
	ppn := parent.GetPanel()
	if rel.X < 0 {
		rel.X = ppn.Rect.Size().X - p.Rect.Size().X - (rel.X + 1)
	}
	if rel.Y < 0 {
		rel.Y = ppn.Rect.Size().Y - p.Rect.Size().Y - (rel.Y + 1)
	}
	sz := V(p.Rect.Size().X, p.Rect.Size().Y)
	p.Rect.Min = rel
	p.Rect.Max = p.Rect.Min.Add(sz)
}

// Close widget and its children.
func (p *Panel) Close() {
	p.removeChildren()
	if p.Parent != nil {
		p.Parent.removeChild(p.Virt)
		p.Parent = nil
	}
}

// Foster reconnects widget to this parent,
// removing it from previous parent if needed.
func (p *Panel) Foster(ch Widget) {
	if ch == nil || p == nil {
		return
	}
	par := ch.GetPanel().Parent
	if p.Equals(par) {
		return
	}
	if par != nil {
		par.removeChild(ch)
	}
	ch.GetPanel().Parent = p.Virt
	p.addChild(ch)
}

func (p *Panel) addChild(ch Widget) {
	for _, c := range p.Children {
		if c.Equals(ch) {
			// already in children.
			return
		}
	}
	p.Children = append(p.Children, ch)
}

func (p *Panel) removeChild(ch Widget) {
	pch := p.Children
	l := len(pch)
	for i, c := range pch {
		if c.Equals(ch) {
			if i < l-1 {
				copy(pch[i:l-1], pch[i+1:])
				// 0 1 2 3 4 5
				// a b c d e f
				//     ^i=2 l=6
				//     [2:5](c to e) replaced by [3:6](d to f)
			}
			pch[l-1] = nil
			p.Children = pch[:l-1]
			break
		}
	}
}

func (p *Panel) removeChildren() {
	for _, c := range p.Children {
		c.GetPanel().Parent = nil
		c.Close()
	}
	p.Children = nil
}

// PrintWidgets prints a tree of widgets for debugging.
func PrintWidgets(w Widget, indent string) {
	if w == nil {
		return
	}
	fmt.Printf("%vText:%v, Widget=%p, Panel=%p, parent=%p, children:%v\n",
		indent, w.GetPanel().Text, w, w.GetPanel(),
		w.GetPanel().Parent, w.GetPanel().Children)
	if w != w.GetPanel().Virt {
		fmt.Printf("!!! WARNING: Widget(%p) != Virt(%p) !!!\n", w, w.GetPanel().Virt)
	}
	for _, ch := range w.GetPanel().Children {
		PrintWidgets(ch, indent+" ")
	}
}
