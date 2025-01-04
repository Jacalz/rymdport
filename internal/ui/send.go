package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

func createSendPage(_ fyne.App, _ fyne.Window) fyne.CanvasObject {
	icon := canvas.NewImageFromResource(theme.UploadIcon())
	icon.SetMinSize(fyne.NewSquareSize(200))

	title := &widget.Label{Text: "Send Data", Alignment: fyne.TextAlignCenter, TextStyle: fyne.TextStyle{Bold: true}}
	description := &widget.Label{Text: "Select data type below or drop files here", Alignment: fyne.TextAlignCenter}

	// Buttons for starting sends.
	file := &widget.Button{Icon: theme.FileTextIcon(), Text: "Send File"}
	folder := &widget.Button{Icon: theme.FolderIcon(), Text: "Send Folder"}
	text := &widget.Button{Icon: theme.DocumentIcon(), Text: "Send Text"}

	buttons := container.NewVBox(icon, &widget.Separator{}, title, description, &widget.Separator{}, file, folder, text)

	return container.NewCenter(buttons)
}
