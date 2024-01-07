package gtfs

import (
	"fmt"
	"net/url"
	"sync"
	"time"

	"github.com/mitchellh/hashstructure/v2"
	starlibtime "go.starlark.net/lib/time"
	"go.starlark.net/starlark"
	"go.starlark.net/starlarkstruct"

	"tidbyt.dev/gtfs"
	"tidbyt.dev/gtfs/model"
)

const (
	ModuleName = "gtfs"
)

// The GTFSManager must be set for module to load
var Manager *gtfs.Manager

var (
	once              sync.Once
	module            starlark.StringDict
	moduleInitialized bool

	mutex sync.RWMutex
	cache map[cacheKey]*GTFS
)

type cacheKey struct {
	appID       string
	staticURL   string
	realtimeURL string
}

func LoadModule() (starlark.StringDict, error) {
	once.Do(func() {
		cache = map[cacheKey]*GTFS{}
		if Manager != nil {
			module = starlark.StringDict{
				ModuleName: &starlarkstruct.Module{
					Name: ModuleName,
					Members: starlark.StringDict{
						"GTFS": starlark.NewBuiltin("GTFS", newGTFS),
					},
				},
			}
			moduleInitialized = true
		}
	})

	if !moduleInitialized {
		return nil, fmt.Errorf("gtfs not initialized")
	}

	return module, nil
}

type GTFS struct {
	static   *gtfs.Static
	realtime *gtfs.Realtime

	routes *starlark.Dict
	stops  *starlark.Dict
	trips  *starlark.Dict

	departures  *starlark.Builtin
	directions  *starlark.Builtin
	nearbyStops *starlark.Builtin
}

func buildHeaders(headers *starlark.Dict) (map[string]string, error) {
	goHeaders := map[string]string{}

	if headers != nil && headers.Len() > 0 {
		for _, kv := range headers.Items() {
			k := kv.Index(0)
			v := kv.Index(1)

			kStr, ok := k.(starlark.String)
			if !ok {
				return nil, fmt.Errorf("key must be string, found %s: %s", k.Type(), k.String())
			}
			vStr, ok := v.(starlark.String)
			if !ok {
				return nil, fmt.Errorf("value be string, found %s: %s", v.Type(), v.String())
			}

			goHeaders[kStr.GoString()] = vStr.GoString()
		}
	}

	return goHeaders, nil
}

func loadGTFS(
	appID string,
	staticURL string,
	staticHeaders map[string]string,
	realtimeURL string,
	realtimeHeaders map[string]string,
) (*GTFS, error) {

	key := cacheKey{appID, staticURL, realtimeURL}

	// From cache, if available.
	mutex.RLock()
	g, found := cache[cacheKey{appID, staticURL, realtimeURL}]
	mutex.RUnlock()
	if found {
		return g, nil
	}

	// Otherwise from Manager. Double-checked locking to avoid
	// redundant loads.
	mutex.Lock()
	defer mutex.Unlock()

	g, found = cache[cacheKey{appID, staticURL, realtimeURL}]
	if found {
		return g, nil
	}

	// Still not in cache, so we gotta load it. Static first
	static, err := Manager.LoadStaticAsync(appID, staticURL, staticHeaders, time.Now())
	if err != nil {
		return nil, err
	}

	// Realtime, if requested
	var realtime *gtfs.Realtime
	if realtimeURL != "" {
		realtime, err = Manager.LoadRealtime(
			appID,
			static,
			realtimeURL,
			realtimeHeaders,
			time.Now(),
		)
		if err != nil {
			return nil, fmt.Errorf("loading realtime feed: %w", err)
		}
	}

	// Here's the object
	g = &GTFS{
		static:      static,
		realtime:    realtime,
		departures:  starlark.NewBuiltin("departures", gtfsDepartures),
		directions:  starlark.NewBuiltin("directions", gtfsDirections),
		nearbyStops: starlark.NewBuiltin("nearby_stops", gtfsNearbyStops),
	}

	// Now load stops
	stops, err := g.static.Reader.Stops()
	if err != nil {
		return nil, fmt.Errorf("getting stops: %w", err)
	}
	g.stops = starlark.NewDict(len(stops))
	for _, s := range stops {
		g.stops.SetKey(starlark.String(s.ID), makeStop(s))
	}
	g.stops.Freeze()

	// Routes
	routes, err := g.static.Reader.Routes()
	if err != nil {
		return nil, fmt.Errorf("getting routes: %w", err)
	}
	g.routes = starlark.NewDict(len(routes))
	for _, r := range routes {
		g.routes.SetKey(starlark.String(r.ID), makeRoute(r))
	}
	g.routes.Freeze()

	// Trips
	trips, err := g.static.Reader.Trips()
	if err != nil {
		return nil, fmt.Errorf("getting trips: %w", err)
	}
	g.trips = starlark.NewDict(len(trips))
	for _, t := range trips {
		g.trips.SetKey(starlark.String(t.ID), makeTrip(t))
	}
	g.trips.Freeze()

	// Write to cache, and we're done
	cache[key] = g

	return g, nil
}

