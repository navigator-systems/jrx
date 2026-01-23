package cli

import (
	"github.com/navigator-systems/jrx/cmd"
	"github.com/navigator-systems/jrx/internal/templates"

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

		cmd.NewCmd(name, template, varsFlag)

		return nil
	},
	Flags: []cli.Flag{
		flagVars,
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

		templates.InitTemplates()
		return nil
	},
}

// Server Command
var serverCmd = &cli.Command{
	Name:    "server",
	Aliases: []string{"s", "serve"},
	Usage:   "Start the JRX web server",
	Action: func(c *cli.Context) error {
		port := c.String("port")
		cmd.ServerCmd(port)
		return nil
	},
	Flags: []cli.Flag{
		flagPort,
	},
}
