package digitalblasphemy

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
)

const userAgent = "https://github.com/fstanis/digitalblaspheme"

var (
	// ErrInvalidCredentials happens when the credentials were rejected by the
	// server.
	ErrInvalidCredentials = errors.New("got unauthorized from server")

	httpClient httpDoer = http.DefaultClient
)

// Credentials holds the username and password used to authorize in the members
// section of the website.
type Credentials struct {
	Username string
	Password string
}

// Validate tries to load a page and thus verify the credentials.
func (c *Credentials) Validate() error {
	resp, err := httpRequest("HEAD", urlMembers, c)
	if err != nil {
		return err
	}
	if resp.StatusCode == http.StatusUnauthorized {
		return ErrInvalidCredentials
	}
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("got status: %q", resp.Status)
	}
	return nil
}

// FetchWallpaper downloads the given wallpaper, optionally with the provided
// credentials.
func FetchWallpaper(wallpaper *Wallpaper, creds *Credentials) (io.ReadCloser, error) {
	return fetch(wallpaper.URL, creds)
}

func fetch(url string, creds *Credentials) (io.ReadCloser, error) {
	resp, err := httpRequest("GET", url, creds)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode == http.StatusUnauthorized {
		return nil, ErrInvalidCredentials
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("got status: %q", resp.Status)
	}
	return resp.Body, nil
}

func httpRequest(method string, url string, creds *Credentials) (*http.Response, error) {
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", userAgent)

	// We don't want to send credentials in the open for non-https URLs.
	if creds != nil && strings.HasPrefix(url, "https://") {
		req.SetBasicAuth(creds.Username, creds.Password)
	}
	return httpClient.Do(req)
}

type httpDoer interface {
	Do(*http.Request) (*http.Response, error)
}
