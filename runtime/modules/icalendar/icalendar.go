package icalendar

import (
	"fmt"
	"strings"
	"sync"
	"tidbyt.dev/pixlet/runtime/modules/icalendar/parser"
	"time"

	godfe "github.com/newm4n/go-dfe"
	"go.starlark.net/starlark"
	"go.starlark.net/starlarkstruct"
)

const (
	ModuleName = "icalendar"
)

var (
	once        sync.Once
	module      starlark.StringDict
	empty       time.Time
	translation *godfe.PatternTranslation
)

func LoadModule() (starlark.StringDict, error) {
	translation = godfe.NewPatternTranslation()
	once.Do(func() {
		module = starlark.StringDict{
			ModuleName: &starlarkstruct.Module{
				Name: ModuleName,
				Members: starlark.StringDict{
					"parseCalendar": starlark.NewBuiltin("findNextEvent", parseCalendar),
				},
			},
		}
	})

	return module, nil
}

/*
* This function returns a list of events with the events metadata
 */
func parseCalendar(thread *starlark.Thread, _ *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
	var (
		rawCalendar starlark.String
	)
	if err := starlark.UnpackArgs(
		"parseCalendar",
		args, kwargs,
		"str", &rawCalendar,
	); err != nil {
		return nil, fmt.Errorf("unpacking arguments for bytes: %s", err)
	}
	calendar := parser.NewParser(strings.NewReader(rawCalendar.GoString()))
	if err := calendar.Parse(); err != nil {
		return nil, fmt.Errorf("parsing calendar: %s", err)
	}

	events := make([]starlark.Value, 0, len(calendar.Events))

	for _, event := range calendar.Events {
		dict := starlark.NewDict(25)

		if err := dict.SetKey(starlark.String("uid"), starlark.String(event.Uid)); err != nil {
			return nil, fmt.Errorf("setting uid: %s", err)
		}
		if err := dict.SetKey(starlark.String("summary"), starlark.String(event.Summary)); err != nil {
			return nil, fmt.Errorf("setting summary: %s", err)
		}
		if err := dict.SetKey(starlark.String("description"), starlark.String(event.Description)); err != nil {
			return nil, fmt.Errorf("setting description: %s", err)
		}
		if err := dict.SetKey(starlark.String("status"), starlark.String(event.Status)); err != nil {
			return nil, fmt.Errorf("setting status: %s", err)
		}
		if err := dict.SetKey(starlark.String("comment"), starlark.String(event.Comment)); err != nil {
			return nil, fmt.Errorf("setting comment: %s", err)
		}
		if err := dict.SetKey(starlark.String("start"), starlark.String(event.Start.Format(time.RFC3339))); err != nil {
			return nil, fmt.Errorf("setting start: %s", err)
		}
		if err := dict.SetKey(starlark.String("end"), starlark.String(event.End.Format(time.RFC3339))); err != nil {
			return nil, fmt.Errorf("setting end: %s", err)
		}
		if err := dict.SetKey(starlark.String("is_recurring"), starlark.Bool(event.IsRecurring)); err != nil {
			return nil, fmt.Errorf("setting is_recurring: %s", err)
		}
		if err := dict.SetKey(starlark.String("location"), starlark.String(event.Location)); err != nil {
			return nil, fmt.Errorf("setting location: %s", err)
		}
		if err := dict.SetKey(starlark.String("duration_in_seconds"), starlark.Float(event.Duration.Seconds())); err != nil {
			return nil, fmt.Errorf("setting duration: %s", err)
		}
		if err := dict.SetKey(starlark.String("url"), starlark.String(event.Url)); err != nil {
			return nil, fmt.Errorf("setting end: %s", err)
		}
		if err := dict.SetKey(starlark.String("sequence"), starlark.MakeInt(event.Sequence)); err != nil {
			return nil, fmt.Errorf("setting end: %s", err)
		}
		if err := dict.SetKey(starlark.String("created_at"), starlark.String(event.Created.Format(time.RFC3339))); err != nil {
			return nil, fmt.Errorf("setting end: %s", err)
		}
		if err := dict.SetKey(starlark.String("updated_at"), starlark.String(event.LastModified.Format(time.RFC3339))); err != nil {
			return nil, fmt.Errorf("setting end: %s", err)
		}

		events = append(events, dict)
	}
	return starlark.NewList(events), nil
}
