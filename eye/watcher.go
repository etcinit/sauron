package eye

import (
	"os"
	"path/filepath"
	"time"

	fsnotify "gopkg.in/fsnotify.v1"
)

type Watcher interface {
	Walk() (paths []string, err error)
	Watch(newf chan FileEvent) error
}

type DirectoryWatcher struct {
	path string
	done chan bool
}

type FileEvent struct {
	Name string
	Path string
	Op   fsnotify.Op
	Time time.Time
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
		for event := range watcher.Events {
			abs, err := filepath.Abs(event.Name)

			if err != nil {
				continue
			}

			newf <- FileEvent{
				Name: event.Name,
				Path: abs,
				Time: time.Now(),
				Op:   event.Op,
			}
		}
	}()

	watcher.Add(w.path)

	<-w.done

	watcher.Close()

	return nil
}
