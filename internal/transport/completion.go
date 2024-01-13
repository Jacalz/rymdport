package transport

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"io"
	"sort"
	"strings"
	"time"

	"github.com/psanford/wormhole-william/rendezvous"
	"github.com/psanford/wormhole-william/wordlist"
	"github.com/psanford/wormhole-william/wormhole"
)

/* The code below is an adapted and improved version based initially on the following two files:
https://github.com/psanford/wormhole-william/blob/master/cmd/completion.go and
https://github.com/psanford/wormhole-william/blob/master/internal/crypto/crypto.go.
*/

// CompleteRecvCode returns completion suggestions for the receiver code.
func (c *Client) CompleteRecvCode(toComplete string) []string {
	separators := strings.Count(toComplete, "-")
	if separators == 0 {
		return c.completeNameplates(toComplete)
	}

	lastPart := strings.LastIndexByte(toComplete, '-')
	completionMatch := toComplete[lastPart+1:] // Word prefix to match completion against.
	prefix := toComplete[:lastPart+1]          // Everything before the match prefix.

	// Even/odd is based on just the number of words, ignore mailbox
	even := (separators-1)%2 == 0

	// Perform binary search for word prefix in alphabetically sorted word list.
	index := sort.Search(256, func(i int) bool {
		pair := wordlist.RawWords[byte(i)]
		prefix := ""
		if even {
			prefix = pair.Even
		} else {
			prefix = pair.Odd
		}

		prefix = prefix[:len(completionMatch)]
		return prefix >= completionMatch
	})

	var candidates []string

	// Search forward for prefix matches from found index.
	for i := index; i < 256; i++ {
		candidate, match := lookupWordMatch(byte(i), completionMatch, even)
		if !match {
			break // Sorted in increasing alphabetical order. No more matches.
		}

		candidates = append(candidates, prefix+candidate)
	}

	// Search backwards for prefix matches from before the found index.
	for i := index - 1; i >= 0; i-- {
		candidate, match := lookupWordMatch(byte(i), completionMatch, even)
		if !match {
			break // Sorted in increasing alphabetical order. No more matches.
		}

		candidates = append(candidates, prefix+candidate)
	}

	return candidates
}

func lookupWordMatch(index byte, prefix string, even bool) (string, bool) {
	pair := wordlist.RawWords[index]
	var candidate string
	if even {
		candidate = pair.Even
	} else {
		candidate = pair.Odd
	}

	if !strings.HasPrefix(candidate, prefix) {
		return "", false
	}

	return candidate, true
}

func (c *Client) completeNameplates(toComplete string) []string {
	nameplates, err := c.activeNameplates()
	if err != nil {
		return nameplates
	}

	var candidates []string
	for _, nameplate := range nameplates {
		if strings.HasPrefix(nameplate, toComplete) {
			candidates = append(candidates, nameplate+"-")
		}
	}

	return candidates
}

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
	defer client.Close(ctx, rendezvous.Happy)

	_, err := client.Connect(ctx)
	if err != nil {
		return nil, err
	}

	return client.ListNameplates(ctx)
}

// randSideID returns a string appropate for use as the Side ID for a client.
func randSideID() string {
	buf := make([]byte, 5)
	_, err := io.ReadFull(rand.Reader, buf)
	if err != nil {
		panic(err)
	}

	return hex.EncodeToString(buf)
}
