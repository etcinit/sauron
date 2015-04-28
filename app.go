package main

import (
	"fmt"
	"os"

	"github.com/codegangsta/cli"
	"github.com/etcinit/sauron/console"
)

func main() {
	// Setup the command line application
	app := cli.NewApp()
	app.Name = "sauron"
	app.Usage = "Utility for monitoring files in a directory"

	// Set version and authorship info
	app.Version = "0.0.2"
	app.Author = "Eduardo Trujillo <ed@chromabits.com>"

	cli.VersionPrinter = func(c *cli.Context) {
		fmt.Fprintf(c.App.Writer, "%v", c.App.Version)
	}

	app.Flags = []cli.Flag{
		cli.BoolFlag{
			Name: "verbose",
		},
	}

	// Setup the default action. This action will be triggered when no
	// subcommand is provided as an argument
	app.Action = console.MainAction

	// Begin
	app.Run(os.Args)
}
