package templategit

type BuildTarget struct {
	Arch  string `toml:"arch"`
	OS    string `toml:"os"`
	Flags string `toml:"flags,omitempty"`
}

type Template struct {
	Name        string                 `toml:"name"`
	Description string                 `toml:"description"`
	Language    string                 `toml:"language"`
	Entry       string                 `toml:"entry"`
	Path        string                 `toml:"path"`
	Version     string                 `toml:"version"`
	Tags        []string               `toml:"tags"`
	Builds      map[string]BuildTarget `toml:"builds,omitempty"`
	Execute     string                 `toml:"execute,omitempty"`
}

type TemplateFile struct {
	Templates map[string]Template `toml:"templates"`
}
