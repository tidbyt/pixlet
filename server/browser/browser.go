// Package browser provides the ability to send WebP images to a browser over
// websockets.
package browser

import (
	_ "embed"
	"fmt"
	"html/template"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"golang.org/x/sync/errgroup"
	"tidbyt.dev/pixlet/server/fanout"
	"tidbyt.dev/pixlet/server/loader"
)

// Browser provides a structure for serving WebP images over websockets to
// a web browser.
type Browser struct {
	addr       string      // The address to listen on.
	title      string      // The title of the HTML document.
	updateChan chan string // A channel of base64 encoded WebP images.
	watch      bool
	fo         *fanout.Fanout
	r          *mux.Router
	tmpl       *template.Template
	loader     *loader.Loader
}

//go:embed preview-mask.png
var previewMask []byte

//go:embed favicon.png
var favicon []byte

//go:embed preview.html
var previewHTML string

// previewData is used to populate the HTML template.
type previewData struct {
	Title string
	WebP  string
	Watch bool
}

// NewBrowser sets up a browser structure. Call Run() to kick off the main loops.
func NewBrowser(addr string, title string, watch bool, updateChan chan string, l *loader.Loader) (*Browser, error) {
	tmpl, err := template.New("preview").Parse(previewHTML)
	if err != nil {
		return nil, err
	}

	b := &Browser{
		updateChan: updateChan,
		addr:       addr,
		fo:         fanout.NewFanout(),
		tmpl:       tmpl,
		title:      title,
		loader:     l,
		watch:      watch,
	}

	r := mux.NewRouter()
	r.HandleFunc("/", b.rootHandler)
	r.HandleFunc("/ws", b.websocketHandler)
	r.HandleFunc("/favicon.png", b.faviconHandler)
	r.HandleFunc("/preview-mask.png", b.previewMaskHandler)
	b.r = r

	return b, nil
}

// Run starts the server process and runs forever in a blocking fashion. The
// main routines include an update watcher to process incomming changes to the
// webp and running the http handlers.
func (b *Browser) Run() error {
	defer b.fo.Quit()

	g := errgroup.Group{}
	g.Go(b.updateWatcher)
	g.Go(b.serveHTTP)

	return g.Wait()
}

func (b *Browser) serveHTTP() error {
	log.Printf("listening at http://%s\n", b.addr)
	return http.ListenAndServe(b.addr, b.r)
}

func (b *Browser) faviconHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "image/png")
	w.Write(favicon)
}

func (b *Browser) previewMaskHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "image/png")
	w.Write(previewMask)
}

func (b *Browser) websocketHandler(w http.ResponseWriter, r *http.Request) {
	if !b.watch {
		return
	}

	var upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("error establishing a new connection %v\n", err)
		return
	}

	b.fo.NewClient(conn)
}

func (b *Browser) updateWatcher() error {
	for {
		select {
		case webp := <-b.updateChan:
			b.fo.Broadcast(webp)
		}
	}
}

func (b *Browser) rootHandler(w http.ResponseWriter, r *http.Request) {
	config := make(map[string]string)
	for k, vals := range r.URL.Query() {
		config[k] = vals[0]
	}

	webp, err := b.loader.LoadApplet(config)
	if err != nil {
		w.WriteHeader(500)
		fmt.Fprintln(w, err)
		return
	}

	data := previewData{
		Title: b.title,
		Watch: b.watch,
		WebP:  webp,
	}

	w.Header().Set("Content-Type", "text/html")
	b.tmpl.Execute(w, data)
}
