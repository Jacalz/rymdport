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
	version = "v2.0.1"
	rurl    = "https://github.com/Jacalz/wormhole-gui/releases/tag/" + version
)

type about struct {
	icon *canvas.Image
}

func newAbout() *about {
	return &about{}
}

func (a *about) buildUI() *fyne.Container {
	a.icon = canvas.NewImageFromResource(assets.AppIcon)
	a.icon.SetMinSize(fyne.NewSize(256, 256))

	return container.NewVBox(
		layout.NewSpacer(),
		container.NewHBox(layout.NewSpacer(), a.icon, layout.NewSpacer()),
		widget.NewHBox(
			layout.NewSpacer(),
			widget.NewLabelWithStyle("wormhole-gui", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
			widget.NewLabelWithStyle("-", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
			widget.NewHyperlinkWithStyle(version, parseURL(rurl), fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
			layout.NewSpacer(),
		),
		layout.NewSpacer(),
	)
}

func (a *about) tabItem() *container.TabItem {
	return container.NewTabItemWithIcon("About", theme.InfoIcon(), a.buildUI())
}

func parseURL(input string) *url.URL {
	link, err := url.Parse(input)
	if err != nil {
		fyne.LogError("Could not parse URL string", err)
	}

	return link
}
