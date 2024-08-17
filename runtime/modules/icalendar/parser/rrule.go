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

	now := time.Now()
	threeMonthsFromNow := now.AddDate(0, 3, 0)

	nextThreeMonthsOfRecurrences := rule.Between(now, threeMonthsFromNow, true)

	var excludedDateTime map[string]*time.Time
	for _, t := range buf.ExcludeDates {
		str := t.Format(time.RFC3339)
		excludedDateTime[str] = t
	}

	var expandedEvents []Event
	for _, rec := range nextThreeMonthsOfRecurrences {
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
