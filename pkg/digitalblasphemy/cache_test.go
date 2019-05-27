package digitalblasphemy

import "testing"

func TestCache(t *testing.T) {
	w := Wallpaper{
		ID:         "cacheexample",
		URL:        "https://example.com/",
		Resolution: "1920x1080",
	}
	cache.WallpaperMap[w] = "test"

	if err := SaveCache(); err != nil {
		t.Fatal(err)
	}
	cache.WallpaperMap = make(map[Wallpaper]string)
	if err := LoadCache(); err != nil {
		t.Fatal(err)
	}
	if cache.WallpaperMap[w] != "test" {
		t.Errorf("expected cache to have value %q, instead got %q", "test", cache.WallpaperMap[w])
	}
}
