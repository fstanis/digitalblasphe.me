package digitalblasphemy

import (
	"crypto/sha256"
	"fmt"
	"io"
	"net/http"
	"strings"
	"testing"
)

type mockHTTPClient struct {
	lastRequest      *http.Request
	responseToReturn *http.Response
}

func (m *mockHTTPClient) Do(req *http.Request) (*http.Response, error) {
	m.lastRequest = req
	return m.responseToReturn, nil
}

func TestValidate(t *testing.T) {
	mock := &mockHTTPClient{
		responseToReturn: &http.Response{
			StatusCode: http.StatusOK,
		},
	}
	httpClient = mock

	creds := &Credentials{"user", "password"}
	if err := creds.Validate(); err != nil {
		t.Errorf("expected no error for Validate, got %v", err)
	}
	req := mock.lastRequest
	if req.Method != "HEAD" {
		t.Errorf("expected Validate to use HEAD method, got %v", req.Method)
	}
	if req.URL.String() != urlMembers {
		t.Errorf("expected Validate to use %q URL, got %q", urlMembers, req.URL.String())
	}
	if !strings.Contains(req.Header.Get("Authorization"), "Basic ") {
		t.Error("expected Validate to use basic authorization")
	}
}

func TestFetch(t *testing.T) {
	httpClient = http.DefaultClient

	reader, err := FetchWallpaper(&Wallpaper{
		ID:         "portals1",
		URL:        "https://secure.digitalblasphemy.com/graphics/HDfree/portals1HDfree.jpg",
		Resolution: "1920x1080",
	}, nil)
	if err != nil {
		t.Fatal(err)
	}
	defer reader.Close()

	hash := getHash(reader)
	if hash != "34f97c90805efd5c5da5db73efd66807cfcaeeb6b7eabd7741c8b5520e5eea86" {
		t.Fatalf("invalid hash: %s", hash)
	}
}

func getHash(r io.ReadCloser) string {
	h := sha256.New()
	if _, err := io.Copy(h, r); err != nil {
		return ""
	}
	return fmt.Sprintf("%x", h.Sum(nil))
}
