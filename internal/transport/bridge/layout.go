package bridge

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

var maxMinSizeHeight float32 // Keeping all instances of the list layout consistent in height

type listLayout struct{}

func (g *listLayout) Layout(objects []fyne.CanvasObject, size fyne.Size) {
	objects[0].Move(fyne.NewPos(0, theme.Padding()))
	objects[0].Resize(fyne.NewSize(size.Height-theme.Padding(), size.Height-2*theme.Padding()))

	cellSize := (size.Width - size.Height - theme.Padding()) / (float32(len(objects) - 1))
	start, end := size.Height, size.Height+cellSize-theme.Padding()
	for _, child := range objects[1:] {
		if _, label := child.(*widget.Label); label {
			child.Move(fyne.NewPos(start, (size.Height-child.MinSize().Height)/2))
		} else {
			child.Move(fyne.NewPos(start, theme.Padding()))
		}

		child.Resize(fyne.NewSize(end-start, size.Height-2*theme.Padding()))

		start = end + theme.Padding()
		end = start + cellSize
	}
}

// MinSize finds the smallest size that satisfies all the child objects.
// Height will stay consistent between each each instance.
func (g *listLayout) MinSize(objects []fyne.CanvasObject) fyne.Size {
	maxMinSizeWidth := float32(0)
	for _, child := range objects {
		if child.Visible() {
			maxMinSizeWidth += child.MinSize().Width
			maxMinSizeHeight = fyne.Max(child.MinSize().Height, maxMinSizeHeight)
		}
	}

	return fyne.NewSize(maxMinSizeWidth, maxMinSizeHeight+2*theme.Padding())
}
