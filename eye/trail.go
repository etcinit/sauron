// Package eye provides the internals behind the Sauron CLI tool.
package eye

import (
	"strconv"
	"time"

	"github.com/ActiveState/tail"
	"github.com/Sirupsen/logrus"
	fsnotify "gopkg.in/fsnotify.v1"
)

// Line contains a log line of a log file.
type Line struct {
	Path string
	Text string
	Time time.Time
	Err  error
}

// LineHandler is a function capable to handle log lines.
type LineHandler func(line Line) error

// Trail represents a log trail that can be followed for new lines. In
// conjuction with a Watcher, a Trail is capable of monitoring existing and new
// files in a directory.
//
// However, unlike the Watcher, a Trail is limited to traditional filesystems.
type Trail struct {
	watcher Watcher
	done    chan bool
	tails   []*tail.Tail
	options *TrailOptions
}

// NewTrail creates a new instance of a Trail.
func NewTrail(watcher Watcher) *Trail {
	return &Trail{
		watcher: watcher,
		done:    make(chan bool),
		options: &TrailOptions{
			Logger: logrus.New(),
		},
	}
}

// NewTrailWithOptions creates a new instance of a Trail with a custom set of
// options. If any option provided is nil, it will be replaced with a safe
// default.
func NewTrailWithOptions(watcher Watcher, options *TrailOptions) *Trail {
	// Create a default set of options.
	defaults := &TrailOptions{
		Logger: logrus.New(),
	}

	// Replace the logger if an alternative is provided.
	if options.Logger != nil {
		defaults.Logger = options.Logger
	}

	return &Trail{
		watcher: watcher,
		done:    make(chan bool),
		options: defaults,
	}
}

// Follow starts following a trail. Every time a file is changed, the affected
// lines will be passed to the handler function to be proccessed. The handler
// function could do something as simple as writing the lines that standard
// output, or do more advanced things like writing to an external log server.
func (t *Trail) Follow(handler LineHandler) error {
	t.options.Logger.Infoln("Sauron is now watching")

	// First, we tail all the files that we already know.
	files, err := t.watcher.Walk()

	if err != nil {
		t.options.Logger.Errorln("Failed to walk directory")

		return err
	}

	for _, file := range files {
		t.followFile(file, handler, false)
	}

	// Second, we watch for new files, and tail them too.
	events := make(chan FileEvent)

	go func() {
		for {
			select {
			case event := <-events:
				switch event.Op {
				case fsnotify.Create:
					t.options.Logger.Infoln("Created: " + event.Path)

					t.followFile(event.Path, handler, true)
				case fsnotify.Remove:
					t.options.Logger.Infoln("Removed: " + event.Path)
				case fsnotify.Rename:
					t.options.Logger.Infoln("Renamed: " + event.Path)
				case fsnotify.Write:
					t.options.Logger.Infoln("Write: " + event.Path)
				default:
					t.options.Logger.Infoln(
						"Event " + strconv.Itoa(int(event.Op)) + ": " + event.Path,
					)
				}
			case <-t.done:
				// Stop the watcher
				t.watcher.End()

				// Stop any tailers
				for _, current := range t.tails {
					current.Stop()
				}

				// Exit the goroutine
				return
			}
		}
	}()

	t.watcher.Watch(events)

	return nil
}

// End stops watching.
func (t *Trail) End() {
	t.options.Logger.Infoln("Stopping...")

	t.done <- true
}

// followFile simply setups the appropriate options for the tail library and
// starts tailing that file. It also repackages events as Line objects for the
// handler function. The isNew parameter tells the function whether the file
// was just created or it already existed when the trail started following.
func (t *Trail) followFile(path string, handler LineHandler, isNew bool) {
	t.options.Logger.Infoln("Following: " + path)

	go func() {
		var current *tail.Tail
		var err error

		if isNew {
			current, err = tail.TailFile(path, tail.Config{
				Follow: true,
				Logger: tail.DiscardingLogger,
			})

			if err != nil {
				return
			}
		} else {
			current, err = tail.TailFile(path, tail.Config{
				Follow:   true,
				Location: &tail.SeekInfo{Offset: 0, Whence: 2},
				Logger:   tail.DiscardingLogger,
			})

			if err != nil {
				return
			}
		}

		t.tails = append(t.tails, current)

		for line := range current.Lines {
			newLine := Line{
				Path: path,
				Text: line.Text,
				Time: line.Time,
				Err:  line.Err,
			}

			handler(newLine)
		}
	}()
}
