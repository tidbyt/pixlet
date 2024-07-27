package parser

import (
	"bufio"
	"fmt"
	"strings"
	"time"
)

const (
	StrictModeFailFeed = iota
	StrictModeFailAttribute
	StrictModeFailEvent
)

const (
	DuplicateModeFailStrict = iota
	DuplicateModeKeepFirst
	DuplicateModeKeepLast
)

type StrictParams struct {
	Mode int
}

type DuplicateParams struct {
	Mode int
}

type DuplicateAttributeError struct {
	Key, Value string
}

func NewDuplicateAttributeError(key string, value string) DuplicateAttributeError {
	return DuplicateAttributeError{key, value}
}

func (err DuplicateAttributeError) Error() string {
	return fmt.Sprintf("duplicate attribute '%s': %s", err.Key, err.Value)
}

type Calendar struct {
	scanner        *bufio.Scanner
	Events         []*Event
	SkipBounds     bool
	Strict         StrictParams
	Duplicate      DuplicateParams
	buffer         *Event
	Start          *time.Time
	End            *time.Time
	Method         string
	AllDayEventsTZ *time.Location
}

func (cal *Calendar) IsInRange(d Event) bool {
	if (d.Start.Before(*cal.Start) && d.End.After(*cal.Start)) ||
		(d.Start.After(*cal.Start) && d.End.Before(*cal.End)) ||
		(d.Start.Before(*cal.End) && d.End.After(*cal.End)) {
		return true
	}
	return false
}

const (
	ContextRoot = iota
	ContextEvent
	ContextUnknown
)

type Context struct {
	value int
}

func (ctx *Context) Nest() *Context {
	return &Context{ctx.value}
}

type RawDate struct {
	Params map[string]string
	Value  string
}

type Line struct {
	Key    string
	Params map[string]string
	Value  string
}

func (l *Line) Is(key, value string) bool {
	if strings.TrimSpace(l.Key) == key && strings.TrimSpace(l.Value) == value {
		return true
	}

	return false
}

func (l *Line) IsKey(key, value string) bool {
	return strings.TrimSpace(l.Key) == key
}

func (l *Line) IsValue(key, value string) bool {
	return strings.TrimSpace(l.Value) == value
}

type Event struct {
	delayed []*Line

	Uid              string
	Summary          string
	Description      string
	Categories       []string
	Start            *time.Time
	End              *time.Time
	RawStart         *RawDate
	RawEnd           *RawDate
	Duration         *time.Duration
	Stamp            *time.Time
	Created          *time.Time
	LastModified     *time.Time
	Location         string
	LatLng           LatLng
	Url              string
	Status           string
	Organizer        *Organizer
	Attendees        []*Attendee
	Attachments      []*Attachment
	IsRecurring      bool
	RecurrenceId     string
	RecurrenceRule   map[string]string
	ExcludeDates     []*time.Time
	Sequence         int
	CustomAttributes []*time.Time
	Valid            bool
	Comment          string
	Class            string
}

type Organizer struct {
	Cn          string
	DirectoryDn string
	Value       string
}

type Attendee struct {
	Cn               string
	DirectoryDn      string
	Status           string
	Value            string
	CustomAttributes map[string]string
}

type Attachment struct {
	Encoding string
	Type     string
	Mime     string
	Filename string
	Value    string
}

type LatLng struct {
	lat  string
	long string
}
