package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/Jacalz/rymdport/v3/internal/transport"
	"github.com/Jacalz/rymdport/v3/internal/transport/bridge"
	"github.com/Jacalz/rymdport/v3/internal/util"
)

type send struct {
	contentPicker   dialog.Dialog // TODO: Switch to non-interface type when using Fyne v2.4.0.
	fileDialog      *dialog.FileDialog
	directoryDialog *dialog.FileDialog

	data *bridge.SendData

	client *transport.Client
	canvas fyne.Canvas
}

func newSendTab(w fyne.Window, c *transport.Client) *container.TabItem {
	send := &send{client: c, canvas: w.Canvas()}

	return &container.TabItem{
		Text:    "Send",
		Icon:    theme.MailSendIcon(),
		Content: send.buildUI(w),
	}
}

func (s *send) buildUI(window fyne.Window) *fyne.Container {
	fileChoice := &widget.Button{Text: "File", Icon: theme.FileIcon(), OnTapped: s.onFileSend}
	directoryChoice := &widget.Button{Text: "Directory", Icon: theme.FolderOpenIcon(), OnTapped: s.onDirSend}
	textChoice := &widget.Button{Text: "Text", Icon: theme.DocumentCreateIcon(), OnTapped: s.onTextSend}
	codeChoice := &widget.Check{Text: "Use a custom code", OnChanged: s.onCustomCode}

	choiceContent := container.NewGridWithColumns(1, fileChoice, directoryChoice, textChoice, codeChoice)
	s.contentPicker = dialog.NewCustom("Pick a content type", "Cancel", choiceContent, window)

	s.data = &bridge.SendData{Client: s.client, Window: window, Canvas: s.canvas}
	contentToSend := &widget.Button{Text: "Add content to send", Icon: theme.ContentAddIcon(), OnTapped: s.contentPicker.Show}

	s.fileDialog = dialog.NewFileOpen(s.data.OnFileSelect, window)
	s.directoryDialog = dialog.NewFolderOpen(s.data.OnDirSelect, window)

	box := container.NewVBox(&widget.Separator{}, contentToSend, &widget.Separator{})
	return container.NewBorder(box, nil, nil, nil, bridge.NewSendList(s.data))
}

func (s *send) onFileSend() {
	s.contentPicker.Hide()
	s.fileDialog.Resize(util.WindowSizeToDialog(s.canvas.Size()))
	s.fileDialog.Show()
}

func (s *send) onDirSend() {
	s.contentPicker.Hide()
	s.fileDialog.Resize(util.WindowSizeToDialog(s.canvas.Size()))
	s.directoryDialog.Show()
}

func (s *send) onTextSend() {
	s.contentPicker.Hide()
	s.data.SendText()
}

func (s *send) onCustomCode(enabled bool) {
	s.client.CustomCode = enabled
}
