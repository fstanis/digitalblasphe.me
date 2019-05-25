package config

import (
	"fmt"

	"github.com/AlecAivazis/survey"
	"github.com/fstanis/screenresolution"

	"github.com/fstanis/digitalblasphe.me/pkg/digitalblasphemy"
)

var qs = []*survey.Question{
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
	{
		Name: "resolution",
		Prompt: &survey.Select{
			Message: "Select your desktop resolution (autodetected):",
			Options: digitalblasphemy.GetValidResolutions(),
			Default: screenresolution.Get().String(),
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

// FromSurvey generates a config file by asking the user interactively.
func FromSurvey() (*Config, error) {
	var isMember bool
	survey.AskOne(&survey.Confirm{
		Message: "Are you a member of digitalblasphemy.com?",
		Default: true,
	}, &isMember, nil)
	if !isMember {
		fmt.Println("")
		conf := &Config{
			UseFree: true,
		}
		return conf, nil
	}

	answers := struct {
		Username   string
		Password   string
		Resolution string
		Type       string
	}{}

	if err := survey.Ask(qs, &answers); err != nil {
		return nil, err
	}
	fmt.Println("")

	conf := &Config{
		UseFree:    false,
		Resolution: answers.Resolution,
		Random:     answers.Type == "random",
	}
	creds := &digitalblasphemy.Credentials{
		Username: answers.Username,
		Password: answers.Password,
	}
	if err := creds.Validate(); err != nil {
		return nil, fmt.Errorf("failed to validate credentils: %v", err)
	}
	err := conf.SaveCredentials(creds)
	if err != nil {
		return nil, err
	}

	return conf, nil
}
