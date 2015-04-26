package console

import (
	"fmt"
	"os"
	"os/signal"

	"github.com/codegangsta/cli"
	"github.com/etcinit/sauron/eye"
)

// MainAction is the main action exceuted when using Sauron.
func MainAction(c *cli.Context) {
	done := make(chan bool)

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

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt)
	go func() {
		for sig := range signalChan {
			if sig == os.Interrupt || sig == os.Kill {
				trail.End()
				done <- true
			}
		}
	}()

	<-done
}
