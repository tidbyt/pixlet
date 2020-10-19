package runtime

import (
	"fmt"

	gocolor "image/color"

	"github.com/lucasb-eyer/go-colorful"
	"github.com/mitchellh/hashstructure"
	"go.starlark.net/starlark"

	"tidbyt.dev/pixlet/render"
)

type Plot struct {
	Widget
	render.Plot
	starlarkWidth  starlark.Int
	starlarkHeight starlark.Int
}

func parseFloatTuple(tuple starlark.Tuple, nullable bool) ([]*float64, error) {
	ret := make([]*float64, 0, tuple.Len())

	for i := 0; i < tuple.Len(); i++ {
		el := tuple.Index(i)

		if _, isNone := el.(starlark.NoneType); isNone {
			if !nullable {
				return nil, fmt.Errorf("element %d is None, expected float", i)
			}
			ret = append(ret, nil)
			continue
		}

		val, ok := el.(starlark.Float)
		if !ok {
			return nil, fmt.Errorf("element %d is %s, expected float", i, el.Type())
		}

		x := new(float64)
		*x = float64(val)
		ret = append(ret, x)
	}

	return ret, nil
}

func newPlot(thread *starlark.Thread, _ *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
	var (
		width, height starlark.Int
		data          *starlark.List
		xLim          starlark.Tuple
		yLim          starlark.Tuple
		color         starlark.String
		colorInverted starlark.String
		fill          starlark.Bool
	)

	if err := starlark.UnpackArgs(
		"Plot",
		args, kwargs,
		"width", &width,
		"height", &height,
		"data", &data,
		"xlim?", &xLim,
		"ylim?", &yLim,
		"color?", &color,
		"color_inverted?", &colorInverted,
		"fill?", &fill,
	); err != nil {
		return nil, fmt.Errorf("unpacking arguments for Plot: %s", err)
	}

	p := Plot{
		starlarkWidth:  width,
		starlarkHeight: height,
	}

	p.Width = int(width.BigInt().Int64())
	p.Height = int(height.BigInt().Int64())
	p.X = make([]float64, 0, data.Len())
	p.Y = make([]float64, 0, data.Len())
	p.Fill = bool(fill)

	if color.Len() != 0 {
		c, err := colorful.Hex(color.GoString())
		if err != nil {
			return nil, fmt.Errorf("color is not a valid hex string: %s", color.String())
		}
		r, g, b := c.RGB255()
		p.Color = &gocolor.RGBA{r, g, b, 0xff}
	}

	if colorInverted.Len() != 0 {
		c, err := colorful.Hex(colorInverted.GoString())
		if err != nil {
			return nil, fmt.Errorf("color_inverted is not a valid hex string: %s", colorInverted.String())
		}
		r, g, b := c.RGB255()
		p.ColorInverted = &gocolor.RGBA{r, g, b, 0xff}
	}

	xLimEl, err := parseFloatTuple(xLim, true)
	if err != nil {
		return nil, fmt.Errorf("parsing xlim: %s", err)
	}
	if len(xLimEl) != 0 {
		if len(xLimEl) != 2 {
			return nil, fmt.Errorf("xlim has len %d, expected 2", len(xLimEl))
		}
		p.XLimMin = xLimEl[0]
		p.XLimMax = xLimEl[1]
	}

	yLimEl, err := parseFloatTuple(yLim, true)
	if err != nil {
		return nil, fmt.Errorf("parsing ylim: %s", err)
	}
	if len(yLimEl) != 0 {
		if len(yLimEl) != 2 {
			return nil, fmt.Errorf("ylim has len %d, expected 2", len(yLimEl))
		}
		p.YLimMin = yLimEl[0]
		p.YLimMax = yLimEl[1]
	}

	var listVal starlark.Value
	iter := data.Iterate()
	defer iter.Done()

	i := 0
	for iter.Next(&listVal) {
		point, ok := listVal.(starlark.Tuple)
		if !ok {
			return nil, fmt.Errorf(
				"expected data to be a list of Tuple but found %s (at index %d)",
				listVal.Type(),
				i,
			)
		}

		if point.Len() != 2 {
			return nil, fmt.Errorf(
				"data tuples must have length 2, found length %d (at index %d)",
				point.Len(),
				i,
			)
		}

		x, xok := point.Index(0).(starlark.Float)
		y, yok := point.Index(1).(starlark.Float)
		if !xok || !yok {
			return nil, fmt.Errorf(
				"data tuples must hold float, found (%s, %s) (at index %d)",
				point.Index(0).Type(),
				point.Index(1).Type(),
				i,
			)
		}

		p.X = append(p.X, float64(x))
		p.Y = append(p.Y, float64(y))
	}

	return p, nil
}

func (p Plot) AsRenderWidget() render.Widget {
	return p.Plot
}

func (p Plot) AttrNames() []string {
	return []string{
		"width",
		"height",
	}
}

func (p Plot) Attr(name string) (starlark.Value, error) {
	switch name {
	case "width":
		return p.starlarkWidth, nil

	case "height":
		return p.starlarkHeight, nil

	default:
		return nil, nil
	}
}

func (p Plot) String() string       { return fmt.Sprintf("Plot(%dx%d)", p.Width, p.Height) }
func (p Plot) Type() string         { return "Plot" }
func (p Plot) Freeze()              {}
func (p Plot) Truth() starlark.Bool { return true }

func (p Plot) Hash() (uint32, error) {
	sum, err := hashstructure.Hash(p, nil)
	return uint32(sum), err
}
