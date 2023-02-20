package render

import (
	"fmt"
	"image/color"
	"strings"
)

func ParseColor(scol string) (color.Color, error) {
	var format string
	var fourBits bool
	var hasAlpha bool

	scol = strings.TrimPrefix(scol, "#")

	switch len(scol) {
	case 3:
		format = "%1x%1x%1x"
		fourBits = true
		hasAlpha = false
	case 4:
		format = "%1x%1x%1x%1x"
		fourBits = true
		hasAlpha = true
	case 6:
		format = "%02x%02x%02x"
		fourBits = false
		hasAlpha = false
	case 8:
		format = "%02x%02x%02x%02x"
		fourBits = false
		hasAlpha = true
	default:
		return color.Gray{0}, fmt.Errorf("color: %v is not a hex-color", scol)
	}

	var r, g, b, a uint8

	if hasAlpha {
		n, err := fmt.Sscanf(scol, format, &r, &g, &b, &a)
		if err != nil {
			return color.Gray{0}, err
		}
		if n != 4 {
			return color.Gray{0}, fmt.Errorf("color: %v is not a hex-color", scol)
		}
	} else {
		n, err := fmt.Sscanf(scol, format, &r, &g, &b)
		if err != nil {
			return color.Gray{0}, err
		}
		if n != 3 {
			return color.Gray{0}, fmt.Errorf("color: %v is not a hex-color", scol)
		}
		if fourBits {
			a = 15
		} else {
			a = 255
		}
	}

	if fourBits {
		r |= r << 4
		g |= g << 4
		b |= b << 4
		a |= a << 4
	}

	return color.NRGBA{r, g, b, a}, nil
}
