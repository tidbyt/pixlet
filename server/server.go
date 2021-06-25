package server

import (
	"encoding/base64"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"sync"

	"github.com/fsnotify/fsnotify"
	"github.com/gorilla/websocket"
	"tidbyt.dev/pixlet/encode"
	"tidbyt.dev/pixlet/runtime"
)

// Server provides functionality to serve Starlark over HTTP. It has
// functionality to watch a file and hot reload the browser on changes.
type Server struct {
	host     string
	port     int
	watch    bool
	filename string
	mutex    sync.RWMutex
	applet   runtime.Applet
}

type websocketEvent struct {
	Type string `json:"type"`
}

// NewServer creates a new server initialized with the applet.
func NewServer(host string, port int, watch bool, filename string) *Server {
	applet := runtime.Applet{}

	return &Server{
		host:     host,
		port:     port,
		watch:    watch,
		filename: filename,
		mutex:    sync.RWMutex{},
		applet:   applet,
	}
}

// Run serves the http server and runs forever in a blocking fashion.
func (s *Server) Run() error {
	err := loadScript(&s.applet, s.filename)
	if err != nil {
		return err
	}
	log.Println("loaded", s.filename)

	http.HandleFunc("/", s.serveRoot)
	http.HandleFunc("/favicon.ico", s.serveFavicon)
	http.HandleFunc("/ws", s.serveWebsocket)

	log.Printf("listening at http://%s:%d\n", s.host, s.port)
	return http.ListenAndServe(fmt.Sprintf("%s:%d", s.host, s.port), nil)

}
func (s *Server) serveFavicon(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(404)
}

func (s *Server) serveWebsocket(w http.ResponseWriter, r *http.Request) {
	if !s.watch {
		return
	}

	var upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("error establishing a new connection %v\n", err)
	}

	go s.fileWatcher(conn)
}

func (s *Server) serveRoot(w http.ResponseWriter, r *http.Request) {
	config := make(map[string]string)
	for k, vals := range r.URL.Query() {
		config[k] = vals[0]
	}

	s.mutex.RLock()
	defer s.mutex.RUnlock()

	roots, err := s.applet.Run(config)
	if err != nil {
		log.Printf("Error running script: %s\n", err)
		return
	}

	webp, err := encode.ScreensFromRoots(roots).EncodeWebP()
	if err != nil {
		fmt.Printf("Error rendering: %s\n", err)
		return
	}

	w.Header().Set("Content-Type", "text/html")
	s.writePreviewHTML(w, s.host, s.port, webp)
}

func (s *Server) writePreviewHTML(w io.Writer, host string, port int, webp []byte) {
	fmt.Fprintln(
		w,
		`<html>
		<style type="text/css">
		img {
			image-rendering: pixelated;
			image-rendering: -moz-crisp-edges;
			image-rendering: crisp-edges;
			width: 100%;
		}
		</style>
		<body bgcolor="black">
		<div style="border: solid 1px white">
		`,
	)

	fmt.Fprintf(
		w,
		`<img src="data:image/webp;base64,%s" />`,
		base64.StdEncoding.EncodeToString(webp),
	)

	fmt.Fprintf(
		w,
		`<script>
			function connect() {
				var conn = new WebSocket("ws://%s:%d/ws");

				conn.onmessage = function(e) {
				  var data = JSON.parse(e.data);
				  switch(data.type) {
					case "update":
					  conn.close(1000, "Reloading page after receiving update");
					  console.log("Reloading page after receiving update");
					  location.reload(true);
					  break;
		
					default:
					  console.log("Don't know how to handle type '${data.type}'");
				  }
				}

				conn.onclose = function(event) {
					console.log("Attempting to reconnect");
					setTimeout(function() {
						connect();
					}, 1000);
				}

				conn.onerror = function(event) {
					console.log("error");
					console.log(event);
				}
			}
			connect()
		</script>`,
		host,
		port,
	)

	fmt.Fprintln(w, `</div></body></html>`)
}

func (s *Server) fileWatcher(conn *websocket.Conn) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Printf("error watching for changes: %v\n", err)
		os.Exit(1)
	}

	watcher.Add(s.filename)

	for {
		event, ok := <-watcher.Events
		if !ok {
			break
		}
		if (event.Op & fsnotify.Rename) != 0 {
			// When Vim saves a file, we get a Rename event followed
			// by silence. Re-adding allows us to capture future
			// events.
			watcher.Remove(event.Name)
			watcher.Add(s.filename)
		} else if (event.Op & (fsnotify.Write | fsnotify.Chmod)) != 0 {
			log.Printf("detected updates for %s, reloading\n", s.filename)

			// Reloading on Write is sufficient for most editors,
			// but with Vim we only get Chmod. No clue why.
			s.mutex.Lock()
			err := loadScript(&s.applet, s.filename)
			s.mutex.Unlock()
			if err != nil {
				log.Printf("error on reload: %v\n", err)
				continue
			}

			// Send update on websocket. If there is an error, break out of the
			// loop. When the browser reloads, it will recreate it's connection
			// and this goroutine will spawn again.
			err = conn.WriteJSON(websocketEvent{Type: "update"})
			if err != nil {
				break
			}
		}
	}
}

func loadScript(applet *runtime.Applet, filename string) error {
	src, err := ioutil.ReadFile(filename)
	if err != nil {
		return fmt.Errorf("failed to read file %s: %v", filename, err)
	}

	runtime.InitCache(runtime.NewInMemoryCache())

	err = applet.Load(filename, src, nil)
	if err != nil {
		return fmt.Errorf("failed to load applet: %v", err)
	}

	return nil
}
