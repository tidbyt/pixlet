package schema_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"tidbyt.dev/pixlet/runtime"
)

var typeaheadSource = `
load("schema.star", "schema")

def assert(success, message = None):
    if not success:
        fail(message or "assertion failed")

def search(pattern):
    return [
        schema.Option(
            display = "Grand Central",
            value = "abc123",
        ),
        schema.Option(
            display = "Penn Station",
            value = "xyz123",
        ),
    ]

t = schema.Typeahead(
    id = "search",
    name = "Search",
    desc = "A list of items that match search.",
    icon = "gear",
    handler = search,
)

assert(t.id == "search")
assert(t.name == "Search")
assert(t.desc == "A list of items that match search.")
assert(t.icon == "gear")
assert(t.handler("")[0].display == "Grand Central")

def main():
    return []

`

func TestTypeahead(t *testing.T) {
	app := &runtime.Applet{}
	err := app.Load("typeahead.star", []byte(typeaheadSource), nil)
	assert.NoError(t, err)

	screens, err := app.Run(map[string]string{})
	assert.NoError(t, err)
	assert.NotNil(t, screens)
}
