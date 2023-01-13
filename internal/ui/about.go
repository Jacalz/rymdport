package ui

import (
	"net/url"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

const version = "v3.3.0-rc2"

type about struct {
	icon        *clickableIcon
	nameLabel   *widget.Label
	spacerLabel *widget.Label
	hyperlink   *widget.Hyperlink

	app fyne.App
}

func newAbout(app fyne.App) *about {
	return &about{app: app}
}

func (a *about) buildUI() *fyne.Container {
	const (
		https  = "https"
		github = "github.com"
	)

	repoURL := &url.URL{Scheme: https, Host: github, Path: "/jacalz/rymdport"}
	a.icon = newClickableIcon(a.app.Icon(), repoURL, a.app)

	a.nameLabel = newBoldLabel("Rymdport")
	a.spacerLabel = newBoldLabel("-")

	releaseURL := &url.URL{
		Scheme: https, Host: github,
		Path: "/jacalz/rymdport/releases/tag/" + version,
	}
	a.hyperlink = &widget.Hyperlink{Text: version, URL: releaseURL, TextStyle: fyne.TextStyle{Bold: true}}

	spacer := &layout.Spacer{}
	return container.NewVBox(
		spacer,
		container.NewHBox(spacer, a.icon, spacer),
		container.NewHBox(
			spacer,
			a.nameLabel,
			a.spacerLabel,
			a.hyperlink,
			spacer,
		),
		spacer,
	)
}

func (a *about) tabItem() *container.TabItem {
	return &container.TabItem{Text: "About", Icon: theme.InfoIcon(), Content: a.buildUI()}
}

type clickableIcon struct {
	widget.BaseWidget
	app  fyne.App
	url  *url.URL
	icon *canvas.Image
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

func (c *clickableIcon) CreateRenderer() fyne.WidgetRenderer {
	c.ExtendBaseWidget(c)
	return widget.NewSimpleRenderer(c.icon)
}

func (c *clickableIcon) MinSize() fyne.Size {
	return fyne.Size{Width: 256, Height: 256}
}

func newClickableIcon(res fyne.Resource, url *url.URL, app fyne.App) *clickableIcon {
	return &clickableIcon{app: app, url: url, icon: canvas.NewImageFromResource(res)}
}
