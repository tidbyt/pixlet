package parser

import (
	"bufio"
	"io"
	"strings"
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

/*
* The iCal mandates that lines longer than 75 octets require a linebreak.
* The format uses progressive whitespace indentation to denote a line is continued on a new line.
* https://icalendar.org/iCalendar-RFC-5545/3-1-content-lines.html
 */
func (cal *Calendar) findNextLine() string {
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

	return l
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
