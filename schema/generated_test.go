package schema_test

import (
	"context"
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
	app, err := runtime.NewApplet("generated.star", []byte(generatedSource))
	assert.NoError(t, err)

	screens, err := app.Run(context.Background())
	assert.NoError(t, err)
	assert.NotNil(t, screens)
}
