package server

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"golang.org/x/sync/errgroup"
	"tidbyt.dev/pixlet/server/browser"
	"tidbyt.dev/pixlet/server/loader"
	"tidbyt.dev/pixlet/tools"
)

// Server provides functionality to serve Starlark over HTTP. It has
// functionality to watch a file and hot reload the browser on changes.
type Server struct {
	watcher *Watcher
	browser *browser.Browser
	loader  *loader.Loader
	watch   bool
}

// NewServer creates a new server initialized with the applet.
func NewServer(host string, port int, watch bool, path string, maxDuration int, timeout int) (*Server, error) {
	fileChanges := make(chan bool, 100)

	// check if path exists, and whether it is a directory or a file
	info, err := os.Stat(path)
	if err != nil {
		return nil, fmt.Errorf("failed to stat %s: %w", path, err)
	}

	var fs fs.FS
	var w *Watcher
	if info.IsDir() {
		fs = os.DirFS(path)
		w = NewWatcher(path, fileChanges)
	} else {
		if !strings.HasSuffix(path, ".star") {
			return nil, fmt.Errorf("script file must have suffix .star: %s", path)
		}

		fs = tools.NewSingleFileFS(path)
		w = NewWatcher(filepath.Dir(path), fileChanges)
	}

	updatesChan := make(chan loader.Update, 100)
	l, err := loader.NewLoader(fs, watch, fileChanges, updatesChan, maxDuration, timeout)
	if err != nil {
		return nil, err
	}

	addr := fmt.Sprintf("%s:%d", host, port)
	b, err := browser.NewBrowser(addr, filepath.Base(path), watch, updatesChan, l)
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
