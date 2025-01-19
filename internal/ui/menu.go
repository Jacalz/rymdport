package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"fyne.io/x/fyne/dialog"
	"github.com/Jacalz/rymdport/v3/internal/util"
)

func openSettingsMenu(a fyne.App) {
	w := a.NewWindow("Settings")
	w.SetContent(widget.NewLabel("Under contstruction..."))
	w.Show()
}

func openAboutMenu(a fyne.App) {
	links := []*widget.Hyperlink{
		{Text: "Repository", URL: util.URLToGitHubProject("")},
		{Text: "Issue Tracker", URL: util.URLToGitHubProject("/issues")},
		{Text: "Wiki", URL: util.URLToGitHubProject("/wiki")},
	}
	w := dialog.NewAboutWindow("Easy encrypted file, folder, and text sharing between devices.", links, a)
	w.Resize(fyne.NewSize(500, 300))
	w.Show()
}

func buildApplicationMenu(a fyne.App) *fyne.Menu {
	menus := []*fyne.MenuItem{
		{Label: "Settings", Icon: theme.SettingsIcon(), Action: func() { openSettingsMenu(a) }},
		{Label: "About", Icon: theme.InfoIcon(), Action: func() { openAboutMenu(a) }},
	}

	return &fyne.Menu{Items: menus}
}
