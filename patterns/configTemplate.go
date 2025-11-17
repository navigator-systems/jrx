package patterns

// Root: Template definition
type RootTemplate struct {
	ProjectName string
	Name        string   `toml:"name"`
	Description string   `toml:"description"`
	Path        string   `toml:"path"`
	Tags        []string `toml:"tags"`
	ProjectInfo ProjectTemplate
	Variables   []VariablesTemplate `toml:"variables"`
}

type ProjectTemplate struct {
	Language        string `toml:"language"`
	LanguageVersion string `toml:"language_version,omitempty"`
	Entry           string `toml:"entry"`
	AppVersion      string `toml:"appversion,omitempty"`
}

// Metadata fields used to substitute inside template files.
type VariablesTemplate struct {
	Key         string `toml:"key"`
	Description string `toml:"description"`
	Default     string `toml:"default,omitempty"`
}

type TemplateFile struct {
	Templates map[string]RootTemplate `toml:"templates"`
}
