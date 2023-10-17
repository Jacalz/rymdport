package util

import "net/url"

const (
	https  = "https"
	github = "github.com"
	repo   = "/jacalz/rymdport"
)

// URLToGitHubProject returns a pre-parsed link to a GitHub site starting from github.com/jacalz/rymdport.
func URLToGitHubProject(subpath string) *url.URL {
	return &url.URL{Scheme: https, Host: github, Path: repo + subpath}
}
