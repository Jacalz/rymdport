package bridge

import (
	"fyne.io/fyne"
	"fyne.io/fyne/theme"
	"fyne.io/fyne/widget"
)

var maxMinSizeHeight float32 // Keeping all instances of the list layout consistent in height

type listLayout struct{}

// Layout is called to pack all child objects into a specified size.
func (g *listLayout) Layout(objects []fyne.CanvasObject, size fyne.Size) {
	length := float32(len(objects) - 1)
	padWidth := length * theme.Padding()
	cellWidth := (size.Width - padWidth - objects[0].MinSize().Width) / length

	y2 := size.Height
	oldx, newx := float32(0), y2

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
				_, isLabel = cont.Objects[0].(*codeDisplay)
			}

			if isLabel { // Proper middle alignment for text labels
				child.Move(fyne.NewPos(oldx, (y2-child.MinSize().Height)/2))
			} else {
				child.Move(fyne.NewPos(oldx, 0))
			}

		} else {
			newx = size.Height + theme.Padding()
			child.Move(fyne.NewPos(0, 0))
		}

		// TODO: Better solution for too long text
		if child.MinSize().Width > cellWidth {
			oldx += (child.MinSize().Width - cellWidth)
		}

		child.Resize(fyne.NewSize(newx-oldx, y2))
	}
}

// MinSize finds the smallest size that satisfies all the child objects.
// Height will stay consistent between each each instance.
func (g *listLayout) MinSize(objects []fyne.CanvasObject) fyne.Size {
	length := float32(len(objects) - 1)
	maxMinSizeWidth := float32(0)
	for _, child := range objects {
		if child.Visible() {
			maxMinSizeWidth = fyne.Max(child.MinSize().Width, maxMinSizeWidth)
			maxMinSizeHeight = fyne.Max(child.MinSize().Height, maxMinSizeHeight)
		}
	}

	return fyne.NewSize((maxMinSizeWidth+theme.Padding())*length, maxMinSizeHeight+2*theme.Padding())
}
