package cmd

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"time"

	pprof_driver "github.com/google/pprof/driver"
	pprof_profile "github.com/google/pprof/profile"
	"github.com/spf13/cobra"
	"go.starlark.net/starlark"

	"tidbyt.dev/pixlet/runtime"
)

var (
	pprof_cmd string
)

func init() {
	ProfileCmd.Flags().StringVarP(
		&pprof_cmd, "pprof", "", "top 10", "Command to call pprof with",
	)
}

var ProfileCmd = &cobra.Command{
	Use:   "profile [script] [<key>=value>]...",
	Short: "Runs script with provided config parameters and prints execution-time profile.",
	Args:  cobra.MinimumNArgs(1),
	Run:   profile,
}

// We save the profile into an in-memory buffer, which is simpler than the tool expects.
// Simple adapter to pipe it through.
type FetchFunc func(src string, duration, timeout time.Duration) (*pprof_profile.Profile, string, error)

func (f FetchFunc) Fetch(src string, duration, timeout time.Duration) (*pprof_profile.Profile, string, error) {
	return f(src, duration, timeout)
}
func MakeFetchFunc(prof *pprof_profile.Profile) FetchFunc {
	return func(src string, duration, timeout time.Duration) (*pprof_profile.Profile, string, error) {
		return prof, "", nil
	}
}

// Calls the pprof program to print the top users of CPU, then exit
type printUI struct{}

var pprof_printed = false

func (u printUI) ReadLine(prompt string) (string, error) {
	if pprof_printed {
		os.Exit(0)
	}
	pprof_printed = true
	return pprof_cmd, nil
}
func (u printUI) Print(args ...interface{})                    {}
func (u printUI) PrintErr(args ...interface{})                 {}
func (u printUI) IsTerminal() bool                             { return false }
func (u printUI) WantBrowser() bool                            { return false }
func (u printUI) SetAutoComplete(complete func(string) string) {}

func profile(cmd *cobra.Command, args []string) {
	script := args[0]

	if !strings.HasSuffix(script, ".star") {
		fmt.Printf("script file must have suffix .star: %s\n", script)
		os.Exit(1)
	}

	config := map[string]string{}
	for _, param := range args[1:] {
		split := strings.Split(param, "=")
		if len(split) != 2 {
			fmt.Printf("parameters must be on form <key>=<value>, found %s\n", param)
			os.Exit(1)
		}
		config[split[0]] = split[1]
	}

	src, err := ioutil.ReadFile(script)
	if err != nil {
		fmt.Printf("failed to read file %s: %v\n", script, err)
		os.Exit(1)
	}

	runtime.InitCache(runtime.NewInMemoryCache())

	applet := runtime.Applet{}
	err = applet.Load(script, src, nil)
	if err != nil {
		fmt.Printf("failed to load applet: %v\n", err)
		os.Exit(1)
	}

	buf := new(bytes.Buffer)
	if err = starlark.StartProfile(buf); err != nil {
		fmt.Printf("Error starting profiler: %s\n", err)
		os.Exit(1)
	}

	_, err = applet.Run(config)
	if err != nil {
		_ = starlark.StopProfile()
		fmt.Printf("Error running script: %s\n", err)
		os.Exit(1)
	}

	if err = starlark.StopProfile(); err != nil {
		fmt.Printf("Error stopping profiler: %s\n", err)
		os.Exit(1)
	}

	profile, err := pprof_profile.ParseData(buf.Bytes())
	if err != nil {
		fmt.Printf("Could not parse pprof profile: %s\n", err)
		os.Exit(1)
	}

	options := &pprof_driver.Options{
		Fetch: MakeFetchFunc(profile),
		UI:    printUI{},
	}
	if err = pprof_driver.PProf(options); err != nil {
		fmt.Printf("Could not start pprof driver: %s\n", err)
		os.Exit(1)
	}
}
