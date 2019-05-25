package config

import (
	"errors"
	"log"

	keyring "github.com/zalando/go-keyring"

	"github.com/fstanis/digitalblasphe.me/pkg/digitalblasphemy"
)

var (
	// ErrNoUsername happens when the config file contains no username.
	ErrNoUsername = errors.New("no username stored in config")

	// ErrInvalidConfig happens when the config file is invalid.
	ErrInvalidConfig = errors.New("invalid config")
)

// Config contains the config file for the app.
type Config struct {
	UseFree    bool
	Username   string
	Resolution string
	Random     bool
}

// LoadConfig loads the config file from the given location.
func LoadConfig() (*Config, error) {
	var conf Config

	if err := load(configPath, &conf); err != nil {
		return nil, err
	}

	if !conf.Validate() {
		return nil, ErrInvalidConfig
	}

	log.Printf("Loaded config from: %s", configPath)

	return &conf, nil
}

// Save saves the config file to the given location.
func (c *Config) Save() error {
	if !c.Validate() {
		return ErrInvalidConfig
	}

	if err := save(configPath, c); err != nil {
		return err
	}

	log.Printf("Saved config to: %s", configPath)

	return nil
}

// Validate checks if the config object is valid.
func (c *Config) Validate() bool {
	if c.UseFree {
		return true
	}
	if c.Username == "" {
		return false
	}

	if !digitalblasphemy.IsValidResolution(c.Resolution) {
		return false
	}

	return true
}

// LoadCredentials constructs the Credentials object using the data from the
// given Config combined with the password stored in the (platform dependent)
// keyring.
func (c *Config) LoadCredentials() (*digitalblasphemy.Credentials, error) {
	if c.Username == "" {
		return nil, ErrNoUsername
	}

	pass, err := keyring.Get(keyringService, c.Username)
	if err != nil {
		return nil, err
	}

	return &digitalblasphemy.Credentials{
		Username: c.Username,
		Password: pass,
	}, nil
}

// SaveCredentials saves the data from the given Credentials object to the
// specified config while also storing the password in the (platform dependent)
// keyring.
func (c *Config) SaveCredentials(creds *digitalblasphemy.Credentials) error {
	c.Username = creds.Username
	return keyring.Set(keyringService, creds.Username, creds.Password)
}
