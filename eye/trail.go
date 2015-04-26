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

type LineHandler func(line Line) error

type Trail struct {
	watcher Watcher
	done    chan bool
}

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
		for event := range events {
			if event.Op == fsnotify.Create {
				t.followFile(event.Path, handler, true)
			}
		}
	}()

	t.watcher.Watch(events)

	<-t.done

	return nil
}

func (t *Trail) followFile(path string, handler LineHandler, isNew bool) {
	go func() {
		var t *tail.Tail

		if isNew {
			t, _ = tail.TailFile(path, tail.Config{
				Follow: true,
				Logger: tail.DiscardingLogger,
			})
		} else {
			t, _ = tail.TailFile(path, tail.Config{
				Follow:   true,
				Location: &tail.SeekInfo{Offset: 0, Whence: 2},
				Logger:   tail.DiscardingLogger,
			})
		}

		for line := range t.Lines {
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
