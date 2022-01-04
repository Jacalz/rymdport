package bridge

import (
	"sync"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/storage"
	"fyne.io/fyne/v2/widget"
	"github.com/Jacalz/wormhole-gui/v2/internal/transport"
	"github.com/Jacalz/wormhole-gui/v2/internal/util"
)

// SendItem is the item that is being sent.
type SendItem struct {
	URI      fyne.URI
	Progress *util.ProgressBar
	Code     string
	Name     string
}

// SendList is a list of progress bars that track send progress.
type SendList struct {
	widget.List

	client *transport.Client

	Items []*SendItem
	lock  sync.RWMutex

	window fyne.Window
}

// Length returns the length of the data.
func (p *SendList) Length() int {
	return len(p.Items)
}

// CreateItem creates a new item in the list.
func (p *SendList) CreateItem() fyne.CanvasObject {
	return container.New(&listLayout{},
		widget.NewFileIcon(nil),
		&widget.Label{Text: "Waiting for filename...", Wrapping: fyne.TextTruncate},
		newCodeDisplay(p.window),
		util.NewProgressBar(),
	)
}

// UpdateItem updates the data in the list.
func (p *SendList) UpdateItem(i int, item fyne.CanvasObject) {
	p.lock.RLock()
	defer p.lock.RUnlock()

	item.(*fyne.Container).Objects[0].(*widget.FileIcon).SetURI(p.Items[i].URI)
	item.(*fyne.Container).Objects[1].(*widget.Label).SetText(p.Items[i].Name)
	item.(*fyne.Container).Objects[2].(*fyne.Container).Objects[0].(*codeDisplay).SetText(p.Items[i].Code)
	p.Items[i].Progress = item.(*fyne.Container).Objects[3].(*util.ProgressBar)
}

// RemoveItem removes the item at the specified index.
func (p *SendList) RemoveItem(i int) {
	p.lock.Lock()
	defer p.lock.Unlock()

	copy(p.Items[i:], p.Items[i+1:])
	p.Items[p.Length()-1] = nil // Make sure that GC run on removed element.
	p.Items = p.Items[:p.Length()-1]
}

// OnSelected handles removing items and stopping send (in the future).
func (p *SendList) OnSelected(i int) {
	if p.Items[i].Progress.Value != p.Items[i].Progress.Max { // TODO: Stop the send instead.
		p.Unselect(i)
		return // We can't stop running sends due to bug in wormhole-gui.
	}

	dialog.ShowConfirm("Remove from list", "Do you wish to remove the item from the list?", func(remove bool) {
		if remove {
			p.RemoveItem(i)
			p.Refresh()
		}
	}, p.window)

	p.Unselect(i)
}

// NewSendItem adds data about a new send to the list and then returns the item.
func (p *SendList) NewSendItem(name string, uri fyne.URI) *SendItem {
	p.lock.Lock()
	defer p.lock.Unlock()

	item := &SendItem{Name: name, Code: "Waiting for code...", URI: uri}
	p.Items = append(p.Items, item)
	return item
}

// OnFileSelect is intended to be passed as callback to a FileOpen dialog.
func (p *SendList) OnFileSelect(file fyne.URIReadCloser, err error) {
	if err != nil {
		fyne.LogError("Error on selecting file to send", err)
		dialog.ShowError(err, p.window)
		return
	} else if file == nil {
		return
	}

	item := p.NewSendItem(file.URI().Name(), file.URI())
	p.Refresh()

	go func() {
		// We want to catch close errors for security reasons.
		defer func() {
			if err = file.Close(); err != nil {
				item.Progress.Failed()
				fyne.LogError("Error on closing file", err)
			}
		}()

		code, result, err := p.client.NewFileSend(file, item.Progress.WithProgress())
		if err != nil {
			fyne.LogError("Error on sending file", err)
			item.Progress.Failed()
			dialog.ShowError(err, p.window)
			return
		}

		item.Code = code
		p.Refresh()

		if res := <-result; res.Error != nil {
			fyne.LogError("Error on sending file", res.Error)
			item.Progress.Failed()
			dialog.ShowError(res.Error, p.window)
			p.client.ShowNotification("File send failed", "An error occurred when sending the file.")
		} else if res.OK {
			p.client.ShowNotification("File send completed", "The file was sent successfully.")
		}
	}()
}

// OnDirSelect is intended to be passed as callback to a FolderOpen dialog.
func (p *SendList) OnDirSelect(dir fyne.ListableURI, err error) {
	if err != nil {
		fyne.LogError("Error on selecting dir to send", err)
		dialog.ShowError(err, p.window)
		return
	} else if dir == nil {
		return
	}

	item := p.NewSendItem(dir.Name(), dir)
	p.Refresh()

	go func() {
		code, result, err := p.client.NewDirSend(dir, item.Progress.WithProgress())
		if err != nil {
			fyne.LogError("Error on sending directory", err)
			item.Progress.Failed()
			dialog.ShowError(err, p.window)
			return
		}

		item.Code = code
		p.Refresh()

		if res := <-result; res.Error != nil {
			fyne.LogError("Error on sending directory", res.Error)
			item.Progress.Failed()
			dialog.ShowError(res.Error, p.window)
			p.client.ShowNotification("Directory send failed", "An error occurred when sending the directory.")
		} else if res.OK {
			p.client.ShowNotification("Directory send completed", "The directory was sent successfully.")
		}
	}()
}

// SendText sends new text.
func (p *SendList) SendText() {
	// The file URI is a hack to get the correct icon.
	item := &SendItem{Name: "Text Snippet", Code: "Waiting for code...", URI: storage.NewFileURI("text")}

	go func() {
		text := <-p.client.ShowTextSendWindow()
		if text == "" {
			return
		}

		p.Items = append(p.Items, item)
		p.Refresh()

		code, result, err := p.client.NewTextSend(text, item.Progress.WithProgress())
		if err != nil {
			fyne.LogError("Error on sending text", err)
			item.Progress.Failed()
			dialog.ShowError(err, p.window)
			return
		}

		item.Code = code
		p.Refresh()

		if res := <-result; res.Error != nil {
			fyne.LogError("Error on sending text", res.Error)
			item.Progress.Failed()
			dialog.ShowError(res.Error, p.window)
			p.client.ShowNotification("Text send failed", "An error occurred when sending the text.")
		} else if res.OK && p.client.Notifications {
			p.client.ShowNotification("Text send completed", "The text was sent successfully.")
		}
	}()
}

// NewSendList greates a list of progress bars.
func NewSendList(window fyne.Window, client *transport.Client) *SendList {
	p := &SendList{client: client, window: window}
	p.List.Length = p.Length
	p.List.CreateItem = p.CreateItem
	p.List.UpdateItem = p.UpdateItem
	p.List.OnSelected = p.OnSelected
	p.ExtendBaseWidget(p)

	return p
}
