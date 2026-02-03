package config

import (
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"
)

// This config is used for jrx to know about the templates repository and the server configuration
// This is not for reading each jrx project, but for the jrx system itself
type JRXConfig struct {
	TemplatesRepo     string   `toml:"templates_repo"`
	TemplatesDefault  string   `toml:"templates_default,omitempty"`
	TemplatesBranch   []string `toml:"templates_branches"`
	TemplatesTag      []string `toml:"templates_tags"`
	TemplatesCacheDir string   `toml:"templates_cache_dir,omitempty"` // Cache directory for templates

	SshKeyPath       string         `toml:"ssh_key_path"`
	SshKeyPassphrase string         `toml:"ssh_key_passphrase,omitempty"`
	ServerPort       string         `toml:"server_port"`
	GitProvider      JRXGitProvider `toml:"git_provider"`
	Database         JRXDataBase    `toml:"data_base"`
}

type JRXDataBase struct {
	Database   string `toml:"database,omitempty"`    // sqlite or postgres
	DBHost     string `toml:"db_host,omitempty"`     // postgres
	DBPort     int    `toml:"db_port,omitempty"`     // postgres
	DBUser     string `toml:"db_user,omitempty"`     // postgres
	DBPassword string `toml:"db_password,omitempty"` // postgres
	DBName     string `toml:"db_name,omitempty"`     // postgres
	DBPath     string `toml:"db_path,omitempty"`     //sqlite
}

type JRXGitProvider struct {
	GithubToken        string   `toml:"github_token,omitempty"`
	GithubURL          string   `toml:"github_url,omitempty"`
	GithubOrganization []string `toml:"github_organization_url,omitempty"`

	GitlabToken string `toml:"gitlab_token,omitempty"`
	GitlabGroup string `toml:"gitlab_group,omitempty"`
}

func ReadJRXConfig() (JRXConfig, error) {
	var jrxConfig JRXConfig
	path := os.Getenv("HOME")
	cfgFile := filepath.Join(path, ".jrxrc")

	_, err := toml.DecodeFile(cfgFile, &jrxConfig)
	if err != nil {
		return jrxConfig, err
	}

	return jrxConfig, nil
}
