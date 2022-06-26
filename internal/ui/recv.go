package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/Jacalz/rymdport/v3/internal/completion"
	"github.com/Jacalz/rymdport/v3/internal/transport"
	"github.com/Jacalz/rymdport/v3/internal/transport/bridge"
	"github.com/Jacalz/rymdport/v3/internal/util"
)

type recv struct {
	codeEntry  *completionEntry
	codeButton *widget.Button

	recvList *bridge.RecvList

	client *transport.Client
	window fyne.Window
	app    fyne.App
}

func newRecv(a fyne.App, w fyne.Window, c *transport.Client) *recv {
	return &recv{app: a, window: w, client: c}
}

func (r *recv) onRecv() {
	if err := r.codeEntry.Validate(); err != nil || r.codeEntry.Text == "" {
		dialog.ShowInformation("Invalid code", "The code is invalid. Please try again.", r.window)
		return
	}

	r.recvList.NewReceive(r.codeEntry.Text)
	r.codeEntry.SetText("")
}

func (r *recv) buildUI() *fyne.Container {
	r.codeEntry = newCompletionEntry(r.client, r.window.Canvas())
	r.codeEntry.OnSubmitted = func(_ string) { r.onRecv() }

	r.codeButton = &widget.Button{Text: "Receive", Icon: theme.DownloadIcon(), OnTapped: r.onRecv}

	r.recvList = bridge.NewRecvList(r.window, r.client)

	box := container.NewVBox(container.NewGridWithColumns(2, r.codeEntry, r.codeButton), &widget.Label{})
	return container.NewBorder(box, nil, nil, nil, r.recvList)
}

func (r *recv) tabItem() *container.TabItem {
	return &container.TabItem{Text: "Receive", Icon: theme.DownloadIcon(), Content: r.buildUI()}
}

type completionEntry struct {
	widget.Entry
	canvas   fyne.Canvas
	complete *completion.TabCompleter
}

// AcceptsTab overrides tab handling to allow tabs as input.
func (c *completionEntry) AcceptsTab() bool {
	return true
}

// TypedKey adapts the key inputs to handle tab completion.
func (c *completionEntry) TypedKey(key *fyne.KeyEvent) {
	switch key.Name {
	case fyne.KeyTab:
		c.next()
	case fyne.KeyEscape:
		c.canvas.Unfocus()
	default:
		c.complete.Reset()
		c.Entry.TypedKey(key)
	}
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
