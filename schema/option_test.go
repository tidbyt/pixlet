package schema_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"tidbyt.dev/pixlet/runtime"
)

var optionSource = `
load("schema.star", "schema")

def assert(success, message=None):
    if not success:
        fail(message or "assertion failed")

s = schema.Option(
	display = "Green",
	value = "#00FF00",
)

assert(s.display == "Green")
assert(s.value == "#00FF00")

def main():
	return []
`

func TestOption(t *testing.T) {
	app, err := runtime.NewApplet("option.star", []byte(optionSource))
	assert.NoError(t, err)

	screens, err := app.Run(context.Background())
	assert.NoError(t, err)
	assert.NotNil(t, screens)
}
