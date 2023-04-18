package cmd

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"go.starlark.net/starlark"

	"tidbyt.dev/pixlet/globals"
	"tidbyt.dev/pixlet/runtime"
	"tidbyt.dev/pixlet/runtime/modules/testing"
)

var (
	allowHttp bool
)

func init() {
	TestCmd.Flags().BoolVar(&allowHttp, "allow-http", false, "Do not error on unstubbed http requests")
	TestCmd.Flags().BoolVarP(&silenceOutput, "silent", "", false, "Silence print statements when rendering app")
	TestCmd.Flags().BoolVarP(&vflag, "verbose", "v", false, "print verbose information about test run")
}

var TestCmd = &cobra.Command{
	Use:   "test [script] [<key>=value>]...",
	Short: "Run tests for a Pixlet script with provided config parameters",
	Args:  cobra.MinimumNArgs(1),
	RunE:  test,
}

func test(cmd *cobra.Command, args []string) error {
	script := args[0]

	globals.Width = width
	globals.Height = height

	if !strings.HasSuffix(script, ".star") {
		return fmt.Errorf("script file must have suffix .star: %s", script)
	}

	testScript := strings.TrimSuffix(script, ".star") + ".test.star"

	if _, err := os.Stat(testScript); err != nil {
		return fmt.Errorf("test file %s not found", testScript)
	}

	config := map[string]string{}
	for _, param := range args[1:] {
		split := strings.Split(param, "=")
		if len(split) != 2 {
			return fmt.Errorf("parameters must be on form <key>=<value>, found %s", param)
		}
		config[split[0]] = split[1]
	}

	src, err := ioutil.ReadFile(script)
	if err != nil {
		return fmt.Errorf("failed to read file %s: %w", script, err)
	}

	testSrc, err := ioutil.ReadFile(testScript)
	if err != nil {
		return fmt.Errorf("failed to read file %s: %w", testScript, err)
	}

	initializers := []runtime.ThreadInitializer{}

	runtime.InitCache(runtime.NewInMemoryCache())

	stubbedHttpTransport := testing.InitHttpStub(allowHttp)
	initializers = append(initializers, func(thread *starlark.Thread) *starlark.Thread {
		stubbedHttpTransport.ClearStubs()
		return thread
	})

	applet := runtime.Applet{}
	err = applet.Load(script, src, runtime.TestingLoader)
	if err != nil {
		return fmt.Errorf("failed to load applet: %w", err)
	}

	err = applet.LoadTests(testScript, testSrc, nil)
	if err != nil {
		return fmt.Errorf("failed to load tests: %w", err)
	}

	// Remove the print function from the starlark thread if the silent flag is
	// passed.
	if silenceOutput {
		initializers = append(initializers, func(thread *starlark.Thread) *starlark.Thread {
			thread.Print = func(thread *starlark.Thread, msg string) {}
			return thread
		})
	}

	initializers = append(initializers, func(thread *starlark.Thread) *starlark.Thread {
		thread.SetLocal("scriptPath", script)
		return thread
	})

	testResult, err := applet.RunTests(config, initializers...)
	if err != nil {
		return fmt.Errorf("error running tests for script: %w", err)
	}

	testFailures := 0
	for _, res := range testResult {
		if res.Success() {
			if vflag {
				fmt.Print(res.FunctionName, ": ")
				fmt.Println("SUCCESS")
			}
		} else {
			testFailures += 1
			fmt.Print(res.FunctionName, ": ")
			fmt.Println("FAILURE")
		}
		for _, err := range res.Errors {
			for _, line := range strings.Split(err.String(), "\n") {
				fmt.Println("    ", line)
			}
			fmt.Println("")
		}
	}

	if testFailures > 0 {
		return fmt.Errorf("%v test failures.", testFailures)
	}

	return nil
}
