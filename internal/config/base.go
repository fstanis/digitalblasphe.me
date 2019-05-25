package config

import (
	"io/ioutil"
	"path/filepath"

	"github.com/shibukawa/configdir"
	yaml "gopkg.in/yaml.v2"
)

const (
	appName        = "digitalblasphe.me"
	configFileName = "config.yaml"
	keyringService = appName
)

var (
	// Directory is the target directory for the config file.
	configDirectory string
	configPath      = configFileName
)

func init() {
	dirs := configdir.New("", appName).QueryFolders(configdir.Global)
	if len(dirs) > 0 {
		dirs[0].MkdirAll()
		configDirectory = dirs[0].Path
		configPath = PathInFolder(configFileName)
	}
}

// PathInFolder returns a path in the same folder as the config file.
func PathInFolder(filename string) string {
	return filepath.Join(configDirectory, filename)
}

func load(filename string, conf interface{}) error {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}

	return yaml.Unmarshal(data, conf)
}

func save(filename string, conf interface{}) error {
	data, err := yaml.Marshal(conf)
	if err != nil {
		return err
	}

	return ioutil.WriteFile(filename, data, 0644)
}
