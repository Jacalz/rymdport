package main

import (
	"fyne.io/fyne"
	"fyne.io/fyne/layout"
	"fyne.io/fyne/theme"
	"fyne.io/fyne/widget"
)

func recieveTab() *widget.TabItem {
	codeEntry := widget.NewEntry()
	codeEntry.SetPlaceHolder("Enter code")

	codeButton := widget.NewButtonWithIcon("Download", theme.MoveDownIcon(), nil)

	codeContainer := fyne.NewContainerWithLayout(layout.NewGridLayout(2), codeEntry, codeButton)

	recieveContent := fyne.NewContainerWithLayout(layout.NewVBoxLayout(), codeContainer)

	return widget.NewTabItemWithIcon("Recieve", theme.MoveDownIcon(), recieveContent)
}
