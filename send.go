package main

import (
	"fyne.io/fyne"
	"fyne.io/fyne/dialog"
	"fyne.io/fyne/layout"
	"fyne.io/fyne/theme"
	"fyne.io/fyne/widget"

	"github.com/Jacalz/wormhole-gui/bridge"
	"github.com/Jacalz/wormhole-gui/widgets"
)

func (ad *appData) sendTab() *widget.TabItem {
	fileChoice := widget.NewButtonWithIcon("File", theme.FileIcon(), nil)
	textChoice := widget.NewButtonWithIcon("Text", theme.DocumentCreateIcon(), nil)
	// TODO: Add support for sending directories when fyne supports it in the file picker.

	choiceContent := widget.NewVBox(fileChoice, textChoice)

	contentPicker := dialog.NewCustom("Pick a content type", "Cancel", choiceContent, ad.Window)
	contentPicker.Hide()

	sendGrid := fyne.NewContainerWithLayout(layout.NewGridLayout(3), widgets.NewBoldLabel("Filename"), widgets.NewBoldLabel("Code"), widgets.NewBoldLabel("Progress"))

	fileChoice.OnTapped = func() {
		go func() {
			contentPicker.Hide()

			dialog.ShowFileOpen(func(file fyne.URIReadCloser, err error) {
				if err != nil {
					fyne.LogError("Error on selecting file to send", err)
					dialog.ShowError(err, ad.Window)
					return
				} else if file == nil {
					return
				}

				code := make(chan string)

				progress := bridge.NewSendProgress()

				go func() {
					err = ad.Bridge.SendFile(file, code, progress.Update)
					if err != nil {
						dialog.ShowError(err, ad.Window)
					} else if ad.Notifications {
						ad.App.SendNotification(fyne.NewNotification("Send completed", "The sending of a file completed successfully"))
					}
				}()

				codeLabel := widgets.NewCodeLabel(code)
				sendGrid.AddObject(widget.NewLabel(file.Name()))
				sendGrid.AddObject(codeLabel)
				sendGrid.AddObject(progress.Widget)
			}, ad.Window)
		}()
	}

	textChoice.OnTapped = func() {
		go func() {
			contentPicker.Hide()
			text := make(chan string)

			ad.Bridge.EnterSendText(ad.App, text)
			t := <-text

			if t == "" {
				return
			}

			progress := bridge.NewSendProgress()

			code := make(chan string)
			go func() {
				err := ad.Bridge.SendText(t, code, progress.Update)
				if err != nil {
					dialog.ShowError(err, ad.Window)
				} else if ad.Notifications {
					ad.App.SendNotification(fyne.NewNotification("Send completed", "The sending of text completed successfully"))
				}

			}()

			codeLabel := widgets.NewCodeLabel(code)
			sendGrid.AddObject(widget.NewLabel("Text Snippet"))
			sendGrid.AddObject(codeLabel)
			sendGrid.AddObject(progress.Widget)
		}()
	}

	send := widget.NewButtonWithIcon("Add content to send", theme.ContentAddIcon(), func() {
		contentPicker.Show()
	})

	sendContent := fyne.NewContainerWithLayout(layout.NewVBoxLayout(), send, widget.NewLabel(""), sendGrid)

	return widget.NewTabItemWithIcon("Send", theme.MailSendIcon(), widget.NewScrollContainer(sendContent))
}
