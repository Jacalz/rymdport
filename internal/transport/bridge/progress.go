package bridge

import (
	"sync"

	"fyne.io/fyne"
	"fyne.io/fyne/container"
	"fyne.io/fyne/widget"
	"github.com/psanford/wormhole-william/wormhole"
)

// sendProgress is contains a widget for displaying wormhole send progress.
type sendProgress struct {
	widget.ProgressBar

	// Update is the SendOption that should be passed to the wormhole client.
	update wormhole.SendOption
	once   sync.Once
}

// UpdateProgress is the function that runs when updating the progress.
func (p *sendProgress) updateProgress(sent int64, total int64) {
	p.once.Do(func() { p.Max = float64(total) })
	p.SetValue(float64(sent))
}

// newSendProgress creates a new fyne progress bar and update function for wormhole send.
func newSendProgress() *sendProgress {
	p := &sendProgress{}
	p.ExtendBaseWidget(p)
	p.update = wormhole.WithProgress(p.updateProgress)

	return p
}

type recvProgress struct {
	widget.ProgressBarInfinite
	done       *widget.ProgressBar
	container  *fyne.Container
	statusText string
}

func (r *recvProgress) status() string {
	return r.statusText
}

func (r *recvProgress) finished() {
	r.Stop()
	r.Hide()
	r.done.Show()
}

func (r *recvProgress) completed() {
	r.done.Value = 1.0
	r.finished()
}

func (r *recvProgress) failed() {
	r.done.Value = 0.0
	r.finished()
}

func (r *recvProgress) setStatus(stat string) {
	switch stat {
	case "start":
		r.Start()
	case "Failed":
		r.failed()
	case "Completed":
		r.completed()
	}

	r.statusText = stat
}

func newRecvProgress() *fyne.Container {
	r := &recvProgress{done: &widget.ProgressBar{}}
	r.ExtendBaseWidget(r)
	r.container = container.NewMax(r, r.done)

	r.done.TextFormatter = r.status
	r.done.Hide()
	return r.container
}
