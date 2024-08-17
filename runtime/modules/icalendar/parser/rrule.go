package parser

import (
	"github.com/teambition/rrule-go"
	"time"
)

func (cal *Calendar) ExpandRecurringEvent(buf *Event) []Event {
	rule, err := rrule.StrToRRule(buf.RawRecurrenceRule)
	if err != nil || buf.Status == "CANCELLED" {
		return []Event{}
	}

	allRecurrences := rule.All()

	var excludedDateTime map[string]*time.Time
	for _, t := range buf.ExcludeDates {
		str := t.Format(time.RFC3339)
		excludedDateTime[str] = t
	}

	var expandedEvents []Event
	for _, rec := range allRecurrences {
		if _, ok := excludedDateTime[rec.Format(time.RFC3339)]; ok {
			continue
		}

		e := *buf
		newEnd := time.Date(rec.Year(), rec.Month(), rec.Day(), buf.End.Hour(), buf.End.Minute(), rec.Second(), buf.End.Nanosecond(), time.UTC)

		e.Start = &rec
		e.End = &newEnd

		expandedEvents = append(expandedEvents, e)
	}

	return expandedEvents
}
