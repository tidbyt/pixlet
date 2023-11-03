package schema_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"tidbyt.dev/pixlet/runtime"
)

var generatedSource = `
load("schema.star", "schema")

def assert(success, message=None):
    if not success:
        fail(message or "assertion failed")

s = schema.Generated(
	id = "foo",
        source = "bar",
        handler = assert,
)

assert(s.id == "foo")
assert(s.source == "bar")
assert(s.handler == assert)

def main():
	return []
`

func TestGenerated(t *testing.T) {
	app := &runtime.Applet{}
	err := app.Load("gid", "generated.star", []byte(generatedSource), nil)
	assert.NoError(t, err)

	screens, err := app.Run(map[string]string{})
	assert.NoError(t, err)
	assert.NotNil(t, screens)
}
