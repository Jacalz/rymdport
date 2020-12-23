package bridge

import (
	"fyne.io/fyne"
	"fyne.io/fyne/widget"
)

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
func NewPressEntry(placeholder string, onReturn func()) *PressEntry {
	p := &PressEntry{OnReturn: onReturn}
	p.Entry.PlaceHolder = placeholder
	p.ExtendBaseWidget(p)
	return p
}
