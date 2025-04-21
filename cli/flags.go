package cli

import (
	"github.com/urfave/cli/v2"
)

var flagOSV = &cli.BoolFlag{
	Name:        "osv",
	Aliases:     []string{"o"},
	Usage:       "Check if the packages have known vulnerabilities",
	Destination: &osvFlag,
}

var flagArch = &cli.StringFlag{
	Name:        "arch",
	Usage:       "Architecture to build for",
	Destination: &archFlag,
}

var flagOS = &cli.StringFlag{
	Name:        "os",
	Usage:       "Operating system to build for",
	Destination: &osFlag,
}
