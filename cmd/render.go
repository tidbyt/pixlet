package cmd

import (
	"context"
	"fmt"
	"image"
	"io/ioutil"
	"log"
	"os"
	"runtime/pprof"
	"strconv"
	"strings"
	"time"

	"github.com/spf13/cobra"
	"go.starlark.net/starlark"

	"tidbyt.dev/pixlet/encode"
	pixletrender "tidbyt.dev/pixlet/render"
	"tidbyt.dev/pixlet/runtime"
	"tidbyt.dev/pixlet/starlarkutil"
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
	Use:   "render [script] [<key>=value>]...",
	Short: "Run a Pixlet script with provided config parameters",
	Args:  cobra.MinimumNArgs(1),
	RunE:  render,
}

func render(cmd *cobra.Command, args []string) error {
	f, err := os.Create("profile.out")
	if err != nil {
		log.Fatal(err)
	}
	pprof.StartCPUProfile(f)
	defer pprof.StopCPUProfile()

	return doTheThing(args)
}

func doTheThing(args []string) error {
	script := args[0]

	pixletrender.FrameWidth = width
	pixletrender.FrameHeight = height

	if !strings.HasSuffix(script, ".star") {
		return fmt.Errorf("script file must have suffix .star: %s", script)
	}

	outPath := strings.TrimSuffix(script, ".star")
	if renderGif {
		outPath += ".gif"
	} else {
		outPath += ".webp"
	}
	if output != "" {
		outPath = output
	}

	config := map[string]string{}
	for _, param := range args[1:] {
		split := strings.Split(param, "=")
		if len(split) < 2 {
			return fmt.Errorf("parameters must be on form <key>=<value>, found %s", param)
		}
		config[split[0]] = strings.Join(split[1:], "=")
	}

	src, err := ioutil.ReadFile(script)
	if err != nil {
		return fmt.Errorf("failed to read file %s: %w", script, err)
	}

	// Remove the print function from the starlark thread if the silent flag is
	// passed.
	initializers := []runtime.ThreadInitializer{}
	if silenceOutput {
		initializers = append(initializers, func(thread *starlark.Thread) *starlark.Thread {
			thread.Print = func(thread *starlark.Thread, msg string) {}
			return thread
		})
	}

	// Timeout?
	if timeout > 0 {
		initializers = append(initializers, func(thread *starlark.Thread) *starlark.Thread {
			ctx, _ := context.WithTimeoutCause(
				context.Background(),
				time.Duration(timeout)*time.Millisecond,
				fmt.Errorf("timeout after %dms", timeout),
			)
			starlarkutil.AttachThreadContext(ctx, thread)
			return thread
		})
	}

	cache := runtime.NewInMemoryCache()
	runtime.InitHTTP(cache)
	runtime.InitCache(cache)

	applet := runtime.Applet{}
	err = applet.LoadWithInitializers("", script, src, nil, initializers...)
	if err != nil {
		return fmt.Errorf("failed to load applet: %w", err)
	}

	var buf []byte
	iterations := os.Getenv("ITERATIONS")
	if iterations == "" {
		iterations = "1"
	}
	n, _ := strconv.Atoi(iterations)

	for i := 0; i < n; i++ {
		roots, err := applet.Run(config, initializers...)
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
	}

	if outPath == "-" {
		_, err = os.Stdout.Write(buf)
	} else {
		err = ioutil.WriteFile(outPath, buf, 0644)
	}

	if err != nil {
		return fmt.Errorf("writing %s: %s", outPath, err)
	}

	return nil
}
