package icalendar_test

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
	"tidbyt.dev/pixlet/runtime"
)

import (
	"context"
)

var randomSrc = `
load("icalendar.star", "icalendar")
raw_string = """
	BEGIN:VCALENDAR
	VERSION:2.0
	CALSCALE:GREGORIAN
	BEGIN:VTIMEZONE
	TZID:America/Phoenix
	LAST-MODIFIED:20231222T233358Z
	TZURL:https://www.tzurl.org/zoneinfo-outlook/America/Phoenix
	X-LIC-LOCATION:America/Phoenix
	BEGIN:STANDARD
	TZNAME:MST
	TZOFFSETFROM:-0700
	TZOFFSETTO:-0700
	DTSTART:19700101T000000
	END:STANDARD
	END:VTIMEZONE
	BEGIN:VEVENT
	DTSTAMP:20240817T225510Z
	UID:1723935287215-82939@ical.marudot.com
	DTSTART;TZID=America/Phoenix:20240801T120000
	RRULE:FREQ=DAILY
	DTEND;TZID=America/Phoenix:20240801T150000
	SUMMARY:Test
	DESCRIPTION:Test Description
	LOCATION:Phoenix
	END:VEVENT
	END:VCALENDAR
"""

def test_icalendar():
	events = icalendar.parse(raw_string)
	for event in events:
		print(event)


def main():
	return []
`

func TestICalendar(t *testing.T) {
	app, err := runtime.NewApplet("icalendar_test.star", []byte(randomSrc))
	require.NoError(t, err)

	screens, err := app.Run(context.Background())
	require.NoError(t, err)
	assert.NotNil(t, screens)
}
