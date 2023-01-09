/*
Copyright 2016 Google LLC

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    https://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package cmd

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"

	"github.com/bazelbuild/buildtools/build"
	"github.com/bazelbuild/buildtools/buildifier/utils"
	"github.com/bazelbuild/buildtools/differ"
	"github.com/bazelbuild/buildtools/warn"
	"github.com/bazelbuild/buildtools/wspace"
)

var (
	vflag        bool
	rflag        bool
	dryRunFlag   bool
	fixFlag      bool
	outputFormat string
)

func runBuildifier(args []string, lint string, mode string, format string, recursive bool, verbose bool) int {
	tf := &utils.TempFile{}
	defer tf.Clean()

	exitCode := 0
	var diagnostics *utils.Diagnostics
	if len(args) == 0 || (len(args) == 1 && (args)[0] == "-") {
		// Read from stdin, write to stdout.
		data, err := io.ReadAll(os.Stdin)
		if err != nil {
			fmt.Fprintf(os.Stderr, "buildifier: reading stdin: %v\n", err)
			return 2
		}
		if mode == "fix" {
			mode = "pipe"
		}
		var fileDiagnostics *utils.FileDiagnostics
		fileDiagnostics, exitCode = processFile("", data, lint, false, tf, mode, verbose)
		diagnostics = utils.NewDiagnostics(fileDiagnostics)
	} else {
		files := args
		if recursive {
			var err error
			files, err = utils.ExpandDirectories(&args)
			if err != nil {
				fmt.Fprintf(os.Stderr, "buildifier: %v\n", err)
				return 3
			}
		}
		diagnostics, exitCode = processFiles(files, lint, tf, mode, verbose)
	}

	diagnosticsOutput := diagnostics.Format(format, verbose)
	if format != "" {
		// Explicitly provided --format means the diagnostics are printed to stdout
		fmt.Printf(diagnosticsOutput)
		// Exit code should be set to 0 so that other tools know they can safely parse the json
		exitCode = 0
	} else {
		// --format is not provided, stdout is reserved for file contents
		fmt.Fprint(os.Stderr, diagnosticsOutput)
	}

	if err := diff.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		return 2
	}

	return exitCode
}

func processFiles(files []string, lint string, tf *utils.TempFile, mode string, verbose bool) (*utils.Diagnostics, int) {
	// Decide how many file reads to run in parallel.
	// At most 100, and at most one per 10 input files.
	nworker := 100
	if n := (len(files) + 9) / 10; nworker > n {
		nworker = n
	}
	runtime.GOMAXPROCS(nworker + 1)

	// Start nworker workers reading stripes of the input
	// argument list and sending the resulting data on
	// separate channels. file[k] is read by worker k%nworker
	// and delivered on ch[k%nworker].
	type result struct {
		file string
		data []byte
		err  error
	}

	ch := make([]chan result, nworker)
	for i := 0; i < nworker; i++ {
		ch[i] = make(chan result, 1)
		go func(i int) {
			for j := i; j < len(files); j += nworker {
				file := files[j]
				data, err := os.ReadFile(file)
				ch[i] <- result{file, data, err}
			}
		}(i)
	}

	exitCode := 0
	fileDiagnostics := []*utils.FileDiagnostics{}

	// Process files. The processing still runs in a single goroutine
	// in sequence. Only the reading of the files has been parallelized.
	// The goal is to optimize for runs where most files are already
	// formatted correctly, so that reading is the bulk of the I/O.
	for i, file := range files {
		res := <-ch[i%nworker]
		if res.file != file {
			fmt.Fprintf(os.Stderr, "buildifier: internal phase error: got %s for %s", res.file, file)
			os.Exit(3)
		}
		if res.err != nil {
			fmt.Fprintf(os.Stderr, "buildifier: %v\n", res.err)
			exitCode = 3
			continue
		}
		fd, newExitCode := processFile(file, res.data, lint, len(files) > 1, tf, mode, verbose)
		if fd != nil {
			fileDiagnostics = append(fileDiagnostics, fd)
		}
		if newExitCode != 0 {
			exitCode = newExitCode
		}
	}
	return utils.NewDiagnostics(fileDiagnostics...), exitCode
}

// diff is the differ to use when *mode == "diff".
var diff *differ.Differ

func defaultWarnings() []string {
	warnings := []string{}
	for _, warning := range warn.AllWarnings {
		if !disabledWarnings[warning] {
			warnings = append(warnings, warning)
		}
	}
	return warnings
}

var disabledWarnings = map[string]bool{
	"function-docstring":        true, // disables docstring warnings
	"function-docstring-header": true, // disables docstring warnings
	"function-docstring-args":   true, // disables docstring warnings
	"function-docstring-return": true, // disables docstring warnings
	"native-android":            true, // disables native android rules
	"native-cc":                 true, // disables native cc rules
	"native-java":               true, // disables native java rules
	"native-proto":              true, // disables native proto rules
	"native-py":                 true, // disables native python rules
}

// processFile processes a single file containing data.
// It has been read from filename and should be written back if fixing.
func processFile(filename string, data []byte, lint string, displayFileNames bool, tf *utils.TempFile, mode string, verbose bool) (*utils.FileDiagnostics, int) {
	var exitCode int

	displayFilename := filename
	parser := utils.GetParser("auto")

	f, err := parser(displayFilename, data)
	if err != nil {
		// Do not use buildifier: prefix on this error.
		// Since it is a parse error, it begins with file:line:
		// and we want that to be the first thing in the error.
		fmt.Fprintf(os.Stderr, "%v\n", err)
		if exitCode < 1 {
			exitCode = 1
		}
		return utils.InvalidFileDiagnostics(displayFilename), exitCode
	}

	if absoluteFilename, err := filepath.Abs(displayFilename); err == nil {
		f.WorkspaceRoot, f.Pkg, f.Label = wspace.SplitFilePath(absoluteFilename)
	}

	enabledWarnings := defaultWarnings()
	warnings := utils.Lint(f, lint, &enabledWarnings, verbose)
	if len(warnings) > 0 {
		exitCode = 4
	}
	fileDiagnostics := utils.NewFileDiagnostics(f.DisplayPath(), warnings)

	ndata := build.Format(f)

	switch mode {
	case "check":
		// check mode: print names of files that need formatting.
		if !bytes.Equal(data, ndata) {
			fileDiagnostics.Formatted = false
			return fileDiagnostics, 4
		}

	case "diff":
		// diff mode: run diff on old and new.
		if bytes.Equal(data, ndata) {
			return fileDiagnostics, exitCode
		}
		outfile, err := tf.WriteTemp(ndata)
		if err != nil {
			fmt.Fprintf(os.Stderr, "buildifier: %v\n", err)
			return fileDiagnostics, 3
		}
		infile := filename
		if filename == "" {
			// data was read from standard filename.
			// Write it to a temporary file so diff can read it.
			infile, err = tf.WriteTemp(data)
			if err != nil {
				fmt.Fprintf(os.Stderr, "buildifier: %v\n", err)
				return fileDiagnostics, 3
			}
		}
		if displayFileNames {
			fmt.Fprintf(os.Stderr, "%v:\n", f.DisplayPath())
		}
		if err := diff.Show(infile, outfile); err != nil {
			fmt.Fprintf(os.Stderr, "%v\n", err)
			return fileDiagnostics, 4
		}

	case "pipe":
		// pipe mode - reading from stdin, writing to stdout.
		// ("pipe" is not from the command line; it is set above in main.)
		os.Stdout.Write(ndata)

	case "fix":
		// fix mode: update files in place as needed.
		if bytes.Equal(data, ndata) {
			return fileDiagnostics, exitCode
		}

		err := os.WriteFile(filename, ndata, 0666)
		if err != nil {
			fmt.Fprintf(os.Stderr, "buildifier: %s\n", err)
			return fileDiagnostics, 3
		}

		if verbose {
			fmt.Fprintf(os.Stderr, "fixed %s\n", f.DisplayPath())
		}
	case "print_if_changed":
		if bytes.Equal(data, ndata) {
			return fileDiagnostics, exitCode
		}

		if _, err := os.Stdout.Write(ndata); err != nil {
			fmt.Fprintf(os.Stderr, "buildifier: error writing output: %v\n", err)
			return fileDiagnostics, 3
		}
	}
	return fileDiagnostics, exitCode
}
