package actions

import (
	"errors"
	"io"
	"os"

	"github.com/ProtonMail/go-autostart"
	"github.com/urfave/cli"

	"github.com/fstanis/digitalblasphe.me/internal/config"
)

// Install copies the executable and makes it run on startup.
var Install = cli.Command{
	Name:    "install",
	Aliases: []string{"i"},
	Usage:   "installs the app to run on startup",
	Action:  install,
	Flags: []cli.Flag{
		cli.BoolFlag{
			Name:  "uninstall, u",
			Usage: "uninstall the app",
		},
		cli.BoolFlag{
			Name:  "force, f",
			Usage: "reinstall if already installed",
		},
	},
}

var autostartApp = &autostart.App{
	Name:        "digitalblasphe.me",
	DisplayName: "digitalblasphe.me",
}

func install(c *cli.Context) error {
	app := autostartApp

	targetPath := config.FilePath(c.App.Name)
	app.Exec = []string{targetPath}

	if isInstalled() {
		if c.Bool("uninstall") {
			os.Remove(targetPath)
			return app.Disable()
		}

		if c.Bool("force") {
			app.Disable()
		} else {
			return errors.New("App is already installed. Use -u to uninstall or -f to force reinstall.")
		}
	}

	if c.Bool("uninstall") {
		return errors.New("App is not installed.")
	}

	sourcePath := os.Args[0]
	if err := copyFile(sourcePath, targetPath); err != nil {
		return err
	}
	return app.Enable()
}

func isInstalled() bool {
	return autostartApp.IsEnabled()
}

func copyFile(sourcePath string, targetPath string) error {
	source, err := os.Open(sourcePath)
	if err != nil {
		return err
	}
	defer source.Close()

	target, err := os.OpenFile(targetPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0755)
	if err != nil {
		return err
	}
	defer target.Close()

	if _, err := io.Copy(target, source); err != nil {
		return err
	}

	return nil
}
