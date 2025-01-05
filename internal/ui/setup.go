package ui

import (
	"github.com/Jacalz/rymdport/v3/internal/ui/components"

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

	navigator := &components.StackNavigator{}
	navigator.OnBack = navigator.Pop
	tabs := &container.AppTabs{
		Items: []*container.TabItem{
			{Text: "Send", Icon: theme.UploadIcon(), Content: createSendPage(navigator)},
			{Text: "Receive", Icon: theme.DownloadIcon(), Content: createRecvPage(navigator)},
		},
	}

	upperRightCorner := container.NewBorder(container.NewBorder(nil, nil, nil, dropdown), nil, nil, nil)
	navigator.Push(container.NewStack(tabs, upperRightCorner), "")
	return navigator
}
