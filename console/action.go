package console

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/signal"

	"github.com/Sirupsen/logrus"
	"github.com/codegangsta/cli"
	"github.com/etcinit/sauron/eye"
)

// MainAction is the main action exceuted when using Sauron.
func MainAction(c *cli.Context) {
	done := make(chan bool)

	// Decide which directory to follow.
	directory := "."
	if c.Args().First() != "" {
		directory = c.Args().First()
	}

	watcher := eye.NewDirectoryWatcher(directory)

	options := &eye.TrailOptions{}

	// Decide whether to output logs.
	if c.Bool("verbose") {
		options.Logger = logrus.New()
	} else {
		options.Logger = log.New(ioutil.Discard, "", 0)
	}

	// Create the new instance of the trail and begin following it.
	trail := eye.NewTrailWithOptions(watcher, options)
	trail.Follow(func(line eye.Line) error {
		fmt.Println(line.Path, line.Text)

		return nil
	})

	// Wait for an interrupt or kill signal.
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
