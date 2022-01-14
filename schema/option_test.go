package schema_test

import (
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
	app := &runtime.Applet{}
	err := app.Load("option.star", []byte(optionSource), nil)
	assert.NoError(t, err)

	screens, err := app.Run(map[string]string{})
	assert.NoError(t, err)
	assert.NotNil(t, screens)
}
