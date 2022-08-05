package cmd

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"go.starlark.net/starlark"

	"tidbyt.dev/pixlet/runtime"
)

var (
	profileOutput string
)

func init() {
	ProfileCmd.Flags().StringVarP(&profileOutput, "profile", "p", "", "Path for gzipped pprof output")
}

var ProfileCmd = &cobra.Command{
	Use:   "profile [script] [<key>=value>]...",
	Short: "Runs script with provided config parameters and prints execution-time profile. See https://github.com/google/pprof for how to read the output.",
	Args:  cobra.MinimumNArgs(1),
	Run:   profile,
}

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

	outPath := strings.TrimSuffix(script, ".star") + "_pprof.gz"
	if profileOutput != "" {
		outPath = profileOutput
	}
	prof, err := os.Create(outPath)
	if err != nil {
		fmt.Printf("Cannot create temp file for output: %s", err)
		os.Exit(1)
	}
	defer prof.Close()

	if err = starlark.StartProfile(prof); err != nil {
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
	prof.Sync()
}
