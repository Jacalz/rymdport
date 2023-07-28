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
}

func createTextRecvWindow(app fyne.App) *textRecvWindow {
	display := &textRecvWindow{
		window:     app.NewWindow("Received Text"),
		textEntry:  &widget.Entry{MultiLine: true, Wrapping: fyne.TextWrapWord},
		copyButton: &widget.Button{Text: "Copy", Icon: theme.ContentCopyIcon()},
		saveButton: &widget.Button{Text: "Save", Icon: theme.DocumentSaveIcon()},
	}

	actionContainer := container.NewGridWithColumns(2, display.copyButton, display.saveButton)
	display.window.SetContent(container.NewBorder(nil, actionContainer, nil, nil, display.textEntry))
	display.window.Resize(fyne.NewSize(400, 300))

	return display
}

// showTextReceiveWindow handles the creation of a window for displaying text content.
func (c *Client) showTextReceiveWindow(received []byte) {
	if c.textRecvWindow == nil {
		c.textRecvWindow = createTextRecvWindow(c.app)
	}

	d := c.textRecvWindow
	text := string(received)

	d.window.SetCloseIntercept(func() {
		d.window.Hide()
		d.textEntry.SetText("")
	})

	d.copyButton.OnTapped = func() {
		d.window.Clipboard().SetContent(text)
	}

	d.saveButton.OnTapped = func() {
		go func() {
			save := dialog.NewFileSave(func(file fyne.URIWriteCloser, err error) { // TODO: Might want to save this instead of recreating each time
				if err != nil {
					fyne.LogError("Error on selecting file to write to", err)
					dialog.ShowError(err, d.window)
					return
				} else if file == nil {
					return
				}

				if _, err := file.Write(received); err != nil {
					fyne.LogError("Error on writing text to the file", err)
					dialog.ShowError(err, d.window)
				}

				if err := file.Close(); err != nil {
					fyne.LogError("Error on closing text file", err)
					dialog.ShowError(err, d.window)
				}
			}, d.window)
			now := time.Now().Format("2006-01-02T15:04")
			save.SetFileName("received-" + now + ".txt")
			save.Resize(util.WindowSizeToDialog(d.window.Canvas().Size()))
			save.Show()
		}()
	}

	d.textEntry.SetText(text)
	d.window.Show()
	d.window.RequestFocus()
	d.window.Canvas().Focus(d.textEntry)
}

type textSendWindow struct {
	textEntry                *widget.Entry
	cancelButton, sendButton *widget.Button
	window                   fyne.Window
}

func createTextSendWindow(app fyne.App) *textSendWindow {
	display := &textSendWindow{
		window:       app.NewWindow("Send Text"),
		textEntry:    &widget.Entry{MultiLine: true, Wrapping: fyne.TextWrapWord},
		cancelButton: &widget.Button{Text: "Cancel", Icon: theme.CancelIcon()},
		sendButton:   &widget.Button{Text: "Send", Icon: theme.MailSendIcon(), Importance: widget.HighImportance},
	}
	display.textEntry.OnSubmitted = func(_ string) {
		display.sendButton.OnTapped()
	}

	actionContainer := container.NewGridWithColumns(2, display.cancelButton, display.sendButton)
	display.window.SetContent(container.NewBorder(nil, actionContainer, nil, nil, display.textEntry))
	display.window.Resize(fyne.NewSize(400, 300))

	return display
}

// ShowTextSendWindow opens a new window for setting up text to send.
func (c *Client) ShowTextSendWindow() chan string {
	if c.textSendWindow == nil {
		c.textSendWindow = createTextSendWindow(c.app)
	}

	d := c.textSendWindow

	text := make(chan string)
	onClose := func() {
		text <- ""
		d.window.Hide()
		d.textEntry.SetText("")
		// don't force focus the main window - we didn't send anything so there's no new code to see
	}

	d.window.SetCloseIntercept(onClose)
	d.cancelButton.OnTapped = onClose
	d.sendButton.OnTapped = func() {
		text <- d.textEntry.Text
		d.window.Hide()
		d.textEntry.SetText("")
		c.mainWindow.RequestFocus() // So the user can see the code for the sent text
	}

	d.window.Show()
	d.window.RequestFocus()
	d.window.Canvas().Focus(d.textEntry)

	return text
}
