package main

import (
	"context"
	"os"

	"fyne.io/fyne"
	"fyne.io/fyne/dialog"
	"fyne.io/fyne/layout"
	"fyne.io/fyne/theme"
	"fyne.io/fyne/widget"

	"github.com/psanford/wormhole-william/wormhole"
)

func sendFile(file fyne.URIReadCloser, s settings, code chan string) error {
	c := wormhole.Client{PassPhraseComponentLength: s.PassPhraseComponentLength}

	defer file.Close()

	f, err := os.Open(file.URI().String()[7:]) // Ignore the file:/ prefix by doing [7:]
	if err != nil {
		fyne.LogError("Error on opening file to send", err)
		return err
	}

	defer f.Close()

	codestr, status, err := c.SendFile(context.Background(), file.Name(), f)
	if err != nil {
		fyne.LogError("Error on sending file", err)
		return err
	}

	code <- codestr

	if stat := <-status; stat.Error != nil {
		fyne.LogError("Error on status of share", err)
		return err
	} else if stat.OK {
		return nil
	}

	return nil
}

func sendText(text string, s settings, code chan string) error {
	c := wormhole.Client{PassPhraseComponentLength: s.PassPhraseComponentLength}

	codestr, status, err := c.SendText(context.Background(), text)
	if err != nil {
		fyne.LogError("Error on sending text", err)
		return err
	}

	code <- codestr

	if stat := <-status; stat.Error != nil {
		fyne.LogError("Error on status of share", err)
		return err
	} else if stat.OK {
		return nil
	}

	return nil
}

func (s *settings) sendTab(w fyne.Window) *widget.TabItem {
	fileChoice := widget.NewButtonWithIcon("File", theme.FileIcon(), nil)
	textChoice := widget.NewButtonWithIcon("Text", theme.DocumentCreateIcon(), nil)
	// TODO: Add support for sending directories when fyne supports it in the file picker.

	choiceContent := widget.NewVBox(fileChoice, textChoice)

	contentPicker := dialog.NewCustom("Pick a content type", "Cancel", choiceContent, w)
	contentPicker.Hide()

	sendGrid := fyne.NewContainerWithLayout(layout.NewGridLayout(2), widget.NewLabel("Filename"), widget.NewLabel("Code"))

	fileChoice.OnTapped = func() {
		contentPicker.Hide()

		dialog.ShowFileOpen(func(file fyne.URIReadCloser, err error) {
			if err != nil {
				fyne.LogError("Error on picking the file", err)
				return
			}

			code := make(chan string)

			go func() {
				err = sendFile(file, *s, code)
				if err != nil {
					dialog.ShowError(err, w)
				}
			}()

			sendGrid.AddObject(widget.NewLabel(file.Name()))
			codeLabel := widget.NewLabel("Waiting for code")

			sendGrid.AddObject(codeLabel)
			go codeLabel.SetText(<-code)
		}, w)
	}

	textEntry := widget.NewMultiLineEntry()
	textEntry.SetPlaceHolder("Enter text to send")

	textChoice.OnTapped = func() {
		contentPicker.Hide()
		dialog.ShowCustomConfirm("Text to send", "Send", "Cancel", textEntry, func(send bool) {
			if send {
				code := make(chan string)

				go func() {
					err := sendText(textEntry.Text, *s, code)
					if err != nil {
						dialog.ShowError(err, w)
					}
				}()

				sendGrid.AddObject(widget.NewLabel("Text Snippet"))
				codeLabel := widget.NewLabel("Waiting for code")
				sendGrid.AddObject(codeLabel)
				go codeLabel.SetText(<-code)
			} else {
				textEntry.SetText("")
			}

		}, w)
	}

	send := widget.NewButtonWithIcon("Add content to send", theme.ContentAddIcon(), func() {
		contentPicker.Show()
	})

	sendContent := fyne.NewContainerWithLayout(layout.NewVBoxLayout(), send, widget.NewLabel(""), sendGrid)

	return widget.NewTabItemWithIcon("Send", theme.MailSendIcon(), widget.NewScrollContainer(sendContent))
}
