package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

// Create sets up the user interface for the application.
func Create(a fyne.App, w fyne.Window) fyne.CanvasObject {
	menu := buildApplicationMenu(a)
	dropdown := &widget.Button{Icon: theme.MenuIcon(), Importance: widget.LowImportance}
	dropdown.OnTapped = func() {
		offset := fyne.Position{Y: dropdown.Size().Height + theme.Padding()}
		widget.ShowPopUpMenuAtRelativePosition(menu, w.Canvas(), offset, dropdown)
	}

	tabs := &container.AppTabs{
		Items: []*container.TabItem{
			{Text: "Send", Icon: theme.UploadIcon(), Content: createSendPage(a, w)},
			{Text: "Receive", Icon: theme.DownloadIcon(), Content: createRecvPage(a, w)},
		},
	}

	upperRightCorner := container.NewBorder(container.NewBorder(nil, nil, nil, dropdown), nil, nil, nil)
	return container.NewStack(tabs, upperRightCorner)
}
