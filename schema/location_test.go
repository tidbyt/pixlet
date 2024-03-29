package schema_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"tidbyt.dev/pixlet/runtime"
)

var locationSource = `
load("schema.star", "schema")

def assert(success, message=None):
    if not success:
        fail(message or "assertion failed")

s = schema.Location(
	id = "location",
	name = "Location",
	desc = "Location for which to display time.",
	icon = "locationDot",
)

assert(s.id == "location")
assert(s.name == "Location")
assert(s.desc == "Location for which to display time.")
assert(s.icon == "locationDot")

def main():
	return []
`

func TestLocation(t *testing.T) {
	app := &runtime.Applet{}
	err := app.Load("lid", "location.star", []byte(locationSource), nil)
	assert.NoError(t, err)

	screens, err := app.Run(map[string]string{})
	assert.NoError(t, err)
	assert.NotNil(t, screens)
}
