package cli

import (
	"github.com/urfave/cli/v2"
)

var (
	osvFlag  bool
	archFlag string
	osFlag   string
	gitOrg   string
)

var projectCmd = &cli.Command{
	Name:    "project",
	Aliases: []string{"p"},
	Usage:   "Manage projects",
	Subcommands: []*cli.Command{
		newCmd,
		buildCmd,
		cleanCmd,
		infoCmd,
	},
}

var templatesCmd = &cli.Command{
	Name:    "templates",
	Aliases: []string{"t"},
	Usage:   "Manage project templates",
	Subcommands: []*cli.Command{
		tmplInfoCmd,
		tmplDownloadCmd,
	},
}
