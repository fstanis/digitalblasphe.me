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
	cacheName           = "digitalblasphemy"
	cacheDataFilename   = "data"
	cachedIndexDuration = time.Hour * 24
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

func init() {
	cache.init()
}

type cacheData struct {
	WallpaperMap map[Wallpaper]string
	Wallpapers   []Wallpaper

	Indexes      map[string]cachedIndex
	FreebieIndex cachedIndex
}

type cachedIndex struct {
	Updated time.Time
	Index   []Wallpaper
}

func (c *cacheData) init() {
	if cache.WallpaperMap == nil {
		cache.WallpaperMap = make(map[Wallpaper]string)
	}
	if cache.Indexes == nil {
		cache.Indexes = make(map[string]cachedIndex)
	}
}

func (c *cacheData) garbageCollect() error {
	if len(c.Wallpapers) > 10 {
		for i := 0; i < 10; i++ {
			wallpaper := c.Wallpapers[i]
			filename := c.WallpaperMap[wallpaper]
			if _, err := os.Stat(filename); os.IsNotExist(err) {
				continue
			}
			if err := os.Remove(filename); err != nil {
				return err
			}
			delete(c.WallpaperMap, wallpaper)
		}
		c.Wallpapers = c.Wallpapers[10:]
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
		return ""
	}
	return filename
}

func (c *cacheData) putWallpaper(w Wallpaper, data io.Reader) (string, error) {
	if filename, exists := c.WallpaperMap[w]; exists {
		if _, err := os.Stat(filename); os.IsExist(err) {
			os.Remove(filename)
		}
	}
	filename := filepath.Join(cacheFolder.Path, fmt.Sprintf("cache%d.jpg", rand.Int()))

	cacheFolder.MkdirAll()
	f, err := os.Create(filename)
	if err != nil {
		return "", err
	}
	defer f.Close()
	io.Copy(f, data)

	c.WallpaperMap[w] = filename
	c.Wallpapers = append(c.Wallpapers, w)

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
