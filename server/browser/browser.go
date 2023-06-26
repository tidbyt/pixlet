// Package browser provides the ability to send WebP images to a browser over
// websockets.
package browser

import (
	_ "embed"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"golang.org/x/sync/errgroup"
	"tidbyt.dev/pixlet/dist"
	"tidbyt.dev/pixlet/server/fanout"
	"tidbyt.dev/pixlet/server/loader"
)

// Browser provides a structure for serving WebP images over websockets to
// a web browser.
type Browser struct {
	addr       string             // The address to listen on.
	title      string             // The title of the HTML document.
	updateChan chan loader.Update // A channel of base64 encoded WebP images.
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
	Title string `json:"title"`
	WebP  string `json:"webp"`
	Watch bool   `json:"-"`
	Err   string `json:"error,omitempty"`
}
type handlerRequest struct {
	ID    string `json:"id"`
	Param string `json:"param"`
}

// NewBrowser sets up a browser structure. Call Run() to kick off the main loops.
func NewBrowser(addr string, title string, watch bool, updateChan chan loader.Update, l *loader.Loader) (*Browser, error) {
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

	// In order for React Router to work, all routes that React Router should
	// manage need to return the root handler.
	r.HandleFunc("/", b.rootHandler)
	r.HandleFunc("/oauth-callback", b.rootHandler)

	// This enables the static directory containing JS and CSS to be available
	// at /static.
	r.PathPrefix("/static").Handler(http.FileServer(http.FS(dist.Static)))

	// In case we broke something or someone prefers the legacy editor, it is
	// still available for now. This will be removed in the future once we
	// have confirmed the new editor is stable.
	r.HandleFunc("/legacy", b.oldRootHandler)
	r.HandleFunc("/ws", b.websocketHandler)
	r.HandleFunc("/favicon.png", b.faviconHandler).Methods("GET")
	r.HandleFunc("/preview-mask.png", b.previewMaskHandler).Methods("GET")

	// API endpoints to support the React frontend.
	r.HandleFunc("/api/v1/preview", b.previewHandler)
	r.HandleFunc("/api/v1/preview.webp", b.imageHandler)
	r.HandleFunc("/api/v1/push", b.pushHandler)
	r.HandleFunc("/api/v1/schema", b.schemaHandler).Methods("GET")
	r.HandleFunc("/api/v1/handlers/{handler}", b.schemaHandlerHandler).Methods("POST")
	r.HandleFunc("/api/v1/ws", b.websocketHandler)
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

func (b *Browser) faviconHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "image/png")
	w.Write(favicon)
}

func (b *Browser) previewMaskHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "image/png")
	w.Write(previewMask)
}

func (b *Browser) schemaHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Write(b.loader.GetSchema())
}

func (b *Browser) schemaHandlerHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	if _, ok := vars["handler"]; !ok {
		w.WriteHeader(404)
		fmt.Fprintln(w, "no handler")
		return
	}

	msg := &handlerRequest{}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(msg)
	if err != nil {
		w.WriteHeader(500)
		fmt.Fprintln(w, err)
		return
	}

	data, err := b.loader.CallSchemaHandler(r.Context(), vars["handler"], msg.Param)
	if err != nil {
		w.WriteHeader(500)
		fmt.Fprintln(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(data))
}

func (b *Browser) imageHandler(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		w.WriteHeader(500)
		http.Error(w, "bad form data", http.StatusBadRequest)
		return
	}

	config := make(map[string]string)
	for k, val := range r.Form {
		config[k] = val[0]
	}

	webp, err := b.loader.LoadApplet(config)
	if err != nil {
		http.Error(w, "loading applet", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "image/webp")

	data, err := base64.StdEncoding.DecodeString(webp)
	if err != nil {
		http.Error(w, "decoding webp", http.StatusInternalServerError)
		return
	}

	w.Write(data)
}

func (b *Browser) previewHandler(w http.ResponseWriter, r *http.Request) {
	// Parse the request form so we can use it as config values.
	if err := r.ParseMultipartForm(100); err != nil {
		log.Printf("form parsing failed: %+v", err)
		http.Error(w, "bad form data", http.StatusBadRequest)
		return
	}
	config := make(map[string]string)
	for k, val := range r.Form {
		config[k] = val[0]
	}

	webp, err := b.loader.LoadApplet(config)
	data := &previewData{
		WebP:  webp,
		Title: b.title,
	}
	if err != nil {
		data.Err = err.Error()
	}

	d, err := json.Marshal(data)
	if err != nil {
		w.WriteHeader(500)
		fmt.Fprintln(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(d)
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
		case up := <-b.updateChan:
			b.fo.Broadcast(
				fanout.WebsocketEvent{
					Type:    fanout.EventTypeWebP,
					Message: up.WebP,
				},
			)

			if up.Err != nil {
				b.fo.Broadcast(
					fanout.WebsocketEvent{
						Type:    fanout.EventTypeErr,
						Message: up.Err.Error(),
					},
				)
			}

			if up.Schema != "" {
				b.fo.Broadcast(
					fanout.WebsocketEvent{
						Type:    fanout.EventTypeSchema,
						Message: up.Schema,
					},
				)
			}
		}
	}
}
func (b *Browser) rootHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	w.Write(dist.Index)
}

func (b *Browser) oldRootHandler(w http.ResponseWriter, r *http.Request) {
	config := make(map[string]string)
	for k, vals := range r.URL.Query() {
		config[k] = vals[0]
	}

	webp, err := b.loader.LoadApplet(config)

	data := previewData{
		Title: b.title,
		Watch: b.watch,
		WebP:  webp,
	}

	if err != nil {
		data.Err = err.Error()
	}

	w.Header().Set("Content-Type", "text/html")
	b.tmpl.Execute(w, data)
}
