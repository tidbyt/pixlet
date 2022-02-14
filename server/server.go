package server

import (
	"fmt"

	"golang.org/x/sync/errgroup"
	"tidbyt.dev/pixlet/server/browser"
	"tidbyt.dev/pixlet/server/loader"
	"tidbyt.dev/pixlet/server/watcher"
)

// Server provides functionality to serve Starlark over HTTP. It has
// functionality to watch a file and hot reload the browser on changes.
type Server struct {
	watcher *watcher.Watcher
	browser *browser.Browser
	loader  *loader.Loader
	watch   bool
}

// NewServer creates a new server initialized with the applet.
func NewServer(host string, port int, watch bool, filename string) (*Server, error) {
	fileChanges := make(chan bool, 100)
	w := watcher.NewWatcher(filename, fileChanges)

	updatesChan := make(chan loader.Update, 100)
	l, err := loader.NewLoader(filename, watch, fileChanges, updatesChan)
	if err != nil {
		return nil, err
	}

	addr := fmt.Sprintf("%s:%d", host, port)
	b, err := browser.NewBrowser(addr, filename, watch, updatesChan, l)
	if err != nil {
		return nil, err
	}

	return &Server{
		watcher: w,
		browser: b,
		loader:  l,
		watch:   watch,
	}, nil
}

// Run serves the http server and runs forever in a blocking fashion.
func (s *Server) Run() error {
	g := errgroup.Group{}

	g.Go(s.loader.Run)
	g.Go(s.browser.Run)
	if s.watch {
		g.Go(s.watcher.Run)
		s.loader.LoadApplet(make(map[string]string))
	}

	return g.Wait()
}
