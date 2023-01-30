package bridge

import (
	"sync"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/storage"
	"fyne.io/fyne/v2/widget"
	"github.com/Jacalz/rymdport/v3/internal/transport"
)

// RecvItem is the item that is being received
type RecvItem struct {
	URI  fyne.URI
	Name string

	Value  int64
	Max    int64
	Status func() string

	list *RecvList
}

func (r *RecvItem) update(delta, total int64) {
	r.Value += delta
	r.Max = total
	r.list.Refresh()
}

func (r *RecvItem) done() {
	r.Value = r.Max
	r.list.Refresh()
}

func (r *RecvItem) failed() {
	r.Status = func() string { return "Failed" }
	r.list.Refresh()
}

// RecvList is a list of progress bars that track send progress.
type RecvList struct {
	widget.List

	client *transport.Client

	Items []*RecvItem
	lock  sync.RWMutex

	window fyne.Window
}

// Length returns the length of the data.
func (p *RecvList) Length() int {
	return len(p.Items)
}

// CreateItem creates a new item in the list.
func (p *RecvList) CreateItem() fyne.CanvasObject {
	return container.New(&listLayout{},
		&widget.FileIcon{URI: nil},
		&widget.Label{Text: "Waiting for filename...", Wrapping: fyne.TextTruncate},
		&widget.ProgressBar{},
	)
}

// UpdateItem updates the data in the list.
func (p *RecvList) UpdateItem(i int, item fyne.CanvasObject) {
	p.lock.RLock()
	defer p.lock.RUnlock()

	container := item.(*fyne.Container)
	container.Objects[0].(*widget.FileIcon).SetURI(p.Items[i].URI)
	container.Objects[1].(*widget.Label).SetText(p.Items[i].Name)

	progress := container.Objects[2].(*widget.ProgressBar)
	progress.Max = float64(p.Items[i].Max)
	progress.Value = float64(p.Items[i].Value)
	progress.TextFormatter = p.Items[i].Status
	progress.Refresh()
}

// NewRecvItem creates a new send item and adds it to the items.
func (p *RecvList) NewRecv() *RecvItem {
	p.lock.Lock()
	defer p.lock.Unlock()

	item := &RecvItem{Name: "Waiting for filename...", Max: 1, list: p}
	p.Items = append(p.Items, item)
	return item
}

// NewReceive adds data about a new send to the list and then returns the channel to update the code.
func (p *RecvList) NewReceive(code string) {
	item := p.NewRecv()
	p.Refresh()

	path := make(chan string)

	go func() {
		name := <-path
		item.URI = storage.NewFileURI(name)
		if name != "text" {
			item.Name = item.URI.Name()
		} else {
			item.Name = "Text Snippet"
		}

		close(path)
		p.Refresh()
	}()

	go func(code string) {
		if err := p.client.NewReceive(code, path, item.update); err != nil {
			p.client.ShowNotification("Receive failed", "An error occurred when receiving the data.")
			item.failed()
			dialog.ShowError(err, p.window)
		} else if item.Name != "Text Snippet" {
			p.client.ShowNotification("Receive completed", "The contents were saved to "+item.URI.Path()+".")
			item.done()
		} else {
			p.client.ShowNotification("Receive completed", "The text was received successfully.")
			item.done()
		}
	}(code)
}

// NewRecvList greates a list of progress bars.
func NewRecvList(window fyne.Window, client *transport.Client) *RecvList {
	p := &RecvList{client: client, window: window}
	p.List.Length = p.Length
	p.List.CreateItem = p.CreateItem
	p.List.UpdateItem = p.UpdateItem
	p.List.OnSelected = p.Unselect
	p.ExtendBaseWidget(p)

	return p
}
