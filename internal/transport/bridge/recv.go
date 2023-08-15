package bridge

import (
	"path/filepath"

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

	list *widget.List
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

// RecvData is a list of progress bars that track send progress.
type RecvData struct {
	Client *transport.Client
	Window fyne.Window

	items []*RecvItem
	list  *widget.List
}

// Length returns the length of the data.
func (d *RecvData) Length() int {
	return len(d.items)
}

// CreateItem creates a new item in the list.
func (d *RecvData) CreateItem() fyne.CanvasObject {
	return container.New(&listLayout{},
		&widget.FileIcon{},
		&widget.Label{Text: "Waiting for filename...", Truncation: fyne.TextTruncateEllipsis},
		&widget.Label{Text: "Waiting for code..."},
		&widget.ProgressBar{},
	)
}

// UpdateItem updates the data in the list.
func (d *RecvData) UpdateItem(i int, item fyne.CanvasObject) {
	container := item.(*fyne.Container)
	container.Objects[0].(*widget.FileIcon).SetURI(d.items[i].URI)
	container.Objects[1].(*widget.Label).SetText(d.items[i].Name)
	container.Objects[2].(*widget.Label).SetText(d.items[i].Code)

	progress := container.Objects[3].(*widget.ProgressBar)
	progress.Max = float64(d.items[i].Max)
	progress.Value = float64(d.items[i].Value)
	progress.TextFormatter = d.items[i].Status
	progress.Refresh()
}

// OnSelected currently just makes sure that we don't persist selection.
func (d *RecvData) OnSelected(i int) {
	d.list.Unselect(i)

	removeLabel := &widget.Label{Text: "This item has completed the transfer and can be removed."}
	removeButton := &widget.Button{Icon: theme.DeleteIcon(), Importance: widget.DangerImportance, Text: "Remove", OnTapped: func() {
		if i < len(d.items)-1 {
			copy(d.items[i:], d.items[i+1:])
		}

		d.items[len(d.items)-1] = nil // Allow the GC to reclaim memory.
		d.items = d.items[:len(d.items)-1]

		d.list.Refresh()
	}}

	removeCard := &widget.Card{Content: container.NewVBox(removeLabel, removeButton)}

	// Only allow failed or completed items to be removed.
	if d.items[i].Value < d.items[i].Max && d.items[i].Status == nil {
		removeLabel.Text = "This item can not be removed yet. The transfer needs to complete first."
		removeButton.Disable()
	}

	dialog.ShowCustom("Information", "Close", removeCard, d.Window)
}

// NewRecvItem creates a new send item and adds it to the items.
func (d *RecvData) NewRecv(code string) *RecvItem {
	item := &RecvItem{Name: "Waiting for filename...", Code: code, Max: 1, list: d.list}
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
		d.list.Refresh()
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
