package transport

import (
	"fyne.io/fyne"
	"fyne.io/fyne/container"
	"fyne.io/fyne/dialog"
	"fyne.io/fyne/theme"
	"fyne.io/fyne/widget"
)

// displayReceivedText handles the creation of a window for displaying text content.
func (c *Client) displayReceivedText(content []byte) {
	c.window.SetTitle("Received Text")
	c.window.SetCloseIntercept(func() {
		c.window.Hide()
	})

	textEntry := &widget.Entry{MultiLine: true, Text: string(content)}

	copyText := widget.NewButtonWithIcon("Copy text", theme.ContentCopyIcon(), func() {
		c.window.Clipboard().SetContent(string(content))
	})

	saveFile := widget.NewButtonWithIcon("Save text to file", theme.MoveDownIcon(), func() {
		go func() {
			dialog.ShowFileSave(func(file fyne.URIWriteCloser, err error) {
				if err != nil {
					fyne.LogError("Error on slecting file to write to", err)
					dialog.ShowError(err, c.window)
					return
				} else if file == nil {
					return
				}

				if _, err := file.Write(content); err != nil {
					fyne.LogError("Error on writing data to the file", err)
					dialog.ShowError(err, c.window)
				}

				if err := file.Close(); err != nil {
					fyne.LogError("Error on writing data to the file", err)
					dialog.ShowError(err, c.window)
				}
			}, c.window)
		}()
	})

	textContainer := container.NewScroll(textEntry)
	actionContainer := container.NewGridWithColumns(2, copyText, saveFile)

	c.window.SetContent(container.NewBorder(nil, actionContainer, nil, nil, textContainer))
	c.window.Show()
}

// EnterSendText opens a new window for setting up text to send.
func (c *Client) EnterSendText() chan string {
	text := make(chan string)
	textEntry := &widget.Entry{MultiLine: true, PlaceHolder: "Enter text to send..."}

	c.window.SetTitle("Send text")
	c.window.SetCloseIntercept(func() {
		text <- ""
		c.window.Hide()
	})

	cancel := widget.NewButtonWithIcon("Cancel", theme.CancelIcon(), func() {
		text <- ""
		c.window.Hide()
	})

	send := &widget.Button{Text: "Send", Icon: theme.MailSendIcon(), OnTapped: func() {
		text <- textEntry.Text
		c.window.Hide()
	}, Importance: widget.HighImportance}

	textContainer := container.NewScroll(textEntry)
	actionContainer := container.NewGridWithColumns(2, cancel, send)

	c.window.SetContent(container.NewBorder(nil, actionContainer, nil, nil, textContainer))
	c.window.Canvas().Focus(textEntry)
	c.window.Show()

	return text
}
