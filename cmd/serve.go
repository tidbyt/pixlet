package cmd

import (
	"github.com/spf13/cobra"

	"tidbyt.dev/pixlet/server"
)

var (
	host  string
	port  int
	watch bool
)

func init() {
	ServeCmd.Flags().StringVarP(&host, "host", "i", "127.0.0.1", "Host interface for serving rendered images")
	ServeCmd.Flags().IntVarP(&port, "port", "p", 8080, "Port for serving rendered images")
	ServeCmd.Flags().BoolVarP(&watch, "watch", "w", false, "Reload scripts on change")
}

var ServeCmd = &cobra.Command{
	Use:   "serve [script]",
	Short: "Serves a starlark render script over HTTP.",
	Args:  cobra.ExactArgs(1),
	RunE:  serve,
}

func serve(cmd *cobra.Command, args []string) error {
	s, err := server.NewServer(host, port, watch, args[0])
	if err != nil {
		return err
	}
	return s.Run()
}
