package digitalblasphemy

import (
	"testing"
)

func TestGetFreebies(t *testing.T) {
	list, err := GetFreebiesIndex()
	if err != nil {
		t.Fatalf("failed to get freebies with error %v", err)
	}
	if len(list) != 6 {
		t.Errorf("excepted 6 items, got %d", len(list))
	}
}
