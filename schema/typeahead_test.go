package schema_test

import (
	"context"
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
	app, err := runtime.NewApplet("typeahead.star", []byte(typeaheadSource))
	assert.NoError(t, err)

	screens, err := app.Run(context.Background())
	assert.NoError(t, err)
	assert.NotNil(t, screens)
}
