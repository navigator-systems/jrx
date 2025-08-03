package patterns

type BuildTarget struct {
	Arch  string `toml:"arch"`
	OS    string `toml:"os"`
	Flags string `toml:"flags,omitempty"`
}

type Template struct {
	ProjectName     string
	Name            string                 `toml:"name"`
	Description     string                 `toml:"description"`
	Language        string                 `toml:"language"`
	LanguageVersion string                 `toml:"language_version,omitempty"`
	Entry           string                 `toml:"entry"`
	Path            string                 `toml:"path"`
	AppVersion      string                 `toml:"appversion,omitempty"`
	Tags            []string               `toml:"tags"`
	Builds          map[string]BuildTarget `toml:"builds,omitempty"`
	Execute         string                 `toml:"execute,omitempty"`
}

type TemplateFile struct {
	Templates map[string]Template `toml:"templates"`
}
