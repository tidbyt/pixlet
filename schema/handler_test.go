package schema_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"tidbyt.dev/pixlet/runtime"
)

var handlerSource = `
load("schema.star", "schema")

def assert(success, message=None):
    if not success:
        fail(message or "assertion failed")

def foobar(param):
    return "derp"

h = schema.Handler(
    handler = foobar,
    type = schema.HandlerType.String,
)

assert(h.handler == foobar)
assert(h.type == schema.HandlerType.String)

def main():
	return []
`

func TestHandler(t *testing.T) {
	app := &runtime.Applet{}
	err := app.Load("handler.star", []byte(handlerSource), nil)
	assert.NoError(t, err)

	screens, err := app.Run(map[string]string{})
	assert.NoError(t, err)
	assert.NotNil(t, screens)
}

func TestHandlerBadParams(t *testing.T) {
	// Handler is a string
	app := &runtime.Applet{}
	err := app.Load("text.star", []byte(`
load("schema.star", "schema")

def foobar(param):
    return "derp"

h = schema.Handler(
    handler = "foobar",
    type = schema.HandlerType.String,
)

def main():
	return []
`), nil)
	assert.Error(t, err)

	// Type is not valid
	app = &runtime.Applet{}
	err = app.Load("text.star", []byte(`
load("schema.star", "schema")

def foobar(param):
    return "derp"

h = schema.Handler(
    handler = foobar,
    type = 42,
)

def main():
	return []
`), nil)
	assert.Error(t, err)

}