func newGTFS(thread *starlark.Thread, _ *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
	var (
		staticURL       starlark.String
		realtimeURL     starlark.String
		staticHeaders   *starlark.Dict
		realtimeHeaders *starlark.Dict
	)

	err := starlark.UnpackArgs(
		"GTFS", args, kwargs,
		"static_url", &staticURL,
		"realtime_url?", &realtimeURL,
		"static_headers?", &staticHeaders,
		"realtime_headers?", &realtimeHeaders,
	)
	if err != nil {
		return nil, fmt.Errorf("unpacking GTFS args: %w", err)
	}

	appID, ok := thread.Local("appID").(string)
	if !ok {
		return nil, fmt.Errorf("appID not in thread local storage")
	}

	// Validate static params
	goStaticURL := staticURL.GoString()
	if goStaticURL == "" {
		return nil, fmt.Errorf("static URL is required")
	}
	goStaticHeaders, err := buildHeaders(staticHeaders)
	if err != nil {
		return nil, fmt.Errorf("static headers: %w", err)
	}

	u, e := url.Parse(staticURL.GoString())
	if e != nil {
		return nil, fmt.Errorf("bad static URL: %w", e)
	}
	if u.Scheme != "http" && u.Scheme != "https" {
		return nil, fmt.Errorf("static URL must be http or https")
	}

	// Validate realtime params
	var goRealtimeHeaders map[string]string
	goRealtimeURL := realtimeURL.GoString()
	if goRealtimeURL != "" {
		u, e = url.Parse(goRealtimeURL)
		if e != nil {
			return nil, fmt.Errorf("bad realtime URL: %w", e)
		}
		if u.Scheme != "http" && u.Scheme != "https" {
			return nil, fmt.Errorf("realtime URL must be http:// or https://")
		}

		if realtimeHeaders != nil {
			goRealtimeHeaders, err = buildHeaders(realtimeHeaders)
			if err != nil {
				return nil, fmt.Errorf("realtime headers: %w", err)
			}
		}
	}

	gtfs, err := loadGTFS(appID, goStaticURL, goStaticHeaders, goRealtimeURL, goRealtimeHeaders)
	if err != nil {
		return nil, fmt.Errorf("loading GTFS: %w", err)
	}

	return gtfs, nil
}

func gtfsDepartures(thread *starlark.Thread, b *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
	g := b.Receiver().(*GTFS)

	var stopID starlark.String
	var when starlibtime.Time
	var window starlibtime.Duration
	var limit starlark.Int
	var routeID starlark.String
	var directionID starlark.Int

	err := starlark.UnpackArgs(
		"departures",
		args, kwargs,
		"stop_id", &stopID,
		"when", &when,
		"window", &window,
		"limit?", &limit,
		"route_id?", &routeID,
		"direction_id?", &directionID,
	)
	if err != nil {
		return nil, fmt.Errorf("unpacking arguments for GTFS.departures: %s", err)
	}

	// TODO: would assigning -1 before UnpackArgs() work as well?

	// These will be 0 if not provided, but we need them to
	// default to -1
	goLimit := -1
	goDirectionID := int8(-1)
	for _, kwarg := range kwargs {
		key := kwarg.Index(0).(starlark.String).GoString()
		if key == "direction_id" {
			goDirectionID = int8(directionID.BigInt().Int64())
		} else if key == "limit" {
			goLimit = int(limit.BigInt().Int64())
		}
	}

	// Get departures. Realtime if available, otherwise Static.
	var departures []model.Departure
	if g.realtime != nil {
		departures, err = g.realtime.Departures(stopID.GoString(), time.Time(when), time.Duration(window), goLimit, routeID.GoString(), goDirectionID, nil)
		if err != nil {
			return nil, fmt.Errorf("getting realtime departures: %w", err)
		}
	} else {
		departures, err = g.static.Departures(stopID.GoString(), time.Time(when), time.Duration(window), goLimit, routeID.GoString(), goDirectionID, nil)
		if err != nil {
			return nil, fmt.Errorf("getting static departures: %w", err)
		}
	}

	// Returns as dicts
	departureList := make([]starlark.Value, 0, len(departures))
	for _, d := range departures {
		departureList = append(departureList, makeDeparture(d))
	}

	return starlark.NewList(departureList), nil
}

