package util

import (
	"net/url"
	"testing"

	"github.com/alecthomas/assert/v2"
)

func TestURLToGitHubProject(t *testing.T) {
	tests := []struct {
		name    string
		subpath string
	}{
		{"Repository root", ""},
		{"Releases", "/releases"},
		{"Release v3.0.0", "/releases/v3.0.0"},
		{"Supported clients", "/wiki/Supported-clients"},
	}

	const basepath = https + "://" + github + repo
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			path := basepath + tt.subpath
			want, err := url.Parse(path)
			assert.NoError(t, err)

			got := URLToGitHubProject(tt.subpath)
			assert.Equal(t, want, got)
		})
	}
}
