package widgets

import (
	"fyne.io/fyne"
	"fyne.io/fyne/theme"
)

type sendLayout struct{}

// Layout is called to pack all child objects into a specified size.
func (g *sendLayout) Layout(objects []fyne.CanvasObject, size fyne.Size) {
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

			child.Move(fyne.NewPos(oldx, 0))
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
func (g *sendLayout) MinSize(objects []fyne.CanvasObject) fyne.Size {
	len := len(objects)
	minSize := fyne.NewSize(0, 0)
	for _, child := range objects {
		if !child.Visible() {
			continue
		}

		minSize = minSize.Max(child.MinSize())
	}

	minContentSize := fyne.NewSize(minSize.Width*len, minSize.Height)
	return minContentSize.Add(fyne.NewSize(theme.Padding()*fyne.Max(len-1, 0), 0))
}

// newSendLayout creates a grid that keep all element in one long grid.
// The first object will be an icon and the rest will be equally sized according to the given size.
func newSendLayout() fyne.Layout {
	return &sendLayout{}
}
