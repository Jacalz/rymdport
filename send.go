package main

import (
	"fyne.io/fyne"
	"fyne.io/fyne/dialog"
	"fyne.io/fyne/layout"
	"fyne.io/fyne/theme"
	"fyne.io/fyne/widget"
)

func sendTab(w fyne.Window) *widget.TabItem {
	fileChoice := widget.NewButtonWithIcon("File", theme.FileIcon(), nil)
	textChoice := widget.NewButtonWithIcon("Text", theme.DocumentCreateIcon(), nil)
	// TODO: Add support for sending directories when fyne supports it in the file picker.

	choiceContent := widget.NewVBox(fileChoice, textChoice)

	contentPicker := dialog.NewCustom("Pick a content type", "Cancel", choiceContent, w)
	contentPicker.Hide()

	fileChoice.OnTapped = func() {
		contentPicker.Hide()

		dialog.ShowFileOpen(func(file fyne.URIReadCloser, err error) {

		}, w)
	}

	textEntry := widget.NewMultiLineEntry()
	textEntry.SetPlaceHolder("Enter text to send")

	textChoice.OnTapped = func() {
		contentPicker.Hide()
		dialog.ShowCustomConfirm("text to send", "Send", "Cancel", textEntry, nil, w)
	}

	send := widget.NewButtonWithIcon("Add content to send", theme.ContentAddIcon(), func() {
		contentPicker.Show()
	})

	sendContent := fyne.NewContainerWithLayout(layout.NewVBoxLayout(), send)

	return widget.NewTabItemWithIcon("Send", theme.MailSendIcon(), widget.NewScrollContainer(sendContent))
}
