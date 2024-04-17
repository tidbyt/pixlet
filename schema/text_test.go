package schema_test

import (
	"context"
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
	default = "foo",
)

assert(s.id == "screen_name")
assert(s.name == "Screen Name")
assert(s.desc == "A text entry for your screen name.")
assert(s.icon == "user")
assert(s.default == "foo")

def main():
	return []
`

func TestText(t *testing.T) {
	app, err := runtime.NewApplet("text.star", []byte(textSource))
	assert.NoError(t, err)

	screens, err := app.Run(context.Background())
	assert.NoError(t, err)
	assert.NotNil(t, screens)
}
