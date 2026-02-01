package cli

import (
	"github.com/urfave/cli/v2"
)

var flagVars = &cli.StringFlag{
	Name:        "vars",
	Aliases:     []string{"v"},
	Usage:       "Variables for template in format key1=value1,key2=value2",
	Destination: &varsFlag,
}

var flagPort = &cli.StringFlag{
	Name:    "port",
	Aliases: []string{"p"},
	Usage:   "Port to run the server on",
	Value:   "8080",
}

var flagGitHubOrg = &cli.StringFlag{
	Name:        "github-organization",
	Aliases:     []string{"g"},
	Usage:       "GitHub organization URL to operate on",
	Destination: &gitHubOrg,
}
