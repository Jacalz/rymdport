package widgets

import (
	"fyne.io/fyne"
	"fyne.io/fyne/dialog"
	"fyne.io/fyne/storage"
	"fyne.io/fyne/widget"
)

var emptyRecvItem = &RecvItem{}

// RecvItem is the item that is being received
type RecvItem struct {
	Filename chan string
	Status   chan string
}

// RecvList is a list of progress bars that track send progress.
type RecvList struct {
	widget.List

	Items []RecvItem
}

// Length returns the length of the data.
func (p *RecvList) Length() int {
	return len(p.Items)
}

// CreateItem creates a new item in the list.
func (p *RecvList) CreateItem() fyne.CanvasObject {
	return fyne.NewContainerWithLayout(&listLayout{}, widget.NewFileIcon(nil), widget.NewLabel("Waiting for filename..."), widget.NewLabel("Checking status..."))
}

// UpdateItem updates the data in the list.
func (p *RecvList) UpdateItem(i int, item fyne.CanvasObject) {
	go func(filename, status chan string) { // Get the channel out of the scope to avoid stalling render thread
		file, stat := <-filename, <-status
		item.(*fyne.Container).Objects[0].(*widget.FileIcon).SetURI(storage.NewURI(file))
		item.(*fyne.Container).Objects[1].(*widget.Label).SetText(file)
		item.(*fyne.Container).Objects[2].(*widget.Label).SetText(stat)
	}(p.Items[i].Filename, p.Items[i].Status)
}

// OnSelectionChanged handles removing items and stopping send (in the future)
func (p *RecvList) OnSelectionChanged(i int) {
	dialog.ShowConfirm("Remove from list", "Do you wish to remove the item from the list?", func(remove bool) {
		if remove {
			// Make sure that GC run on removed element
			copy(p.Items[i:], p.Items[i+1:])
			p.Items[p.Length()-1] = *emptyRecvItem
			p.Items = p.Items[:p.Length()-1]

			p.Refresh()
		}
	}, fyne.CurrentApp().Driver().AllWindows()[0])
}

// NewRecvItem adds data about a new send to the list and then returns the channel to update the code.
func (p *RecvList) NewRecvItem() (file, status chan string) {
	p.Items = append(p.Items, RecvItem{Filename: make(chan string), Status: make(chan string)})
	p.Refresh()

	return p.Items[p.Length()-1].Filename, p.Items[p.Length()-1].Status
}

// NewRecvList greates a list of progress bars.
func NewRecvList() *RecvList {
	p := &RecvList{}
	p.List.Length = p.Length
	p.List.CreateItem = p.CreateItem
	p.List.UpdateItem = p.UpdateItem
	p.List.OnSelectionChanged = p.OnSelectionChanged
	p.ExtendBaseWidget(p)

	return p
}
