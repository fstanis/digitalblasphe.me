package digitalblasphemy

import (
	"os"
	"reflect"
	"testing"

	"github.com/PuerkitoBio/goquery"
)

func loadDocument(filename string) (*goquery.Document, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return goquery.NewDocumentFromReader(f)
}

func TestParseIndex(t *testing.T) {
	doc, err := loadDocument("test/index.html")
	if err != nil {
		t.Fatal(err)
	}
	got := parseIndex(doc, "1024x768")
	want := []string{"second1024st.jpg"}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %v, want %v", got, want)
	}
}

func TestParseFreebies(t *testing.T) {
	doc, err := loadDocument("test/freebies.html")
	if err != nil {
		t.Fatal(err)
	}
	got := parseFreebies(doc)
	want := []string{"standingstones1", "acumen1", "skygatespring", "moonbeamsea1", "portals1", "moonshadow1"}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %v, want %v", got, want)
	}
}
