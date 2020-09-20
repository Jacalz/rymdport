package widgets

import (
	"fyne.io/fyne"
	"fyne.io/fyne/widget"
)

// SendItem is the item that is being sent.
type SendItem struct {
	Progress *SendProgress
	Code     chan string
	URI      fyne.URI
}

// ProgressList is a list of progress bars that track send progress.
type ProgressList struct {
	widget.List

	Items []SendItem
}

// Length returns the length of the data.
func (p *ProgressList) Length() int {
	return len(p.Items)
}

// CreateItem creates a new item in the list.
func (p *ProgressList) CreateItem() fyne.CanvasObject {
	return fyne.NewContainerWithLayout(newSendLayout(), widget.NewIcon(nil), widget.NewLabel("Waiting for filename..."), newCodeDisplay(), NewSendProgress())
}

// UpdateItem updates the data in the list.
func (p *ProgressList) UpdateItem(i int, item fyne.CanvasObject) {
	item.(*fyne.Container).Objects[0].(*widget.Icon).SetResource(iconFromURI(p.Items[i].URI))
	item.(*fyne.Container).Objects[1].(*widget.Label).SetText(p.Items[i].URI.Name())
	item.(*fyne.Container).Objects[2].(*fyne.Container).Objects[0].(*CodeDisplay).waitForCode(p.Items[i].Code)
	p.Items[i].Progress = item.(*fyne.Container).Objects[3].(*SendProgress)
}

// NewSendItem adds data about a new send to the list and then returns the channel to update the code.
func (p *ProgressList) NewSendItem(URI fyne.URI) chan string {
	p.Items = append(p.Items, SendItem{Progress: NewSendProgress(), URI: URI, Code: make(chan string)})
	p.Refresh()

	return p.Items[p.Length()-1].Code
}

// NewProgressList greates a list of progress bars.
func NewProgressList() *ProgressList {
	p := &ProgressList{}
	p.List.Length = p.Length
	p.List.CreateItem = p.CreateItem
	p.List.UpdateItem = p.UpdateItem
	p.ExtendBaseWidget(p)

	return p
}
