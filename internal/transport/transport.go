// Package transport handles sending and receiving using wormhole-william
package transport

import (
	"github.com/rymdport/wormhole/wormhole"
)

// Client defines the client for handling sending and receiving using wormhole-william
type Client struct {
	wormhole.Client
}
