package widgets

import (
	"fyne.io/fyne"
	"fyne.io/fyne/widget"
)

// NewBoldLabel returns a new label with bold text.
func NewBoldLabel(text string) *widget.Label {
	return widget.NewLabelWithStyle(text, fyne.TextAlignLeading, fyne.TextStyle{Bold: true})
}

// CodeLabel is a label extended to copy the code with a menu popup on rightclick.
type CodeLabel struct {
	widget.Label
	popUp *widget.PopUp
}

// TappedSecondary adds rightclick for showing a menu and copy code.
func (cl *CodeLabel) TappedSecondary(pe *fyne.PointEvent) {
	copyCode := fyne.NewMenuItem("Copy code", func() {
		clipboard := fyne.CurrentApp().Driver().AllWindows()[0].Clipboard()
		clipboard.SetContent(cl.Text)
	})

	pos := fyne.CurrentApp().Driver().AbsolutePositionForObject(cl).Add(fyne.NewPos(pe.Position.X, pe.Position.Y))
	c := fyne.CurrentApp().Driver().CanvasForObject(cl)
	menu := fyne.NewMenu("", copyCode)

	cl.popUp = widget.NewPopUpMenuAtPosition(menu, c, pos)
}

// NewCodeLabel creates a new code label.
func NewCodeLabel(code chan string) *CodeLabel {
	c := &CodeLabel{}
	c.ExtendBaseWidget(c)
	c.SetText("Waiting for code...")
	go c.SetText(<-code)
	return c
}

// PressEntry is an extended entry for running a function on pressing enter.
type PressEntry struct {
	widget.Entry
	OnReturn func()
}

// TypedKey handles the key presses inside our UsernameEntry and uses Action to press the linked button.
func (p *PressEntry) TypedKey(ev *fyne.KeyEvent) {
	switch ev.Name {
	case fyne.KeyReturn:
		if p.OnReturn != nil {
			p.OnReturn()
		}
	default:
		p.Entry.TypedKey(ev)
	}
}

// NewPressEntry returns a new entry that runs a function on pressing return.
func NewPressEntry(placeholder string) *PressEntry {
	p := &PressEntry{}
	p.ExtendBaseWidget(p)
	p.SetPlaceHolder(placeholder)
	return p
}
