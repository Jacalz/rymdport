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

func buildRecvView(nav *components.StackNavigator, code string) fyne.CanvasObject {
	image := &canvas.Image{FillMode: canvas.ImageFillContain}
	image.SetMinSize(fyne.NewSquareSize(150))

	recvType := 0
	switch recvType {
	case 0:
		image.Resource = theme.FileImageIcon()
	case 1:
		image.Resource = theme.FolderIcon()
	case 2:
		image.Resource = theme.DocumentIcon()
	}

	name := "Something"

	codeStyle := widget.RichTextStyleSubHeading
	codeStyle.Alignment = fyne.TextAlignCenter
	nameStyle := widget.RichTextStyleInline
	nameStyle.Alignment = fyne.TextAlignCenter
	text := widget.NewRichText(&widget.TextSegment{Style: codeStyle, Text: code}, &widget.TextSegment{Style: nameStyle, Text: name})

	progress := widget.NewProgressBar()
	progress.SetValue(0.5)

	cancel := &widget.Button{Text: "Cancel", OnTapped: nav.Pop, Importance: widget.WarningImportance}

	return container.NewCenter(
		container.NewVBox(
			image,
			text,
			progress,
			&widget.Separator{},
			container.NewCenter(cancel),
		),
	)
}

func createRecvPage(nav *components.StackNavigator) fyne.CanvasObject {
	icon := canvas.NewImageFromResource(theme.DownloadIcon())
	icon.FillMode = canvas.ImageFillContain
	icon.SetMinSize(fyne.NewSquareSize(200))

	description := &widget.Label{Text: "Enter a code below to start receiving data.", Alignment: fyne.TextAlignCenter}

	code := &widget.Entry{PlaceHolder: "Code from sender", Validator: util.CodeValidator}
	code.OnSubmitted = func(input string) {
		if code.Validator(input) != nil {
			return
		}

		nav.Push(buildRecvView(nav, input), "Receiving Data")
	}

	start := &widget.Button{
		Text:       "Start Receive",
		Icon:       theme.DownloadIcon(),
		Importance: widget.HighImportance,
		OnTapped:   func() { code.OnSubmitted(code.Text) },
	}

	content := container.NewVBox(icon, description, &widget.Separator{}, code, &widget.Separator{}, container.NewCenter(start))

	return container.NewCenter(content)
}
