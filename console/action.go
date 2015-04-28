// Package console contains Sauron's CLI commands.
package console

import (
	"fmt"
	"io/ioutil"
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
		log := logrus.New()

		log.Out = ioutil.Discard

		options.Logger = log
	}

	// Create the new instance of the trail and begin following it.
	trail := eye.NewTrailWithOptions(watcher, options)
	err := trail.Follow(func(line eye.Line) error {
		output := ""

		if c.BoolT("prefix-path") {
			output += "[" + line.Path + "] "
		}

		if c.Bool("prefix-time") {
			output += "[" + line.Time.Format("Jan 2, 2006 at 3:04pm (MST)") + "] "
		}

		output += line.Text

		fmt.Println(output)

		return nil
	})

	if err != nil {
		return
	}

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
