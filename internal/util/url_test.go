package util

import (
	"net/url"
	"reflect"
	"testing"
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
			if err != nil {
				t.Errorf("Failed to parse the path %s: %v", path, err)
				return
			}

			if got := URLToGitHubProject(tt.subpath); !reflect.DeepEqual(got, want) {
				t.Errorf("URLToGitHubProject() = %v, want %v", got, want)
			}
		})
	}
}
