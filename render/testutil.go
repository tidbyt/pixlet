package render

import (
	"fmt"
	"image"
	"image/color"
	"strings"

	"tidbyt.dev/pixlet/render/canvas"
)

var DefaultPalette = map[string]color.RGBA{
	"r": {0xff, 0, 0, 0xff},
	"g": {0, 0xff, 0, 0xff},
	"b": {0, 0, 0xff, 0xff},
	"w": {0xff, 0xff, 0xff, 0xff},
	".": {0, 0, 0, 0},
	"x": {0, 0, 0, 0xff},
}

type ImageChecker struct {
	Palette map[string]color.RGBA
}

func (ic ImageChecker) Check(expected []string, actual image.Image) error {
	var runes [][]string

	for _, str := range expected {
		runes = append(runes, strings.Split(str, ""))
	}

	if len(expected) != actual.Bounds().Dy() {
		ic.PrintDiff(expected, actual)
		return fmt.Errorf("expected %d rows, found %d", len(expected), actual.Bounds().Dy())
	}

	for y := 0; y < actual.Bounds().Dy(); y++ {
		if len(runes[y]) != actual.Bounds().Dx() {
			ic.PrintDiff(expected, actual)
			return fmt.Errorf(
				"row %d: expected %d columns, found %d",
				y, len(runes[0]), actual.Bounds().Dx())
		}
		for x := 0; x < actual.Bounds().Dx(); x++ {
			var actualColorRGBA color.RGBA
			actualColor := actual.At(x, y)
			if nrgba, ok := actualColor.(color.NRGBA); ok {
				r, g, b, a := nrgba.RGBA()
				actualColorRGBA = color.RGBA{uint8(r), uint8(g), uint8(b), uint8(a)}
			} else {
				actualColorRGBA = actualColor.(color.RGBA)
			}

			if actualColorRGBA != ic.Palette[string(runes[y][x])] {
				ic.PrintDiff(expected, actual)
				return fmt.Errorf("color differs at %d,%d", x, y)
			}
		}
	}

	return nil
}

func (ic ImageChecker) PrintDiff(expected []string, actual image.Image) {
	fmt.Println("EXPECTED")
	for _, row := range expected {
		fmt.Println(row)
	}

	fmt.Println("ACTUAL")
	ic.PrintImage(actual)
}

func (ic ImageChecker) PrintImage(im image.Image) {
	color2Ascii := map[color.RGBA]string{}
	for t, rgba := range ic.Palette {
		color2Ascii[rgba] = t
	}

	for y := 0; y < im.Bounds().Dy(); y++ {
		for x := 0; x < im.Bounds().Dx(); x++ {
			ascii := color2Ascii[im.At(x, y).(color.RGBA)]
			if ascii == "" {
				ascii = "?"
			}
			fmt.Printf(ascii)
		}
		fmt.Printf("\n")
	}
}

func printExpectedActual(expected []string, actual image.Image) {
	ic := ImageChecker{Palette: DefaultPalette}
	ic.PrintDiff(expected, actual)
}

func checkImage(expected []string, actual image.Image) error {
	ic := ImageChecker{Palette: DefaultPalette}
	return ic.Check(expected, actual)
}

func CheckImage(expected []string, actual image.Image) error {
	ic := ImageChecker{Palette: DefaultPalette}
	return ic.Check(expected, actual)
}

func PaintWidget(w Widget, bounds image.Rectangle, frameIdx int) image.Image {
	pb := w.PaintBounds(bounds, frameIdx)
	dc := canvas.NewGGCanvas(pb.Dx(), pb.Dy())
	w.Paint(dc, bounds, frameIdx)
	return dc.Image()
}
