package render

import (
	"fmt"
	"image"
	"image/color"
)

var DefaultPalette = map[string]color.RGBA{
	"r": color.RGBA{0xff, 0, 0, 0xff},
	"g": color.RGBA{0, 0xff, 0, 0xff},
	"b": color.RGBA{0, 0, 0xff, 0xff},
	"w": color.RGBA{0xff, 0xff, 0xff, 0xff},
	".": color.RGBA{0, 0, 0, 0},
	"x": color.RGBA{0, 0, 0, 0xff},
}

type ImageChecker struct {
	palette map[string]color.RGBA
}

func (ic ImageChecker) Check(expected []string, actual image.Image) error {
	if len(expected) != actual.Bounds().Dy() {
		ic.PrintDiff(expected, actual)
		return fmt.Errorf("expected %d rows, found %d", len(expected), actual.Bounds().Dy())
	}

	for y := 0; y < actual.Bounds().Dy(); y++ {
		if len(expected[y]) != actual.Bounds().Dx() {
			ic.PrintDiff(expected, actual)
			return fmt.Errorf(
				"row %d: expected %d columns, found %d",
				y, len(expected[0]), actual.Bounds().Dx())
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

			if actualColorRGBA != ic.palette[string(expected[y][x])] {
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
	for t, rgba := range ic.palette {
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
	ic := ImageChecker{palette: DefaultPalette}
	ic.PrintDiff(expected, actual)
}

func checkImage(expected []string, actual image.Image) error {
	ic := ImageChecker{palette: DefaultPalette}
	return ic.Check(expected, actual)
}

func CheckImage(expected []string, actual image.Image) error {
	ic := ImageChecker{palette: DefaultPalette}
	return ic.Check(expected, actual)
}
