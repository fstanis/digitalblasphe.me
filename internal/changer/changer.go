package changer

import (
	"math/rand"

	setbackground "github.com/fstanis/setbackground/go"

	"github.com/fstanis/digitalblasphe.me/internal/config"
	"github.com/fstanis/digitalblasphe.me/pkg/digitalblasphemy"
)

// Apply changes the desktop background based on information from the config
// file.
func Apply(conf *config.Config) error {
	switch {
	case conf.UseFree:
		return changeToFree()
	case conf.Random:
		return changeToRandom(conf)
	default:
		return changeToLatest(conf)
	}
}

func changeToFree() error {
	index, err := digitalblasphemy.GetFreebiesIndex()
	if err != nil {
		return err
	}

	return downloadAndSetBackground(index[rand.Intn(len(index))], nil)
}

func changeToRandom(conf *config.Config) error {
	index, creds, err := getIndex(conf)
	if err != nil {
		return err
	}

	return downloadAndSetBackground(index[rand.Intn(len(index))], creds)
}

func changeToLatest(conf *config.Config) error {
	index, creds, err := getIndex(conf)
	if err != nil {
		return err
	}

	return downloadAndSetBackground(index[0], creds)
}

func getIndex(conf *config.Config) ([]digitalblasphemy.Wallpaper, *digitalblasphemy.Credentials, error) {
	creds, err := conf.LoadCredentials()
	if err != nil {
		return nil, nil, err
	}
	index, err := digitalblasphemy.GetIndex(conf.Resolution, creds)
	if err != nil {
		return nil, nil, err
	}
	return index, creds, nil
}

func downloadAndSetBackground(wallpaper digitalblasphemy.Wallpaper, creds *digitalblasphemy.Credentials) error {
	filename, err := digitalblasphemy.FetchWallpaper(wallpaper, creds)
	if err != nil {
		return nil
	}

	return setbackground.SetBackground(filename, setbackground.StyleCenter)
}
