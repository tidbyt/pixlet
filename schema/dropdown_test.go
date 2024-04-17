package schema_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"tidbyt.dev/pixlet/runtime"
)

var dropdownSource = `
load("schema.star", "schema")

def assert(success, message=None):
    if not success:
        fail(message or "assertion failed")

options = [
	schema.Option(
		display = "Green",
		value = "#00FF00",
	),
	schema.Option(
		display = "Red",
		value = "#FF0000",
	),
]
	
s = schema.Dropdown(
	id = "colors",
	name = "Text Color",
	desc = "The color of text to be displayed.", 
	icon = "brush",
	default = options[0].value,
	options = options,
)

assert(s.id == "colors")
assert(s.name == "Text Color")
assert(s.desc == "The color of text to be displayed.")
assert(s.icon == "brush")
assert(s.default == "#00FF00")

assert(s.options[0].display == "Green")
assert(s.options[0].value == "#00FF00")

assert(s.options[1].display == "Red")
assert(s.options[1].value == "#FF0000")

def main():
	return []
`

func TestDropdown(t *testing.T) {
	app, err := runtime.NewApplet("dropdown.star", []byte(dropdownSource))
	assert.NoError(t, err)

	screens, err := app.Run(context.Background())
	assert.NoError(t, err)
	assert.NotNil(t, screens)
}
