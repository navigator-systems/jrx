package cli

import (
	"github.com/navigator-systems/jrx/cmd"
	"github.com/navigator-systems/jrx/server"

	"github.com/urfave/cli/v2"
)

var (
	osvFlag  bool
	archFlag string
	osFlag   string
	ciFlag   string
)

var newCmd = &cli.Command{
	Name:    "new",
	Aliases: []string{"n"},
	Usage:   "Create a new project",
	Action: func(c *cli.Context) error {
		name := c.Args().Get(0)
		cmd.NewCmd(name)
		return nil
	},
}

var buildCmd = &cli.Command{
	Name:    "build",
	Aliases: []string{"b"},
	Usage:   "Build and compile a project",
	Action: func(c *cli.Context) error {
		name := c.Args().Get(0)
		cmd.BuildCmd(name, archFlag, osFlag)
		return nil
	},
	Flags: []cli.Flag{
		flagArch,
		flagOS,
	},
}

var modCmd = &cli.Command{
	Name:    "mod",
	Aliases: []string{"m"},
	Usage:   "Start a simple go.mod file",
	Action: func(c *cli.Context) error {
		name := c.Args().Get(0)
		cmd.ModCmd(name)
		return nil
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

var cleanCmd = &cli.Command{
	Name:    "clean",
	Aliases: []string{"c"},
	Usage:   "Clean the project binaries",
	Action: func(c *cli.Context) error {
		name := c.Args().Get(0)
		cmd.CleanCmd(name)
		return nil
	},
}

var ciCmd = &cli.Command{
	Name:  "ci",
	Usage: "add a CI template (Jenkins, GitHub Actions or Gitlab Template) to the project",
	Action: func(c *cli.Context) error {
		name := c.Args().Get(0)
		cmd.AddCICmd(name, ciFlag)
		return nil
	},
	Flags: []cli.Flag{
		templateCI,
	},
}

var serverCmd = &cli.Command{
	Name:    "server",
	Aliases: []string{"s"},
	Usage:   "Start a simple web server",
	Action: func(c *cli.Context) error {
		server.StartServer()
		return nil
	},
}
