package eye

import (
	"time"

	fsnotify "gopkg.in/fsnotify.v1"
)

// A Watcher is capable of providing a list of the current files and notify
// about changes in a directory. An implementation might decide whether it
// supports multiple directory levels (recursive) or just one level. Files also
// do not need to be in a traditional filesystem.
type Watcher interface {
	Walk() (paths []string, err error)
	Watch(newf chan FileEvent) error
	End()
}

// FileEvent represents an event affecting a single file.
type FileEvent struct {
	// Name of the file affected.
	Name string
	// Path to the file, including the filename.
	Path string
	// Operation that triggerred the event.
	Op fsnotify.Op
	// Time at which the event occured.
	Time time.Time
}
