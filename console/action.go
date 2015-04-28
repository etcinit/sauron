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

	options := &eye.TrailOptions{}

	// Decide whether to output logs.
	if c.Bool("verbose") {
		options.Logger = logrus.New()
	} else {
		log := logrus.New()

		log.Out = ioutil.Discard

		options.Logger = log
	}

	// Decide which directories to follow.
	directories := []string{"."}
	if c.Args().First() != "" {
		directories = []string{c.Args().First()}
	}

	// Add any additional directories provided.
	for _, other := range c.Args().Tail() {
		directories = append(directories, other)
	}

	var trails []*eye.Trail

	for _, directory := range directories {
		watcher, err := eye.NewDirectoryWatcher(directory)

		if err != nil {
			logrus.Errorln(err)
			return
		}

		// Create the new instance of the trail and begin following it.
		trail := eye.NewTrailWithOptions(watcher, options)
		err = trail.Follow(getHandler(c))

		if err != nil {
			logrus.Errorln(err)
			return
		}

		trails = append(trails, trail)
	}

	// Wait for an interrupt or kill signal.
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt)
	go func() {
		for sig := range signalChan {
			if sig == os.Interrupt || sig == os.Kill {
				for _, trail := range trails {
					trail.End()
				}
				done <- true
			}
		}
	}()

	<-done
}

// getHandler builds the handler function to be used while following a trail.
func getHandler(c *cli.Context) eye.LineHandler {
	return func(line eye.Line) error {
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
	}
}
