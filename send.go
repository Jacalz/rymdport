package main

import (
	"context"
	"os"
	"sync"

	"fyne.io/fyne"
	"fyne.io/fyne/dialog"
	"fyne.io/fyne/layout"
	"fyne.io/fyne/theme"
	"fyne.io/fyne/widget"

	"github.com/Jacalz/wormhole-gui/widgets"
	"github.com/psanford/wormhole-william/wormhole"
)

func sendFile(file fyne.URIReadCloser, s settings, code chan string, progress wormhole.SendOption) error {
	c := wormhole.Client{PassPhraseComponentLength: s.ComponentLength}

	defer file.Close()

	f, err := os.Open(file.URI().String()[7:]) // Ignore the file:/ prefix by doing [7:]
	if err != nil {
		fyne.LogError("Error on opening file to send", err)
		return err
	}

	defer f.Close() // #nosec - We are not writing to the file.

	codestr, status, err := c.SendFile(context.Background(), file.Name(), f, progress)
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

func sendText(text string, s settings, code chan string, progress wormhole.SendOption) error {
	c := wormhole.Client{PassPhraseComponentLength: s.ComponentLength}

	codestr, status, err := c.SendText(context.Background(), text, progress)
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

func (s *settings) sendTab(a fyne.App, w fyne.Window) *widget.TabItem {
	fileChoice := widget.NewButtonWithIcon("File", theme.FileIcon(), nil)
	textChoice := widget.NewButtonWithIcon("Text", theme.DocumentCreateIcon(), nil)
	// TODO: Add support for sending directories when fyne supports it in the file picker.

	choiceContent := widget.NewVBox(fileChoice, textChoice)

	contentPicker := dialog.NewCustom("Pick a content type", "Cancel", choiceContent, w)
	contentPicker.Hide()

	sendGrid := fyne.NewContainerWithLayout(layout.NewGridLayout(3), widgets.NewBoldLabel("Filename"), widgets.NewBoldLabel("Code"), widgets.NewBoldLabel("Progress"))

	fileChoice.OnTapped = func() {
		go func() {
			contentPicker.Hide()

			dialog.ShowFileOpen(func(file fyne.URIReadCloser, err error) {
				if err != nil {
					fyne.LogError("Error on picking the file", err)
					return
				}

				code := make(chan string)

				var once sync.Once
				progress := widget.NewProgressBar()
				update := wormhole.WithProgress(func(sent int64, total int64) {
					once.Do(func() { progress.Max = float64(total) })
					progress.SetValue(float64(sent))
				})

				go func() {
					err = sendFile(file, *s, code, update)
					if err != nil {
						dialog.ShowError(err, w)
					} else if s.Notifications {
						a.SendNotification(fyne.NewNotification("Send completed", "The sending of a file completed successfully"))
					}
				}()

				codeLabel := widgets.NewCodeLabel(code)
				sendGrid.AddObject(widget.NewLabel(file.Name()))
				sendGrid.AddObject(codeLabel)
				sendGrid.AddObject(progress)
			}, w)
		}()
	}

	textEntry := widget.NewMultiLineEntry()
	textEntry.SetPlaceHolder("Enter text to send")

	textChoice.OnTapped = func() {
		go func() {
			contentPicker.Hide()
			textEntry.SetText("")

			dialog.ShowCustomConfirm("Text to send", "Send", "Cancel", textEntry, func(send bool) {
				if send {
					text := textEntry.Text

					var once sync.Once
					progress := widget.NewProgressBar()
					update := wormhole.WithProgress(func(sent int64, total int64) {
						once.Do(func() { progress.Max = float64(total) })
						progress.SetValue(float64(sent))
					})

					code := make(chan string)
					go func() {
						err := sendText(text, *s, code, update)
						if err != nil {
							dialog.ShowError(err, w)
						} else if s.Notifications {
							a.SendNotification(fyne.NewNotification("Send completed", "The sending of text completed successfully"))
						}

					}()

					codeLabel := widgets.NewCodeLabel(code)
					sendGrid.AddObject(widget.NewLabel("Text Snippet"))
					sendGrid.AddObject(codeLabel)
					sendGrid.AddObject(progress)
				}
				textEntry.SetText("")
			}, w)
		}()
	}

	send := widget.NewButtonWithIcon("Add content to send", theme.ContentAddIcon(), func() {
		contentPicker.Show()
	})

	sendContent := fyne.NewContainerWithLayout(layout.NewVBoxLayout(), send, widget.NewLabel(""), sendGrid)

	return widget.NewTabItemWithIcon("Send", theme.MailSendIcon(), widget.NewScrollContainer(sendContent))
}
