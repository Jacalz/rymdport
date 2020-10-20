package widgets

import (
	"fyne.io/fyne"
	"fyne.io/fyne/dialog"
	"fyne.io/fyne/widget"
)

var emptySendItem = &SendItem{}

// SendItem is the item that is being sent.
type SendItem struct {
	Progress *SendProgress
	Code     chan string
	URI      fyne.URI
}

// SendList is a list of progress bars that track send progress.
type SendList struct {
	widget.List

	Items []SendItem
}

// Length returns the length of the data.
func (p *SendList) Length() int {
	return len(p.Items)
}

// CreateItem creates a new item in the list.
func (p *SendList) CreateItem() fyne.CanvasObject {
	return fyne.NewContainerWithLayout(&listLayout{}, widget.NewFileIcon(nil), widget.NewLabel("Waiting for filename..."), newCodeDisplay(), NewSendProgress())
}

// UpdateItem updates the data in the list.
func (p *SendList) UpdateItem(i int, item fyne.CanvasObject) {
	item.(*fyne.Container).Objects[0].(*widget.FileIcon).SetURI(p.Items[i].URI)
	item.(*fyne.Container).Objects[1].(*widget.Label).SetText(p.Items[i].URI.Name())
	item.(*fyne.Container).Objects[2].(*fyne.Container).Objects[0].(*CodeDisplay).waitForCode(p.Items[i].Code)
	p.Items[i].Progress = item.(*fyne.Container).Objects[3].(*SendProgress)
}

// OnSelectionChanged handles removing items and stopping send (in the future)
func (p *SendList) OnSelectionChanged(i int) {
	if p.Items[i].Progress.Value != p.Items[i].Progress.Max { // TODO: Stop the send instead.
		return // We can't stop running sends due to bug in wormhole-gui.
	}

	dialog.ShowConfirm("Remove from list", "Do you wish to remove the item from the list?", func(remove bool) {
		if remove {
			// Make sure that GC run on removed element
			copy(p.Items[i:], p.Items[i+1:])
			p.Items[p.Length()-1] = *emptySendItem
			p.Items = p.Items[:p.Length()-1]

			p.Refresh()
		}
	}, fyne.CurrentApp().Driver().AllWindows()[0])
}

// NewSendItem adds data about a new send to the list and then returns the channel to update the code.
func (p *SendList) NewSendItem(URI fyne.URI) chan string {
	p.Items = append(p.Items, SendItem{Progress: NewSendProgress(), URI: URI, Code: make(chan string)})
	p.Refresh()

	return p.Items[p.Length()-1].Code
}

// NewSendList greates a list of progress bars.
func NewSendList() *SendList {
	p := &SendList{}
	p.List.Length = p.Length
	p.List.CreateItem = p.CreateItem
	p.List.UpdateItem = p.UpdateItem
	p.List.OnSelectionChanged = p.OnSelectionChanged
	p.ExtendBaseWidget(p)

	return p
}
