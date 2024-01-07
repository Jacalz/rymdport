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

func (r *RecvItem) setPath(path string) {
	r.URI = storage.NewFileURI(path)
	r.Name = r.URI.Name()
	r.refresh(r.index)
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
	data.setUpRecvInfoDialog()
	return list
}

// RecvData is a list of progress bars that track send progress.
type RecvData struct {
	Client *transport.Client
	Window fyne.Window

	items []*RecvItem
	info  recvInfoDialog

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

	d.info.button.OnTapped = func() {
		d.remove(i)
		d.info.dialog.Hide()
	}

	if d.info.button.Disabled() {
		d.info.label.Text = "This item can be removed.\nThe transfer has completed."
		d.info.button.Enable()
	}

	// Only allow failed or completed items to be removed.
	item := d.items[i]
	if item.Value < item.Max && item.Status == nil {
		d.info.label.Text = "This item can't be removed yet.\nThe transfer needs to complete first."
		d.info.button.Disable()
	}

	d.info.dialog.Show()
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

	go func(code string) {
		if err := d.Client.NewReceive(code, item.setPath, item.update); err != nil {
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

func (d *RecvData) remove(index int) {
	// Make sure that no updates happen while we modify the slice.
	d.deleting.Store(true)

	if index < len(d.items)-1 {
		copy(d.items[index:], d.items[index+1:])
	}

	d.items[len(d.items)-1] = nil // Allow the GC to reclaim memory.
	d.items = d.items[:len(d.items)-1]

	// Update the moved items to have the correct index.
	for j := index; j < len(d.items); j++ {
		d.items[j].index = j
	}

	// Refresh the whole list.
	d.list.Refresh()

	// Allow individual objects to be refreshed again.
	d.deleting.Store(false)
}

func (d *RecvData) setUpRecvInfoDialog() {
	d.info.label = &widget.Label{Text: "This item can be removed.\nThe transfer has completed."}
	d.info.button = &widget.Button{Icon: theme.DeleteIcon(), Importance: widget.DangerImportance, Text: "Remove"}
	removeCard := &widget.Card{Content: container.NewVBox(d.info.label, d.info.button)}
	d.info.dialog = dialog.NewCustom("Information", "Close", removeCard, d.Window)
}

type recvInfoDialog struct {
	dialog *dialog.CustomDialog
	button *widget.Button
	label  *widget.Label
}
