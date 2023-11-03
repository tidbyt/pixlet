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
	Short: "Run a Pixlet script and print its execution-time profile",
	Args:  cobra.MinimumNArgs(1),
	RunE:  profile,
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

func profile(cmd *cobra.Command, args []string) error {
	script := args[0]

	if !strings.HasSuffix(script, ".star") {
		return fmt.Errorf("script file must have suffix .star: %s", script)
	}

	config := map[string]string{}
	for _, param := range args[1:] {
		split := strings.Split(param, "=")
		if len(split) != 2 {
			return fmt.Errorf("parameters must be on form <key>=<value>, found %s", param)
		}
		config[split[0]] = split[1]
	}

	profile, err := ProfileApp(script, config)
	if err != nil {
		return err
	}

	options := &pprof_driver.Options{
		Fetch: MakeFetchFunc(profile),
		UI:    printUI{},
	}
	if err = pprof_driver.PProf(options); err != nil {
		return fmt.Errorf("could not start pprof driver: %w", err)
	}

	return nil
}

func ProfileApp(script string, config map[string]string) (*pprof_profile.Profile, error) {
	src, err := ioutil.ReadFile(script)
	if err != nil {
		return nil, fmt.Errorf("failed to read file %s: %w", script, err)
	}

	cache := runtime.NewInMemoryCache()
	runtime.InitHTTP(cache)
	runtime.InitCache(cache)

	// Remove the print function from the starlark thread.
	initializers := []runtime.ThreadInitializer{}
	initializers = append(initializers, func(thread *starlark.Thread) *starlark.Thread {
		thread.Print = func(thread *starlark.Thread, msg string) {}
		return thread
	})

	applet := runtime.Applet{}
	err = applet.LoadWithInitializers("", script, src, nil, initializers...)
	if err != nil {
		return nil, fmt.Errorf("failed to load applet: %w", err)
	}

	buf := new(bytes.Buffer)
	if err = starlark.StartProfile(buf); err != nil {
		return nil, fmt.Errorf("error starting profiler: %w", err)
	}

	_, err = applet.Run(config, initializers...)
	if err != nil {
		_ = starlark.StopProfile()
		return nil, fmt.Errorf("error running script: %w", err)
	}

	if err = starlark.StopProfile(); err != nil {
		return nil, fmt.Errorf("error stopping profiler: %w", err)
	}

	profile, err := pprof_profile.ParseData(buf.Bytes())
	if err != nil {
		return nil, fmt.Errorf("could not parse pprof profile: %w", err)
	}

	return profile, nil
}
