package bridge

import (
	"fyne.io/fyne"
	"fyne.io/fyne/container"
	"fyne.io/fyne/dialog"
	"fyne.io/fyne/widget"
	"github.com/Jacalz/wormhole-gui/internal/transport"
)

// SendItem is the item that is being sent.
type SendItem struct {
	URI      fyne.URI
	Progress *sendProgress
	Code     string
	Name     string
}

// SendList is a list of progress bars that track send progress.
type SendList struct {
	widget.List

	client *transport.Client

	Items []*SendItem
}

// Length returns the length of the data.
func (p *SendList) Length() int {
	return len(p.Items)
}

// CreateItem creates a new item in the list.
func (p *SendList) CreateItem() fyne.CanvasObject {
	return container.New(&listLayout{}, widget.NewFileIcon(nil), widget.NewLabel("Waiting for filename..."), newCodeDisplay(), newSendProgress())
}

// UpdateItem updates the data in the list.
func (p *SendList) UpdateItem(i int, item fyne.CanvasObject) {
	item.(*fyne.Container).Objects[0].(*widget.FileIcon).SetURI(p.Items[i].URI)
	item.(*fyne.Container).Objects[1].(*widget.Label).SetText(p.Items[i].Name)
	item.(*fyne.Container).Objects[2].(*fyne.Container).Objects[0].(*codeDisplay).SetText(p.Items[i].Code)
	p.Items[i].Progress = item.(*fyne.Container).Objects[3].(*sendProgress)
}

// RemoveItem removes the item at the specified index.
func (p *SendList) RemoveItem(i int) {
	copy(p.Items[i:], p.Items[i+1:])
	p.Items[p.Length()-1] = nil // Make sure that GC run on removed element
	p.Items = p.Items[:p.Length()-1]
	p.Refresh()
}

// OnSelected handles removing items and stopping send (in the future)
func (p *SendList) OnSelected(i int) {
	if p.Items[i].Progress.Value != p.Items[i].Progress.Max { // TODO: Stop the send instead.
		return // We can't stop running sends due to bug in wormhole-gui.
	}

	dialog.ShowConfirm("Remove from list", "Do you wish to remove the item from the list?", func(remove bool) {
		if remove {
			p.RemoveItem(i)
			p.Refresh()
		}
	}, fyne.CurrentApp().Driver().AllWindows()[0])

	p.Unselect(i)
}

// NewSendItem adds data about a new send to the list and then returns the channel to update the code.
func (p *SendList) NewSendItem(name string, uri fyne.URI) {
	p.Items = append(p.Items, &SendItem{Name: name, Code: "Waiting for code...", URI: uri})
	p.Refresh()
}

// OnFileSelect is intended to be passed as callback to a FileOpen dialog.
func (p *SendList) OnFileSelect(file fyne.URIReadCloser, err error) {
	if err != nil {
		fyne.LogError("Error on selecting file to send", err)
		dialog.ShowError(err, fyne.CurrentApp().Driver().AllWindows()[0])
		return
	} else if file == nil {
		return
	}

	p.NewSendItem(file.URI().Name(), file.URI())

	go func(i int) {
		code, result, err := p.client.NewFileSend(file, p.Items[i].Progress.update)
		if err != nil {
			fyne.LogError("Error on sending file", err)
			dialog.ShowError(err, fyne.CurrentApp().Driver().AllWindows()[0])
			return
		}

		p.Items[i].Code = code
		p.Refresh()

		if res := <-result; res.Error != nil {
			fyne.LogError("Error on sending file", res.Error)
			dialog.ShowError(res.Error, fyne.CurrentApp().Driver().AllWindows()[0])
			p.client.ShowNotification("File send failed", "An error occurred when sending the file.")
		} else if res.OK {
			p.client.ShowNotification("File send completed", "The file was sent successfully.")
		}

		if err = file.Close(); err != nil {
			fyne.LogError("Error on closing file", err)
		}
	}(p.Length() - 1)
}

// OnDirSelect is intended to be passed as callback to a FolderOpen dialog.
func (p *SendList) OnDirSelect(dir fyne.ListableURI, err error) {
	if err != nil {
		fyne.LogError("Error on selecting dir to send", err)
		dialog.ShowError(err, fyne.CurrentApp().Driver().AllWindows()[0])
		return
	} else if dir == nil {
		return
	}

	p.NewSendItem(dir.Name(), dir)

	go func(i int) {
		code, result, err := p.client.NewDirSend(dir, p.Items[i].Progress.update)
		if err != nil {
			fyne.LogError("Error on sending directory", err)
			dialog.ShowError(err, fyne.CurrentApp().Driver().AllWindows()[0])
			return
		}

		p.Items[i].Code = code
		p.Refresh()

		if res := <-result; res.Error != nil {
			fyne.LogError("Error on sending directory", res.Error)
			dialog.ShowError(res.Error, fyne.CurrentApp().Driver().AllWindows()[0])
			p.client.ShowNotification("Directory send failed", "An error occurred when sending the directory.")
		} else if res.OK {
			fyne.CurrentApp().SendNotification(fyne.NewNotification("Directory send completed", "The directory was sent successfully."))
		}
	}(p.Length() - 1)
}

// SendText sends new text.
func (p *SendList) SendText() {
	if text := <-p.client.ShowTextSendWindow(); text != "" {
		p.NewSendItem("Text Snippet", nil)

		go func(i int) {
			code, result, err := p.client.NewTextSend(text, p.Items[i].Progress.update)
			if err != nil {
				fyne.LogError("Error on sending text", err)
				dialog.ShowError(err, fyne.CurrentApp().Driver().AllWindows()[0])
				return
			}

			p.Items[i].Code = code
			p.Refresh()

			if res := <-result; res.Error != nil {
				fyne.LogError("Error on sending text", res.Error)
				dialog.ShowError(res.Error, fyne.CurrentApp().Driver().AllWindows()[0])
				p.client.ShowNotification("Text send failed", "An error occurred when sending the text.")
			} else if res.OK && p.client.Notifications {
				fyne.CurrentApp().SendNotification(fyne.NewNotification("Text send completed", "The text was sent successfully."))
			}
		}(p.Length() - 1)
	}
}

// NewSendList greates a list of progress bars.
func NewSendList(client *transport.Client) *SendList {
	p := &SendList{client: client}
	p.List.Length = p.Length
	p.List.CreateItem = p.CreateItem
	p.List.UpdateItem = p.UpdateItem
	p.List.OnSelected = p.OnSelected
	p.ExtendBaseWidget(p)

	return p
}
