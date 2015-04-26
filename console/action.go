package console

import (
	"fmt"

	"github.com/codegangsta/cli"
	"github.com/etcinit/sauron/eye"
)

// MainAction is the main action exceuted when using Sauron.
func MainAction(c *cli.Context) {
	directory := "."

	if c.Args().First() != "" {
		directory = c.Args().First()
	}

	watcher := eye.NewDirectoryWatcher(directory)

	trail := eye.NewTrail(watcher)
	trail.Follow(func(line eye.Line) error {
		fmt.Println(line.Path, line.Text)

		return nil
	})
}
