package eye

import (
	"time"

	"github.com/ActiveState/tail"
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

// Trail represents a log trail.
type Trail struct {
	watcher Watcher
	done    chan bool
	tails   []*tail.Tail
}

// NewTrail creates a new instance of a Trail.
func NewTrail(watcher Watcher) *Trail {
	return &Trail{
		watcher: watcher,
		done:    make(chan bool),
	}
}

func (t *Trail) Follow(handler LineHandler) error {
	// First, we tail all the files that we already know.
	files, err := t.watcher.Walk()

	if err != nil {
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
				if event.Op == fsnotify.Create {
					t.followFile(event.Path, handler, true)
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
	t.done <- true
}

func (t *Trail) followFile(path string, handler LineHandler, isNew bool) {
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
