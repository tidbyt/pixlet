package schema_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"tidbyt.dev/pixlet/runtime"
)

var moduleSource = `
load("schema.star", "schema")

def main():
	return []
`

var schemaSource = `
load("schema.star", "schema")

def assert(success, message=None):
    if not success:
        fail(message or "assertion failed")

s = schema.Schema(
	version = "1",
	fields = [
		schema.Toggle(
			id = "display_weather",
			name = "Display Weather",
			desc = "A toggle to determine if the weather should be displayed.",
			icon = "cloud",
		),
	],
)

assert(s.version == "1")
assert(s.fields[0].name == "Display Weather")

def main():
	return []
`

func TestStarlarkSchema(t *testing.T) {
	app := &runtime.Applet{}
	err := app.Load("starlark.star", []byte(schemaSource), nil)
	assert.NoError(t, err)

	screens, err := app.Run(map[string]string{})
	assert.NoError(t, err)
	assert.NotNil(t, screens)
}
func TestSchemaModuleLoads(t *testing.T) {
	app := &runtime.Applet{}
	err := app.Load("source.star", []byte(moduleSource), nil)
	assert.NoError(t, err)

	screens, err := app.Run(map[string]string{})
	assert.NoError(t, err)
	assert.NotNil(t, screens)
}
