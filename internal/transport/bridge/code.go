package bridge

import (
	"regexp"
	"time"

	"fyne.io/fyne"
	"fyne.io/fyne/container"
	"fyne.io/fyne/theme"
	"fyne.io/fyne/widget"
)

var codeValidator = regexp.MustCompile(`^\d\d?(-\w{2,12}){2,6}$`)

// codeDisplay is a label extended to copy the code with a menu popup on rightclick.
type codeDisplay struct {
	widget.Label
	button *widget.Button
}

func (c *codeDisplay) copyOnPress() {
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
	c := &codeDisplay{button: &widget.Button{Icon: theme.ContentCopyIcon(), Importance: widget.LowImportance}}
	c.ExtendBaseWidget(c)

	c.Text = "Waiting for code..."
	c.button.OnTapped = c.copyOnPress

	return container.NewHBox(c, c.button)
}
