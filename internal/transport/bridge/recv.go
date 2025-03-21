package bridge

import (
	"path/filepath"
	"slices"
	"sync"
	"sync/atomic"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/storage"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/Jacalz/rymdport/v3/internal/transport"
	"github.com/Jacalz/rymdport/v3/internal/util"
	"github.com/rymdport/wormhole/wormhole"
)

// RecvItem is the item that is being received
type RecvItem struct {
	URI  fyne.URI
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

	items      []*RecvItem
	info       recvInfoDialog
	textWindow textRecvWindow

	deleting atomic.Bool
	lock     sync.Mutex
	list     *widget.List
}

// NewRecvList greates a list of progress bars.
func (d *RecvData) NewRecvList() *widget.List {
	d.list = &widget.List{
		Length:     d.Length,
		CreateItem: d.CreateItem,
		UpdateItem: d.UpdateItem,
		OnSelected: d.OnSelected,
	}
	d.setUpInfoDialog()
	return d.list
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
	container.Objects[1].(*widget.Label).SetText(item.URI.Name())
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

func (d *RecvData) NewFailedRecv(code string) {
	d.lock.Lock()
	item := &RecvItem{URI: storage.NewFileURI("Failed"), Code: code, Max: 1, refresh: d.refresh, index: len(d.items)}
	item.Status = func() string { return "failed" }
	d.items = append(d.items, item)
	d.lock.Unlock()

	d.list.Refresh()
}

// NewRecv creates a new item to be displayed and returns it plus the path to the file.
// A text receive will have a path of an empty string.
func (d *RecvData) NewRecv(code string, msg *wormhole.IncomingMessage) (item *RecvItem, path string) {
	d.lock.Lock()
	item = &RecvItem{Code: code, Max: 1, refresh: d.refresh, index: len(d.items)}

	if msg.Type == wormhole.TransferText {
		item.URI = storage.NewFileURI("Text Snippet")
	} else {
		path = filepath.Join(d.Client.DownloadPath, msg.Name)
		item.URI = storage.NewFileURI(path)

		if msg.Type == wormhole.TransferDirectory {
			item.URI = &folderURI{item.URI}
		}
	}

	d.items = append(d.items, item)
	d.lock.Unlock()

	d.list.Refresh()
	return item, path
}

// NewReceive adds data about a new send to the list and then returns the channel to update the code.
func (d *RecvData) NewReceive(code string) {
	go func(code string) {
		msg, err := d.Client.NewReceive(code)
		if err != nil {
			d.Client.ShowNotification("Receive failed", "An error occurred when receiving the data.")
			d.NewFailedRecv(code)
			dialog.ShowError(err, d.Window)
			return
		}

		item, path := d.NewRecv(code, msg)

		if msg.Type == wormhole.TransferText {
			d.showTextWindow(msg.ReadText())
			item.done()

			d.Client.ShowNotification("Receive completed", "The text was received successfully.")
			return
		}

		err = d.Client.SaveToDisk(msg, path, item.update)
		if err != nil {
			item.failed()
			d.Client.ShowNotification("Receive failed", "An error occurred when receiving the data.")
			dialog.ShowError(err, d.Window)
			return
		}

		d.Client.ShowNotification("Receive completed", "The contents were saved to "+filepath.Dir(item.URI.Path())+".")
		item.done()
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

	d.items = slices.Delete(d.items, index, index+1)

	// Update the moved items to have the correct index.
	for j := index; j < len(d.items); j++ {
		d.items[j].index = j
	}

	// Refresh the whole list.
	d.list.Refresh()

	// Allow individual objects to be refreshed again.
	d.deleting.Store(false)
}

func (d *RecvData) setUpInfoDialog() {
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

type textRecvWindow struct {
	textLabel              *widget.Label
	copyButton, saveButton *widget.Button
	window                 fyne.Window
	fileSaveDialog         *dialog.FileDialog
}

func (r *textRecvWindow) copy() {
	r.window.Clipboard().SetContent(r.textLabel.Text)
}

func (r *textRecvWindow) interceptClose() {
	r.window.Hide()
	r.textLabel.SetText("")
}

func (r *textRecvWindow) saveFileToDisk(file fyne.URIWriteCloser, err error) {
	if err != nil {
		fyne.LogError("Error on selecting file to write to", err)
		dialog.ShowError(err, r.window)
		return
	} else if file == nil {
		return
	}

	if _, err := file.Write([]byte(r.textLabel.Text)); err != nil {
		fyne.LogError("Error on writing text to the file", err)
		dialog.ShowError(err, r.window)
	}

	if err := file.Close(); err != nil {
		fyne.LogError("Error on closing text file", err)
		dialog.ShowError(err, r.window)
	}
}

func (r *textRecvWindow) save() {
	now := time.Now().Format("2006-01-02T15:04") // TODO: Might want to use AppendFormat and strings.Builder
	r.fileSaveDialog.SetFileName("received-" + now + ".txt")
	r.fileSaveDialog.Resize(util.WindowSizeToDialog(r.window.Canvas().Size()))
	r.fileSaveDialog.Show()
}

func (d *RecvData) createTextWindow() {
	window := d.Client.App.NewWindow("Received Text")
	window.SetCloseIntercept(d.textWindow.interceptClose)

	d.textWindow = textRecvWindow{
		window:         window,
		textLabel:      &widget.Label{},
		copyButton:     &widget.Button{Text: "Copy", Icon: theme.ContentCopyIcon(), OnTapped: d.textWindow.copy},
		saveButton:     &widget.Button{Text: "Save", Icon: theme.DocumentSaveIcon(), OnTapped: d.textWindow.save},
		fileSaveDialog: dialog.NewFileSave(d.textWindow.saveFileToDisk, window),
	}

	actionContainer := container.NewGridWithColumns(2, d.textWindow.copyButton, d.textWindow.saveButton)
	window.SetContent(container.NewBorder(nil, actionContainer, nil, nil, container.NewScroll(d.textWindow.textLabel)))
	window.Resize(fyne.NewSize(400, 300))
}

// showTextWindow handles the creation of a window for displaying text content.
func (d *RecvData) showTextWindow(received string) {
	if d.textWindow.window == nil {
		d.createTextWindow()
	}

	d.textWindow.textLabel.SetText(received)

	win := d.textWindow.window
	win.Show()
	win.RequestFocus()
}

var _ fyne.URIWithIcon = (*folderURI)(nil)

// folderURI is used to force folders to use an icon even before the actual folder is created.
// This avoids issues with the icon initially appearing as a file.
type folderURI struct {
	fyne.URI
}

func (c *folderURI) Icon() fyne.Resource {
	return theme.FolderIcon()
}
