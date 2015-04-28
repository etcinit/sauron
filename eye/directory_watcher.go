package eye

import (
	"errors"
	"os"
	"path/filepath"
	"time"

	fsnotify "gopkg.in/fsnotify.v1"
)

// DirectoryWatcher is an implementation of a Watcher capable of monitoring
// for changes on a directory recursively.
type DirectoryWatcher struct {
	path string
	done chan bool
}

// NewDirectoryWatcher creates a new instance of a DirectoryWatcher.
func NewDirectoryWatcher(path string) (*DirectoryWatcher, error) {
	fileInfo, err := os.Stat(path)

	if err != nil {
		return nil, err
	}

	if !fileInfo.IsDir() {
		return nil, errors.New("Unable to watch. Cannot watch a file.")
	}

	return &DirectoryWatcher{
		path: path,
	}, nil
}

// Walk returns a list of all the files within the target directory.
func (w *DirectoryWatcher) Walk() (paths []string, err error) {
	visit := func(path string, f os.FileInfo, err error) error {
		if f.IsDir() {
			return nil
		}

		abs, err := filepath.Abs(path)

		if err != nil {
			return err
		}

		paths = append(paths, abs)

		return nil
	}

	err = filepath.Walk(w.path, visit)

	return
}

// Watch starts watching for filesystem events.
func (w *DirectoryWatcher) Watch(newf chan FileEvent) error {
	watcher, err := fsnotify.NewWatcher()

	if err != nil {
		return err
	}

	w.done = make(chan bool)

	go func() {
		for {
			select {
			case event := <-watcher.Events:
				if abs, err := filepath.Abs(event.Name); err == nil {
					newf <- FileEvent{
						Name: event.Name,
						Path: abs,
						Time: time.Now(),
						Op:   event.Op,
					}
				}
			case <-w.done:
				watcher.Close()
			}
		}
	}()

	watcher.Add(w.path)

	return nil
}

// End stops the watching operation.
func (w *DirectoryWatcher) End() {
	w.done <- true
}
