package config

import (
	"fmt"
	"path/filepath"

	"github.com/BurntSushi/toml"
)

type Cfg struct {
	ProjectName string             `toml:"name"`
	Version     string             `toml:"version,omitempty"`
	Authors     []string           `toml:"authors,omitempty"`
	Team        string             `toml:"team,omitempty"`
	Packages    map[string]Package `toml:"packages,omitempty"`
}

type Package struct {
	Type    string `toml:"type,omitempty"`
	Command string `toml:"command,omitempty"`
}

func ReadCfgFile(path string) (Cfg, error) {
	var jrxConfig Cfg
	cfgFile := filepath.Join(path, "jrx.toml")
	_, err := toml.DecodeFile(cfgFile, &jrxConfig)
	if err != nil {
		fmt.Println("Error reading config file:", err)
		return jrxConfig, err
	}
	return jrxConfig, nil
}
