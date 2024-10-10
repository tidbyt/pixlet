package members

import (
	cases "golang.org/x/text/cases"
	"golang.org/x/text/language"
	"strings"
	"time"
)

const (
	TimeStart = iota
	TimeEnd
)

func ParseTime(s string, params map[string]string, ty int, allday bool, allDayTZ *time.Location) (*time.Time, error) {
	// Date field is YYYYMMDD
	// Reference: https://icalendar.org/iCalendar-RFC-5545/3-3-4-date.html
	if params["VALUE"] == "DATE" || len(s) == 8 {
		t, err := time.Parse("20060102", s)
		if err != nil {
			return nil, err
		}

		if ty == TimeStart {
			t = time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, allDayTZ)

		} else if ty == TimeEnd {
			if allday {
				t = time.Date(t.Year(), t.Month(), t.Day(), 23, 59, 59, 999, allDayTZ)
			} else {
				t = time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, allDayTZ)
			}
		}
		return &t, err
	}

	// Z indicates we're using UTC
	if strings.HasSuffix(s, "Z") {
		format := "20060102T150405Z"
		tz, _ := time.LoadLocation("UTC")
		t, err := time.ParseInLocation(format, s, tz)
		if err != nil {
			return nil, err
		}
		return &t, err
	} else if params["TZID"] != "" {
		format := "20060102T150405"
		tz, err := time.LoadLocation(params["TZID"])
		if err != nil {
			// If there's an error we can assume that the timezones are in Window's format.
			// This is especially common due to Microsoft Exchange iCalendar files
			unixTz, err := ConvertTimeZoneWindowsToLinux(params["TZID"])
			if err != nil {
				return nil, err
			}
			tz, err := time.LoadLocation(*unixTz)
			if err != nil {
				return nil, err
			}
			t, err := time.ParseInLocation(format, s, tz)
			if err != nil {
				return nil, err
			}
			return &t, err
		}

		t, err := time.ParseInLocation(format, s, tz)
		if err != nil {
			return nil, err
		}
		return &t, err
	}

	// Default to local time if Z is not in use, there's no TZID and there is no Date
	format := "20060102T150405"
	tz := time.Local

	t, err := time.ParseInLocation(format, s, tz)
	if err != nil {
		return nil, err
	}
	return &t, err
}

func ParseDuration(s string) (*time.Duration, error) {
	dur, err := time.ParseDuration(s)
	if err != nil {
		return nil, err
	}

	return &dur, err
}

func LoadTimezone(tzid string) (*time.Location, error) {
	tz, err := time.LoadLocation(tzid)
	if err != nil {
		return tz, nil
	}

	tokens := strings.Split(tzid, "_")
	for idx, token := range tokens {
		t := strings.ToLower(token)
		if t != "of" && t != "es" {
			tokens[idx] = cases.Title(language.English).String(token)
		} else {
			tokens[idx] = t
		}
	}

	return time.LoadLocation(strings.Join(tokens, "_"))
}
