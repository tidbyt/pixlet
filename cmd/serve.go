package cmd

import (
	"github.com/spf13/cobra"

	"tidbyt.dev/pixlet/server"
)

var (
	host  string
	port  int
	watch bool
	serveGif bool
)

func init() {
	ServeCmd.Flags().StringVarP(&host, "host", "i", "127.0.0.1", "Host interface for serving rendered images")
	ServeCmd.Flags().IntVarP(&port, "port", "p", 8080, "Port for serving rendered images")
	ServeCmd.Flags().BoolVarP(&watch, "watch", "w", true, "Reload scripts on change. Does not recurse sub-directories.")
	ServeCmd.Flags().IntVarP(&maxDuration, "max_duration", "d", 15000, "Maximum allowed animation duration (ms)")
	ServeCmd.Flags().IntVarP(&timeout, "timeout", "", 30000, "Timeout for execution (ms)")
	ServeCmd.Flags().BoolVarP(&serveGif, "gif", "", false, "Generate GIF instead of WebP")
}

var ServeCmd = &cobra.Command{
	Use:   "serve [path]",
	Short: "Serve a Pixlet app in a web server",
	Args:  cobra.ExactArgs(1),
	RunE:  serve,
	Long: `Serve a Pixlet app in a web server.

The path argument should be the path to the Pixlet program to run. The
program can be a single file with the .star extension, or a directory
containing multiple Starlark files and resources.`,
}

func serve(cmd *cobra.Command, args []string) error {
	s, err := server.NewServer(host, port, watch, args[0], maxDuration, timeout, serveGif)
	if err != nil {
		return err
	}
	return s.Run()
}
