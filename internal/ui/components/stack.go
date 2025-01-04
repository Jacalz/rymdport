package components

import (
	"fyne.io/fyne/v2"
)

// StackNavigator represents a stack-based navigation manager
type StackNavigator struct {
	stack   []fyne.CanvasObject
	current int
	OnBack  func()
}

// NewNavigator creates a new Navigator instance.
func NewNavigator(initialPage fyne.CanvasObject) *StackNavigator {
	return &StackNavigator{stack: []fyne.CanvasObject{initialPage}}
}

// Next moves the view to the next view in the stack without adding contents.
// This allows viewes to move forwards through views without recreating each time.
func (n *StackNavigator) Next() {
	if n.current == len(n.stack)-1 {
		return
	}

	n.current++
}

// Previous moves the view to the previous view in the stack without removing contents.
// This allows viewes to move backwards through views without recreating each time.
func (n *StackNavigator) Previous() {
	if n.current == 0 || len(n.stack) < 1 {

	}

	n.current--
}

// Push adds a new page to the stack and displays it.
func (n *StackNavigator) Push(page fyne.CanvasObject) {
	n.stack = append(n.stack, page)
	n.Next()
}

// Pop removes the current page and returns to the previous one.
func (n *StackNavigator) Pop() {
	if len(n.stack) <= 1 {
		return // Prevent popping the last page
	}

	n.stack[len(n.stack)-1] = nil
	n.stack = n.stack[:len(n.stack)-1]
	n.Previous()
}
