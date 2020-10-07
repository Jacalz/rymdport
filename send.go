package main

import (
	"fyne.io/fyne"
	"fyne.io/fyne/container"
	"fyne.io/fyne/dialog"
	"fyne.io/fyne/layout"
	"fyne.io/fyne/storage"
	"fyne.io/fyne/theme"
	"fyne.io/fyne/widget"

	"github.com/Jacalz/wormhole-gui/bridge/widgets"
)

func (ad *appData) sendTab() *container.TabItem {
	fileChoice := widget.NewButtonWithIcon("File", theme.FileIcon(), nil)
	textChoice := widget.NewButtonWithIcon("Text", theme.DocumentCreateIcon(), nil)
	// TODO: Add support for sending directories when fyne supports it in the file picker.

	choiceContent := fyne.NewContainerWithLayout(layout.NewGridLayout(1), fileChoice, textChoice)

	contentPicker := dialog.NewCustom("Pick a content type", "Cancel", choiceContent, ad.Window)
	contentPicker.Hide()

	sendList := widgets.NewProgressList()

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

				code := sendList.NewSendItem(file.URI())

				go func() {
					err = ad.Bridge.SendFile(file, code, sendList.Items[sendList.Length()-1].Progress.Update)
					if err != nil {
						dialog.ShowError(err, ad.Window)
					} else if ad.Notifications {
						ad.App.SendNotification(fyne.NewNotification("Send completed", "The file was sent successfully"))
					}
				}()
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

			code := sendList.NewSendItem(storage.NewURI("Text Snippet"))

			go func() {
				err := ad.Bridge.SendText(t, code)
				if err != nil {
					dialog.ShowError(err, ad.Window)
				} else if ad.Notifications {
					ad.App.SendNotification(fyne.NewNotification("Send completed", "The sending of text completed successfully"))
				} else {
					sendList.Items[sendList.Length()-1].Progress.SetValue(1)
				}
			}()
		}()
	}

	send := widget.NewButtonWithIcon("Add content to send", theme.ContentAddIcon(), func() {
		contentPicker.Show()
	})

	box := widget.NewVBox(send, widget.NewLabel(""))

	sendContent := fyne.NewContainerWithLayout(layout.NewBorderLayout(box, nil, nil, nil), box, sendList)

	return widget.NewTabItemWithIcon("Send", theme.MailSendIcon(), sendContent)
}
