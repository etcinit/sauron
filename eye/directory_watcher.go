package eye

import (
	"os"
	"path/filepath"
	"time"

	fsnotify "gopkg.in/fsnotify.v1"
)

type DirectoryWatcher struct {
	path string
	done chan bool
}

func NewDirectoryWatcher(path string) *DirectoryWatcher {
	//if string(path[len(path)-1]) != "/" {
	//	path += "/"
	//}

	return &DirectoryWatcher{
		path: path,
	}
}

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

func (w *DirectoryWatcher) End() {
	w.done <- true
}
