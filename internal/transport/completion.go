package transport

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"io"
	"strings"
	"time"

	"github.com/psanford/wormhole-william/rendezvous"
	"github.com/psanford/wormhole-william/wordlist"
	"github.com/psanford/wormhole-william/wormhole"
)

/* The code below is largely based on the following two files:
https://github.com/psanford/wormhole-william/blob/master/cmd/completion.go and
https://github.com/psanford/wormhole-william/blob/master/internal/crypto/crypto.go.
*/

// CompleteCode returns completion suggestions for the receiver code.
func (c *Client) CompleteRecvCode(toComplete string) []string {
	parts := strings.Split(toComplete, "-")
	if len(parts) < 2 {
		nameplates, err := c.activeNameplates()
		if err != nil {
			return nil
		} else if len(parts) == 0 {
			return nameplates
		}

		var candidates []string
		for _, nameplate := range nameplates {
			if strings.HasPrefix(nameplate, parts[0]) {
				candidates = append(candidates, nameplate+"-")
			}
		}

		return candidates
	}

	currentCompletion := parts[len(parts)-1]
	prefix := parts[:len(parts)-1]

	// Even/odd is based on just the number of words, slice off the mailbox
	parts = parts[1:]
	even := len(parts)%2 == 0

	var candidates []string
	for _, pair := range wordlist.RawWords {
		var candidateWord string
		if even {
			candidateWord = pair.Even
		} else {
			candidateWord = pair.Odd
		}
		if strings.HasPrefix(candidateWord, currentCompletion) {
			guessParts := append(prefix, candidateWord)
			candidates = append(candidates, strings.Join(guessParts, "-"))
		}
	}

	return candidates
}

// TODO: This part is slow. We might want to cache the result.
func (c *Client) activeNameplates() ([]string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 4*time.Second)
	defer cancel()

	url := c.RendezvousURL
	if url == "" {
		url = wormhole.DefaultRendezvousURL
	}

	appID := c.AppID
	if appID == "" {
		appID = wormhole.WormholeCLIAppID
	}

	client := rendezvous.NewClient(url, randSideID(), appID)

	mood := rendezvous.Happy
	defer client.Close(ctx, mood)

	_, err := client.Connect(ctx)
	if err != nil {
		return nil, err
	}

	return client.ListNameplates(ctx)
}

// randSideID returns a string appropate for use as the Side ID for a client.
func randSideID() string {
	return randHex(5)
}

// randHex generates hex encoded secure random bytes of the byteCount length.
func randHex(byteCount int) string {
	buf := make([]byte, byteCount)
	_, err := io.ReadFull(rand.Reader, buf)
	if err != nil {
		panic(err)
	}

	return hex.EncodeToString(buf)
}
