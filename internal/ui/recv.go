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
	data      bridge.RecvData

	window fyne.Window
}

func newRecvTab(w fyne.Window, c *transport.Client) *container.TabItem {
	recv := &recv{window: w}

	return &container.TabItem{
		Text:    "Receive",
		Icon:    theme.DownloadIcon(),
		Content: recv.buildUI(c),
	}
}

func (r *recv) buildUI(client *transport.Client) *fyne.Container {
	r.codeEntry = newCompletionEntry(client, r.window.Canvas(), client.App)
	r.codeEntry.OnSubmitted = func(_ string) { r.onRecv() }

	codeButton := &widget.Button{Text: "Receive", Icon: theme.DownloadIcon(), OnTapped: r.onRecv}

	r.data = bridge.RecvData{Client: client, Window: r.window}

	box := container.NewVBox(&widget.Separator{}, container.NewGridWithColumns(2, r.codeEntry, codeButton), &widget.Separator{})
	return container.NewBorder(box, nil, nil, nil, r.data.NewRecvList())
}

func (r *recv) onRecv() {
	if err := r.codeEntry.Validate(); err != nil || r.codeEntry.Text == "" {
		dialog.ShowInformation("Invalid code", "The code is invalid. Please try again.", r.window)
		return
	}

	r.data.NewReceive(r.codeEntry.Text)
	r.codeEntry.SetText("")
}

type completionEntry struct {
	widget.Entry
	driver   desktop.Driver
	canvas   fyne.Canvas
	complete completion.TabCompleter
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
		if c.driver.CurrentKeyModifiers()&fyne.KeyModifierShift != 0 {
			c.setCompletion(c.complete.Previous)
			return
		}

		c.setCompletion(c.complete.Next)
	case fyne.KeyEscape:
		c.canvas.Unfocus()
	default:
		c.complete.Reset()
		c.Entry.TypedKey(key)
	}
}

func (c *completionEntry) setCompletion(lookup func(string) string) {
	text := lookup(c.Text)
	c.CursorColumn = len(text)
	c.SetText(text)
}

func newCompletionEntry(client *transport.Client, canvas fyne.Canvas, app fyne.App) *completionEntry {
	entry := &completionEntry{
		canvas: canvas,
		driver: app.Driver().(desktop.Driver),
		Entry: widget.Entry{
			PlaceHolder: "Enter code", Scroll: container.ScrollHorizontalOnly, Validator: util.CodeValidator,
		},
	}
	entry.complete.Generate = client.GenerateCodeCompletion
	entry.ExtendBaseWidget(entry)
	return entry
}
