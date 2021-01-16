package ui

import (
	"fyne.io/fyne"
	"fyne.io/fyne/container"
	"fyne.io/fyne/dialog"
	"fyne.io/fyne/theme"
	"fyne.io/fyne/widget"
	"github.com/Jacalz/wormhole-gui/internal/transport"
	"github.com/Jacalz/wormhole-gui/internal/transport/bridge"
)

type send struct {
	contentPicker dialog.Dialog

	fileChoice      *widget.Button
	fileDialog      *dialog.FileDialog
	directoryChoice *widget.Button
	directoryDialog *dialog.FileDialog
	textChoice      *widget.Button

	contentToSend *widget.Button
	sendList      *bridge.SendList

	client      *transport.Client
	appSettings *AppSettings
	window      fyne.Window
	app         fyne.App
}

func newSend(a fyne.App, w fyne.Window, c *transport.Client, as *AppSettings) *send {
	return &send{app: a, window: w, client: c, appSettings: as}
}

func (s *send) onFileSend() {
	s.contentPicker.Hide()
	s.fileDialog.Show()
}

func (s *send) onDirSend() {
	s.contentPicker.Hide()
	s.directoryDialog.Show()
}

func (s *send) onTextSend() {
	s.contentPicker.Hide()
	s.sendList.SendText()
}

func (s *send) buildUI() *fyne.Container {
	s.fileChoice = &widget.Button{Text: "File", Icon: theme.FileIcon(), OnTapped: s.onFileSend}
	s.directoryChoice = &widget.Button{Text: "Directory", Icon: theme.FolderOpenIcon(), OnTapped: s.onDirSend}
	s.textChoice = &widget.Button{Text: "Text", Icon: theme.DocumentCreateIcon(), OnTapped: s.onTextSend}

	choiceContent := container.NewGridWithColumns(1, s.fileChoice, s.directoryChoice, s.textChoice)
	s.contentPicker = dialog.NewCustom("Pick a content type", "Cancel", choiceContent, s.window)

	s.sendList = bridge.NewSendList(s.client)
	s.contentToSend = &widget.Button{Text: "Add content to send", Icon: theme.ContentAddIcon(), OnTapped: s.contentPicker.Show}

	s.fileDialog = dialog.NewFileOpen(s.sendList.OnFileSelect, s.window)
	s.directoryDialog = dialog.NewFolderOpen(s.sendList.OnDirSelect, s.window)

	box := container.NewVBox(s.contentToSend, &widget.Label{})
	return container.NewBorder(box, nil, nil, nil, s.sendList)
}

func (s *send) tabItem() *container.TabItem {
	return &container.TabItem{Text: "Send", Icon: theme.MailSendIcon(), Content: s.buildUI()}
}
