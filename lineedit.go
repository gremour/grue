package grue

// LineEdit is a widget to input text.
type LineEdit struct {
	*Panel
	CursorPos  int
	TextOffset int

	OnTextChanged     func()
	OnEditingFinished func()
}

// NewLineEdit creates new line edit.
func NewLineEdit(parent Widget, b Base) *LineEdit {
	le := &LineEdit{
		Panel: NewPanel(nil, b),
	}
	InitWidget(parent, le)

	le.OnMouseDown = func(bt Button) {
		if bt != MouseButtonLeft {
			return
		}
		le.CursorPos = len(le.Text)
		le.TextOffset = 0
	}
	le.OnKeys = le.onKeys

	return le
}

// Paint draws the widget without children.
func (le *LineEdit) Paint() {
	r := le.GlobalRect()
	theme := le.MyTheme()
	tcur, _ := theme.Drawers[ThemeLineEdit]
	tcol := theme.TextColor
	txt := le.Text
	if le.Text == "" {
		tcol = theme.PlaceholderColor
		txt = le.PlaceholderText
	}
	tcur.Draw(le.Surface, r)
	editMode := le.Equals(le.Surface.Focus())
	if editMode {
		curRight := theme.Pad
		if le.Text != "" && le.CursorPos >= le.TextOffset {
			cursorAfterStr := le.Text[le.TextOffset : le.CursorPos-le.TextOffset]
			curRight += le.Surface.GetTextRect(cursorAfterStr, theme.TitleFont).W()
		}
		le.DrawImageAndText("", le.Text, tcol, 0, AlignLeft)
		theme.CursorDrawer.Draw(le.Surface, r.Min.Add(V(curRight, theme.Pad)), le.Rect.H()-theme.Pad*2)
	} else {
		le.DrawImageAndText("", txt, tcol, 0, AlignLeft)
	}
}

func (le *LineEdit) onKeys() bool {
	switch {
	case le.Surface.JustPressed(KeyEnter):
		if le.OnEditingFinished != nil {
			le.OnEditingFinished()
		}
		le.Surface.SetFocus(nil)
	case le.Surface.JustPressed(KeyBackspace) || le.Surface.Repeated(KeyBackspace):
		if le.CursorPos == 0 {
			break
		}
		if le.Text != "" {
			le.Text = le.Text[:le.CursorPos-1] + le.Text[le.CursorPos:len(le.Text)]
			le.CursorPos--
		}
		if le.OnTextChanged != nil {
			le.OnTextChanged()
		}
	case le.Surface.JustPressed(KeyDelete) || le.Surface.Repeated(KeyDelete):
		if le.CursorPos >= len(le.Text) {
			break
		}
		if le.Text != "" {
			le.Text = le.Text[:le.CursorPos] + le.Text[le.CursorPos+1:len(le.Text)]
		}
		if le.OnTextChanged != nil {
			le.OnTextChanged()
		}
	case le.Surface.JustPressed(KeyLeft) || le.Surface.Repeated(KeyLeft):
		le.CursorPos--
		if le.CursorPos < 0 {
			le.CursorPos = 0
		} else if le.CursorPos < le.TextOffset {
			le.TextOffset = le.CursorPos
		}
	case le.Surface.JustPressed(KeyRight) || le.Surface.Repeated(KeyRight):
		le.CursorPos++
		if le.CursorPos > len(le.Text) {
			le.CursorPos = len(le.Text)
		}
	case le.Surface.JustPressed(KeyHome):
		le.CursorPos = 0
	case le.Surface.JustPressed(KeyEnd):
		le.CursorPos = len(le.Text)
	default:
		t := le.Surface.KeysInput()
		if len(t) > 0 {
			le.Text = le.Text[:le.CursorPos] + t + le.Text[le.CursorPos:len(le.Text)]
			le.CursorPos += len(t)
		}
		if le.OnTextChanged != nil {
			le.OnTextChanged()
		}
	}
	return true
}
