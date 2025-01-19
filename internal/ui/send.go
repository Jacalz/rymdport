package ui

import (
	"github.com/Jacalz/rymdport/v3/internal/ui/components"
	qrcode "github.com/rymdport/go-qrcode"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

func buildSendView(nav *components.StackNavigator) fyne.CanvasObject {
	code := "123-example-code"

	qr, err := qrcode.New("wormhole-transfer:"+code, qrcode.High)
	if err != nil {
		fyne.LogError("Failed to encode qr code", err)
	}

	qr.DisableBorder = true
	qr.BackgroundColor = theme.Color(theme.ColorNameBackground)
	qr.ForegroundColor = theme.Color(theme.ColorNameForeground)

	const size = 150
	image := &canvas.Image{Image: qr.Image(size), FillMode: canvas.ImageFillOriginal, ScaleMode: canvas.ImageScalePixels}
	image.SetMinSize(fyne.NewSquareSize(float32(size)))

	codeLabel := &widget.Label{Text: code, TextStyle: fyne.TextStyle{Bold: true}, Alignment: fyne.TextAlignCenter}

	progress := &widget.ProgressBar{Max: 100}

	cancel := &widget.Button{Text: "Cancel", OnTapped: func() { nav.Pop() }}

	return container.NewCenter(
		container.NewVBox(
			image,
			codeLabel,
			progress,
			&widget.Separator{},
			container.NewCenter(cancel),
		),
	)
}

func createSendPage(nav *components.StackNavigator) fyne.CanvasObject {
	icon := canvas.NewImageFromResource(theme.UploadIcon())
	icon.FillMode = canvas.ImageFillContain
	icon.SetMinSize(fyne.NewSquareSize(200))

	description := &widget.Label{Text: "Select data type below or drop files here.", Alignment: fyne.TextAlignCenter}

	file := &widget.Button{
		Icon:       theme.FileTextIcon(),
		Text:       "Send File",
		Importance: widget.HighImportance,
		OnTapped:   func() { nav.Push(buildSendView(nav), "Sending File") },
	}
	folder := &widget.Button{
		Icon:       theme.FolderIcon(),
		Text:       "Send Folder",
		Importance: widget.HighImportance,
		OnTapped:   func() { nav.Push(buildSendView(nav), "Sending Folder") },
	}
	text := &widget.Button{
		Icon:       theme.DocumentIcon(),
		Text:       "Send Text",
		Importance: widget.HighImportance,
		OnTapped:   func() { nav.Push(buildSendView(nav), "Sending Text") },
	}

	buttons := container.NewCenter(container.NewHBox(file, &widget.Separator{}, folder, &widget.Separator{}, text))
	content := container.NewVBox(icon, description, &widget.Separator{}, buttons)
	return container.NewCenter(content)
}
