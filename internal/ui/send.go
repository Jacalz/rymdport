package ui

import (
	"fyne.io/fyne"
	"fyne.io/fyne/container"
	"fyne.io/fyne/dialog"
	"fyne.io/fyne/storage"
	"fyne.io/fyne/theme"
	"fyne.io/fyne/widget"
	"github.com/Jacalz/wormhole-gui/internal/bridge"
	"github.com/Jacalz/wormhole-gui/internal/bridge/widgets"
)

type send struct {
	contentPicker   dialog.Dialog
	fileChoice      *widget.Button
	directoryChoice *widget.Button
	textChoice      *widget.Button

	contentToSend *widget.Button
	sendList      *widgets.SendList

	bridge      *bridge.Bridge
	appSettings *AppSettings
	window      fyne.Window
	app         fyne.App
}

func newSend(a fyne.App, w fyne.Window, b *bridge.Bridge, as *AppSettings) *send {
	return &send{app: a, window: w, bridge: b, appSettings: as}
}

func (s *send) onFileSend() {
	go func() {
		s.contentPicker.Hide()

		dialog.ShowFileOpen(func(file fyne.URIReadCloser, err error) {
			if err != nil {
				fyne.LogError("Error on selecting file to send", err)
				dialog.ShowError(err, s.window)
				return
			} else if file == nil {
				return
			}

			code := s.sendList.NewSendItem(file.URI())

			go func(i int) {
				err = s.bridge.SendFile(file, code, s.sendList.Items[i].Progress.Update)
				if err != nil {
					s.sendList.RemoveItem(i)
					dialog.ShowError(err, s.window)
				} else if s.appSettings.Notifications {
					s.app.SendNotification(fyne.NewNotification("Send completed", "The file was sent successfully"))
				}
			}(s.sendList.Length() - 1)
		}, s.window)
	}()
}

func (s *send) onDirSend() {
	go func() {
		s.contentPicker.Hide()

		dialog.ShowFolderOpen(func(dir fyne.ListableURI, err error) {
			if err != nil {
				fyne.LogError("Error on selecting directory to send", err)
				dialog.ShowError(err, s.window)
				return
			} else if dir == nil {
				return
			}

			code := s.sendList.NewSendItem(dir)

			go func(i int) {
				if err := s.bridge.SendDir(dir, code, s.sendList.Items[i].Progress.Update); err != nil {
					s.sendList.RemoveItem(i)
					dialog.ShowError(err, s.window)
				} else if s.appSettings.Notifications {
					s.app.SendNotification(fyne.NewNotification("Send completed", "The directory was sent successfully"))
				}
			}(s.sendList.Length() - 1)
		}, s.window)
	}()
}

func (s *send) onTextSend() {
	go func() {
		s.contentPicker.Hide()
		text := make(chan string)

		s.bridge.EnterSendText(s.app, text)
		t := <-text

		if t == "" {
			return
		}

		code := s.sendList.NewSendItem(storage.NewURI("Text Snippet"))

		go func(i int) {
			err := s.bridge.SendText(t, code)
			if err != nil {
				s.sendList.RemoveItem(i)
				dialog.ShowError(err, s.window)
			} else if s.appSettings.Notifications {
				s.app.SendNotification(fyne.NewNotification("Send completed", "The sending of text completed successfully"))
			} else {
				s.sendList.Items[i].Progress.SetValue(1)
			}
		}(s.sendList.Length() - 1)
	}()
}

func (s *send) onContentToSend() {
	s.contentPicker.Show()
}

func (s *send) buildUI() *fyne.Container {
	s.fileChoice = &widget.Button{Text: "File", Icon: theme.FileIcon(), OnTapped: s.onFileSend}
	s.directoryChoice = &widget.Button{Text: "Directory", Icon: theme.FolderIcon(), OnTapped: s.onDirSend}
	s.textChoice = &widget.Button{Text: "Text", Icon: theme.DocumentCreateIcon(), OnTapped: s.onTextSend}

	choiceContent := container.NewGridWithColumns(1, s.fileChoice, s.directoryChoice, s.textChoice)
	s.contentPicker = dialog.NewCustom("Pick a content type", "Cancel", choiceContent, s.window)
	s.contentPicker.Hide() // Bug in Fyne API. Can be remove after Fyne 2.0 and later.

	s.sendList = widgets.NewSendList()
	s.contentToSend = &widget.Button{Text: "Add content to send", Icon: theme.ContentAddIcon(), OnTapped: s.onContentToSend}

	box := container.NewVBox(s.contentToSend, &widget.Label{})
	return container.NewBorder(box, nil, nil, nil, s.sendList)
}

func (s *send) tabItem() *container.TabItem {
	return container.NewTabItemWithIcon("Send", theme.MailSendIcon(), s.buildUI())
}
