package ops

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"
)

type Config struct {
	ProjectName string           `toml:"name"`
	Version     string           `toml:"version,omitempty"`
	Authors     []string         `toml:"authors,omitempty"`
	Builds      map[string]Build `toml:"builds,omitempty"`
}

type Build struct {
	Arch  string `toml:"arch"`
	OS    string `toml:"os"`
	Flags string `toml:"flags"`
}

func ReadCfgFile(path string) (Config, error) {
	var jrxConfig Config
	cfgFile := filepath.Join(path, "jrx.toml")
	_, err := toml.DecodeFile(cfgFile, &jrxConfig)
	if err != nil {
		fmt.Println("Error reading config file:", err)
		return jrxConfig, err
	}
	return jrxConfig, nil
}

func CheckIfCfgExists(path string) bool {
	cfgFile := filepath.Join(path, "jrx.toml")

	_, err := os.Stat(cfgFile)
	return !os.IsNotExist(err)
}
