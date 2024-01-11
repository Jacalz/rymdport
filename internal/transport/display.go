package transport

import (
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/Jacalz/rymdport/v3/internal/util"
)

type textRecvWindow struct {
	textEntry              *widget.Entry
	copyButton, saveButton *widget.Button
	window                 fyne.Window
	received               []byte
}

func (r *textRecvWindow) copy() {
	r.window.Clipboard().SetContent(string(r.received))
}

func (r *textRecvWindow) interceptClose() {
	r.window.Hide()
	r.textEntry.SetText("")
}

func (r *textRecvWindow) save() {
	save := dialog.NewFileSave(func(file fyne.URIWriteCloser, err error) { // TODO: Might want to save this instead of recreating each time
		if err != nil {
			fyne.LogError("Error on selecting file to write to", err)
			dialog.ShowError(err, r.window)
			return
		} else if file == nil {
			return
		}

		if _, err := file.Write(r.received); err != nil {
			fyne.LogError("Error on writing text to the file", err)
			dialog.ShowError(err, r.window)
		}

		if err := file.Close(); err != nil {
			fyne.LogError("Error on closing text file", err)
			dialog.ShowError(err, r.window)
		}
	}, r.window)
	now := time.Now().Format("2006-01-02T15:04") // TODO: Might want to use AppendFormat and strings.Builder
	save.SetFileName("received-" + now + ".txt")
	save.Resize(util.WindowSizeToDialog(r.window.Canvas().Size()))
	save.Show()
}

func (c *Client) createTextRecvWindow() {
	window := c.App.NewWindow("Received Text")
	window.SetCloseIntercept(c.textRecvWindow.interceptClose)

	c.textRecvWindow = textRecvWindow{
		window:     window,
		textEntry:  &widget.Entry{MultiLine: true, Wrapping: fyne.TextWrapWord},
		copyButton: &widget.Button{Text: "Copy", Icon: theme.ContentCopyIcon(), OnTapped: c.textRecvWindow.copy},
		saveButton: &widget.Button{Text: "Save", Icon: theme.DocumentSaveIcon(), OnTapped: c.textRecvWindow.save},
	}

	actionContainer := container.NewGridWithColumns(2, c.textRecvWindow.copyButton, c.textRecvWindow.saveButton)
	window.SetContent(container.NewBorder(nil, actionContainer, nil, nil, c.textRecvWindow.textEntry))
	window.Resize(fyne.NewSize(400, 300))
}

// showTextReceiveWindow handles the creation of a window for displaying text content.
func (c *Client) showTextReceiveWindow(received []byte) {
	if c.textRecvWindow.window == nil {
		c.createTextRecvWindow()
	} // else: Might want to request window focus?

	c.textRecvWindow.received = received
	c.textRecvWindow.textEntry.SetText(string(received))

	win := c.textRecvWindow.window
	win.Show()
	win.RequestFocus()
	win.Canvas().Focus(c.textRecvWindow.textEntry)
}
