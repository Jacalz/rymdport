package widgets

import (
	"regexp"
	"time"

	"fyne.io/fyne"
	"fyne.io/fyne/container"
	"fyne.io/fyne/theme"
	"fyne.io/fyne/widget"
)

var codeValidator = regexp.MustCompile(`^\d\d?(-\w{2,12}){2,6}$`)

// NewBoldLabel returns a new label with bold text.
func NewBoldLabel(text string) *widget.Label {
	return widget.NewLabelWithStyle(text, fyne.TextAlignLeading, fyne.TextStyle{Bold: true})
}

// CodeDisplay is a label extended to copy the code with a menu popup on rightclick.
type CodeDisplay struct {
	widget.Label
	button *widget.Button
}

func (c *CodeDisplay) copyOnPress() {
	clipboard := fyne.CurrentApp().Driver().AllWindows()[0].Clipboard()

	if codeValidator.MatchString(c.Text) {
		c.button.SetIcon(theme.ConfirmIcon())
		clipboard.SetContent(c.Text)
	} else {
		c.button.SetIcon(theme.CancelIcon())
	}

	time.Sleep(500 * time.Millisecond)
	c.button.SetIcon(theme.ContentCopyIcon())
}

func newCodeDisplay() *fyne.Container {
	c := &CodeDisplay{button: &widget.Button{Icon: theme.ContentCopyIcon(), Importance: widget.LowImportance}}
	c.ExtendBaseWidget(c)

	c.SetText("Waiting for code...")
	c.button.OnTapped = c.copyOnPress

	return container.NewHBox(c, c.button)
}
