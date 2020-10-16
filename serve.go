package main

import (
	"encoding/base64"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/spf13/cobra"

	"github.com/tidbyt/pixlet/runtime"
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

func serve(cmd *cobra.Command, args []string) {
	filename := args[0]

	src, err := ioutil.ReadFile(filename)
	if err != nil {
		fmt.Printf("failed to read file %s: %v\n", filename, err)
		return
	}

	runtime.InitCache(runtime.NewInMemoryCache())

	applet := runtime.Applet{}
	err = applet.Load(filename, src, nil)
	if err != nil {
		fmt.Printf("failed to load applet: %v\n", err)
		return
	}

	http.HandleFunc("/favicon.ico", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(404)
	})

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		config := make(map[string]string)
		for k, vals := range r.URL.Query() {
			config[k] = vals[0]
		}

		screens, err := applet.Run(config)
		if err != nil {
			log.Printf("Error running script: %s\n", err)
			return
		}

		webp, err := screens.RenderWebP()
		if err != nil {
			fmt.Println("Error rendering: %s\n", err)
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
