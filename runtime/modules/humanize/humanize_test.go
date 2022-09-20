package humanize_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"tidbyt.dev/pixlet/runtime"
)

var humanSource = `
load("time.star", "time")
load("humanize.star", "humanize")

def assert(success, message=None):
    if not success:
        fail(message or "assertion failed")

# Set up
now = time.now()
tomorrow = now + time.parse_duration("26h")

# Call methods.
humanized_time_past = humanize.time(now - time.parse_duration("48h"))
humanized_time_future = humanize.time(tomorrow)
humanized_rel_time = humanize.relative_time(now, tomorrow)
humanized_date_format = humanize.time_format("yyyy-MM-dd")
humanized_date_format_date = humanize.time_format("yyyy-MM-dd", now)
humanized_day_of_week = humanize.day_of_week(now)
humanized_size = humanize.bytes(1401946112)
humanized_size_iec = humanize.bytes(1401946112, True)
humanized_size_parsed = humanize.parse_bytes("42 MB")
humanized_comma_int = humanize.comma(123456)
humanized_comma_float = humanize.comma(123456.78)
humanized_float = humanize.float("#,###.##", 123456.78002)
humanized_int = humanize.int("#,###.", 123456)
humanized_ordinal = humanize.ordinal(1)
humanized_ftoa = humanize.ftoa(3.1450000)
humanized_ftoa_digits = humanize.ftoa(3.1450000, 2)
humanized_ftoa_digits_z = humanize.ftoa(3.1450000, 0)
humanized_plural = humanize.plural(42, "object")
humanized_plural_test = humanize.plural(1, "star", "")
humanized_plural_word = humanize.plural_word(1, "star", "")
humanized_word_series = humanize.word_series(["foo", "bar", "baz"], "and")
humanized_word_series_oxford = humanize.oxford_word_series(["foo", "bar", "baz"], "and")
iso_date = now.format(humanized_date_format)
humanized_url_encode = humanize.url_encode("bar baz")
humanized_url_decode = humanize.url_decode("http://example.com/foo=bar+baz")

# Assert.
assert(humanized_time_past == "2 days ago")
assert(humanized_time_future == "1 day from now")
assert(humanized_rel_time == "1 day ")
assert(humanized_date_format == "2006-01-02")
assert(humanized_date_format_date == iso_date)
weekday = now.format("Monday")
weekday_names = ["Sunday", "Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday"]
assert(weekday_names[humanized_day_of_week] == weekday)
assert(humanized_size == "1.4 GB")
assert(humanized_size_iec == "1.3 GiB")
assert(humanized_size_parsed == 42000000)
assert(humanized_comma_int == "123,456")
assert(humanized_comma_float == "123,456.78")
assert(humanized_float == "123,456.78")
assert(humanized_int == "123,456")
assert(humanized_ordinal == "1st")
assert(humanized_ftoa == "3.145")
assert(humanized_ftoa_digits == "3.14")
assert(humanized_ftoa_digits_z == "3")
assert(humanized_plural == "42 objects")
assert(humanized_plural_test == "1 star")
assert(humanized_plural_word == "star")
assert(humanized_word_series == "foo, bar and baz")
assert(humanized_word_series_oxford == "foo, bar, and baz")
assert(humanized_url_encode == "bar+baz")
assert(humanized_url_decode == "http://example.com/foo=bar baz")

def main():
	return []
`

func TestHumanize(t *testing.T) {
	app := &runtime.Applet{}
	err := app.Load("human.star", []byte(humanSource), nil)
	assert.NoError(t, err)

	screens, err := app.Run(map[string]string{})
	assert.NoError(t, err)
	assert.NotNil(t, screens)
}
