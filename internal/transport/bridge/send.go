package bridge

import (
	"path/filepath"
	"sync/atomic"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/storage"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/Jacalz/rymdport/v3/internal/transport"
	"github.com/Jacalz/rymdport/v3/internal/util"
	"github.com/psanford/wormhole-william/wormhole"
	"github.com/rymdport/go-qrcode"
)

// SendItem is the item that is being sent.
type SendItem struct {
	URI  fyne.URI
	Code string

	Value  int64
	Max    int64
	Status func() string

	// Allow the list to only refresh a single object.
	refresh func(int)
	index   int
}

func (s *SendItem) update(sent, total int64) {
	s.Value = sent
	s.Max = total
	s.refresh(s.index)
}

func (s *SendItem) failed() {
	s.Status = func() string { return "Failed" }
	s.refresh(s.index)
}

// SendData is a list of progress bars that track send progress.
type SendData struct {
	Client *transport.Client
	Window fyne.Window
	Canvas fyne.Canvas

	items []*SendItem

	deleting atomic.Bool
	list     *widget.List
}

// Length returns the length of the data.
func (d *SendData) Length() int {
	return len(d.items)
}

// CreateItem creates a new item in the list.
func (d *SendData) CreateItem() fyne.CanvasObject {
	return container.New(listLayout{},
		&widget.FileIcon{},
		&widget.Label{Text: "Waiting for filename...", Truncation: fyne.TextTruncateEllipsis},
		newCodeDisplay(d.Window),
		&widget.ProgressBar{},
	)
}

// UpdateItem updates the data in the list.
func (d *SendData) UpdateItem(i int, object fyne.CanvasObject) {
	container := object.(*fyne.Container)

	item := d.items[i]
	container.Objects[0].(*widget.FileIcon).SetURI(item.URI)
	container.Objects[1].(*widget.Label).SetText(item.URI.Name())
	container.Objects[2].(*fyne.Container).Objects[0].(*widget.Label).SetText(item.Code)

	progress := container.Objects[3].(*widget.ProgressBar)
	progress.Max = float64(item.Max)
	progress.Value = float64(item.Value)
	progress.TextFormatter = item.Status
	progress.Refresh()
}

// OnSelected currently just makes sure that we don't persist selection.
func (d *SendData) OnSelected(i int) {
	d.list.Unselect(i)

	code, err := qrcode.New("wormhole-transfer:"+d.items[i].Code, qrcode.High)
	if err != nil {
		fyne.LogError("Failed to encode qr code", err)
		return
	}

	code.BackgroundColor = theme.OverlayBackgroundColor()
	code.ForegroundColor = theme.ForegroundColor()

	qrcode := canvas.NewImageFromImage(code.Image(100))
	qrcode.FillMode = canvas.ImageFillOriginal
	qrcode.ScaleMode = canvas.ImageScalePixels
	qrcode.SetMinSize(fyne.NewSize(100, 100))

	supportedClientsUrl := util.URLToGitHubProject("/wiki/Supported-clients")
	qrCodeInfo := widget.NewRichText(&widget.TextSegment{
		Style: widget.RichTextStyleInline,
		Text:  "A list of supported apps can be found ",
	}, &widget.HyperlinkSegment{
		Text: "here",
		URL:  supportedClientsUrl,
	}, &widget.TextSegment{
		Style: widget.RichTextStyleInline,
		Text:  ".",
	})
	qrCard := &widget.Card{Image: qrcode, Content: container.NewCenter(qrCodeInfo)}

	var infoDialog *dialog.CustomDialog
	removeLabel := &widget.Label{Text: "This item can be removed.\nThe transfer has completed."}
	removeButton := &widget.Button{Icon: theme.DeleteIcon(), Importance: widget.DangerImportance, Text: "Remove", OnTapped: func() {
		d.remove(i)
		infoDialog.Hide()
		infoDialog = nil
	}}

	// Only allow failed or completed items to be removed.
	if d.items[i].Value < d.items[i].Max && d.items[i].Status == nil {
		removeLabel.Text = "This item can not be removed yet.\nThe transfer needs to complete first."
		removeButton.Disable()
	} else {
		qrcode.Image = nil
		qrcode.Resource = theme.InfoIcon()
		qrcode.ScaleMode = canvas.ImageScaleSmooth
		qrcode.Refresh()

		qrCard.Content = &widget.Label{Text: "This transfer is not active.\nCan't show a QR code."}
	}

	removeCard := &widget.Card{Content: container.NewVBox(removeLabel, removeButton)}

	infoDialog = dialog.NewCustom("Information", "Close", container.NewGridWithColumns(2, qrCard, removeCard), d.Window)
	infoDialog.Show()
}

// NewSend adds data about a new send to the list and then returns the item.
func (d *SendData) NewSend(uri fyne.URI) *SendItem {
	item := &SendItem{Code: "Waiting for code...", URI: uri, Max: 1, refresh: d.refresh, index: len(d.items)}
	d.items = append(d.items, item)
	return item
}

