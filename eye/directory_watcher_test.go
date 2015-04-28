package eye

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWalk(t *testing.T) {
	watcher := DirectoryWatcher{path: "../_resources/"}

	files, _ := watcher.Walk()

	assert.True(t, len(files) == 2)
}
