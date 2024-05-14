package schema_test

import (
	"context"
	"testing"
	"testing/fstest"

	"github.com/stretchr/testify/assert"
	"tidbyt.dev/pixlet/runtime"
)

var soundSource = `
load("schema.star", "schema")
load("sound.mp3", "file")

def assert(success, message=None):
    if not success:
        fail(message or "assertion failed")

s = schema.Sound(
	id = "sound1",
	title = "Sneezing Elephant",
	file = file,
)

assert(s.id == "sound1")
assert(s.title == "Sneezing Elephant")
assert(s.file == file)
assert(s.file.readall() == "sound data")

def main():
	return []
`

func TestSound(t *testing.T) {
	vfs := fstest.MapFS{
		"sound.mp3":  &fstest.MapFile{Data: []byte("sound data")},
		"sound.star": &fstest.MapFile{Data: []byte(soundSource)},
	}
	app, err := runtime.NewAppletFromFS("sound", vfs)
	assert.NoError(t, err)

	screens, err := app.Run(context.Background())
	assert.NoError(t, err)
	assert.NotNil(t, screens)
}
