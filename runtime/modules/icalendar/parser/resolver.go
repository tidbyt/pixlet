package parser

import (
	"fmt"
	"tidbyt.dev/pixlet/runtime/modules/icalendar/parser/members"
	"time"
)

func resolve[T comparable](cal *Calendar, l *Line, dst *T, resolve func(cal *Calendar, line *Line) (T, T, error), post func(cal *Calendar, out T)) error {
	value, empty, err := resolve(cal, l)
	if err != nil {
		return err
	}

	if dst != nil && *dst != empty {
		if cal.Duplicate.Mode == DuplicateModeFailStrict {
			return NewDuplicateAttribute(l.Key, l.Value)
		}
	}

	// If the value is empty or the duplicate mode allows further processing, set the value
	if *dst == empty || cal.Duplicate.Mode == DuplicateModeKeepLast {
		*dst = value
		if post != nil && dst != nil {
			post(cal, *dst)
		}
	}

	return nil
}

func resolveString(cal *Calendar, l *Line) (string, string, error) {
	return l.Value, "", nil
}

func resolveLatLng(gc *Calendar, l *Line) (*LatLng, *LatLng, error) {
	lat, long, err := members.ParseLatLng(l.Value)
	if err != nil {
		return nil, nil, err
	}

	return &LatLng{lat, long}, nil, nil
}

func resolveDate(cal *Calendar, l *Line) (*time.Time, *time.Time, error) {
	d, err := members.ParseTime(l.Value, l.Params, members.TimeStart, false, cal.AllDayEventsTZ)
	if err != nil {
		return nil, nil, fmt.Errorf("could not parse: %s", err)
	}

	return d, nil, nil
}

func resolveDateEnd(cal *Calendar, l *Line) (*time.Time, *time.Time, error) {
	d, err := members.ParseTime(l.Value, l.Params, members.TimeEnd, false, cal.AllDayEventsTZ)
	if err != nil {
		return nil, nil, fmt.Errorf("could not parse: %s", err)
	}

	return d, nil, nil
}

func resolveOrganizer(cal *Calendar, l *Line) (*Organizer, *Organizer, error) {
	o := Organizer{
		Cn:          l.Params["CN"],
		DirectoryDn: l.Params["DIR"],
		Value:       l.Value,
	}

	return &o, nil, nil
}

func resolveDuration(cal *Calendar, l *Line) (*time.Duration, *time.Duration, error) {
	d, err := members.ParseDuration(l.Value)
	if err != nil {
		return nil, nil, fmt.Errorf("could not parse: %s", err)
	}

	return d, nil, nil
}
