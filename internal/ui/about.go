package ui

import (
	"net/url"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/Jacalz/rymdport/v3/internal/util"
)

func newAboutTab(app fyne.App) *container.TabItem {
	const version = "v3.9.2"

	repoURL := util.URLToGitHubProject("")
	icon := newClickableIcon(app.Icon(), repoURL, app)

	nameLabel := newBoldLabel("Rymdport")
	spacerLabel := newBoldLabel("-")

	releaseURL := util.URLToGitHubProject("/releases/tag/" + version)
	hyperlink := &widget.Hyperlink{Text: version, URL: releaseURL, TextStyle: fyne.TextStyle{Bold: true}}

	spacer := &layout.Spacer{}
	content := container.NewVBox(
		spacer,
		container.NewHBox(spacer, icon, spacer),
		container.NewHBox(
			spacer,
			nameLabel,
			spacerLabel,
			hyperlink,
			spacer,
		),
		spacer,
	)

	return &container.TabItem{
		Text:    "About",
		Icon:    theme.InfoIcon(),
		Content: content,
	}
}

type clickableIcon struct {
	widget.Icon
	app fyne.App
	url *url.URL
}

func (c *clickableIcon) Tapped(_ *fyne.PointEvent) {
	err := c.app.OpenURL(c.url)
	if err != nil {
		fyne.LogError("Failed to open repository: ", err)
	}
}

func (c *clickableIcon) Cursor() desktop.Cursor {
	return desktop.PointerCursor
}

func (c *clickableIcon) MinSize() fyne.Size {
	return fyne.Size{Width: 256, Height: 256}
}

func newClickableIcon(res fyne.Resource, url *url.URL, app fyne.App) *clickableIcon {
	return &clickableIcon{Icon: widget.Icon{Resource: res}, app: app, url: url}
}
