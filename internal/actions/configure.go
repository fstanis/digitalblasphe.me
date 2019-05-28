package actions

import (
	"fmt"
	"time"

	"github.com/AlecAivazis/survey"
	"github.com/fstanis/screenresolution"
	"github.com/urfave/cli"

	"github.com/fstanis/digitalblasphe.me/internal/config"
	"github.com/fstanis/digitalblasphe.me/pkg/digitalblasphemy"
)

// Configure asks the user a series of questions to build a config file.
var Configure = cli.Command{
	Name:    "configure",
	Aliases: []string{"c"},
	Usage:   "run the initial configuration",
	Action:  configure,
}

func configure(c *cli.Context) error {
	conf, err := configFromSurvey()
	if err != nil {
		return err
	}
	conf.Save()

	Config = conf

	if prompt("Would you like to run now?", true) {
		return c.App.Command("apply").Run(c)
	}

	return nil
}

func prompt(question string, def bool) bool {
	var result bool
	survey.AskOne(&survey.Confirm{
		Message: question,
		Default: def,
	}, &result, nil)
	return result
}

func configFromSurvey() (*config.Config, error) {
	isMember := askIsMember()
	conf := &config.Config{
		UseFree: true,
	}
	if isMember {
		if err := askCredentials(conf); err != nil {
			return nil, err
		}
		if err := askPreferences(conf); err != nil {
			return nil, err
		}
	}
	if conf.Random {
		conf.Frequency = askFrequency()
	}

	return conf, nil
}

func askIsMember() bool {
	return prompt("Are you a member of digitalblasphemy.com?", true)
}

var credentialsQuestions = []*survey.Question{
	{
		Name:     "username",
		Prompt:   &survey.Input{Message: "digitalblasphemy.com username:"},
		Validate: survey.Required,
	},
	{
		Name:     "password",
		Prompt:   &survey.Password{Message: "digitalblasphemy.com password:"},
		Validate: survey.Required,
	},
}

func askCredentials(conf *config.Config) error {
	answers := struct {
		Username string
		Password string
	}{}
	if err := survey.Ask(credentialsQuestions, &answers); err != nil {
		return err
	}

	creds := &digitalblasphemy.Credentials{
		Username: answers.Username,
		Password: answers.Password,
	}
	if err := creds.Validate(); err != nil {
		fmt.Println("Invalid credentils")
		return askCredentials(conf)
	}
	conf.UseFree = false
	return conf.SaveCredentials(creds)
}

var preferencesQuestions = []*survey.Question{
	{
		Name: "resolution",
		Prompt: &survey.Select{
			Message: "Select your desktop resolution (autodetected):",
			Options: digitalblasphemy.GetValidResolutions(),
			Default: screenresolution.GetPrimary().String(),
		},
	},
	{
		Name: "type",
		Prompt: &survey.Select{
			Message: "Would you like to get random or latest wallpapers?",
			Options: []string{
				"latest",
				"random",
			},
			Default: "latest",
		},
	},
}

func askPreferences(conf *config.Config) error {
	answers := struct {
		Resolution string
		Type       string
	}{}

	if err := survey.Ask(preferencesQuestions, &answers); err != nil {
		return err
	}

	conf.Resolution = answers.Resolution
	conf.Random = answers.Type == "random"
	return nil
}

var frequencies = map[string]time.Duration{
	"15 minutes":      15 * time.Minute,
	"30 minutes":      30 * time.Minute,
	"1 hour":          1 * time.Hour,
	"2 hours":         2 * time.Hour,
	"6 hours":         6 * time.Hour,
	"only on startup": 0,
}

func askFrequency() time.Duration {
	var freq string
	survey.AskOne(&survey.Select{
		Message: "How often would you like to update the wallpaper?",
		Options: func() []string {
			result := make([]string, 0, len(frequencies))
			for key := range frequencies {
				result = append(result, key)
			}
			return result
		}(),
		Default: "1 hour",
	}, &freq, nil)
	return frequencies[freq]
}
