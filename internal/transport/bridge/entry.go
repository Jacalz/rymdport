package bridge

import (
	"fyne.io/fyne"
	"fyne.io/fyne/widget"
)

// PressEntry is an extended entry for running a function on pressing enter.
type PressEntry struct {
	widget.Entry
	OnEnter func()
}

// TypedKey handles the key presses inside our UsernameEntry and uses Action to press the linked button.
func (p *PressEntry) TypedKey(ev *fyne.KeyEvent) {
	switch ev.Name {
	case fyne.KeyReturn, fyne.KeyEnter: // fyne.KeyReturn is the enter/return key on the keyboard, fyne.KeyEnter is on the NumPad.
		if p.OnEnter != nil {
			p.OnEnter()
		}
	default:
		p.Entry.TypedKey(ev)
	}
}

// NewPressEntry returns a new entry that runs a function on pressing return.
func NewPressEntry(placeholder string, onEnter func()) *PressEntry {
	p := &PressEntry{OnEnter: onEnter}
	p.Entry.PlaceHolder = placeholder
	p.ExtendBaseWidget(p)
	return p
}
