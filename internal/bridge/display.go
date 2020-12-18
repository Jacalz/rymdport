package bridge

import (
	"fyne.io/fyne"
	"fyne.io/fyne/container"
	"fyne.io/fyne/dialog"
	"fyne.io/fyne/theme"
	"fyne.io/fyne/widget"
)

// displayReceivedText handles the creation of a window for displaying text content.
func (b *Bridge) displayReceivedText(content []byte) {
	b.window.SetTitle("Received Text")
	b.window.SetCloseIntercept(func() {
		b.window.Hide()
	})

	textEntry := &widget.Entry{MultiLine: true, Text: string(content)}

	copyText := widget.NewButtonWithIcon("Copy text", theme.ContentCopyIcon(), func() {
		fyne.CurrentApp().Driver().AllWindows()[0].Clipboard().SetContent(string(content))
	})

	saveFile := widget.NewButtonWithIcon("Save text to file", theme.MoveDownIcon(), func() {
		go func() {
			dialog.ShowFileSave(func(file fyne.URIWriteCloser, err error) {
				if err != nil {
					fyne.LogError("Error on slecting file to write to", err)
					dialog.ShowError(err, b.window)
					return
				} else if file == nil {
					return
				}

				if _, err := file.Write(content); err != nil {
					fyne.LogError("Error on writing data to the file", err)
					dialog.ShowError(err, b.window)
				}

				if err := file.Close(); err != nil {
					fyne.LogError("Error on writing data to the file", err)
					dialog.ShowError(err, b.window)
				}
			}, b.window)
		}()
	})

	textContainer := container.NewScroll(textEntry)
	actionContainer := container.NewGridWithColumns(2, copyText, saveFile)

	b.window.SetContent(container.NewBorder(nil, actionContainer, nil, nil, textContainer))
	b.window.Show()
}

// EnterSendText opens a new window for setting up text to send.
func (b *Bridge) EnterSendText() chan string {
	text := make(chan string)
	textEntry := &widget.Entry{MultiLine: true, PlaceHolder: "Enter text to send..."}

	b.window.SetTitle("Send text")
	b.window.SetCloseIntercept(func() {
		text <- ""
		b.window.Hide()
	})

	cancel := widget.NewButtonWithIcon("Cancel", theme.CancelIcon(), func() {
		text <- ""
		b.window.Hide()
	})

	send := &widget.Button{Text: "Send", Icon: theme.MailSendIcon(), OnTapped: func() {
		text <- textEntry.Text
		b.window.Hide()
	}, Importance: widget.HighImportance}

	textContainer := container.NewScroll(textEntry)
	actionContainer := container.NewGridWithColumns(2, cancel, send)

	b.window.SetContent(container.NewBorder(nil, actionContainer, nil, nil, textContainer))
	b.window.Canvas().Focus(textEntry)
	b.window.Show()

	return text
}
