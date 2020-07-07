package bridge

import (
	"fyne.io/fyne"
	"fyne.io/fyne/dialog"
	"fyne.io/fyne/layout"
	"fyne.io/fyne/theme"
	"fyne.io/fyne/widget"
)

// displayRecievedText handles the creation of a window for displaying text content.
func displayRecievedText(a fyne.App, content string) {
	w := a.NewWindow("Recieved text")

	textEntry := widget.NewMultiLineEntry()
	textEntry.SetText(content)

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

	textContainer := widget.NewScrollContainer(textEntry)
	textContainer.SetMinSize(fyne.NewSize(400, 400))

	actionContainer := fyne.NewContainerWithLayout(layout.NewGridLayout(2), copyText, saveFile)

	w.SetContent(fyne.NewContainerWithLayout(layout.NewVBoxLayout(), textContainer, actionContainer))
	w.Show()
}

// EnterSendText opens a new window for setting up text to send.
func (b *Bridge) EnterSendText(a fyne.App, text chan string) {
	w := a.NewWindow("Enter text to send")

	textEntry := widget.NewMultiLineEntry()
	textContainer := widget.NewScrollContainer(textEntry)
	textContainer.SetMinSize(fyne.NewSize(400, 400))

	w.Canvas().Focus(textEntry)

	cancel := widget.NewButtonWithIcon("Cancel", theme.CancelIcon(), func() {
		textEntry.Text = ""
		w.Close()
	})

	send := widget.NewButtonWithIcon("Send", theme.MailSendIcon(), func() {
		w.Close()
	})

	w.SetOnClosed(func() {
		text <- textEntry.Text
	})

	actionContainer := fyne.NewContainerWithLayout(layout.NewGridLayout(2), cancel, send)

	w.SetContent(fyne.NewContainerWithLayout(layout.NewVBoxLayout(), textContainer, actionContainer))
	w.Show()
}
