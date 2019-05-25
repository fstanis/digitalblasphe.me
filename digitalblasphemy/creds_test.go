package digitalblasphemy

import (
	"crypto/sha256"
	"fmt"
	"io"
	"testing"
)

func TestValidate(t *testing.T) {
	invalidCreds := &Credentials{"invalid", "invalid"}
	if invalidCreds.Validate() != ErrInvalidCredentials {
		t.Fatal("invalid credentials should return ErrInvalidCredentials")
	}
}

func TestFetch(t *testing.T) {
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
