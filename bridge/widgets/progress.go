package widgets

import (
	"sync"

	"fyne.io/fyne/widget"
	"github.com/psanford/wormhole-william/wormhole"
)

// SendProgress is contains a widget for displaying wormhole send progress.
type SendProgress struct {
	widget.ProgressBar

	// Update is the SendOption that should be passed to the wormhole client.
	Update wormhole.SendOption
	once   sync.Once
}

// UpdateProgress is the function that runs when updating the progress.
func (p *SendProgress) UpdateProgress(sent int64, total int64) {
	p.once.Do(func() { p.Max = float64(total) })
	p.SetValue(float64(sent))
}

// NewSendProgress creates a new fyne progress bar and update function for wormhole send.
func NewSendProgress() *SendProgress {
	p := &SendProgress{}
	p.ExtendBaseWidget(p)

	p.Update = wormhole.WithProgress(p.UpdateProgress)

	return p
}
