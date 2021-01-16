package ui

import (
	"net/url"

	"fyne.io/fyne"
	"fyne.io/fyne/canvas"
	"fyne.io/fyne/container"
	"fyne.io/fyne/layout"
	"fyne.io/fyne/theme"
	"fyne.io/fyne/widget"
	"github.com/Jacalz/wormhole-gui/internal/assets"
)

const (
	version    = "v2.1.0"
	releaseURL = "https://github.com/Jacalz/wormhole-gui/releases/tag/" + version
)

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

	a.nameLabel = newBoldLabel("wormhole-gui")
	a.spacerLabel = newBoldLabel("-")
	a.hyperlink = &widget.Hyperlink{Text: version, URL: parseURL(releaseURL), TextStyle: fyne.TextStyle{Bold: true}}

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

func parseURL(input string) *url.URL {
	link, err := url.Parse(input)
	if err != nil {
		fyne.LogError("Could not parse URL string", err)
	}

	return link
}
