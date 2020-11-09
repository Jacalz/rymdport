package bridge

import (
	"fyne.io/fyne"
	"fyne.io/fyne/container"
	"fyne.io/fyne/dialog"
	"fyne.io/fyne/layout"
	"fyne.io/fyne/theme"
	"fyne.io/fyne/widget"
)

// displayRecievedText handles the creation of a window for displaying text content.
func displayRecievedText(a fyne.App, content string) {
	w := a.NewWindow("Received text")

	textEntry := &widget.Entry{MultiLine: true, Text: content}

	copyText := widget.NewButtonWithIcon("Copy text", theme.ContentCopyIcon(), func() {
		fyne.CurrentApp().Driver().AllWindows()[0].Clipboard().SetContent(content)
	})

	saveFile := widget.NewButtonWithIcon("Save text to file", theme.MoveDownIcon(), func() {
		go func() {
			dialog.ShowFileSave(func(file fyne.URIWriteCloser, err error) {
				if err != nil {
					fyne.LogError("Error on slecting file to write to", err)
					dialog.ShowError(err, w)
					return
				} else if file == nil {
					return
				}

				err = writeFile(file, []byte(content))
				if err != nil {
					dialog.ShowError(err, w)
					return
				}
			}, w)
		}()
	})

	textContainer := container.NewScroll(textEntry)
	actionContainer := fyne.NewContainerWithLayout(layout.NewGridLayout(2), copyText, saveFile)

	w.Resize(fyne.NewSize(400, 300))
	w.SetContent(fyne.NewContainerWithLayout(layout.NewBorderLayout(nil, actionContainer, nil, nil), actionContainer, textContainer))
	w.Show()
}

// EnterSendText opens a new window for setting up text to send.
func EnterSendText() chan string {
	w := fyne.CurrentApp().NewWindow("Send text")

	text := make(chan string)
	textEntry := widget.NewMultiLineEntry()

	cancel := widget.NewButtonWithIcon("Cancel", theme.CancelIcon(), func() {
		text <- ""
		w.Close()
	})

	send := &widget.Button{Text: "Send", Icon: theme.MailSendIcon(), OnTapped: func() {
		text <- textEntry.Text
		w.Close()
	}, Importance: widget.HighImportance}

	textContainer := container.NewScroll(textEntry)
	actionContainer := fyne.NewContainerWithLayout(layout.NewGridLayout(2), cancel, send)

	w.Resize(fyne.NewSize(400, 300))
	w.SetContent(container.NewBorder(nil, actionContainer, nil, nil, textContainer))
	w.Canvas().Focus(textEntry)
	w.Show()

	return text
}
