package main

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
	"github.com/spf13/cobra"

	"tidbyt.dev/pixlet/runtime"
)

var (
	port  int
	watch bool
)

func init() {
	rootCmd.AddCommand(serveCmd)
	serveCmd.Flags().IntVarP(&port, "port", "p", 8080, "Port for serving rendered images")
	serveCmd.Flags().BoolVarP(&watch, "watch", "w", false, "Reload scripts on change")
}

var serveCmd = &cobra.Command{
	Use:   "serve [script]",
	Short: "Serves a starlark render script over HTTP.",
	Args:  cobra.ExactArgs(1),
	Run:   serve,
}

func loadScript(applet *runtime.Applet, filename string) error {
	src, err := ioutil.ReadFile(filename)
	if err != nil {
		return fmt.Errorf("failed to read file %s: %v\n", filename, err)
	}

	runtime.InitCache(runtime.NewInMemoryCache())

	err = applet.Load(filename, src, nil)
	if err != nil {
		return fmt.Errorf("failed to load applet: %v\n", err)
	}

	return nil
}

func serve(cmd *cobra.Command, args []string) {
	applet := runtime.Applet{}
	err := loadScript(&applet, args[0])
	if err != nil {
		fmt.Printf("error: %v\n", err)
		os.Exit(1)
	}

	mutex := sync.RWMutex{}

	// if --watch/-w: monitor script and reload on file change
	if watch {
		watcher, err := fsnotify.NewWatcher()
		if err != nil {
			fmt.Printf("error watching for changes: %v\n", err)
			os.Exit(1)
		}

		watcher.Add(args[0])

		go func() {
			for {
				event, ok := <-watcher.Events
				if !ok {
					break
				}

				if (event.Op & fsnotify.Write) != 0 {
					mutex.Lock()
					err := loadScript(&applet, args[0])
					mutex.Unlock()
					if err != nil {
						fmt.Printf("Error on reload: %v\n", err)
					} else {
						fmt.Printf("Reload\n")
					}
				}
			}
		}()

		fmt.Printf("Loaded %s, watching for changes\n", args[0])
	}

	http.HandleFunc("/favicon.ico", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(404)
	})

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		config := make(map[string]string)
		for k, vals := range r.URL.Query() {
			config[k] = vals[0]
		}

		mutex.RLock()
		defer mutex.RUnlock()

		screens, err := applet.Run(config)
		if err != nil {
			log.Printf("Error running script: %s\n", err)
			return
		}

		webp, err := screens.RenderWebP()
		if err != nil {
			fmt.Printf("Error rendering: %s\n", err)
			return
		}

		w.Header().Set("Content-Type", "text/html")
		writePreviewHTML(w, webp)
	})
	fmt.Printf("listening on tcp/%d\n", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf("127.0.0.1:%d", port), nil))
}

func writePreviewHTML(w io.Writer, webp []byte) {
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

	fmt.Fprintln(w, `</div></body></html>`)
}
