package bridge

import (
	"path/filepath"
	"sync/atomic"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/storage"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/Jacalz/rymdport/v3/internal/transport"
)

// RecvItem is the item that is being received
type RecvItem struct {
	URI  fyne.URI
	Name string
	Code string

	Value  int64
	Max    int64
	Status func() string

	// Allow the list to only refresh a single object.
	refresh func(int)
	index   int
}

func (r *RecvItem) update(delta, total int64) {
	r.Value += delta
	r.Max = total
	r.refresh(r.index)
}

func (r *RecvItem) done() {
	r.Value = r.Max
	r.refresh(r.index)
}

func (r *RecvItem) failed() {
	r.Status = func() string { return "Failed" }
	r.refresh(r.index)
}

// RecvData is a list of progress bars that track send progress.
type RecvData struct {
	Client *transport.Client
	Window fyne.Window

	items []*RecvItem

	deleting atomic.Bool
	list     *widget.List
}

// Length returns the length of the data.
func (d *RecvData) Length() int {
	return len(d.items)
}

// CreateItem creates a new item in the list.
func (d *RecvData) CreateItem() fyne.CanvasObject {
	return container.New(listLayout{},
		&widget.FileIcon{},
		&widget.Label{Text: "Waiting for filename...", Truncation: fyne.TextTruncateEllipsis},
		&widget.Label{Text: "Waiting for code...", Truncation: fyne.TextTruncateEllipsis},
		&widget.ProgressBar{},
	)
}

// UpdateItem updates the data in the list.
func (d *RecvData) UpdateItem(i int, object fyne.CanvasObject) {
	item := d.items[i]
	container := object.(*fyne.Container)
	container.Objects[0].(*widget.FileIcon).SetURI(item.URI)
	container.Objects[1].(*widget.Label).SetText(item.Name)
	container.Objects[2].(*widget.Label).SetText(item.Code)

	progress := container.Objects[3].(*widget.ProgressBar)
	progress.Max = float64(item.Max)
	progress.Value = float64(item.Value)
	progress.TextFormatter = item.Status
	progress.Refresh()
}

// OnSelected currently just makes sure that we don't persist selection.
func (d *RecvData) OnSelected(i int) {
	d.list.Unselect(i)

	removeLabel := &widget.Label{Text: "This item has completed the transfer and can be removed."}
	removeButton := &widget.Button{Icon: theme.DeleteIcon(), Importance: widget.DangerImportance, Text: "Remove", OnTapped: func() {
		// Make sure that no updates happen while we modify the slice.
		d.deleting.Store(true)

		if i < len(d.items)-1 {
			copy(d.items[i:], d.items[i+1:])
		}

		d.items[len(d.items)-1] = nil // Allow the GC to reclaim memory.
		d.items = d.items[:len(d.items)-1]

		// Update the moved items to have the correct index.
		for j := i; j < len(d.items); j++ {
			d.items[j].index = j
		}

		// Refresh the whole list.
		d.list.Refresh()

		// Allow individual objects to be refreshed again.
		d.deleting.Store(false)
	}}

	removeCard := &widget.Card{Content: container.NewVBox(removeLabel, removeButton)}

	// Only allow failed or completed items to be removed.
	if d.items[i].Value < d.items[i].Max && d.items[i].Status == nil {
		removeLabel.Text = "This item can not be removed yet. The transfer needs to complete first."
		removeButton.Disable()
	}

	dialog.ShowCustom("Information", "Close", removeCard, d.Window)
}

// NewRecv creates a new send item and adds it to the items.
func (d *RecvData) NewRecv(code string) *RecvItem {
	item := &RecvItem{Name: "Waiting for filename...", Code: code, Max: 1, refresh: d.refresh, index: len(d.items)}
	d.items = append(d.items, item)
	return item
}

// NewReceive adds data about a new send to the list and then returns the channel to update the code.
func (d *RecvData) NewReceive(code string) {
	item := d.NewRecv(code)
	d.list.Refresh()

	path := make(chan string)

	go func() {
		item.URI = storage.NewFileURI(<-path)
		item.Name = item.URI.Name()
		close(path)
		d.refresh(item.index)
	}()

	go func(code string) {
		if err := d.Client.NewReceive(code, path, item.update); err != nil {
			d.Client.ShowNotification("Receive failed", "An error occurred when receiving the data.")
			item.failed()
			dialog.ShowError(err, d.Window)
		} else if item.Name != "Text Snippet" {
			d.Client.ShowNotification("Receive completed", "The contents were saved to "+filepath.Dir(item.URI.Path())+".")
			item.done()
		} else {
			d.Client.ShowNotification("Receive completed", "The text was received successfully.")
			item.done()
		}
	}(code)
}

func (d *RecvData) refresh(index int) {
	if d.deleting.Load() {
		return // Don't update if we are deleting.
	}

	d.list.RefreshItem(index)
}

// NewRecvList greates a list of progress bars.
func NewRecvList(data *RecvData) *widget.List {
	list := &widget.List{
		Length:     data.Length,
		CreateItem: data.CreateItem,
		UpdateItem: data.UpdateItem,
		OnSelected: data.OnSelected,
	}
	data.list = list
	return list
}
