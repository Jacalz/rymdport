package main

import (
	"net/url"

	"fyne.io/fyne"
	"fyne.io/fyne/canvas"
	"fyne.io/fyne/layout"
	"fyne.io/fyne/theme"
	"fyne.io/fyne/widget"
	"github.com/Jacalz/wormhole-gui/assets"
)

const version string = "v2.0.0"

func parseURL(input string) *url.URL {
	link, err := url.Parse(input)
	if err != nil {
		fyne.LogError("Could not parse URL string", err)
	}

	return link
}

func aboutTab() *widget.TabItem {
	logo := canvas.NewImageFromResource(assets.AppIcon)
	logo.SetMinSize(fyne.NewSize(256, 256))

	content := fyne.NewContainerWithLayout(layout.NewVBoxLayout(),
		layout.NewSpacer(),
		widget.NewHBox(layout.NewSpacer(), logo, layout.NewSpacer()),
		widget.NewHBox(
			layout.NewSpacer(),
			widget.NewLabelWithStyle("wormhole-gui", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
			widget.NewLabelWithStyle("-", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
			widget.NewHyperlinkWithStyle(version, parseURL("https://github.com/Jacalz/wormhole-gui/releases/tag/"+version), fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
			layout.NewSpacer(),
		),
		layout.NewSpacer(),
	)

	return widget.NewTabItemWithIcon("About", theme.InfoIcon(), content)
}
