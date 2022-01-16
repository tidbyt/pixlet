package sunrise_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"tidbyt.dev/pixlet/runtime"
)

var sunSource = `
load("time.star", "time")
load("sunrise.star", "sunrise")

def assert(success, message=None):
    if not success:
        fail(message or "assertion failed")


# Setup.
format = "2006-01-02T15:04:05"
input = time.parse_time("2022-01-15T22:40:24", format = format)
expectedRise = time.parse_time("2022-01-15T12:17:29", format = format)
expectedSet = time.parse_time("2022-01-15T21:52:30", format = format)
lat = 40.6781784
lng = -73.9441579

# Call methods.
rise = sunrise.sunrise(lat, lng, input)
set = sunrise.sunset(lat, lng, input)

# Assert.
assert(rise == expectedRise)
assert(set == expectedSet)

def main():
	return []
`

func TestSunrise(t *testing.T) {
	app := &runtime.Applet{}
	err := app.Load("sun.star", []byte(sunSource), nil)
	assert.NoError(t, err)

	screens, err := app.Run(map[string]string{})
	assert.NoError(t, err)
	assert.NotNil(t, screens)
}
