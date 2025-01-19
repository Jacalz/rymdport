package components

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/widget"

	"github.com/Jacalz/rymdport/v3/completion"
)

// NewCompletionEntry returns a new entry widget that allows completing on tab.
func NewCompletionEntry(driver desktop.Driver, generator func(string) []string) *CompletionEntry {
	entry := &CompletionEntry{driver: driver}
	entry.completer.Generate = generator
	entry.ExtendBaseWidget(entry)
	return entry
}

// CompletionEntry allows using tab and shift+tab to
// move forwards and backwards from a set of completions.
type CompletionEntry struct {
	widget.Entry
	driver    desktop.Driver
	completer completion.TabCompleter
}

// AcceptsTab overrides tab handling to allow tabs as input.
func (c *CompletionEntry) AcceptsTab() bool {
	return true
}

// TypedKey adapts the key inputs to handle tab completion.
func (c *CompletionEntry) TypedKey(key *fyne.KeyEvent) {
	switch key.Name {
	case desktop.KeyShiftLeft, desktop.KeyShiftRight:
	case fyne.KeyTab:
		completed := ""

		if c.driver.CurrentKeyModifiers()&fyne.KeyModifierShift != 0 {
			completed = c.completer.Previous(c.Text)
		} else {
			completed = c.completer.Next(c.Text)
		}

		c.CursorColumn = len(completed)
		c.SetText(completed)
	default:
		c.completer.Reset()
		c.Entry.TypedKey(key)
	}
}
