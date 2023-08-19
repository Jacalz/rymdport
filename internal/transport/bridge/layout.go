package bridge

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

type listLayout struct{}

func (g *listLayout) Layout(objects []fyne.CanvasObject, size fyne.Size) {
	padding := theme.InnerPadding()
	doublePadding := 2 * padding

	objects[0].Move(fyne.NewPos(0, padding))
	objects[0].Resize(fyne.NewSize(size.Height-padding, size.Height-doublePadding))

	cellSize := (size.Width - size.Height - doublePadding) / (float32(len(objects) - 1))
	start, end := size.Height, size.Height+cellSize-padding
	for _, child := range objects[1:] {
		if _, label := child.(*widget.Label); label {
			child.Move(fyne.NewPos(start, (size.Height-child.MinSize().Height)/2))
		} else {
			child.Move(fyne.NewPos(start, padding))
		}

		child.Resize(fyne.NewSize(end-start, size.Height-doublePadding))

		start = end + padding
		end = start + cellSize
	}
}

// MinSize finds the smallest size that satisfies all the child objects.
// Height will stay consistent between each each instance.
func (g *listLayout) MinSize(objects []fyne.CanvasObject) fyne.Size {
	doublePadding := 2 * theme.InnerPadding()

	maxMinSizeWidth := float32(0)
	maxMinSizeHeight := theme.IconInlineSize() + doublePadding // Default button height with icon
	for _, child := range objects {
		if child.Visible() {
			min := child.MinSize()
			maxMinSizeWidth += min.Width
			maxMinSizeHeight = fyne.Max(min.Height, maxMinSizeHeight)
		}
	}

	return fyne.NewSize(maxMinSizeWidth, maxMinSizeHeight+doublePadding)
}
