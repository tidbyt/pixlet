package cmd

import (
	"fmt"
	"image"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"tidbyt.dev/pixlet/encode"
	"tidbyt.dev/pixlet/runtime"
)

var (
	output    string
	magnify   int
	renderGif bool
)

func init() {
	RenderCmd.Flags().StringVarP(&output, "output", "o", "", "Path for rendered image")
	RenderCmd.Flags().BoolVarP(&renderGif, "gif", "", false, "Generate GIF instead of WebP")
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
	Short: "Runs script with provided config parameters.",
	Args:  cobra.MinimumNArgs(1),
	Run:   render,
}

func render(cmd *cobra.Command, args []string) {
	script := args[0]

	if !strings.HasSuffix(script, ".star") {
		fmt.Printf("script file must have suffix .star: %s\n", script)
		os.Exit(1)
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

	roots, err := applet.Run(config)
	if err != nil {
		log.Printf("Error running script: %s\n", err)
		os.Exit(1)
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
				for xx := 0; xx < 10; xx++ {
					for yy := 0; yy < 10; yy++ {
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
		fmt.Printf("Error rendering: %s\n", err)
		os.Exit(1)
	}

	if outPath == "-" {
		_, err = os.Stdout.Write(buf)
	} else {
		err = ioutil.WriteFile(outPath, buf, 0644)
	}

	if err != nil {
		fmt.Printf("Writing %s: %s", outPath, err)
		os.Exit(1)
	}
}
