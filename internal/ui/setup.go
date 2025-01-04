package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/Jacalz/rymdport/v3/internal/util"

	"fyne.io/x/fyne/dialog"
)

// Create sets up the user interface for the application.
func Create(a fyne.App, w fyne.Window) fyne.CanvasObject {
	dropdown := &widget.Button{Icon: theme.MenuIcon(), Importance: widget.LowImportance}
	dropdown.OnTapped = func() {
		widget.ShowPopUpMenuAtRelativePosition(
			&fyne.Menu{Items: []*fyne.MenuItem{{Label: "About", Icon: theme.InfoIcon(), Action: func() {
				links := []*widget.Hyperlink{
					{Text: "Repository", URL: util.URLToGitHubProject("")},
					{Text: "Issue Tracker", URL: util.URLToGitHubProject("/issues")},
					{Text: "Wiki", URL: util.URLToGitHubProject("/wiki")},
				}
				abwin := dialog.NewAboutWindow("Easy encrypted file, folder, and text sharing between devices.", links, a)
				abwin.Resize(fyne.NewSize(500, 300))
				abwin.Show()
			}}}}, w.Canvas(), fyne.Position{Y: dropdown.Size().Height}, dropdown)
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
