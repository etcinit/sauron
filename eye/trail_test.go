package eye

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/mock"
)

type MockedWatcher struct {
	mock.Mock
	called chan bool
}

func (m MockedWatcher) Walk() (paths []string, err error) {
	args := m.Called()

	m.called <- true

	return args.Get(0).([]string), args.Error(1)
}

func (m MockedWatcher) Watch(newf chan FileEvent) error {
	m.called <- true

	return nil
}

func TestNewTrail(t *testing.T) {
	watcher := NewDirectoryWatcher("../_resources")

	NewTrail(watcher)
}

func TestFollow(t *testing.T) {
	watcher := MockedWatcher{
		called: make(chan bool),
	}

	path, _ := filepath.Abs("../_resources/example.log")
	watcher.On("Walk").Return([]string{path}, nil)

	trail := NewTrail(watcher)

	done := make(chan bool)

	go func() {
		trail.Follow(func(line Line) error {
			return nil
		})

		done <- true
	}()

	<-watcher.called
	<-watcher.called

	trail.End()

	<-done
}
