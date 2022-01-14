package schema_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"tidbyt.dev/pixlet/runtime"
)

var textSource = `
load("schema.star", "schema")

def assert(success, message=None):
    if not success:
        fail(message or "assertion failed")

s = schema.Text(
	id = "screen_name",
	name = "Screen Name",
	desc = "A text entry for your screen name.",
	icon = "user",
)

assert(s.id == "screen_name")
assert(s.name == "Screen Name")
assert(s.desc == "A text entry for your screen name.")
assert(s.icon == "user")

def main():
	return []
`

func TestText(t *testing.T) {
	app := &runtime.Applet{}
	err := app.Load("text.star", []byte(textSource), nil)
	assert.NoError(t, err)

	screens, err := app.Run(map[string]string{})
	assert.NoError(t, err)
	assert.NotNil(t, screens)
}
