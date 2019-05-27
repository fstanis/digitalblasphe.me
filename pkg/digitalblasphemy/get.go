package digitalblasphemy

import (
	"errors"

	"github.com/PuerkitoBio/goquery"
)

var (
	// ErrInvalidResolution happens when the requested resolution doesn't exist
	// on the website.
	ErrInvalidResolution = errors.New("invalid resolution")

	// ErrNoWallpapersFound happens when the server response is parsed successfully,
	// but no valid items are found.
	ErrNoWallpapersFound = errors.New("zero wallpapers returned from server")
)

// Wallpaper is a single wallpaper image that can be downloaded from the
// website.
type Wallpaper struct {
	ID         string
	URL        string
	Resolution string
}

// GetIndex returns the list of all wallpaper URLs for the given resolution.
func GetIndex(resolution string, creds *Credentials) ([]Wallpaper, error) {
	index, ok := indexURLForResolution[resolution]
	if !ok {
		return nil, ErrInvalidResolution
	}

	if index := cache.getIndex(resolution); index != nil {
		return index, nil
	}

	url := index + indexURLSort
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
	if len(list) == 0 {
		return nil, ErrNoWallpapersFound
	}

	result := make([]Wallpaper, len(list))
	for i, filename := range list {
		submatches := urlRegexpForResolution[resolution].FindStringSubmatch(filename)
		result[i] = Wallpaper{
			ID:         submatches[1],
			URL:        index + filename,
			Resolution: resolution,
		}
	}

	cache.setIndex(resolution, result)
	return result, nil
}

// GetFreebiesIndex returns the list of all wallpaper URLs that are available
// for free.
func GetFreebiesIndex() ([]Wallpaper, error) {
	if index := cache.getIndexFreebies(); index != nil {
		return index, nil
	}

	reader, err := fetch(urlFreebies, nil)
	if err != nil {
		return nil, err
	}
	defer reader.Close()

	doc, err := goquery.NewDocumentFromReader(reader)
	if err != nil {
		return nil, err
	}
	list := parseFreebies(doc)
	if len(list) == 0 {
		return nil, ErrNoWallpapersFound
	}

	result := make([]Wallpaper, len(list))
	for i, id := range list {
		result[i] = Wallpaper{
			ID:         id,
			URL:        makeFreebieURL(id),
			Resolution: "1920x1080",
		}
	}

	cache.setIndexFreebies(result)
	return result, nil
}

// IsValidResolution checks if the given resolution is valid.
func IsValidResolution(resolution string) bool {
	_, ok := indexURLForResolution[resolution]
	return ok
}

// GetValidResolutions returns the list of all valid resolutions.
func GetValidResolutions() []string {
	result := make([]string, 0, len(indexURLForResolution))
	for resolution := range indexURLForResolution {
		result = append(result, resolution)
	}
	return result
}
