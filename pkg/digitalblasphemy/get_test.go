package digitalblasphemy

import (
	"net/http"
	"os"
	"reflect"
	"strings"
	"testing"
)

func TestGetIndex(t *testing.T) {
	data, err := os.Open("test/index.html")
	if err != nil {
		t.Fatal(err)
	}

	mock := &mockHTTPClient{
		responseToReturn: &http.Response{
			StatusCode: http.StatusOK,
			Body:       data,
		},
	}
	httpClient = mock

	creds := &Credentials{"user", "password"}
	list, err := GetIndex("1024x768", creds)
	if err != nil {
		t.Fatalf("got error trying to get index: %v", err)
	}

	req := mock.lastRequest
	if !strings.Contains(req.Header.Get("Authorization"), "Basic ") {
		t.Error("expected GetIndex to use basic authorization")
	}
	if !strings.HasSuffix(req.URL.String(), indexURLSort) {
		t.Error("expected index URL to use sort")
	}

	if len(list) != 1 {
		t.Errorf("expected 1 item in index, got %d", len(list))
	}
	want := Wallpaper{
		ID:         "second",
		URL:        "https://secure.digitalblasphemy.com/content/jpgs/1024st/second1024st.jpg",
		Resolution: "1024x768",
	}
	if !reflect.DeepEqual(list[0], want) {
		t.Errorf("expected wallpaper %+v, got %+v", want, list[0])
	}

	mock.lastRequest = nil
	_, err = GetIndex("1024x768", creds)
	if err != nil {
		t.Fatal(err)
	}
	if mock.lastRequest != nil {
		t.Error("expected second request for index to not make a new HTTP request")
	}
}

func TestGetFreebiesIndex(t *testing.T) {
	data, err := os.Open("test/freebies.html")
	if err != nil {
		t.Fatal(err)
	}

	mock := &mockHTTPClient{
		responseToReturn: &http.Response{
			StatusCode: http.StatusOK,
			Body:       data,
		},
	}
	httpClient = mock

	list, err := GetFreebiesIndex()
	if err != nil {
		t.Fatal(err)
	}

	req := mock.lastRequest
	if req.Header.Get("Authorization") != "" {
		t.Error("expected GetFreebiesIndex to not use authorization")
	}
	if req.URL.String() != urlFreebies {
		t.Errorf("expected URL to be %q, got %q", urlFreebies, req.URL.String())
	}

	if len(list) != 6 {
		t.Errorf("excepted 6 items in index, got %d", len(list))
	}

	mock.lastRequest = nil
	_, err = GetFreebiesIndex()
	if err != nil {
		t.Fatal(err)
	}
	if mock.lastRequest != nil {
		t.Error("expected second request for index to not make a new HTTP request")
	}
}