// OnFileSelect is intended to be passed as callback to a FileOpen dialog.
func (d *SendData) OnFileSelect(file fyne.URIReadCloser, err error) {
	if err != nil {
		fyne.LogError("Error on selecting file to send", err)
		dialog.ShowError(err, d.Window)
		return
	} else if file == nil {
		return
	}

	item := d.NewSend(file.URI())
	d.list.Refresh()

	go func() {
		// We want to catch close errors for security reasons.
		defer func() {
			if err = file.Close(); err != nil {
				item.failed()
				fyne.LogError("Error on closing file", err)
			}
		}()

		code, result, err := d.Client.NewFileSend(file, wormhole.WithProgress(item.update), d.getCustomCode())
		if err != nil {
			fyne.LogError("Error on sending file", err)
			item.failed()
			dialog.ShowError(err, d.Window)
			return
		}

		item.Code = code
		d.refresh(item.index)

		if res := <-result; res.Error != nil {
			fyne.LogError("Error on sending file", res.Error)
			item.failed()
			dialog.ShowError(res.Error, d.Window)
			d.Client.ShowNotification("File send failed", "An error occurred when sending the file.")
		} else if res.OK {
			d.Client.ShowNotification("File send completed", "The file was sent successfully.")
		}
	}()
}

// OnDirSelect is intended to be passed as callback to a FolderOpen dialog.
func (d *SendData) OnDirSelect(dir fyne.ListableURI, err error) {
	if err != nil {
		fyne.LogError("Error on selecting dir to send", err)
		dialog.ShowError(err, d.Window)
		return
	} else if dir == nil {
		return
	}

	item := d.NewSend(dir)
	d.list.Refresh()

	go func() {
		code, result, err := d.Client.NewDirSend(dir, wormhole.WithProgress(item.update), d.getCustomCode())
		if err != nil {
			fyne.LogError("Error on sending directory", err)
			item.failed()
			dialog.ShowError(err, d.Window)
			return
		}

		item.Code = code
		d.refresh(item.index)

		if res := <-result; res.Error != nil {
			fyne.LogError("Error on sending directory", res.Error)
			item.failed()
			dialog.ShowError(res.Error, d.Window)
			d.Client.ShowNotification("Directory send failed", "An error occurred when sending the directory.")
		} else if res.OK {
			d.Client.ShowNotification("Directory send completed", "The directory was sent successfully.")
		}
	}()
}

// NewSendFromFiles creates a directory from the files and sends it as a directory send.
func (d *SendData) NewSendFromFiles(uris []fyne.URI) {
	parentDir := storage.NewFileURI(filepath.Dir(uris[0].Path()))
	item := d.NewSend(parentDir)
	d.list.Refresh()

	go func() {
		code, result, err := d.Client.NewMultipleFileSend(uris, wormhole.WithProgress(item.update), d.getCustomCode())
		if err != nil {
			fyne.LogError("Error on sending directory", err)
			item.failed()
			dialog.ShowError(err, d.Window)
			return
		}

		item.Code = code
		d.refresh(item.index)

		if res := <-result; res.Error != nil {
			fyne.LogError("Error on sending directory", res.Error)
			item.failed()
			dialog.ShowError(res.Error, d.Window)
			d.Client.ShowNotification("Directory send failed", "An error occurred when sending the directory.")
		} else if res.OK {
			d.Client.ShowNotification("Directory send completed", "The directory was sent successfully.")
		}
	}()
}

// SendText sends new text.
func (d *SendData) SendText() {
	go func() {
		text := d.Client.ShowTextSendWindow()
		if text == "" {
			return
		}

		d.Window.RequestFocus() // Refocus the main window
		item := d.NewSend(storage.NewFileURI("Text Snippet"))
		d.list.Refresh()

		code, result, err := d.Client.NewTextSend(text, wormhole.WithProgress(item.update), d.getCustomCode())
		if err != nil {
			fyne.LogError("Error on sending text", err)
			item.failed()
			dialog.ShowError(err, d.Window)
			return
		}

		item.Code = code
		d.refresh(item.index)

		if res := <-result; res.Error != nil {
			fyne.LogError("Error on sending text", res.Error)
			item.failed()
			dialog.ShowError(res.Error, d.Window)
			d.Client.ShowNotification("Text send failed", "An error occurred when sending the text.")
		} else if res.OK && d.Client.Notifications {
			d.Client.ShowNotification("Text send completed", "The text was sent successfully.")
		}
	}()
}

// getCustomCode returns "" if the user has custom codes disabled.
// Otherwise, it will ask the user for a code.
func (d *SendData) getCustomCode() string {
	if !d.Client.CustomCode {
		return ""
	}

	code := make(chan string)
	codeEntry := &widget.Entry{
		PlaceHolder: "123-example-code",
		Scroll:      container.ScrollBoth,
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
	}, d.Window)
	form.Resize(fyne.Size{Width: d.Canvas.Size().Width * 0.8})
	codeEntry.OnSubmitted = func(_ string) { form.Submit() }
	form.Show()
	d.Canvas.Focus(codeEntry)

	return <-code
}

func (d *SendData) refresh(index int) {
	if d.deleting.Load() {
		return // Don't update if we are deleting.
	}

	d.list.RefreshItem(index)
}

func (d *SendData) remove(index int) {
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

// NewSendList greates a list of progress bars.
func NewSendList(data *SendData) *widget.List {
	list := &widget.List{
		Length:     data.Length,
		CreateItem: data.CreateItem,
		UpdateItem: data.UpdateItem,
		OnSelected: data.OnSelected,
	}
	data.list = list
	return list
}
