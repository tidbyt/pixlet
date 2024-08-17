package parser

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"strconv"
	"strings"
	"tidbyt.dev/pixlet/runtime/modules/icalendar/parser/members"
	"time"
)

func NewParser(r io.Reader) *Calendar {
	return &Calendar{
		scanner: bufio.NewScanner(r),
		Events:  make([]*Event, 0),
		Strict: StrictParams{
			Mode: StrictModeFailFeed,
		},
		Duplicate: DuplicateParams{
			DuplicateModeFailStrict,
		},
		SkipBounds:     false,
		AllDayEventsTZ: time.UTC,
	}
}

func (cal *Calendar) Parse() error {
	if cal.Start == nil {
		start := time.Now().Add(-1 * 24 * time.Hour)
		cal.Start = &start
	}
	if cal.End == nil {
		end := time.Now().Add(3 * 30 * 24 * time.Hour)
		cal.End = &end
	}

	cal.scanner.Scan()

	rInstances := make([]Event, 0)
	ctx := &Context{Value: ContextRoot}

	for {
		l, err, done := cal.parseLine()
		if err != nil {
			if done {
				break
			}
			continue
		}

		if l.IsValue("VCALENDAR") {
			continue
		}

		if ctx.Value == ContextRoot && l.Is("BEGIN", "VEVENT") {
			ctx = ctx.Nest(ContextEvent)
			cal.buffer = &Event{Valid: true, delayed: make([]*Line, 0)}
		} else if ctx.Value == ContextEvent && l.Is("END", "VEVENT") {
			if ctx.Previous == nil {
				return fmt.Errorf("got an END:* without matching BEGIN:*")
			}
			ctx = ctx.Previous

			for _, d := range cal.buffer.delayed {
				cal.parseEvent(d)
			}

			if cal.buffer.RawStart.Value == cal.buffer.RawEnd.Value {
				if value, ok := cal.buffer.RawEnd.Params["VALUE"]; ok && value == "DATE" {
					cal.buffer.End, err = members.ParseTime(cal.buffer.RawEnd.Value, cal.buffer.RawEnd.Params, members.TimeEnd, true, cal.AllDayEventsTZ)
				}
			}

			if cal.buffer.End == nil && cal.buffer.RawStart.Params["VALUE"] == "DATE" {
				d := (*cal.buffer.Start).Add(24 * time.Hour)
				cal.buffer.End = &d
			}

			if err := cal.checkEvent(); err != nil {
				switch cal.Strict.Mode {
				case StrictModeFailFeed:
					return fmt.Errorf("calender error: %s", err)
				case StrictModeFailEvent:
					continue
				}
			}

			if cal.buffer.Start == nil || cal.buffer.End == nil {
				continue
			}

			if cal.buffer.IsRecurring {
				rInstances = append(rInstances, cal.ExpandRecurringEvent(cal.buffer)...)
			} else {
				if cal.buffer.End == nil || cal.buffer.Start == nil {
					continue
				}
				if !cal.SkipBounds && !cal.IsInRange(*cal.buffer) {
					continue
				}
				if cal.Strict.Mode == StrictModeFailEvent && !cal.buffer.Valid {
					continue
				}
				cal.Events = append(cal.Events, cal.buffer)
			}

		} else if l.IsKey("BEGIN") {
			ctx = ctx.Nest(ContextUnknown)

		} else if l.IsKey("END") {
			if ctx.Previous == nil {
				return fmt.Errorf("got an END:* without matching BEGIN:*")
			}

			ctx = ctx.Previous
		} else if ctx.Value == ContextEvent {
			if err := cal.parseEvent(l); err != nil {
				var duplicateAttributeError DuplicateAttributeError
				if errors.As(err, &duplicateAttributeError) {
					switch cal.Duplicate.Mode {
					case DuplicateModeFailStrict:
						switch cal.Strict.Mode {
						case StrictModeFailFeed:
							return fmt.Errorf("gocal error: %s", err)
						case StrictModeFailEvent:
							cal.buffer.Valid = false
							continue
						case StrictModeFailAttribute:
							cal.buffer.Valid = false
							continue
						}
					}
				}

				return fmt.Errorf(fmt.Sprintf("gocal error: %s", err))

			}
		} else {
			continue
		}

		if done {
			break
		}
	}

	for _, i := range rInstances {
		if !cal.IsRecurringInstanceOverridden(&i) && cal.IsInRange(i) {
			cal.Events = append(cal.Events, &i)
		}
	}

	return nil

}

