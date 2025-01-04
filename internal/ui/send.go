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
	icon.FillMode = canvas.ImageFillContain
	icon.SetMinSize(fyne.NewSquareSize(200))

	description := &widget.Label{Text: "Select data type below or drop files here.", Alignment: fyne.TextAlignCenter}

	file := &widget.Button{Icon: theme.FileTextIcon(), Text: "Send File", Importance: widget.HighImportance}
	folder := &widget.Button{Icon: theme.FolderIcon(), Text: "Send Folder", Importance: widget.HighImportance}
	text := &widget.Button{Icon: theme.DocumentIcon(), Text: "Send Text", Importance: widget.HighImportance}

	buttons := container.NewCenter(container.NewHBox(file, &widget.Separator{}, folder, &widget.Separator{}, text))
	content := container.NewVBox(icon, description, &widget.Separator{}, buttons)
	return container.NewCenter(content)
}
