package schema_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"tidbyt.dev/pixlet/runtime"
)

var photoSelectSource = `
load("schema.star", "schema")

def assert(success, message=None):
    if not success:
        fail(message or "assertion failed")

t = schema.PhotoSelect(
	id = "photo",
	name = "Add Photo",
	desc = "A photo.",
	icon = "gear",
)

assert(t.id == "photo")
assert(t.name == "Add Photo")
assert(t.desc == "A photo.")
assert(t.icon == "gear")

def main():
	return []
`

func TestPhotoSelect(t *testing.T) {
	app := &runtime.Applet{}
	err := app.Load("photo_select.star", []byte(photoSelectSource), nil)
	assert.NoError(t, err)

	screens, err := app.Run(map[string]string{})
	assert.NoError(t, err)
	assert.NotNil(t, screens)
}
