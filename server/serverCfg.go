package server

import (
	"log"

	"github.com/BurntSushi/toml"
)

type ServerConfig struct {
	Port           string `toml:"port"`
	TemplateRepo   string `toml:"template_repo"`
	TemplateBranch string `toml:"template_branch"`
	TemplateFile   string `toml:"template_file"`
}

func GetConfig() (*ServerConfig, error) {
	var serverCfg ServerConfig
	if _, err := toml.DecodeFile("jrxcfg.toml", &serverCfg); err != nil {
		log.Fatal("Error decoding template file:", err)
		return &serverCfg, err
	}

	return &serverCfg, nil

}
