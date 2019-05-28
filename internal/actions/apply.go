package actions

import (
	"errors"
	"fmt"
	"log"
	"math/rand"
	"time"

	setbackground "github.com/fstanis/setbackground/go"
	"github.com/urfave/cli"

	"github.com/fstanis/digitalblasphe.me/internal/config"
	"github.com/fstanis/digitalblasphe.me/pkg/digitalblasphemy"
)

const minimumFrequency = time.Minute * 10

// Apply changes the desktop background based on information from the config
// file.
var Apply = cli.Command{
	Name:    "apply",
	Aliases: []string{"a"},
	Usage:   "change the desktop background",
	Action:  apply,
	Flags: []cli.Flag{
		cli.BoolFlag{
			Name:  "now, n",
			Usage: "ignore configured frequency and exit after updating the background",
		},
	},
}

var currentWallpaper string

func apply(c *cli.Context) error {
	if Config == nil {
		return errors.New("Config not found")
	}

	digitalblasphemy.LoadCache()
	currentWallpaper = digitalblasphemy.GetCachedCurrent()
	fmt.Println(currentWallpaper)
	dur := Config.Frequency
	if dur > 0 && dur < time.Hour*6 && !Config.Random && !Config.UseFree {
		dur = time.Hour * 6
	}
	if c.Bool("now") || dur == 0 {
		return applyAction()
	}

	t := time.NewTimer(dur)
	for {
		applyAction()
		<-t.C
		t.Reset(dur)
	}
}

func applyAction() error {
	defer digitalblasphemy.SaveCache()

	switch {
	case Config.UseFree:
		return changeToFree()
	case Config.Random:
		return changeToRandom(Config)
	default:
		return changeToLatest(Config)
	}
}

func changeToFree() error {
	index, err := digitalblasphemy.GetFreebiesIndex()
	if err != nil {
		return err
	}

	i := rand.Intn(len(index))
	if index[i].ID == currentWallpaper {
		i = (i + 1) % len(index)
	}

	return downloadAndSetBackground(index[i], nil)
}

func changeToRandom(conf *config.Config) error {
	index, creds, err := getIndex(conf)
	if err != nil {
		return err
	}

	i := rand.Intn(len(index))
	if index[i].ID == currentWallpaper {
		i = (i + 1) % len(index)
	}

	return downloadAndSetBackground(index[i], creds)
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

	if err := setbackground.SetBackground(filename, setbackground.StyleCenter); err != nil {
		return err
	}

	currentWallpaper = wallpaper.ID
	digitalblasphemy.SetCachedCurrent(currentWallpaper)
	log.Printf("Changed wallpaper to %s\n", wallpaper.ID)
	return nil
}
