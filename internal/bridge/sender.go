package bridge

import (
	"context"
	"os"

	"fyne.io/fyne"
	"github.com/psanford/wormhole-william/wormhole"
)

// SendFile takes the chosen file and sends it using wormhole-william.
func (b *Bridge) SendFile(file fyne.URIReadCloser, code chan string, progress wormhole.SendOption) error {
	c := wormhole.Client{PassPhraseComponentLength: b.ComponentLength}

	defer file.Close()

	f, err := os.Open(file.URI().String()[7:])
	if err != nil {
		fyne.LogError("Error on opening file", err)
		return err
	}

	defer f.Close() // #nosec - We are not writing to the file.

	codestr, status, err := c.SendFile(context.Background(), file.Name(), f, progress)
	if err != nil {
		fyne.LogError("Error on sending file", err)
		return err
	}

	code <- codestr

	if stat := <-status; stat.Error != nil {
		fyne.LogError("Error on status of share", err)
		return err
	} else if stat.OK {
		return nil
	}

	return nil
}

// SendText takes a text input and sends the text using wormhole-william.
func (b *Bridge) SendText(text string, code chan string) error {
	c := wormhole.Client{PassPhraseComponentLength: b.ComponentLength}

	codestr, status, err := c.SendText(context.Background(), text) // TODO: Check why progress doesn't work for sending text.
	if err != nil {
		fyne.LogError("Error on sending text", err)
		return err
	}

	code <- codestr

	if stat := <-status; stat.Error != nil {
		fyne.LogError("Error on status of share", err)
		return err
	} else if stat.OK {
		return nil
	}

	return nil
}
