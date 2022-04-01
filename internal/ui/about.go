package ui

import (
	"net/url"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/Jacalz/rymdport/v3/internal/assets"
)

const version = "v3.0.1"

var releaseURL = &url.URL{
	Scheme: "https",
	Host:   "github.com",
	Path:   "/jacalz/rymdport/releases/tag/" + version,
}

type about struct {
	icon        *canvas.Image
	nameLabel   *widget.Label
	spacerLabel *widget.Label
	hyperlink   *widget.Hyperlink
}

func newAbout() *about {
	return &about{}
}

func (a *about) buildUI() *fyne.Container {
	a.icon = canvas.NewImageFromResource(assets.AppIcon)
	a.icon.SetMinSize(fyne.NewSize(256, 256))

	a.nameLabel = newBoldLabel("Rymdport")
	a.spacerLabel = newBoldLabel("-")
	a.hyperlink = &widget.Hyperlink{Text: version, URL: releaseURL, TextStyle: fyne.TextStyle{Bold: true}}

	return container.NewVBox(
		layout.NewSpacer(),
		container.NewHBox(layout.NewSpacer(), a.icon, layout.NewSpacer()),
		container.NewHBox(
			layout.NewSpacer(),
			a.nameLabel,
			a.spacerLabel,
			a.hyperlink,
			layout.NewSpacer(),
		),
		layout.NewSpacer(),
	)
}

func (a *about) tabItem() *container.TabItem {
	return &container.TabItem{Text: "About", Icon: theme.InfoIcon(), Content: a.buildUI()}
}
