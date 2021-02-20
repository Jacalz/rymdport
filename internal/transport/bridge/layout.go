package bridge

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

var maxMinSizeHeight float32 // Keeping all instances of the list layout consistent in height

type listLayout struct{}

// Layout is called to pack all child objects into a specified size.
func (g *listLayout) Layout(objects []fyne.CanvasObject, size fyne.Size) {
	length := float32(len(objects) - 1)
	padWidth := length * theme.Padding()
	cellWidth := (size.Width - padWidth - objects[0].MinSize().Width) / length

	start, end := float32(0), size.Height
	for i, child := range objects {
		if !child.Visible() {
			continue
		}

		if i == 0 {
			child.Move(fyne.NewPos(0, 0))
			child.Resize(fyne.NewSize(size.Height, size.Height))
			start = end + theme.Padding()
			continue
		}

		if i != 1 {
			start += cellWidth + theme.Padding()
		}

		end = start + cellWidth

		_, isLabel := child.(*widget.Label)
		if cont, ok := child.(*fyne.Container); ok {
			_, isLabel = cont.Objects[0].(*codeDisplay)
		}

		if isLabel { // Proper vertical alignment for text labels
			child.Move(fyne.NewPos(start, (size.Height-child.MinSize().Height)/2))
		} else {
			child.Move(fyne.NewPos(start, 0))
		}

		// TODO: Better solution for too long text
		if width := child.MinSize().Width; width > cellWidth {
			start += width - cellWidth
		}

		child.Resize(fyne.NewSize(end-start, size.Height))
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

	return fyne.NewSize(maxMinSizeWidth, maxMinSizeHeight+theme.Padding())
}
