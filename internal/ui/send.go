package ui

import (
	"github.com/Jacalz/rymdport/v3/internal/ui/components"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

func buildSendView() fyne.CanvasObject {
	return widget.NewLabel("Sending will be implemented soon...")
}

func createSendPage(navigator *components.StackNavigator) fyne.CanvasObject {
	icon := canvas.NewImageFromResource(theme.UploadIcon())
	icon.FillMode = canvas.ImageFillContain
	icon.SetMinSize(fyne.NewSquareSize(200))

	description := &widget.Label{Text: "Select data type below or drop files here.", Alignment: fyne.TextAlignCenter}

	sendView := buildSendView()
	file := &widget.Button{
		Icon:       theme.FileTextIcon(),
		Text:       "Send File",
		Importance: widget.HighImportance,
		OnTapped:   func() { navigator.Push(sendView, "Sending File") },
	}
	folder := &widget.Button{
		Icon:       theme.FolderIcon(),
		Text:       "Send Folder",
		Importance: widget.HighImportance,
		OnTapped:   func() { navigator.Push(sendView, "Sending Folder") },
	}
	text := &widget.Button{
		Icon:       theme.DocumentIcon(),
		Text:       "Send Text",
		Importance: widget.HighImportance,
		OnTapped:   func() { navigator.Push(sendView, "Sending Text") },
	}

	buttons := container.NewCenter(container.NewHBox(file, &widget.Separator{}, folder, &widget.Separator{}, text))
	content := container.NewVBox(icon, description, &widget.Separator{}, buttons)
	return container.NewCenter(content)
}
