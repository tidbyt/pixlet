package parser

import (
	"github.com/teambition/rrule-go"
	"time"
)

// @TODO: Remove cancelled events; event can be cancelled at the top level or have individual recurrences cancelled

func (cal *Calendar) ExpandRecurringEvent(buf *Event, calendar *Calendar) []Event {
	rule, err := rrule.StrToRRule(buf.RawRecurrenceRule)
	if err != nil || buf.Status == "CANCELLED" {
		return []Event{}
	}

	allRecurrences := rule.All()

	var expandedEvents []Event
	for _, rec := range allRecurrences {
		e := *buf

		newEnd := time.Date(rec.Year(), rec.Month(), rec.Day(), buf.End.Hour(), buf.End.Minute(), rec.Second(), buf.End.Nanosecond(), time.UTC)

		e.Start = &rec
		e.End = &newEnd

		expandedEvents = append(expandedEvents, e)
	}

	return expandedEvents
}
