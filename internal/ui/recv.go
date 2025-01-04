package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/Jacalz/rymdport/v3/internal/util"
)

func createRecvPage(_ fyne.App, _ fyne.Window) fyne.CanvasObject {
	icon := canvas.NewImageFromResource(theme.DownloadIcon())
	icon.FillMode = canvas.ImageFillContain
	icon.SetMinSize(fyne.NewSquareSize(200))

	title := &widget.Label{Text: "Receive Data", Alignment: fyne.TextAlignCenter, TextStyle: fyne.TextStyle{Bold: true}}
	description := &widget.Label{Text: "Enter a code below to start receiving data", Alignment: fyne.TextAlignCenter}

	code := &widget.Entry{PlaceHolder: "Code from sender", Validator: util.CodeValidator}

	content := container.NewVBox(icon, &widget.Separator{}, title, description, &widget.Separator{}, code)

	return container.NewCenter(content)
}
