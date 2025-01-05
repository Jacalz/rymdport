package components

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

var _ fyne.Widget = (*StackNavigator)(nil)

// StackNavigator represents a stack-based navigation manager
type StackNavigator struct {
	widget.BaseWidget
	stack  []fyne.CanvasObject
	titles []string
	OnBack func()
}

// NewNavigator creates a new Navigator instance.
func NewNavigator(initial fyne.CanvasObject) *StackNavigator {
	return &StackNavigator{stack: []fyne.CanvasObject{initial}, titles: []string{""}}
}

// Push adds a new page to the stack and displays it.
func (n *StackNavigator) Push(page fyne.CanvasObject, title string) {
	n.stack = append(n.stack, page)
	n.titles = append(n.titles, title)
	n.Refresh()
}

// Pop removes the current page and returns to the previous one.
func (n *StackNavigator) Pop() {
	if len(n.stack) <= 1 {
		return // Prevent popping the last page.
	}

	n.stack[len(n.stack)-1] = nil
	n.stack = n.stack[:len(n.stack)-1]
	n.titles = n.titles[:len(n.titles)-1]

	n.Refresh()
}

func (n *StackNavigator) MinSize() fyne.Size {
	n.ExtendBaseWidget(n)
	return n.BaseWidget.MinSize()
}

// CreateRenderer creats the stackNavigatorRenderer.
func (n *StackNavigator) CreateRenderer() fyne.WidgetRenderer {
	renderer := &stackNavigatorRenderer{
		parent: n,
		backButton: widget.Button{
			Icon:       theme.NavigateBackIcon(),
			Text:       "Go back",
			Importance: widget.LowImportance,
			OnTapped:   n.OnBack,
		},
		titleLabel: widget.Label{
			Text:      n.titles[len(n.titles)-1],
			TextStyle: fyne.TextStyle{Bold: true},
			Alignment: fyne.TextAlignCenter,
		},
	}

	renderer.backButton.Hidden = len(n.stack) == 1
	renderer.titleLabel.Hidden = renderer.backButton.Hidden
	renderer.separator.Hidden = renderer.backButton.Hidden

	renderer.objects = []fyne.CanvasObject{&renderer.backButton, &renderer.titleLabel, &renderer.separator, n.stack[len(n.stack)-1]}
	return renderer
}

var _ fyne.WidgetRenderer = (*stackNavigatorRenderer)(nil)

type stackNavigatorRenderer struct {
	parent  *StackNavigator
	objects []fyne.CanvasObject

	backButton widget.Button
	titleLabel widget.Label
	separator  widget.Separator
}

func (r *stackNavigatorRenderer) Destroy() {
}

// Layout is a hook that is called if the widget needs to be laid out.
// This should never call [Refresh].
func (r *stackNavigatorRenderer) Layout(size fyne.Size) {
	contentStartsAt := float32(0)
	if len(r.parent.stack) > 1 {
		r.backButton.Move(fyne.Position{})
		buttonSize := r.backButton.MinSize()
		r.backButton.Resize(buttonSize)

		labelSize := r.titleLabel.MinSize()
		r.titleLabel.Move(fyne.NewPos((size.Width-labelSize.Width)/2, 0))
		r.titleLabel.Resize(labelSize)

		contentStartsAt = buttonSize.Height + theme.Padding()

		r.separator.Move(fyne.Position{Y: contentStartsAt})
		r.separator.Resize(fyne.NewSize(size.Width, theme.SeparatorThicknessSize()))

	}

	r.objects[3].Move(fyne.NewPos(0, contentStartsAt))
	r.objects[3].Resize(size.SubtractWidthHeight(0, contentStartsAt))
}

// MinSize returns the minimum size of the widget that is rendered by this renderer.
func (r *stackNavigatorRenderer) MinSize() fyne.Size {
	minSize := r.objects[3].MinSize()
	if len(r.parent.stack) > 1 {
		return minSize.AddWidthHeight(0, r.backButton.MinSize().Height+theme.Padding())
	}

	return minSize
}

// Objects returns all objects that should be drawn.
func (r *stackNavigatorRenderer) Objects() []fyne.CanvasObject {
	return r.objects
}

// Refresh is a hook that is called if the widget has updated and needs to be redrawn.
// This might trigger a [Layout].
func (r *stackNavigatorRenderer) Refresh() {
	r.titleLabel.Text = r.parent.titles[len(r.parent.titles)-1]
	r.titleLabel.Hidden = len(r.parent.stack) == 1
	r.titleLabel.Refresh()

	r.backButton.Hidden = r.titleLabel.Hidden
	r.backButton.Refresh()

	r.separator.Hidden = r.titleLabel.Hidden
	r.separator.Refresh()

	r.objects[3] = r.parent.stack[len(r.parent.stack)-1]

	canvas.Refresh(r.parent)
}
