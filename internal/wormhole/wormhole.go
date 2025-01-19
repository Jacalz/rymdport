// Package wormhole is a helper around the wormhole package.
package wormhole

import (
	"github.com/rymdport/wormhole/wormhole"
)

// NewClient creates a new client object.
func NewClient() *Client {
	return &Client{}
}

// Client defines the client for handling sending and receiving.
type Client struct {
	connection wormhole.Client
}
