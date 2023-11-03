package main

// Renders webp images from example snippets for the widget
// documentation.

import (
	"fmt"
	"image"
	"io/ioutil"
	"strings"

	"tidbyt.dev/pixlet/encode"
	"tidbyt.dev/pixlet/runtime"
)

const Magnification = 7

func Magnify(input image.Image) (image.Image, error) {
	in, ok := input.(*image.RGBA)
	if !ok {
		panic("not RGBA")
	}

	out := image.NewRGBA(image.Rect(0, 0, in.Bounds().Dx()*Magnification, in.Bounds().Dy()*Magnification))
	for x := 0; x < in.Bounds().Dx(); x++ {
		for y := 0; y < in.Bounds().Dy(); y++ {
			for xx := 0; xx < 10; xx++ {
				for yy := 0; yy < 10; yy++ {
					out.SetRGBA(
						x*Magnification+xx,
						y*Magnification+yy,
						in.RGBAAt(x, y),
					)
				}
			}
		}
	}

	return out, nil
}

func main() {
	files, err := ioutil.ReadDir(".")
	if err != nil {
		panic(err)
	}

	examples := map[string]string{}
	for _, f := range files {
		if !strings.HasSuffix(f.Name(), ".star") {
			continue
		}

		content, err := ioutil.ReadFile(f.Name())
		if err != nil {
			panic(err)
		}

		examples[strings.TrimSuffix(f.Name(), ".star")] = string(content)
	}

	for name, snippet := range examples {
		src := fmt.Sprintf(`
load("render.star", "render")
def main():
    w = %s
    return render.Root(child=w)
`, snippet)

		app := runtime.Applet{}
		err = app.Load(fmt.Sprintf("id-%s", name), name, []byte(src), nil)
		if err != nil {
			panic(err)
		}

		roots, err := app.Run(nil)
		if err != nil {
			panic(err)
		}

		gif, err := encode.ScreensFromRoots(roots).EncodeGIF(15000, Magnify)
		if err != nil {
			panic(err)
		}

		err = ioutil.WriteFile(fmt.Sprintf("img/widget_%s.gif", name), gif, 0644)
		if err != nil {
			panic(err)
		}
	}
}
