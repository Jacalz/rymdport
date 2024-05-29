package transport

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"io"
	"slices"
	"strings"
	"time"

	"github.com/rymdport/wormhole/rendezvous"
	"github.com/rymdport/wormhole/wordlist"
	"github.com/rymdport/wormhole/wormhole"
)

/* The code below is an adapted and improved version based initially on the following two files:
https://github.com/psanford/wormhole-william/blob/master/cmd/completion.go and
https://github.com/psanford/wormhole-william/blob/master/internal/crypto/crypto.go.
*/

// GenerateCodeCompletion returns completion suggestions for the receiver code.
func (c *Client) GenerateCodeCompletion(toComplete string) []string {
	separators := strings.Count(toComplete, "-")
	if separators == 0 {
		return c.completeNameplates(toComplete)
	}

	lastPart := strings.LastIndexByte(toComplete, '-')
	completionMatch := toComplete[lastPart+1:] // Word prefix to match completion against.
	prefix := toComplete[:lastPart+1]          // Everything before the match prefix.

	// Even/odd is based on just the number of words, ignore mailbox.
	even := (separators-1)%2 == 0

	words := wordlist.Odd
	if even {
		words = wordlist.Even
	}

	// Default is to match everything. No binary search needed.
	index := 0

	// Perform binary search for word prefix in alphabetically sorted word list
	// only if there is a completion to look for.
	if len(completionMatch) > 0 {
		index, _ = slices.BinarySearch(words[:], completionMatch)
	}

	var candidates []string

	// Search forward for the other prefix matches.
	for i := index; i < 256; i++ {
		candidate, match := lookupWordMatch(byte(i), completionMatch, even)
		if !match {
			break // Sorted in increasing alphabetical order. No more matches.
		}

		candidates = append(candidates, prefix+candidate)
	}

	return candidates
}

// lookupWordMatch looks up a word at a specific index and even/odd setting.
// It also returns information about if the given word matches a completion prefix.
func lookupWordMatch(index byte, prefix string, even bool) (string, bool) {
	word := wordlist.Odd[index]
	if even {
		word = wordlist.Even[index]
	}

	return word, strings.HasPrefix(word, prefix)
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
