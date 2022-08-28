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

def abs(x):
	if x > 0:
		return x
	return -x

# Setup.
format = "2006-01-02T15:04:05"
input = time.parse_time("2022-01-15T22:40:24", format = format)
expectedRise = time.parse_time("2022-01-15T12:17:29", format = format)
expectedSet = time.parse_time("2022-01-15T21:52:30", format = format)
lat = 40.6781784
lng = -73.9441579

# Sunrise occurs when center of the sun is 50 arc minutes below horizon
# due to atmospheric refraction and the angular distance to the top.
# https://en.wikipedia.org/wiki/Sunrise#Angle
sunriseElevation = -50.0 / 60.0

# Call methods.
rise = sunrise.sunrise(lat, lng, input)
set = sunrise.sunset(lat, lng, input)
elevation = sunrise.elevation(lat, lng, expectedSet)
morning, evening = sunrise.elevation_time(lat, lng, sunriseElevation, input)

# Assert.
assert(rise == expectedRise)
assert(set == expectedSet)
assert(abs(elevation - sunriseElevation) < 0.005)
assert(abs(expectedRise.unix - morning.unix) < 2)
assert(abs(evening.unix - expectedSet.unix) < 2)

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
