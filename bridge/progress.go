package bridge

import (
	"sync"

	"fyne.io/fyne/widget"
	"github.com/psanford/wormhole-william/wormhole"
)

// SendProgress is contains a widget for displaying wormhole send progress.
type SendProgress struct {
	Widget *widget.ProgressBar

	// Update is the SendOption that should be passed to the wormhole client.
	Update wormhole.SendOption
	once   sync.Once
}

// NewSendProgress creates a new fyne progress bar and update function for wormhole send.
func NewSendProgress() *SendProgress {
	p := &SendProgress{Widget: widget.NewProgressBar()}

	p.Update = wormhole.WithProgress(func(sent int64, total int64) {
		p.once.Do(func() { p.Widget.Max = float64(total) })
		p.Widget.SetValue(float64(sent))
	})

	return p
}
