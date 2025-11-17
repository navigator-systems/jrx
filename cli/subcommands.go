package cli

import (
	"github.com/navigator-systems/jrx/cmd"
	"github.com/navigator-systems/jrx/patterns"

	"github.com/urfave/cli/v2"
)

// project SubCommands

var newCmd = &cli.Command{
	Name:    "new",
	Aliases: []string{"n"},
	Usage:   "Create a new project",
	Action: func(c *cli.Context) error {
		name := c.Args().Get(0)
		template := c.Args().Get(1)
		cmd.NewCmd(name, template, gitOrg)
		return nil
	},
	Flags: []cli.Flag{
		flagGitOrg,
	},
}


var infoCmd = &cli.Command{
	Name:    "info",
	Aliases: []string{"i"},

	Usage: "Get information from the project",
	Action: func(c *cli.Context) error {
		name := c.Args().Get(0)
		cmd.InfoCmd(name, osvFlag)
		return nil
	},
	Flags: []cli.Flag{
		flagOSV,
	},
}


// Templates SubCommands
var tmplInfoCmd = &cli.Command{
	Name: "list",

	Usage: "Get information about templates",
	Action: func(c *cli.Context) error {
		cmd.TmplInfoCmd()
		return nil
	},
}

var tmplDownloadCmd = &cli.Command{
	Name:  "download",
	Usage: "Download the templates for a new project",
	Action: func(c *cli.Context) error {

		patterns.InitTemplates()
		return nil
	},
}
