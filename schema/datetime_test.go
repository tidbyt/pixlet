package schema_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"tidbyt.dev/pixlet/runtime"
)

var dateTimeSource = `
load("schema.star", "schema")

def assert(success, message=None):
    if not success:
        fail(message or "assertion failed")

t = schema.DateTime(
	id = "event_name",
	name = "Event Name",
	desc = "The time of the event.",
	icon = "gear",
)

assert(t.id == "event_name")
assert(t.name == "Event Name")
assert(t.desc == "The time of the event.")
assert(t.icon == "gear")

def main():
	return []
`

func TestDateTime(t *testing.T) {
	app, err := runtime.NewApplet("date_time.star", []byte(dateTimeSource))
	assert.NoError(t, err)

	screens, err := app.Run(context.Background())
	assert.NoError(t, err)
	assert.NotNil(t, screens)
}
