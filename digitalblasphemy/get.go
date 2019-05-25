package digitalblasphemy

import (
	"errors"
	"fmt"

	"github.com/PuerkitoBio/goquery"
)

// ErrInvalidResolution happens when the requested resolution doesn't exist
// on the website.
var ErrInvalidResolution = errors.New("invalid resolution")

// Wallpaper is a single wallpaper image that can be downloaded from the
// website.
type Wallpaper struct {
	ID         string
	URL        string
	Resolution string
}

// GetIndex returns the list of all wallpaper URLs for the given resolution.
func GetIndex(resolution string, creds *Credentials) ([]*Wallpaper, error) {
	index, ok := indexURLForResolution[resolution]
	if !ok {
		return nil, ErrInvalidResolution
	}

	url := index + indexURLSort
	fmt.Println(url)
	reader, err := fetch(url, creds)
	if err != nil {
		return nil, err
	}
	defer reader.Close()

	doc, err := goquery.NewDocumentFromReader(reader)
	if err != nil {
		return nil, err
	}

	list := parseIndex(doc, resolution)
	result := make([]*Wallpaper, len(list))
	for i, filename := range list {
		submatches := urlRegexpForResolution[resolution].FindStringSubmatch(filename)
		result[i] = &Wallpaper{
			ID:         submatches[1],
			URL:        index + filename,
			Resolution: resolution,
		}
	}
	return result, nil
}

// GetFreebiesIndex returns the list of all wallpaper URLs that are available
// for free.
func GetFreebiesIndex() ([]*Wallpaper, error) {
	doc, err := goquery.NewDocument(urlFreebies)
	if err != nil {
		return nil, err
	}
	list := parseFreebies(doc)

	result := make([]*Wallpaper, len(list))
	for i, id := range list {
		result[i] = &Wallpaper{
			ID:         id,
			URL:        makeFreebieURL(id),
			Resolution: "1920x1080",
		}
	}
	return result, nil
}
