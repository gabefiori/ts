package config

import (
	"encoding/json"
	"os"

	"github.com/mitchellh/go-homedir"
)

type Config struct {
	File
	Cli
}

type Cli struct {
	List   bool
	Filter string
	Path   string
}

type File struct {
	Selector []string `json:"selector"`
	Targets  []Target `json:"targets"`
}

type Target struct {
	Path  string `json:"path"`
	Depth uint8  `json:"depth"`
}

func Load(cli Cli) (*Config, error) {
	path, err := homedir.Expand(cli.Path)

	if err != nil {
		return nil, err
	}

	fb, err := os.ReadFile(path)

	if err != nil {
		return nil, err
	}

	var file File

	if err := json.Unmarshal(fb, &file); err != nil {
		return nil, err
	}

	return &Config{file, cli}, nil
}
