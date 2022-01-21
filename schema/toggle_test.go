package schema_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"tidbyt.dev/pixlet/runtime"
)

var toggleSource = `
load("schema.star", "schema")

def assert(success, message=None):
    if not success:
        fail(message or "assertion failed")

t = schema.Toggle(
	id = "display_weather",
	name = "Display Weather",
	desc = "A toggle to determine if the weather should be displayed.",
	icon = "cloud",
	default = True,
)

assert(t.id == "display_weather")
assert(t.name == "Display Weather")
assert(t.desc == "A toggle to determine if the weather should be displayed.")
assert(t.icon == "cloud")
assert(t.default == True)

def main():
	return []
`

func TestToggle(t *testing.T) {
	app := &runtime.Applet{}
	err := app.Load("toggle.star", []byte(toggleSource), nil)
	assert.NoError(t, err)

	screens, err := app.Run(map[string]string{})
	assert.NoError(t, err)
	assert.NotNil(t, screens)
}
