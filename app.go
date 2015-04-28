// Package Sauron provides tools for monitoring changes on existing and new
// files inside a directory.
//
// After installation, the CLI tool should be available as:
//
// 	sauron
//
// Running the command without any parameters will cause sauron to watch the
// current directory for changes. Any new lines appended to any files will be
// printed out.
//
// It is possible to specify which directory to watch:
//
// 	sauron /var/log/hhvm
//
// For more detailed output, use the verbose option:
//
// 	sauron --verbose /var/log/nginx
//
// Most of the code for this tool is available as a standalone package,
// checkout the https://github.com/etcinit/sauron/tree/master/eye package.
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
	app.Version = "0.0.3"
	app.Author = "Eduardo Trujillo <ed@chromabits.com>"

	cli.VersionPrinter = func(c *cli.Context) {
		fmt.Fprintf(c.App.Writer, "%v", c.App.Version)
	}

	app.Flags = []cli.Flag{
		cli.BoolFlag{
			Name:  "verbose",
			Usage: "print verbose output along with logs",
		},
		cli.BoolTFlag{
			Name:  "prefix-path",
			Usage: "prefix file path to every output line",
		},
		cli.BoolFlag{
			Name:  "prefix-time",
			Usage: "prefix time to every output line",
		},
	}

	// Setup the default action. This action will be triggered when no
	// subcommand is provided as an argument
	app.Action = console.MainAction

	// Begin
	app.Run(os.Args)
}
