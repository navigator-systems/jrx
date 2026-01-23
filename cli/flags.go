package cli

import (
	"github.com/urfave/cli/v2"
)

var flagGitOrg = &cli.StringFlag{
	Name:        "repository",
	Usage:       "Git repository to operate on",
	Destination: &gitOrg,
}

var flagVars = &cli.StringFlag{
	Name:        "vars",
	Aliases:     []string{"v"},
	Usage:       "Variables for template in format key1=value1,key2=value2",
	Destination: &varsFlag,
}
