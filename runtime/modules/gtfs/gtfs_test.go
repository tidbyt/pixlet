package gtfs_test

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"tidbyt.dev/gtfs"
	gtfs_storage "tidbyt.dev/gtfs/storage"
	"tidbyt.dev/gtfs/testutil"
	"tidbyt.dev/pixlet/runtime"
	pixlet_gtfs "tidbyt.dev/pixlet/runtime/modules/gtfs"
)

var gtfsSrc = `
load("gtfs.star", "gtfs")
load("time.star", "time")

def test_gtfs():
    g = gtfs.GTFS("%s")
    if len(g.stops) != 2:
        fail("stops")
    if g.stops["s1"].name != "Test Stop":
	fail("stop s1 name")
    if g.stops["s2"].name != "Test Stop 2":
	fail("stop s2 name")


    now = time.now()
    departures = g.departures(
        "s1",
        when=time.time(year=now.year, month=now.month, day=now.day, hour=1),
        window=time.parse_duration("2h"),
    )

    if len(departures) != 1:
	fail("departures")
    if departures[0].stop_id != "s1":
	fail("departure stop_id != s1")
    if departures[0].time.format("15:04:05") != "01:30:00":
	fail("departure time != 01:30:00")

test_gtfs()

def main():
    return []
`

func TestGTFS(t *testing.T) {

	// Test server serving a static GTFS feed, which shows 1
	// departure from s1 at 01:30, every day, starting about 10
	// days ago, and going on for another 10 days.
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		past := time.Now().Add(-240 * time.Hour).Format("20060102")
		future := time.Now().Add(240 * time.Hour).Format("20060102")
		w.Write(testutil.BuildZip(t, map[string][]string{
			"agency.txt": []string{
				"agency_id,agency_name,agency_url,agency_timezone,agency_lang,agency_phone",
				"a1,Test Agency,http://example.com,UTC,en,555-555-5555",
			},
			"calendar.txt": []string{
				"service_id,monday,tuesday,wednesday,thursday,friday,saturday,sunday,start_date,end_date",
				"s1,1,1,1,1,1,0,0," + past + "," + future,
			},
			"routes.txt": []string{
				"route_id,agency_id,route_short_name,route_long_name,route_desc,route_type,route_url,route_color,route_text_color",
				"r1,a1,one,Test Route,Test Route Description,3,http://example.com,FFFF00,FF0000",
			},
			"trips.txt": []string{
				"route_id,service_id,trip_id,trip_headsign,direction_id",
				"r1,s1,t1,Test Trip,0",
			},
			"stops.txt": []string{
				"stop_id,stop_name,stop_desc,stop_lat,stop_lon,location_type,parent_station",
				"s1,Test Stop,Test Stop Description,40.1234,-74.1234,0,",
				"s2,Test Stop 2,Test Stop 2 Description,40.1235,-74.1235,0,",
			},
			"stop_times.txt": []string{
				"trip_id,arrival_time,departure_time,stop_id,stop_sequence",
				"t1,01:00:00,01:30:00,s1,1",
				"t1,02:00:00,02:30:00,s2,2",
			},
		}))
	}))

	src := fmt.Sprintf(gtfsSrc, ts.URL+"/gtfs.zip")

	// A Manager must be set for the module to load
	s, err := gtfs_storage.NewSQLiteStorage()
	require.NoError(t, err)
	pixlet_gtfs.Manager = gtfs.NewManager(s)

	// At first, runs should fail, as no feed data has been
	// downloaded.
	app := &runtime.Applet{}
	err = app.Load("gtfs", "gtfs_test.star", []byte(src), nil)
	assert.ErrorContains(t, err, "no active feed")
	err = app.Load("gtfs", "gtfs_test.star", []byte(src), nil)
	assert.ErrorContains(t, err, "no active feed")
	err = app.Load("gtfs", "gtfs_test.star", []byte(src), nil)
	assert.ErrorContains(t, err, "no active feed")

	// Calling Manager.Refresh() will download the feed
	err = pixlet_gtfs.Manager.Refresh(context.Background())
	require.NoError(t, err)

	// Now the script will succeed
	err = app.Load("gtfs", "gtfs_test.star", []byte(src), nil)
	assert.NoError(t, err)
	err = app.Load("gtfs", "gtfs_test.star", []byte(src), nil)
	assert.NoError(t, err)
	err = app.Load("gtfs", "gtfs_test.star", []byte(src), nil)
	assert.NoError(t, err)

}
