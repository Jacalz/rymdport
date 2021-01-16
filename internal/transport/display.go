package transport

import (
	"fyne.io/fyne"
	"fyne.io/fyne/container"
	"fyne.io/fyne/dialog"
	"fyne.io/fyne/theme"
	"fyne.io/fyne/widget"
)

type textDisplay struct {
	textEntry               *widget.Entry
	leftButton, rightButton *widget.Button
	window                  fyne.Window
}

func (d *textDisplay) Refresh() {
	d.textEntry.Refresh()
	d.leftButton.Refresh()
	d.rightButton.Refresh()
}

func createTextWindow() *textDisplay {
	display := &textDisplay{
		window:      fyne.CurrentApp().NewWindow(""),
		textEntry:   &widget.Entry{MultiLine: true},
		leftButton:  &widget.Button{},
		rightButton: &widget.Button{},
	}

	actionContainer := container.NewGridWithColumns(2, display.leftButton, display.rightButton)
	display.window.SetContent(container.NewBorder(nil, actionContainer, nil, nil, display.textEntry))
	display.window.Resize(fyne.NewSize(400, 300))

	return display
}

// showTextReceiveWindow handles the creation of a window for displaying text content.
func (c *Client) showTextReceiveWindow(text string) {
	d := c.display

	d.window.SetTitle("Received Text")
	d.window.SetCloseIntercept(func() {
		d.window.Hide()
	})

	d.leftButton.Text = "Copy"
	d.leftButton.Icon = theme.ContentCopyIcon()
	d.leftButton.OnTapped = func() {
		d.window.Clipboard().SetContent(text)
	}

	d.rightButton.Text = "Save"
	d.rightButton.Icon = theme.DocumentSaveIcon()
	d.rightButton.Importance = widget.MediumImportance
	d.rightButton.OnTapped = func() {
		go func() {
			save := dialog.NewFileSave(func(file fyne.URIWriteCloser, err error) { // TODO: Might want to save this instead of recreating each time
				if err != nil {
					fyne.LogError("Error on selecting file to write to", err)
					dialog.ShowError(err, d.window)
					return
				} else if file == nil {
					return
				}

				if _, err := file.Write([]byte(text)); err != nil {
					fyne.LogError("Error on writing data to the file", err)
					dialog.ShowError(err, d.window)
				}

				if err := file.Close(); err != nil {
					fyne.LogError("Error on writing data to the file", err)
					dialog.ShowError(err, d.window)
				}
			}, d.window)
			save.SetFileName("received.txt")
			save.Show()
		}()
	}

	d.textEntry.OnSubmitted = nil
	d.textEntry.Text = text

	d.Refresh() // Update all contents in the window
	d.window.Show()
}

// ShowTextSendWindow opens a new window for setting up text to send.
func (c *Client) ShowTextSendWindow() chan string {
	text := make(chan string)
	d := c.display

	onClose := func() {
		text <- ""
		d.window.Hide()
	}

	d.window.SetTitle("Send text")
	d.window.SetCloseIntercept(onClose)

	d.leftButton.Text = "Cancel"
	d.leftButton.Icon = theme.CancelIcon()
	d.leftButton.OnTapped = func() {
		onClose()
	}

	d.rightButton.Text = "Send"
	d.rightButton.Icon = theme.MailSendIcon()
	d.rightButton.Importance = widget.HighImportance
	d.rightButton.OnTapped = func() {
		text <- d.textEntry.Text
		d.window.Hide()
	}

	d.textEntry.OnSubmitted = func(_ string) { d.rightButton.OnTapped() }
	d.textEntry.Text = ""

	d.Refresh() // Update all contents in the window
	d.window.Canvas().Focus(d.textEntry)
	d.window.Show()

	return text
}
