package schema_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"tidbyt.dev/pixlet/runtime"
)

var locationBasedSource = `
load("encoding/json.star", "json")
load("schema.star", "schema")

DEFAULT_LOCATION = """
{
	"lat": "40.6781784",
	"lng": "-73.9441579",
	"description": "Brooklyn, NY, USA",
	"locality": "Brooklyn",
	"place_id": "ChIJCSF8lBZEwokRhngABHRcdoI",
	"timezone": "America/New_York"
}
"""

def assert(success, message = None):
    if not success:
        fail(message or "assertion failed")

def get_stations(location):
    loc = json.decode(location)
    lat, lng = float(loc["lat"]), float(loc["lng"])

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

t = schema.LocationBased(
    id = "station",
    name = "Train Station",
    desc = "A list of train stations based on a location.",
    icon = "train",
    handler = get_stations,
)

assert(t.id == "station")
assert(t.name == "Train Station")
assert(t.desc == "A list of train stations based on a location.")
assert(t.icon == "train")
assert(t.handler(DEFAULT_LOCATION)[0].display == "Grand Central")

def main():
    return []

`

func TestLocationBased(t *testing.T) {
	app := &runtime.Applet{}
	err := app.Load("location_based.star", []byte(locationBasedSource), nil)
	assert.NoError(t, err)

	screens, err := app.Run(map[string]string{})
	assert.NoError(t, err)
	assert.NotNil(t, screens)
}
