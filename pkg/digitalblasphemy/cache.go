package digitalblasphemy

import (
	"encoding/gob"
	"fmt"
	"io"
	"math/rand"
	"os"
	"path/filepath"
	"time"

	"github.com/shibukawa/configdir"
)

const (
	cacheName               = "digitalblasphemy"
	cacheDataFilename       = "data"
	cachedIndexDuration     = time.Hour * 24
	cachedWallpaperDuration = time.Hour * 24 * 30
)

var (
	cacheFolder = configdir.New("", cacheName).QueryCacheFolder()

	cache cacheData
)

// LoadCache loads the cache from disk.
func LoadCache() error {
	return cache.load()
}

// SaveCache saves the cache on disk.
func SaveCache() error {
	return cache.save()
}

// GetCachedCurrent gets the current wallpaper ID from the cache.
func GetCachedCurrent() string {
	return cache.CurrentID
}

// SetCachedCurrent sets the current wallpaper ID in the cache.
func SetCachedCurrent(id string) {
	cache.CurrentID = id
}

func init() {
	cache.init()
}

type cacheData struct {
	WallpaperMap map[Wallpaper]string

	Indexes      map[string]cachedIndex
	FreebieIndex cachedIndex

	CurrentID string
}

type cachedIndex struct {
	Updated time.Time
	Index   []Wallpaper
}

func (c *cacheData) init() {
	if c.WallpaperMap == nil {
		c.WallpaperMap = make(map[Wallpaper]string)
	}
	if c.Indexes == nil {
		c.Indexes = make(map[string]cachedIndex)
	}
}

func (c *cacheData) garbageCollect() error {
	matches, err := filepath.Glob(filepath.Join(cacheFolder.Path, "*.jpg"))
	if err != nil {
		return err
	}

	removed := make(map[string]bool)
	for _, filename := range matches {
		stat, err := os.Stat(filename)
		if err != nil || time.Since(stat.ModTime()) > cachedWallpaperDuration {
			os.Remove(filename)
			removed[filename] = true
		}
	}
	for key, filename := range c.WallpaperMap {
		if removed[filename] {
			delete(c.WallpaperMap, key)
		}
	}
	return nil
}

func (c *cacheData) save() error {
	cacheFolder.MkdirAll()
	f, err := cacheFolder.Create(cacheDataFilename)
	if err != nil {
		return err
	}
	defer f.Close()
	enc := gob.NewEncoder(f)
	if err := enc.Encode(*c); err != nil {
		return err
	}
	return nil
}

func (c *cacheData) load() error {
	f, err := cacheFolder.Open(cacheDataFilename)
	if err != nil {
		return nil
	}
	defer f.Close()
	enc := gob.NewDecoder(f)
	if err := enc.Decode(c); err != nil {
		return err
	}
	return nil
}

func (c *cacheData) getWallpaper(w Wallpaper) string {
	filename, exists := c.WallpaperMap[w]
	if !exists {
		return ""
	}
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		delete(c.WallpaperMap, w)
		return ""
	}
	return filename
}

func (c *cacheData) putWallpaper(w Wallpaper, data io.Reader) (string, error) {
	c.garbageCollect()
	filename := filepath.Join(cacheFolder.Path, fmt.Sprintf("cache%d.jpg", rand.Int()))

	cacheFolder.MkdirAll()
	f, err := os.Create(filename)
	if err != nil {
		return "", err
	}
	defer f.Close()
	io.Copy(f, data)

	c.WallpaperMap[w] = filename

	return filename, nil
}

func (c *cacheData) getIndex(resolution string) []Wallpaper {
	data, exists := c.Indexes[resolution]
	if !exists {
		return nil
	}
	return cachedIndexIfValid(data)
}

func (c *cacheData) setIndex(resolution string, index []Wallpaper) {
	c.Indexes[resolution] = cachedIndex{
		Updated: time.Now(),
		Index:   index,
	}
}

func (c *cacheData) getIndexFreebies() []Wallpaper {
	return cachedIndexIfValid(c.FreebieIndex)
}

func (c *cacheData) setIndexFreebies(index []Wallpaper) {
	c.FreebieIndex = cachedIndex{
		Updated: time.Now(),
		Index:   index,
	}
}

func cachedIndexIfValid(data cachedIndex) []Wallpaper {
	if time.Since(data.Updated) > cachedIndexDuration {
		return nil
	}
	return data.Index
}
