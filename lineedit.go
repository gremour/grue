package grue

// LineEdit is a widget to input text.
type LineEdit struct {
	*Panel
	CursorPos  int
	TextOffset int
	TextLimit  int

	OnTextChanged     func()
	OnEditingFinished func()
}

// NewLineEdit creates new line edit.
func NewLineEdit(parent Widget, b Base) *LineEdit {
	le := &LineEdit{
		Panel:     NewPanel(nil, b),
		TextLimit: 1000,
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
	tcur.Draw(le.Surface, r)
	editMode := le.Equals(le.Surface.Focus())

	tcol := theme.EditTextColor
	if tcol == nil {
		tcol = theme.TextColor
	}
	text := le.Text
	offs := le.TextOffset
	if le.Text == "" {
		tcol = theme.PlaceholderColor
		if tcol == nil {
			tcol = theme.TextColor
		}
		text = le.PlaceholderText
		offs = 0
	}
	text = le.Surface.FitText(text[offs:], theme.TitleFont, le.Rect.W()-theme.Pad*2)
	le.DrawImageAndText("", text, tcol, 0, AlignLeft, Vec{})

	if editMode {
		curRight := le.cursorPos()
		theme.CursorDrawer.Draw(le.Surface, r.Min.Add(V(curRight, theme.Pad)), le.Rect.H()-theme.Pad*2)
	}
}

func (le *LineEdit) onKeys() bool {
	switch {
	case le.Surface.JustPressed(KeyEnter):
		if le.OnEditingFinished != nil {
			le.OnEditingFinished()
		}
		le.Surface.SetFocus(nil)
		le.TextOffset = 0
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
		if len(le.Text) >= le.TextLimit {
			break
		}
		t := le.Surface.KeysInput()
		if len(t) > 0 {
			le.Text = le.Text[:le.CursorPos] + t + le.Text[le.CursorPos:len(le.Text)]
			le.CursorPos += len(t)
		}
		if le.OnTextChanged != nil {
			le.OnTextChanged()
		}
	}
	le.updateTextOffest()
	return true
}

// Return the horizontal position of cursor with current TextOffset
func (le *LineEdit) cursorPos() float64 {
	curRight := le.MyTheme().Pad
	if le.Text == "" || le.CursorPos == le.TextOffset {
		return curRight
	}
	if le.CursorPos > le.TextOffset {
		textPiece := le.Text[le.TextOffset:le.CursorPos]
		curRight += le.Surface.GetTextRect(textPiece, le.MyTheme().TitleFont).W()
	} else {
		if le.TextOffset > len(le.Text) {
			le.TextOffset = len(le.Text)
		}
		textPiece := le.Text[le.CursorPos:le.TextOffset]
		curRight -= le.Surface.GetTextRect(textPiece, le.MyTheme().TitleFont).W()
	}
	return curRight
}

// Update TextOffset so cursor is shown.
func (le *LineEdit) updateTextOffest() {
	pad := le.MyTheme().Pad
	curPos := le.cursorPos()
	for curPos < le.Rect.W()/2+pad && le.TextOffset > 0 {
		le.TextOffset--
		curPos = le.cursorPos()
	}
	for curPos >= le.Rect.W()-pad && le.TextOffset <= len(le.Text) {
		le.TextOffset++
		curPos = le.cursorPos()
	}
}
