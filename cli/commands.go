package cli

import (
	"github.com/urfave/cli/v2"
)

var (
	gitOrg   string
	varsFlag string
)

var projectCmd = &cli.Command{
	Name:    "project",
	Aliases: []string{"p"},
	Usage:   "Manage projects",
	Subcommands: []*cli.Command{
		newCmd,
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
