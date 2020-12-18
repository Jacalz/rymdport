package widgets

import (
	"fyne.io/fyne"
	"fyne.io/fyne/theme"
	"fyne.io/fyne/widget"
)

var maxMinSizeHeight int // Keeping all instances of the list layout consistent in height

type listLayout struct{}

// Layout is called to pack all child objects into a specified size.
func (g *listLayout) Layout(objects []fyne.CanvasObject, size fyne.Size) {
	padWidth := (len(objects) - 1) * theme.Padding()
	cellWidth := (size.Width - padWidth - objects[0].MinSize().Width) / (len(objects) - 1)

	y2 := size.Height

	oldx, newx := 0, y2
	for i, child := range objects {
		if !child.Visible() {
			continue
		}

		if i != 0 {
			if i == 1 {
				oldx += newx + theme.Padding()
			} else {
				oldx += cellWidth + theme.Padding()
			}

			newx = oldx + cellWidth

			_, isLabel := child.(*widget.Label)
			if cont, ok := child.(*fyne.Container); ok {
				_, isLabel = cont.Objects[0].(*CodeDisplay)
			}

			if isLabel {
				child.Move(fyne.NewPos(oldx, (y2-child.MinSize().Height)/2))
			} else {
				child.Move(fyne.NewPos(oldx, 0))
			}

		} else {
			newx = size.Height + theme.Padding()
			child.Move(fyne.NewPos(0, 0))
		}

		if child.MinSize().Width > cellWidth {
			oldx += (child.MinSize().Width - cellWidth)
		}

		child.Resize(fyne.NewSize(newx-oldx, y2))
	}
}

// MinSize finds the smallest size that satisfies all the child objects.
// Height will stay consistent between each each instance.
func (g *listLayout) MinSize(objects []fyne.CanvasObject) fyne.Size {
	len := len(objects) - 1
	maxMinSizeWidth := 0
	for _, child := range objects {
		if child.Visible() {
			maxMinSizeWidth = fyne.Max(child.MinSize().Width, maxMinSizeWidth)
			maxMinSizeHeight = fyne.Max(child.MinSize().Height, maxMinSizeHeight)
		}
	}

	return fyne.NewSize((maxMinSizeWidth+theme.Padding())*len, maxMinSizeHeight+2*theme.Padding())
}
