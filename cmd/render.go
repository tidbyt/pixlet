package cmd

import (
	"fmt"
	"image"
	"io/ioutil"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"go.starlark.net/starlark"

	"tidbyt.dev/pixlet/encode"
	"tidbyt.dev/pixlet/runtime"
)

var (
	output        string
	magnify       int
	renderGif     bool
	silenceOutput bool
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
}

var RenderCmd = &cobra.Command{
	Use:   "render [script] [<key>=value>]...",
	Short: "Run a Pixlet script with provided config parameters",
	Args:  cobra.MinimumNArgs(1),
	RunE:  render,
}

func render(cmd *cobra.Command, args []string) error {
	script := args[0]

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
		if len(split) != 2 {
			return fmt.Errorf("parameters must be on form <key>=<value>, found %s", param)
		}
		config[split[0]] = split[1]
	}

	src, err := ioutil.ReadFile(script)
	if err != nil {
		return fmt.Errorf("failed to read file %s: %w", script, err)
	}

	runtime.InitCache(runtime.NewInMemoryCache())

	applet := runtime.Applet{}
	err = applet.Load(script, src, nil)
	if err != nil {
		return fmt.Errorf("failed to load applet: %w", err)
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

	var buf []byte

	if renderGif {
		buf, err = screens.EncodeGIF(filter)
	} else {
		buf, err = screens.EncodeWebP(filter)
	}
	if err != nil {
		return fmt.Errorf("error rendering: %w", err)
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
