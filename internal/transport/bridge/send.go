package bridge

import (
	"sync"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/storage"
	"fyne.io/fyne/v2/widget"
	"github.com/Jacalz/rymdport/v3/internal/transport"
	"github.com/Jacalz/rymdport/v3/internal/util"
	"github.com/psanford/wormhole-william/wormhole"
)

// SendItem is the item that is being sent.
type SendItem struct {
	URI  fyne.URI
	Code string

	Value  int64
	Max    int64
	Status func() string

	list *SendList
}

func (s *SendItem) update(sent, total int64) {
	s.Value = sent
	s.Max = total
	s.list.Refresh()
}

func (s *SendItem) failed() {
	s.Status = func() string { return "Failed" }
}

// SendList is a list of progress bars that track send progress.
type SendList struct {
	widget.List

	client *transport.Client

	Items []*SendItem
	lock  sync.RWMutex

	window fyne.Window
	canvas fyne.Canvas
}

// Length returns the length of the data.
func (p *SendList) Length() int {
	return len(p.Items)
}

// CreateItem creates a new item in the list.
func (p *SendList) CreateItem() fyne.CanvasObject {
	return container.New(&listLayout{},
		&widget.FileIcon{},
		&widget.Label{Text: "Waiting for filename...", Wrapping: fyne.TextTruncate},
		newCodeDisplay(p.window),
		&widget.ProgressBar{},
	)
}

// UpdateItem updates the data in the list.
func (p *SendList) UpdateItem(i int, item fyne.CanvasObject) {
	p.lock.RLock()
	defer p.lock.RUnlock()

	container := item.(*fyne.Container)

	container.Objects[0].(*widget.FileIcon).SetURI(p.Items[i].URI)
	container.Objects[1].(*widget.Label).SetText(p.Items[i].URI.Name())
	container.Objects[2].(*fyne.Container).Objects[0].(*codeDisplay).SetText(p.Items[i].Code)

	progress := container.Objects[3].(*widget.ProgressBar)
	progress.Max = float64(p.Items[i].Max)
	progress.Value = float64(p.Items[i].Value)
	progress.TextFormatter = p.Items[i].Status
	progress.Refresh()
}

// NewSend adds data about a new send to the list and then returns the item.
func (p *SendList) NewSend(uri fyne.URI) *SendItem {
	p.lock.Lock()
	defer p.lock.Unlock()

	item := &SendItem{Code: "Waiting for code...", URI: uri, list: p, Max: 1}
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

	item := p.NewSend(file.URI())
	p.Refresh()

	go func() {
		// We want to catch close errors for security reasons.
		defer func() {
			if err = file.Close(); err != nil {
				item.failed()
				fyne.LogError("Error on closing file", err)
			}
		}()

		code, result, err := p.client.NewFileSend(file, wormhole.WithProgress(item.update), p.getCustomCode())
		if err != nil {
			fyne.LogError("Error on sending file", err)
			item.failed()
			dialog.ShowError(err, p.window)
			return
		}

		item.Code = code
		p.Refresh()

		if res := <-result; res.Error != nil {
			fyne.LogError("Error on sending file", res.Error)
			item.failed()
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

	item := p.NewSend(dir)
	p.Refresh()

	go func() {
		code, result, err := p.client.NewDirSend(dir, wormhole.WithProgress(item.update), p.getCustomCode())
		if err != nil {
			fyne.LogError("Error on sending directory", err)
			item.failed()
			dialog.ShowError(err, p.window)
			return
		}

		item.Code = code
		p.Refresh()

		if res := <-result; res.Error != nil {
			fyne.LogError("Error on sending directory", res.Error)
			item.failed()
			dialog.ShowError(res.Error, p.window)
			p.client.ShowNotification("Directory send failed", "An error occurred when sending the directory.")
		} else if res.OK {
			p.client.ShowNotification("Directory send completed", "The directory was sent successfully.")
		}
	}()
}

// SendText sends new text.
func (p *SendList) SendText() {
	go func() {
		text := <-p.client.ShowTextSendWindow()
		if text == "" {
			return
		}

		item := p.NewSend(storage.NewFileURI("Text Snippet"))
		p.Refresh()

		code, result, err := p.client.NewTextSend(text, wormhole.WithProgress(item.update), p.getCustomCode())
		if err != nil {
			fyne.LogError("Error on sending text", err)
			item.failed()
			dialog.ShowError(err, p.window)
			return
		}

		item.Code = code
		p.Refresh()

		if res := <-result; res.Error != nil {
			fyne.LogError("Error on sending text", res.Error)
			item.failed()
			dialog.ShowError(res.Error, p.window)
			p.client.ShowNotification("Text send failed", "An error occurred when sending the text.")
		} else if res.OK && p.client.Notifications {
			p.client.ShowNotification("Text send completed", "The text was sent successfully.")
		}
	}()
}

// getCustomCode returns "" if the user has custom codes disabled.
// Otherwise, it will ask the user for a code.
func (p *SendList) getCustomCode() string {
	if !p.client.CustomCode {
		return ""
	}

	code := make(chan string)
	codeEntry := &widget.Entry{
		PlaceHolder: "123-example-code",
		Wrapping:    fyne.TextTruncate,
		Validator:   util.CodeValidator,
	}

	form := dialog.NewForm("Create custom code", "Confirm", "Cancel", []*widget.FormItem{
		{
			Text: "Code", Widget: codeEntry,
			HintText: "A code beginning with a number, followed by groups of letters separated with \"-\".",
		},
	}, func(submitted bool) {
		if !submitted || codeEntry.Text == codeEntry.PlaceHolder {
			code <- ""
		} else {
			code <- codeEntry.Text
		}

		close(code)
	}, p.window)
	form.Resize(fyne.Size{Width: p.canvas.Size().Width * 0.8})
	form.Show()
	p.canvas.Focus(codeEntry)

	return <-code
}

// NewSendList greates a list of progress bars.
func NewSendList(window fyne.Window, client *transport.Client) *SendList {
	p := &SendList{client: client, window: window, canvas: window.Canvas()}
	p.List.Length = p.Length
	p.List.CreateItem = p.CreateItem
	p.List.UpdateItem = p.UpdateItem
	p.List.OnSelected = p.Unselect
	p.ExtendBaseWidget(p)

	return p
}
