//go:build release

package updater

import (
	"crypto/ed25519"
	"time"

	"fyne.io/fyne/v2"
	"github.com/fynelabs/fyneselfupdate"
	"github.com/fynelabs/selfupdate"
)

// Enabled specifies if the binary is built with update checking enabled.
const Enabled = true

// Enable turns on automatic application updates.
func Enable(a fyne.App, w fyne.Window) {
	publicKey := ed25519.PublicKey{165, 235, 49, 149, 238, 5, 192, 143, 38, 96, 124, 12, 235, 116, 94, 72, 244, 97, 230, 44, 60, 135, 85, 202, 98, 103, 233, 119, 63, 57, 83, 106}

	// The public key above matches the signature of the below file served by our CDN
	httpSource := selfupdate.NewHTTPSource(nil, "https://geoffrey-artefacts.fynelabs.com/self-update/22/22a3af17-e614-4b44-bce2-8b32ab7463ac/{{.OS}}-{{.Arch}}/{{.Executable}}{{.Ext}}")

	config := fyneselfupdate.NewConfigWithTimeout(a, w, time.Minute, httpSource, selfupdate.Schedule{FetchOnStart: true}, publicKey)

	_, err := selfupdate.Manage(config)
	if err != nil {
		fyne.LogError("Error while setting up update manager: ", err)
	}
}