/*
* The iCal mandates that lines longer than 75 octets require a linebreak.
* The format uses progressive whitespace indentation to denote a line is continued on a new line.
* https://icalendar.org/iCalendar-RFC-5545/3-1-content-lines.html
 */
func (cal *Calendar) findNextLine() (bool, string) {
	l := cal.scanner.Text()
	done := !cal.scanner.Scan()

	if !done {
		for strings.HasPrefix(cal.scanner.Text(), " ") {
			l = l + strings.TrimPrefix(cal.scanner.Text(), " ")
			if done = !cal.scanner.Scan(); done {
				break
			}
		}
	}

	return done, l
}

// splitLineTokens assures that property parameters that are quoted due to containing special
// characters (like COLON, SEMICOLON, COMMA) are not split.
// See RFC5545, 3.1.1.
func splitLineTokens(line string) []string {
	// go's Split is highly optimized -> use, unless we cannot
	if idxQuote := strings.Index(line, `"`); idxQuote == -1 {
		return strings.SplitN(line, ":", 2)
	} else if idxColon := strings.Index(line, ":"); idxQuote > idxColon {
		return []string{line[0:idxColon], line[idxColon+1:]}
	}

	// otherwise, we need to do it ourselves, let's keep it simple at least:
	quoted := false
	size := len(line)
	for idx, char := range []byte(line) {
		if char == '"' {
			quoted = !quoted
		} else if char == ':' && !quoted && idx+1 < size {
			return []string{line[0:idx], line[idx+1:]}
		}
	}
	return []string{line}
}

func (cal *Calendar) parseLine() (*Line, error, bool) {
	done, line := cal.findNextLine()
	if done {
		return nil, nil, done
	}
	tokens := splitLineTokens(line)
	if len(tokens) < 2 {
		return nil, fmt.Errorf("could not parse item: %s", line), done
	}

	attr, params := members.ParseParameters(tokens[0])

	return &Line{Key: attr, Params: params, Value: members.UnescapeString(strings.TrimPrefix(tokens[1], " "))}, nil, done
}

