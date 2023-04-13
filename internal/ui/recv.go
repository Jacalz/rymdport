package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/Jacalz/rymdport/v3/completion"
	"github.com/Jacalz/rymdport/v3/internal/transport"
	"github.com/Jacalz/rymdport/v3/internal/transport/bridge"
	"github.com/Jacalz/rymdport/v3/internal/util"
)

type recv struct {
	codeEntry *completionEntry
	data      *bridge.RecvData

	client *transport.Client
	window fyne.Window
}

func newRecv(w fyne.Window, c *transport.Client) *recv {
	return &recv{window: w, client: c}
}

func (r *recv) onRecv() {
	if err := r.codeEntry.Validate(); err != nil || r.codeEntry.Text == "" {
		dialog.ShowInformation("Invalid code", "The code is invalid. Please try again.", r.window)
		return
	}

	r.data.NewReceive(r.codeEntry.Text)
	r.codeEntry.SetText("")
}

func (r *recv) buildUI() *fyne.Container {
	r.codeEntry = newCompletionEntry(r.client, r.window.Canvas())
	r.codeEntry.OnSubmitted = func(_ string) { r.onRecv() }

	codeButton := &widget.Button{Text: "Receive", Icon: theme.DownloadIcon(), OnTapped: r.onRecv}

	r.data = &bridge.RecvData{Client: r.client, Window: r.window}

	box := container.NewVBox(&widget.Separator{}, container.NewGridWithColumns(2, r.codeEntry, codeButton), &widget.Separator{})
	return container.NewBorder(box, nil, nil, nil, bridge.NewRecvList(r.data))
}

func (r *recv) tabItem() *container.TabItem {
	return &container.TabItem{Text: "Receive", Icon: theme.DownloadIcon(), Content: r.buildUI()}
}

type completionEntry struct {
	widget.Entry
	canvas    fyne.Canvas
	backwards bool
	complete  *completion.TabCompleter
}

// AcceptsTab overrides tab handling to allow tabs as input.
func (c *completionEntry) AcceptsTab() bool {
	return true
}

// TypedKey adapts the key inputs to handle tab completion.
func (c *completionEntry) TypedKey(key *fyne.KeyEvent) {
	switch key.Name {
	case desktop.KeyShiftLeft, desktop.KeyShiftRight:
	case fyne.KeyTab:
		if c.backwards {
			c.previous()
			return
		}

		c.next()
	case fyne.KeyEscape:
		c.canvas.Unfocus()
	default:
		c.complete.Reset()
		c.Entry.TypedKey(key)
	}
}

// KeyDown handles pressing of shift for going backwards in completion.
func (c *completionEntry) KeyDown(key *fyne.KeyEvent) {
	if key.Name == desktop.KeyShiftLeft || key.Name == desktop.KeyShiftRight {
		c.backwards = true
	}

	c.Entry.KeyDown(key)
}

// KeyUp handles releasing of shift for going backwards in completion.
func (c *completionEntry) KeyUp(key *fyne.KeyEvent) {
	if key.Name == desktop.KeyShiftLeft || key.Name == desktop.KeyShiftRight {
		c.backwards = false
	}

	c.Entry.KeyDown(key)
}

func (c *completionEntry) previous() {
	previous := c.complete.Previous(c.Text)
	c.CursorColumn = len(previous)
	c.SetText(previous)
}

func (c *completionEntry) next() {
	next := c.complete.Next(c.Text)
	c.CursorColumn = len(next)
	c.SetText(next)
}

func newCompletionEntry(client *transport.Client, canvas fyne.Canvas) *completionEntry {
	entry := &completionEntry{
		complete: &completion.TabCompleter{Generate: client.CompleteRecvCode},
		canvas:   canvas,
		Entry: widget.Entry{
			PlaceHolder: "Enter code", Wrapping: fyne.TextTruncate, Validator: util.CodeValidator,
		},
	}
	entry.ExtendBaseWidget(entry)
	return entry
}
