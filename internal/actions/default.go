package actions

import (
	"github.com/urfave/cli"

	"github.com/fstanis/digitalblasphe.me/internal/config"
)

// Config is the config file to use.
var Config *config.Config

// Default is the action used when no action is selected.
func Default(c *cli.Context) error {
	if Config != nil {
		return c.App.Command("apply").Run(c)
	}

	if !isInstalled() {
		if prompt("Would you like to run the app on startup?", true) {
			if err := c.App.Command("install").Run(c); err != nil {
				return err
			}
		}
	}
	return c.App.Command("configure").Run(c)
}
