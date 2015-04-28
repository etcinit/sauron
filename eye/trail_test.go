package eye

import (
	"errors"
	"path/filepath"
	"testing"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	fsnotify "gopkg.in/fsnotify.v1"
)

type MockedWatcher struct {
	mock.Mock
}

func (m *MockedWatcher) Walk() (paths []string, err error) {
	args := m.Called()

	//m.called <- true

	return args.Get(0).([]string), args.Error(1)
}

func (m *MockedWatcher) Watch(newf chan FileEvent) error {
	args := m.Called(newf)

	m.TestData()["watchChannel"] = newf

	return args.Error(0)
}

func (m *MockedWatcher) End() {}

func TestNewTrail(t *testing.T) {
	watcher := NewDirectoryWatcher("../_resources")

	NewTrail(watcher)
}

func TestNewTrailWithOptions(t *testing.T) {
	watcher := NewDirectoryWatcher("../_resources")

	options := &TrailOptions{}

	NewTrailWithOptions(watcher, options)

	options.Logger = logrus.New()

	NewTrailWithOptions(watcher, options)
}

func TestFollow(t *testing.T) {
	watcher := MockedWatcher{}

	path, _ := filepath.Abs("../_resources/example.log")
	watcher.On("Walk").Return([]string{path}, nil)
	watcher.On("Watch", mock.AnythingOfType("chan eye.FileEvent")).Return(nil)

	trail := NewTrail(&watcher)

	trail.Follow(func(line Line) error {
		return nil
	})

	path, _ = filepath.Abs("../_resources/error.log")
	watcher.TestData()["watchChannel"].(chan FileEvent) <- FileEvent{
		Name: "../_resources/error.log",
		Path: path,
		Time: time.Now(),
		Op:   fsnotify.Create,
	}

	trail.End()

	watcher.AssertExpectations(t)
}

func TestFollowWalkError(t *testing.T) {
	watcher := MockedWatcher{}

	watcher.On("Walk").Return([]string{}, errors.New("oops"))

	trail := NewTrail(&watcher)
	err := trail.Follow(func(line Line) error {
		return nil
	})

	assert.NotNil(t, err)

	watcher.AssertExpectations(t)
}
