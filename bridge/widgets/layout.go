package widgets

import (
	"math"

	"fyne.io/fyne"
	"fyne.io/fyne/theme"
)

type sendLayout struct{}

// Get the leading (top or left) edge of a grid cell.
func getLeading(size float64, offset int) int {
	ret := (size + float64(theme.Padding())) * float64(offset)
	return int(math.Round(ret))
}

// Get the trailing (bottom or right) edge of a grid cell.
func getTrailing(size float64, offset int) int {
	return getLeading(size, offset+1) - theme.Padding()
}

// Layout is called to pack all child objects into a specified size.
func (g *sendLayout) Layout(objects []fyne.CanvasObject, size fyne.Size) {
	padWidth := (len(objects) - 1) * theme.Padding()
	cellWidth := float64(size.Width-padWidth-objects[0].MinSize().Width) / float64(len(objects))

	offset := 0
	for i, child := range objects {
		if !child.Visible() {
			continue
		}

		x1, y1, x2, y2 := 0, 0, 0, 0
		if i == 0 {
			x2 = size.Height + theme.Padding()
			y2 = x2
			offset = x2 * 2

			child.Move(fyne.NewPos(x1, y1))
		} else {
			x1 = getLeading(cellWidth, i)
			y1 = getLeading(float64(size.Height), 0)
			x2 = getTrailing(cellWidth, i)
			y2 = getTrailing(float64(size.Height), 0)

			if i == len(objects)-1 {
				offset = -offset / 6
			}

			child.Move(fyne.NewPos(x1-offset, y1))
		}

		// TODO: Make sure to make cell larger if cellWidth < minSize

		child.Resize(fyne.NewSize(x2-x1, y2-y1))
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
