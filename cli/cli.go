package cli

import (
	"log"
	"os"

	"github.com/urfave/cli/v2"
)

func InitCli() {
	app := &cli.App{
		Name:  "jrx",
		Usage: "Just a simple go wrapper CLI",
		Commands: []*cli.Command{
			newCmd,
			buildCmd,
			cleanCmd,
			modCmd,
			infoCmd,
		},
	}
	if err := app.Run(os.Args); err != nil {
		log.Fatalln(err)
	}
}
