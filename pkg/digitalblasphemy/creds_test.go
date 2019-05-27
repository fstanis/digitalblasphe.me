package digitalblasphemy

import (
	"bytes"
	"io/ioutil"
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
	mock := &mockHTTPClient{
		responseToReturn: &http.Response{
			StatusCode: http.StatusOK,
			Body:       ioutil.NopCloser(bytes.NewBuffer([]byte("example"))),
		},
	}
	httpClient = mock

	w := Wallpaper{
		ID:         "example",
		URL:        "https://example.com/",
		Resolution: "1920x1080",
	}

	filename, err := FetchWallpaper(w, nil)
	if err != nil {
		t.Fatal(err)
	}
	if mock.lastRequest.URL.String() != w.URL {
		t.Errorf("expected URL to be %q, got %q", w.URL, mock.lastRequest.URL.String())
	}
	mock.lastRequest = nil

	filename2, err := FetchWallpaper(w, nil)
	if err != nil {
		t.Fatal(err)
	}
	if mock.lastRequest != nil {
		t.Error("expected second Fetch to make no HTTP request")
	}
	if filename != filename2 {
		t.Errorf("expected %q to be the same as %q", filename, filename2)
	}
}
