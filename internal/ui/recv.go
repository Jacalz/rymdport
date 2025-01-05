package ui

import (
	"github.com/Jacalz/rymdport/v3/internal/ui/components"
	"github.com/Jacalz/rymdport/v3/internal/util"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

func buildRecvView() fyne.CanvasObject {
	return widget.NewLabel("Receiving will be implemented soon...")
}

func createRecvPage(navigator *components.StackNavigator) fyne.CanvasObject {
	icon := canvas.NewImageFromResource(theme.DownloadIcon())
	icon.FillMode = canvas.ImageFillContain
	icon.SetMinSize(fyne.NewSquareSize(200))

	description := &widget.Label{Text: "Enter a code below to start receiving data.", Alignment: fyne.TextAlignCenter}

	recvView := buildRecvView()
	code := &widget.Entry{PlaceHolder: "Code from sender", Validator: util.CodeValidator}
	start := &widget.Button{
		Text:       "Start Receive",
		Icon:       theme.DownloadIcon(),
		Importance: widget.HighImportance,
		OnTapped:   func() { navigator.Push(recvView, "Receiving Data") },
	}

	content := container.NewVBox(icon, description, &widget.Separator{}, code, &widget.Separator{}, container.NewCenter(start))

	return container.NewCenter(content)
}
