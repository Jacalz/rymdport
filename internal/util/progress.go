package util

import (
	"sync"

	"fyne.io/fyne/v2/widget"
	"github.com/psanford/wormhole-william/wormhole"
)

// ProgressBar is contains a widget for displaying wormhole send progress.
type ProgressBar struct {
	widget.ProgressBar
}

// WithProgress returns a send option to update the progress.
func (p *ProgressBar) WithProgress() wormhole.SendOption {
	once := sync.Once{}
	return wormhole.WithProgress(func(sent, total int64) {
		once.Do(func() { p.Max = float64(total) })
		p.SetValue(float64(sent))
	})
}

// Done sets the value to max to indicate that it is finished.
func (p *ProgressBar) Done() {
	p.SetValue(p.Max)
}

// Failed sets the text to indicate a failure.
func (p *ProgressBar) Failed() {
	p.TextFormatter = func() string { return "Failed" }
	p.Refresh()
}

// NewProgressBar creates a new fyne progress bar and update function for wormhole send.
func NewProgressBar() *ProgressBar {
	p := &ProgressBar{}
	p.ExtendBaseWidget(p)
	return p
}
