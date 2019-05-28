package config

import (
	"log"
	"path/filepath"

	"github.com/shibukawa/configdir"
	yaml "gopkg.in/yaml.v2"
)

const (
	appName        = "digitalblasphe.me"
	keyringService = appName
)

var (
	// Filename is the filename of the config file.
	Filename = "config.yaml"

	configFolder *configdir.Config
)

func init() {
	dirs := configdir.New("", appName).QueryFolders(configdir.Global)
	if len(dirs) > 0 {
		dirs[0].MkdirAll()
		configFolder = dirs[0]
	} else {
		log.Fatal("Failed to find global config folder")
	}
}

// FilePath returns a path in the same directory as the config file.
func FilePath(filename string) string {
	return filepath.Join(configFolder.Path, filename)
}

func load(conf interface{}) error {
	data, err := configFolder.ReadFile(Filename)
	if err != nil {
		return err
	}

	return yaml.Unmarshal(data, conf)
}

func save(conf interface{}) error {
	data, err := yaml.Marshal(conf)
	if err != nil {
		return err
	}

	return configFolder.WriteFile(Filename, data)
}
