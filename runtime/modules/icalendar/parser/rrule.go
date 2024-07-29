package parser

import (
	"github.com/teambition/rrule-go"
	"time"
)

func (cal *Calendar) ExpandRecurringEvent(buf *Event) []Event {
	rule, err := rrule.StrToRRule(buf.RawRecurrenceRule)
	if err != nil {
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
