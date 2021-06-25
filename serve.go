package main

import (
	"log"

	"github.com/spf13/cobra"

	"tidbyt.dev/pixlet/server"
)

var (
	host  string
	port  int
	watch bool
)

func init() {
	rootCmd.AddCommand(serveCmd)
	serveCmd.Flags().StringVarP(&host, "host", "i", "127.0.0.1", "Host interface for serving rendered images")
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
	s := server.NewServer(host, port, watch, args[0])
	log.Fatal(s.Run())
}