func (cal *Calendar) parseEvent(l *Line) error {

	if cal.buffer == nil {
		return nil
	}

	switch l.Key {
	case "UID":
		if err := resolve(cal, l, &cal.buffer.Uid, resolveString, nil); err != nil {
			return err
		}
	case "SUMMARY":
		if err := resolve(cal, l, &cal.buffer.Summary, resolveString, nil); err != nil {
			return err
		}

	case "DESCRIPTION":
		if err := resolve(cal, l, &cal.buffer.Description, resolveString, nil); err != nil {
			return err
		}
	case "DTSTART":
		if err := resolve(cal, l, &cal.buffer.Start, resolveDate, func(cal *Calendar, out *time.Time) {
			cal.buffer.RawStart = &RawDate{Value: l.Value, Params: l.Params}
		}); err != nil {
			return err
		}

	case "DTEND":
		if err := resolve(cal, l, &cal.buffer.End, resolveDateEnd, func(cal *Calendar, out *time.Time) {
			cal.buffer.RawEnd = &RawDate{Value: l.Value, Params: l.Params}
		}); err != nil {
			return err
		}
	case "DURATION":
		/*
		* Duration should be parsed in conjunction with DTSTART
		* If DTSTART has not been processed, we add to delayed attributes for processing last
		 */
		if cal.buffer.Start == nil {
			cal.buffer.delayed = append(cal.buffer.delayed, l)
			return nil
		}

		if err := resolve(cal, l, &cal.buffer.Duration, resolveDuration, func(cal *Calendar, out *time.Duration) {
			if out != nil {
				cal.buffer.Duration = out
				end := cal.buffer.Start.Add(*out)
				cal.buffer.End = &end
			}
		}); err != nil {
			return err
		}

	case "DTSTAMP":
		if err := resolve(cal, l, &cal.buffer.Stamp, resolveDate, nil); err != nil {
			return err
		}

	case "CREATED":
		if err := resolve(cal, l, &cal.buffer.Created, resolveDate, nil); err != nil {
			return err
		}

	case "LAST-MODIFIED":
		if err := resolve(cal, l, &cal.buffer.LastModified, resolveDate, nil); err != nil {
			return err
		}
	case "RRULE":
		if len(cal.buffer.RecurrenceRule) != 0 {
			return NewDuplicateAttribute(l.Key, l.Value)
		}
		if cal.buffer.RecurrenceRule == nil || cal.Duplicate.Mode == DuplicateModeKeepLast {
			var err error

			cal.buffer.IsRecurring = true

			if cal.buffer.RecurrenceRule, err = members.ParseRecurrenceRule(l.Value); err != nil {
				return err
			}
		}
	case "RECURRENCE-ID":
		if err := resolve(cal, l, &cal.buffer.RecurrenceId, resolveString, nil); err != nil {
			return err
		}
	case "EXDATE":
		/*
		*	Reference: https://icalendar.org/iCalendar-RFC-5545/3-8-5-1-exception-date-times.html
		*	Several parameters are allowed.  We should pass parameters we have
		 */

		d, err := members.ParseTime(l.Value, l.Params, members.TimeStart, false, cal.AllDayEventsTZ)
		if err == nil {
			cal.buffer.ExcludeDates = append(cal.buffer.ExcludeDates, d)
		}

	case "SEQUENCE":
		cal.buffer.Sequence, _ = strconv.Atoi(l.Value)
	case "LOCATION":
		if err := resolve(cal, l, &cal.buffer.Location, resolveString, nil); err != nil {
			return err
		}

	case "STATUS":
		if err := resolve(cal, l, &cal.buffer.Status, resolveString, nil); err != nil {
			return err
		}
	case "ORGANIZER":
		if err := resolve(cal, l, &cal.buffer.Organizer, resolveOrganizer, nil); err != nil {
			return err
		}

	case "ATTENDEE":
		attendee := &Attendee{
			Value: l.Value,
		}

		for key, val := range l.Params {
			key := strings.ToUpper(key)
			switch key {
			case "CN":
				attendee.Cn = val
			case "DIR":
				attendee.DirectoryDn = val

			case "PARTSTAT":
				attendee.Status = val

			default:
				if strings.HasPrefix(key, "X-") {
					if attendee.CustomAttributes == nil {
						attendee.CustomAttributes = make(map[string]string)
					}
					attendee.CustomAttributes[key] = val
				}
			}
		}
		cal.buffer.Attendees = append(cal.buffer.Attendees, attendee)
	case "ATTACH":
		cal.buffer.Attachments = append(cal.buffer.Attachments, &Attachment{
			Type:     l.Params["VALUE"],
			Encoding: l.Params["ENCODING"],
			Mime:     l.Params["FMTTYPE"],
			Filename: l.Params["FILENAME"],
			Value:    l.Value,
		})

	case "GEO":
		if err := resolve(cal, l, &cal.buffer.LatLng, resolveLatLng, nil); err != nil {
			return err
		}
	case "CATEGORIES":
		cal.buffer.Categories = strings.Split(l.Value, ",")
	case "URL":
		cal.buffer.Url = l.Value
	case "COMMENT":
		cal.buffer.Comment = l.Value
	case "CLASS":
		cal.buffer.Class = l.Value
	default:
		key := strings.ToUpper(l.Key)
		if strings.HasPrefix(key, "X-") {
			if cal.buffer.CustomAttributes == nil {
				cal.buffer.CustomAttributes = make(map[string]string)
			}
			cal.buffer.CustomAttributes[key] = l.Value
		}

	}

	return nil
}

func (cal *Calendar) checkEvent() error {
	if cal.buffer.Uid == "" {
		cal.buffer.Valid = false
		return fmt.Errorf("could not parse event without UID")
	}
	if cal.buffer.Start == nil {
		cal.buffer.Valid = false
		return fmt.Errorf("could not parse event without DTSTART")
	}
	if cal.buffer.Stamp == nil {
		cal.buffer.Valid = false
		return fmt.Errorf("could not parse event without DTSTAMP")
	}
	if cal.buffer.RawEnd.Value != "" && cal.buffer.Duration != nil {
		return fmt.Errorf("only one of DTEND and DURATION must be provided")
	}

	return nil
}