func gtfsDirections(thread *starlark.Thread, b *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
	g := b.Receiver().(*GTFS)

	var stopID starlark.String

	if err := starlark.UnpackArgs(
		"directions",
		args, kwargs,
		"stopID", &stopID,
	); err != nil {
		return nil, fmt.Errorf("unpacking arguments for GTFS.directions: %s", err)
	}

	directions, err := g.static.RouteDirections(stopID.GoString())
	if err != nil {
		return nil, fmt.Errorf("getting directions: %w", err)
	}

	directionList := make([]starlark.Value, 0, len(directions))
	for _, rd := range directions {
		headsigns := make([]starlark.Value, 0, len(rd.Headsigns))
		for _, hs := range rd.Headsigns {
			headsigns = append(headsigns, starlark.String(hs))
		}

		directionList = append(directionList, Object{
			name: "RouteDirection",
			entries: map[string]starlark.Value{
				"stop_id":      starlark.String(rd.StopID),
				"route_id":     starlark.String(rd.RouteID),
				"direction_id": starlark.MakeInt(int(rd.DirectionID)),
				"headsigns":    starlark.NewList(headsigns),
			},
		})
	}

	return starlark.NewList(directionList), nil
}

func gtfsNearbyStops(thread *starlark.Thread, b *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
	g := b.Receiver().(*GTFS)

	var lat starlark.Float
	var lon starlark.Float
	var limit starlark.Int

	if err := starlark.UnpackArgs(
		"departures",
		args, kwargs,
		"lat", &lat,
		"lon", &lon,
		"limit", &limit,
	); err != nil {
		return nil, fmt.Errorf("unpacking arguments for GTFS.departures: %s", err)
	}

	stops, err := g.static.NearbyStops(float64(lat), float64(lon), int(limit.BigInt().Int64()), nil)
	if err != nil {
		return nil, fmt.Errorf("getting nearby stops: %w", err)
	}

	stopList := make([]starlark.Value, 0, len(stops))
	for _, s := range stops {
		stopList = append(stopList, makeStop(s))
	}

	return starlark.NewList(stopList), nil
}

func (g *GTFS) AttrNames() []string {
	return []string{
		"routes",
		"stops",
		"trips",
		"departures",
		"directions",
		"nearby_stops",
	}
}

func (g *GTFS) Attr(name string) (starlark.Value, error) {
	switch name {
	case "routes":
		return g.routes, nil
	case "stops":
		return g.stops, nil
	case "trips":
		return g.trips, nil
	case "departures":
		return g.departures.BindReceiver(g), nil
	case "directions":
		return g.directions.BindReceiver(g), nil
	case "nearby_stops":
		return g.nearbyStops.BindReceiver(g), nil
	default:
		return nil, nil
	}
}

func (g *GTFS) Hash() (uint32, error) {
	// TODO: Is this reasonable?
	sum, err := hashstructure.Hash(g, hashstructure.FormatV2, nil)
	return uint32(sum), err
}

func (g *GTFS) String() string {
	// TODO: include url and such in here
	return fmt.Sprintf("GTFS()")
}

func (g *GTFS) Type() string         { return "GTFS" }
func (g *GTFS) Freeze()              {}
func (g *GTFS) Truth() starlark.Bool { return true }
