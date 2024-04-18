package cmd

import (
	"context"
	"fmt"
	"image"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"testing/fstest"
	"time"

	"github.com/spf13/cobra"

	"tidbyt.dev/pixlet/encode"
	"tidbyt.dev/pixlet/globals"
	"tidbyt.dev/pixlet/runtime"
)

var (
	output        string
	magnify       int
	renderGif     bool
	maxDuration   int
	silenceOutput bool
	width         int
	height        int
	timeout       int
)

func init() {
	RenderCmd.Flags().StringVarP(&output, "output", "o", "", "Path for rendered image")
	RenderCmd.Flags().BoolVarP(&renderGif, "gif", "", false, "Generate GIF instead of WebP")
	RenderCmd.Flags().BoolVarP(&silenceOutput, "silent", "", false, "Silence print statements when rendering app")
	RenderCmd.Flags().IntVarP(
		&magnify,
		"magnify",
		"m",
		1,
		"Increase image dimension by a factor (useful for debugging)",
	)
	RenderCmd.Flags().IntVarP(
		&width,
		"width",
		"w",
		64,
		"Set width",
	)
	RenderCmd.Flags().IntVarP(
		&height,
		"height",
		"t",
		32,
		"Set height",
	)
	RenderCmd.Flags().IntVarP(
		&maxDuration,
		"max_duration",
		"d",
		15000,
		"Maximum allowed animation duration (ms)",
	)
	RenderCmd.Flags().IntVarP(
		&timeout,
		"timeout",
		"",
		30000,
		"Timeout for execution (ms)",
	)
}

var RenderCmd = &cobra.Command{
	Use:   "render [path] [<key>=value>]...",
	Short: "Run a Pixlet program with provided config parameters",
	Args:  cobra.MinimumNArgs(1),
	RunE:  render,
	Long: `Render a Pixlet program with provided config parameters.

The path argument should be the path to the Pixlet program to run. The
program can be a single file with the .star extension, or a directory
containing multiple Starlark files and resources.
	`,
}

func render(cmd *cobra.Command, args []string) error {
	path := args[0]

	// check if path exists, and whether it is a directory or a file
	info, err := os.Stat(path)
	if err != nil {
		return fmt.Errorf("failed to stat %s: %w", path, err)
	}

	var fs fs.FS
	var outPath string
	if info.IsDir() {
		fs = os.DirFS(path)
		outPath = filepath.Join(path, filepath.Base(path))
	} else {
		if !strings.HasSuffix(path, ".star") {
			return fmt.Errorf("script file must have suffix .star: %s", path)
		}

		src, err := os.ReadFile(path)
		if err != nil {
			return fmt.Errorf("failed to read file %s: %w", path, err)
		}

		fs = fstest.MapFS{
			filepath.Base(path): {Data: src},
		}

		outPath = strings.TrimSuffix(path, ".star")
	}

	if renderGif {
		outPath += ".gif"
	} else {
		outPath += ".webp"
	}
	if output != "" {
		outPath = output
	}

	globals.Width = width
	globals.Height = height

	config := map[string]string{}
	for _, param := range args[1:] {
		split := strings.Split(param, "=")
		if len(split) < 2 {
			return fmt.Errorf("parameters must be on form <key>=<value>, found %s", param)
		}
		config[split[0]] = strings.Join(split[1:], "=")
	}

	// Remove the print function from the starlark thread if the silent flag is
	// passed.
	var opts []runtime.AppletOption
	if silenceOutput {
		opts = append(opts, runtime.WithPrintDisabled())
	}

	ctx := context.Background()
	if timeout > 0 {
		ctx, _ = context.WithTimeoutCause(
			ctx,
			time.Duration(timeout)*time.Millisecond,
			fmt.Errorf("timeout after %dms", timeout),
		)
	}

	cache := runtime.NewInMemoryCache()
	runtime.InitHTTP(cache)
	runtime.InitCache(cache)

	applet, err := runtime.NewAppletFromFS(filepath.Base(path), fs, opts...)
	if err != nil {
		return fmt.Errorf("failed to load applet: %w", err)
	}

	roots, err := applet.RunWithConfig(ctx, config)
	if err != nil {
		return fmt.Errorf("error running script: %w", err)
	}
	screens := encode.ScreensFromRoots(roots)

	filter := func(input image.Image) (image.Image, error) {
		if magnify <= 1 {
			return input, nil
		}
		in, ok := input.(*image.RGBA)
		if !ok {
			return nil, fmt.Errorf("image not RGBA, very weird")
		}

		out := image.NewRGBA(
			image.Rect(
				0, 0,
				in.Bounds().Dx()*magnify,
				in.Bounds().Dy()*magnify),
		)
		for x := 0; x < in.Bounds().Dx(); x++ {
			for y := 0; y < in.Bounds().Dy(); y++ {
				for xx := 0; xx < magnify; xx++ {
					for yy := 0; yy < magnify; yy++ {
						out.SetRGBA(
							x*magnify+xx,
							y*magnify+yy,
							in.RGBAAt(x, y),
						)
					}
				}
			}
		}

		return out, nil
	}

	var buf []byte

	if screens.ShowFullAnimation {
		maxDuration = 0
	}

	if renderGif {
		buf, err = screens.EncodeGIF(maxDuration, filter)
	} else {
		buf, err = screens.EncodeWebP(maxDuration, filter)
	}
	if err != nil {
		return fmt.Errorf("error rendering: %w", err)
	}

	if outPath == "-" {
		_, err = os.Stdout.Write(buf)
	} else {
		err = os.WriteFile(outPath, buf, 0644)
	}

	if err != nil {
		return fmt.Errorf("writing %s: %s", outPath, err)
	}

	return nil
}
