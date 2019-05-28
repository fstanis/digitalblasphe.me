package main

import (
	"log"
	"math/rand"
	"os"
	"time"

	"github.com/urfave/cli"

	"github.com/fstanis/digitalblasphe.me/internal/actions"
	"github.com/fstanis/digitalblasphe.me/internal/config"
)

func main() {
	rand.Seed(time.Now().UnixNano())

	app := cli.NewApp()
	app.Name = "digitalblasphe.me"
	app.Usage = "changes the desktop wallpaper by downloading it from digitalblasphemy.com"
	app.Action = actions.Default
	app.Commands = []cli.Command{
		actions.Apply,
		actions.Configure,
		actions.Install,
	}

	if conf, err := config.Load(); err == nil {
		actions.Config = conf
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
