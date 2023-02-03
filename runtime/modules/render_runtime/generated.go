package render_runtime

// Code generated by runtime/gen. DO NOT EDIT.

import (
	"fmt"
	"sync"

	"github.com/mitchellh/hashstructure/v2"
	"go.starlark.net/starlark"
	"go.starlark.net/starlarkstruct"

	"tidbyt.dev/pixlet/render"
)

type RenderModule struct {
	once   sync.Once
	module starlark.StringDict
}

var renderModule = RenderModule{}

func LoadRenderModule() (starlark.StringDict, error) {
	renderModule.once.Do(func() {
		fnt := starlark.NewDict(len(render.Font))
		for k := range render.Font {
			fnt.SetKey(starlark.String(k), starlark.String(k))
		}
		fnt.Freeze()

		renderModule.module = starlark.StringDict{
			"render": &starlarkstruct.Module{
				Name: "render",
				Members: starlark.StringDict{
					"fonts": fnt,

					"Animation": starlark.NewBuiltin("Animation", newAnimation),

					"Box": starlark.NewBuiltin("Box", newBox),

					"Circle": starlark.NewBuiltin("Circle", newCircle),

					"Column": starlark.NewBuiltin("Column", newColumn),

					"Image": starlark.NewBuiltin("Image", newImage),

					"Marquee": starlark.NewBuiltin("Marquee", newMarquee),

					"Padding": starlark.NewBuiltin("Padding", newPadding),

					"PieChart": starlark.NewBuiltin("PieChart", newPieChart),

					"Plot": starlark.NewBuiltin("Plot", newPlot),

					"Root": starlark.NewBuiltin("Root", newRoot),

					"Row": starlark.NewBuiltin("Row", newRow),

					"Sequence": starlark.NewBuiltin("Sequence", newSequence),

					"Stack": starlark.NewBuiltin("Stack", newStack),

					"Text": starlark.NewBuiltin("Text", newText),

					"WrappedText": starlark.NewBuiltin("WrappedText", newWrappedText),
				},
			},
		}
	})

	return renderModule.module, nil
}

type Rootable interface {
	AsRenderRoot() render.Root
}

type Widget interface {
	AsRenderWidget() render.Widget
}
type Animation struct {
	Widget

	render.Animation

	starlarkChildren *starlark.List
}

func newAnimation(
	thread *starlark.Thread,
	_ *starlark.Builtin,
	args starlark.Tuple,
	kwargs []starlark.Tuple,
) (starlark.Value, error) {

	var (
		children *starlark.List
	)

	if err := starlark.UnpackArgs(
		"Animation",
		args, kwargs,
		"children?", &children,
	); err != nil {
		return nil, fmt.Errorf("unpacking arguments for Animation: %s", err)
	}

	w := &Animation{}

	var childrenVal starlark.Value
	childrenIter := children.Iterate()
	defer childrenIter.Done()
	for i := 0; childrenIter.Next(&childrenVal); {
		if _, isNone := childrenVal.(starlark.NoneType); isNone {
			continue
		}

		childrenChild, ok := childrenVal.(Widget)
		if !ok {
			return nil, fmt.Errorf(
				"expected children to be a list of Widget but found: %s (at index %d)",
				childrenVal.Type(),
				i,
			)
		}

		w.Children = append(w.Children, childrenChild.AsRenderWidget())
	}
	w.starlarkChildren = children

	return w, nil
}

func (w *Animation) AsRenderWidget() render.Widget {
	return &w.Animation
}

func (w *Animation) AttrNames() []string {
	return []string{
		"children",
	}
}

func (w *Animation) Attr(name string) (starlark.Value, error) {
	switch name {

	case "children":

		return w.starlarkChildren, nil

	default:
		return nil, nil
	}
}

func (w *Animation) String() string       { return "Animation(...)" }
func (w *Animation) Type() string         { return "Animation" }
func (w *Animation) Freeze()              {}
func (w *Animation) Truth() starlark.Bool { return true }

func (w *Animation) Hash() (uint32, error) {
	sum, err := hashstructure.Hash(w, hashstructure.FormatV2, nil)
	return uint32(sum), err
}

type Box struct {
	Widget

	render.Box

	starlarkChild starlark.Value

	starlarkColor starlark.String
}

func newBox(
	thread *starlark.Thread,
	_ *starlark.Builtin,
	args starlark.Tuple,
	kwargs []starlark.Tuple,
) (starlark.Value, error) {

	var (
		child   starlark.Value
		width   starlark.Int
		height  starlark.Int
		padding starlark.Int
		color   starlark.String
	)

	if err := starlark.UnpackArgs(
		"Box",
		args, kwargs,
		"child?", &child,
		"width?", &width,
		"height?", &height,
		"padding?", &padding,
		"color?", &color,
	); err != nil {
		return nil, fmt.Errorf("unpacking arguments for Box: %s", err)
	}

	w := &Box{}

	if child != nil {
		childWidget, ok := child.(Widget)
		if !ok {
			return nil, fmt.Errorf(
				"invalid type for child: %s (expected Widget)",
				child.Type(),
			)
		}
		w.Child = childWidget.AsRenderWidget()
		w.starlarkChild = child
	}

	w.Width = int(width.BigInt().Int64())

	w.Height = int(height.BigInt().Int64())

	w.Padding = int(padding.BigInt().Int64())

	w.starlarkColor = color
	if color.Len() > 0 {
		c, err := render.ParseColor(color.GoString())
		if err != nil {
			return nil, fmt.Errorf("color is not a valid hex string: %s", color.String())
		}
		w.Color = c
	}

	return w, nil
}

func (w *Box) AsRenderWidget() render.Widget {
	return &w.Box
}

func (w *Box) AttrNames() []string {
	return []string{
		"child", "width", "height", "padding", "color",
	}
}

func (w *Box) Attr(name string) (starlark.Value, error) {
	switch name {

	case "child":

		return w.starlarkChild, nil

	case "width":

		return starlark.MakeInt(int(w.Width)), nil

	case "height":

		return starlark.MakeInt(int(w.Height)), nil

	case "padding":

		return starlark.MakeInt(int(w.Padding)), nil

	case "color":

		return w.starlarkColor, nil

	default:
		return nil, nil
	}
}

func (w *Box) String() string       { return "Box(...)" }
func (w *Box) Type() string         { return "Box" }
func (w *Box) Freeze()              {}
func (w *Box) Truth() starlark.Bool { return true }

func (w *Box) Hash() (uint32, error) {
	sum, err := hashstructure.Hash(w, hashstructure.FormatV2, nil)
	return uint32(sum), err
}

type Circle struct {
	Widget

	render.Circle

	starlarkColor starlark.String

	starlarkChild starlark.Value
}

func newCircle(
	thread *starlark.Thread,
	_ *starlark.Builtin,
	args starlark.Tuple,
	kwargs []starlark.Tuple,
) (starlark.Value, error) {

	var (
		color    starlark.String
		diameter starlark.Int
		child    starlark.Value
	)

	if err := starlark.UnpackArgs(
		"Circle",
		args, kwargs,
		"color", &color,
		"diameter", &diameter,
		"child?", &child,
	); err != nil {
		return nil, fmt.Errorf("unpacking arguments for Circle: %s", err)
	}

	w := &Circle{}

	w.starlarkColor = color
	if color.Len() > 0 {
		c, err := render.ParseColor(color.GoString())
		if err != nil {
			return nil, fmt.Errorf("color is not a valid hex string: %s", color.String())
		}
		w.Color = c
	}

	w.Diameter = int(diameter.BigInt().Int64())

	if child != nil {
		childWidget, ok := child.(Widget)
		if !ok {
			return nil, fmt.Errorf(
				"invalid type for child: %s (expected Widget)",
				child.Type(),
			)
		}
		w.Child = childWidget.AsRenderWidget()
		w.starlarkChild = child
	}

	return w, nil
}

func (w *Circle) AsRenderWidget() render.Widget {
	return &w.Circle
}

func (w *Circle) AttrNames() []string {
	return []string{
		"color", "diameter", "child",
	}
}

func (w *Circle) Attr(name string) (starlark.Value, error) {
	switch name {

	case "color":

		return w.starlarkColor, nil

	case "diameter":

		return starlark.MakeInt(int(w.Diameter)), nil

	case "child":

		return w.starlarkChild, nil

	default:
		return nil, nil
	}
}

func (w *Circle) String() string       { return "Circle(...)" }
func (w *Circle) Type() string         { return "Circle" }
func (w *Circle) Freeze()              {}
func (w *Circle) Truth() starlark.Bool { return true }

func (w *Circle) Hash() (uint32, error) {
	sum, err := hashstructure.Hash(w, hashstructure.FormatV2, nil)
	return uint32(sum), err
}

type Column struct {
	Widget

	render.Column

	starlarkChildren *starlark.List
}

func newColumn(
	thread *starlark.Thread,
	_ *starlark.Builtin,
	args starlark.Tuple,
	kwargs []starlark.Tuple,
) (starlark.Value, error) {

	var (
		children    *starlark.List
		main_align  starlark.String
		cross_align starlark.String
		expanded    starlark.Bool
	)

	if err := starlark.UnpackArgs(
		"Column",
		args, kwargs,
		"children", &children,
		"main_align?", &main_align,
		"cross_align?", &cross_align,
		"expanded?", &expanded,
	); err != nil {
		return nil, fmt.Errorf("unpacking arguments for Column: %s", err)
	}

	w := &Column{}

	var childrenVal starlark.Value
	childrenIter := children.Iterate()
	defer childrenIter.Done()
	for i := 0; childrenIter.Next(&childrenVal); {
		if _, isNone := childrenVal.(starlark.NoneType); isNone {
			continue
		}

		childrenChild, ok := childrenVal.(Widget)
		if !ok {
			return nil, fmt.Errorf(
				"expected children to be a list of Widget but found: %s (at index %d)",
				childrenVal.Type(),
				i,
			)
		}

		w.Children = append(w.Children, childrenChild.AsRenderWidget())
	}
	w.starlarkChildren = children

	w.MainAlign = main_align.GoString()

	w.CrossAlign = cross_align.GoString()

	w.Expanded = bool(expanded)

	return w, nil
}

func (w *Column) AsRenderWidget() render.Widget {
	return &w.Column
}

func (w *Column) AttrNames() []string {
	return []string{
		"children", "main_align", "cross_align", "expanded",
	}
}

func (w *Column) Attr(name string) (starlark.Value, error) {
	switch name {

	case "children":

		return w.starlarkChildren, nil

	case "main_align":

		return starlark.String(w.MainAlign), nil

	case "cross_align":

		return starlark.String(w.CrossAlign), nil

	case "expanded":

		return starlark.Bool(w.Expanded), nil

	default:
		return nil, nil
	}
}

func (w *Column) String() string       { return "Column(...)" }
func (w *Column) Type() string         { return "Column" }
func (w *Column) Freeze()              {}
func (w *Column) Truth() starlark.Bool { return true }

func (w *Column) Hash() (uint32, error) {
	sum, err := hashstructure.Hash(w, hashstructure.FormatV2, nil)
	return uint32(sum), err
}

type Image struct {
	Widget

	render.Image

	size *starlark.Builtin
}

func newImage(
	thread *starlark.Thread,
	_ *starlark.Builtin,
	args starlark.Tuple,
	kwargs []starlark.Tuple,
) (starlark.Value, error) {

	var (
		src    starlark.String
		width  starlark.Int
		height starlark.Int
	)

	if err := starlark.UnpackArgs(
		"Image",
		args, kwargs,
		"src", &src,
		"width?", &width,
		"height?", &height,
	); err != nil {
		return nil, fmt.Errorf("unpacking arguments for Image: %s", err)
	}

	w := &Image{}

	w.Src = src.GoString()

	w.Width = int(width.BigInt().Int64())

	w.Height = int(height.BigInt().Int64())

	w.size = starlark.NewBuiltin("size", imageSize)

	if err := w.Init(); err != nil {
		return nil, err
	}

	return w, nil
}

func (w *Image) AsRenderWidget() render.Widget {
	return &w.Image
}

func (w *Image) AttrNames() []string {
	return []string{
		"src", "width", "height", "delay",
	}
}

func (w *Image) Attr(name string) (starlark.Value, error) {
	switch name {

	case "src":

		return starlark.String(w.Src), nil

	case "width":

		return starlark.MakeInt(int(w.Width)), nil

	case "height":

		return starlark.MakeInt(int(w.Height)), nil

	case "delay":

		return starlark.MakeInt(int(w.Delay)), nil

	case "size":
		return w.size.BindReceiver(w), nil

	default:
		return nil, nil
	}
}

func (w *Image) String() string       { return "Image(...)" }
func (w *Image) Type() string         { return "Image" }
func (w *Image) Freeze()              {}
func (w *Image) Truth() starlark.Bool { return true }

func (w *Image) Hash() (uint32, error) {
	sum, err := hashstructure.Hash(w, hashstructure.FormatV2, nil)
	return uint32(sum), err
}

func imageSize(
	thread *starlark.Thread,
	b *starlark.Builtin,
	args starlark.Tuple,
	kwargs []starlark.Tuple) (starlark.Value, error) {

	w := b.Receiver().(*Image)
	width, height := w.Size()

	return starlark.Tuple([]starlark.Value{
		starlark.MakeInt(width),
		starlark.MakeInt(height),
	}), nil
}

type Marquee struct {
	Widget

	render.Marquee

	starlarkChild starlark.Value
}

func newMarquee(
	thread *starlark.Thread,
	_ *starlark.Builtin,
	args starlark.Tuple,
	kwargs []starlark.Tuple,
) (starlark.Value, error) {

	var (
		child            starlark.Value
		width            starlark.Int
		height           starlark.Int
		offset_start     starlark.Int
		offset_end       starlark.Int
		scroll_direction starlark.String
		align            starlark.String
		delay            starlark.Int
	)

	if err := starlark.UnpackArgs(
		"Marquee",
		args, kwargs,
		"child", &child,
		"width?", &width,
		"height?", &height,
		"offset_start?", &offset_start,
		"offset_end?", &offset_end,
		"scroll_direction?", &scroll_direction,
		"align?", &align,
		"delay?", &delay,
	); err != nil {
		return nil, fmt.Errorf("unpacking arguments for Marquee: %s", err)
	}

	w := &Marquee{}

	if child != nil {
		childWidget, ok := child.(Widget)
		if !ok {
			return nil, fmt.Errorf(
				"invalid type for child: %s (expected Widget)",
				child.Type(),
			)
		}
		w.Child = childWidget.AsRenderWidget()
		w.starlarkChild = child
	}

	w.Width = int(width.BigInt().Int64())

	w.Height = int(height.BigInt().Int64())

	w.OffsetStart = int(offset_start.BigInt().Int64())

	w.OffsetEnd = int(offset_end.BigInt().Int64())

	w.ScrollDirection = scroll_direction.GoString()

	w.Align = align.GoString()

	w.Delay = int(delay.BigInt().Int64())

	return w, nil
}

func (w *Marquee) AsRenderWidget() render.Widget {
	return &w.Marquee
}

func (w *Marquee) AttrNames() []string {
	return []string{
		"child", "width", "height", "offset_start", "offset_end", "scroll_direction", "align", "delay",
	}
}

func (w *Marquee) Attr(name string) (starlark.Value, error) {
	switch name {

	case "child":

		return w.starlarkChild, nil

	case "width":

		return starlark.MakeInt(int(w.Width)), nil

	case "height":

		return starlark.MakeInt(int(w.Height)), nil

	case "offset_start":

		return starlark.MakeInt(int(w.OffsetStart)), nil

	case "offset_end":

		return starlark.MakeInt(int(w.OffsetEnd)), nil

	case "scroll_direction":

		return starlark.String(w.ScrollDirection), nil

	case "align":

		return starlark.String(w.Align), nil

	case "delay":

		return starlark.MakeInt(int(w.Delay)), nil

	default:
		return nil, nil
	}
}

func (w *Marquee) String() string       { return "Marquee(...)" }
func (w *Marquee) Type() string         { return "Marquee" }
func (w *Marquee) Freeze()              {}
func (w *Marquee) Truth() starlark.Bool { return true }

func (w *Marquee) Hash() (uint32, error) {
	sum, err := hashstructure.Hash(w, hashstructure.FormatV2, nil)
	return uint32(sum), err
}

type Padding struct {
	Widget

	render.Padding

	starlarkChild starlark.Value

	starlarkPad starlark.Value

	starlarkColor starlark.String
}

func newPadding(
	thread *starlark.Thread,
	_ *starlark.Builtin,
	args starlark.Tuple,
	kwargs []starlark.Tuple,
) (starlark.Value, error) {

	var (
		child    starlark.Value
		pad      starlark.Value
		expanded starlark.Bool
		color    starlark.String
	)

	if err := starlark.UnpackArgs(
		"Padding",
		args, kwargs,
		"child", &child,
		"pad?", &pad,
		"expanded?", &expanded,
		"color?", &color,
	); err != nil {
		return nil, fmt.Errorf("unpacking arguments for Padding: %s", err)
	}

	w := &Padding{}

	if child != nil {
		childWidget, ok := child.(Widget)
		if !ok {
			return nil, fmt.Errorf(
				"invalid type for child: %s (expected Widget)",
				child.Type(),
			)
		}
		w.Child = childWidget.AsRenderWidget()
		w.starlarkChild = child
	}

	w.starlarkPad = pad
	switch padVal := pad.(type) {
	case starlark.Int:
		padInt := int(padVal.BigInt().Int64())
		w.Pad.Left = padInt
		w.Pad.Top = padInt
		w.Pad.Right = padInt
		w.Pad.Bottom = padInt
	case starlark.Tuple:
		padList := []starlark.Value(padVal)
		if len(padList) != 4 {
			return nil, fmt.Errorf(
				"pad tuple must hold 4 elements (left, top, right, bottom), found %d",
				len(padList),
			)
		}
		padListInt := make([]starlark.Int, 4)
		for i := 0; i < 4; i++ {
			pi, ok := padList[i].(starlark.Int)
			if !ok {
				return nil, fmt.Errorf("pad element %d is not int", i)
			}
			padListInt[i] = pi
		}
		w.Pad.Left = int(padListInt[0].BigInt().Int64())
		w.Pad.Top = int(padListInt[1].BigInt().Int64())
		w.Pad.Right = int(padListInt[2].BigInt().Int64())
		w.Pad.Bottom = int(padListInt[3].BigInt().Int64())
	default:
		return nil, fmt.Errorf("pad must be int or 4-tuple of int")
	}

	w.Expanded = bool(expanded)

	w.starlarkColor = color
	if color.Len() > 0 {
		c, err := render.ParseColor(color.GoString())
		if err != nil {
			return nil, fmt.Errorf("color is not a valid hex string: %s", color.String())
		}
		w.Color = c
	}

	return w, nil
}

func (w *Padding) AsRenderWidget() render.Widget {
	return &w.Padding
}

func (w *Padding) AttrNames() []string {
	return []string{
		"child", "pad", "expanded", "color",
	}
}

func (w *Padding) Attr(name string) (starlark.Value, error) {
	switch name {

	case "child":

		return w.starlarkChild, nil

	case "pad":

		return w.starlarkPad, nil

	case "expanded":

		return starlark.Bool(w.Expanded), nil

	case "color":

		return w.starlarkColor, nil

	default:
		return nil, nil
	}
}

func (w *Padding) String() string       { return "Padding(...)" }
func (w *Padding) Type() string         { return "Padding" }
func (w *Padding) Freeze()              {}
func (w *Padding) Truth() starlark.Bool { return true }

func (w *Padding) Hash() (uint32, error) {
	sum, err := hashstructure.Hash(w, hashstructure.FormatV2, nil)
	return uint32(sum), err
}

type PieChart struct {
	Widget

	render.PieChart

	starlarkColors *starlark.List

	starlarkWeights *starlark.List
}

func newPieChart(
	thread *starlark.Thread,
	_ *starlark.Builtin,
	args starlark.Tuple,
	kwargs []starlark.Tuple,
) (starlark.Value, error) {

	var (
		colors   *starlark.List
		weights  *starlark.List
		diameter starlark.Int
	)

	if err := starlark.UnpackArgs(
		"PieChart",
		args, kwargs,
		"colors", &colors,
		"weights", &weights,
		"diameter", &diameter,
	); err != nil {
		return nil, fmt.Errorf("unpacking arguments for PieChart: %s", err)
	}

	w := &PieChart{}

	w.starlarkColors = colors
	if val, err := ColorSeriesFromStarlark(colors); err == nil {
		w.Colors = val
	} else {
		return nil, err
	}

	w.starlarkWeights = weights
	if val, err := WeightsFromStarlark(weights); err == nil {
		w.Weights = val
	} else {
		return nil, err
	}

	w.Diameter = int(diameter.BigInt().Int64())

	return w, nil
}

func (w *PieChart) AsRenderWidget() render.Widget {
	return &w.PieChart
}

func (w *PieChart) AttrNames() []string {
	return []string{
		"colors", "weights", "diameter",
	}
}

func (w *PieChart) Attr(name string) (starlark.Value, error) {
	switch name {

	case "colors":

		return w.starlarkColors, nil

	case "weights":

		return w.starlarkWeights, nil

	case "diameter":

		return starlark.MakeInt(int(w.Diameter)), nil

	default:
		return nil, nil
	}
}

func (w *PieChart) String() string       { return "PieChart(...)" }
func (w *PieChart) Type() string         { return "PieChart" }
func (w *PieChart) Freeze()              {}
func (w *PieChart) Truth() starlark.Bool { return true }

func (w *PieChart) Hash() (uint32, error) {
	sum, err := hashstructure.Hash(w, hashstructure.FormatV2, nil)
	return uint32(sum), err
}

type Plot struct {
	Widget

	render.Plot

	starlarkData *starlark.List

	starlarkColor starlark.String

	starlarkColorInverted starlark.String

	starlarkXLim starlark.Tuple

	starlarkYLim starlark.Tuple

	starlarkFillColor starlark.String

	starlarkFillColorInverted starlark.String
}

func newPlot(
	thread *starlark.Thread,
	_ *starlark.Builtin,
	args starlark.Tuple,
	kwargs []starlark.Tuple,
) (starlark.Value, error) {

	var (
		data                *starlark.List
		width               starlark.Int
		height              starlark.Int
		color               starlark.String
		color_inverted      starlark.String
		x_lim               starlark.Tuple
		y_lim               starlark.Tuple
		fill                starlark.Bool
		chart_type          starlark.String
		fill_color          starlark.String
		fill_color_inverted starlark.String
	)

	if err := starlark.UnpackArgs(
		"Plot",
		args, kwargs,
		"data", &data,
		"width", &width,
		"height", &height,
		"color?", &color,
		"color_inverted?", &color_inverted,
		"x_lim?", &x_lim,
		"y_lim?", &y_lim,
		"fill?", &fill,
		"chart_type?", &chart_type,
		"fill_color?", &fill_color,
		"fill_color_inverted?", &fill_color_inverted,
	); err != nil {
		return nil, fmt.Errorf("unpacking arguments for Plot: %s", err)
	}

	w := &Plot{}

	w.starlarkData = data
	if val, err := DataSeriesFromStarlark(data); err == nil {
		w.Data = val
	} else {
		return nil, err
	}

	w.Width = int(width.BigInt().Int64())

	w.Height = int(height.BigInt().Int64())

	w.starlarkColor = color
	if color.Len() > 0 {
		c, err := render.ParseColor(color.GoString())
		if err != nil {
			return nil, fmt.Errorf("color is not a valid hex string: %s", color.String())
		}
		w.Color = c
	}

	w.starlarkColorInverted = color_inverted
	if color_inverted.Len() > 0 {
		c, err := render.ParseColor(color_inverted.GoString())
		if err != nil {
			return nil, fmt.Errorf("color_inverted is not a valid hex string: %s", color_inverted.String())
		}
		w.ColorInverted = c
	}

	w.starlarkXLim = x_lim
	if val, err := DataPointFromStarlark(x_lim); err == nil {
		w.XLim = val
	} else {
		return nil, err
	}

	w.starlarkYLim = y_lim
	if val, err := DataPointFromStarlark(y_lim); err == nil {
		w.YLim = val
	} else {
		return nil, err
	}

	w.Fill = bool(fill)

	w.ChartType = chart_type.GoString()

	w.starlarkFillColor = fill_color
	if fill_color.Len() > 0 {
		c, err := render.ParseColor(fill_color.GoString())
		if err != nil {
			return nil, fmt.Errorf("fill_color is not a valid hex string: %s", fill_color.String())
		}
		w.FillColor = c
	}

	w.starlarkFillColorInverted = fill_color_inverted
	if fill_color_inverted.Len() > 0 {
		c, err := render.ParseColor(fill_color_inverted.GoString())
		if err != nil {
			return nil, fmt.Errorf("fill_color_inverted is not a valid hex string: %s", fill_color_inverted.String())
		}
		w.FillColorInverted = c
	}

	return w, nil
}

func (w *Plot) AsRenderWidget() render.Widget {
	return &w.Plot
}

func (w *Plot) AttrNames() []string {
	return []string{
		"data", "width", "height", "color", "color_inverted", "x_lim", "y_lim", "fill", "chart_type", "fill_color", "fill_color_inverted",
	}
}

func (w *Plot) Attr(name string) (starlark.Value, error) {
	switch name {

	case "data":

		return w.starlarkData, nil

	case "width":

		return starlark.MakeInt(int(w.Width)), nil

	case "height":

		return starlark.MakeInt(int(w.Height)), nil

	case "color":

		return w.starlarkColor, nil

	case "color_inverted":

		return w.starlarkColorInverted, nil

	case "x_lim":

		return w.starlarkXLim, nil

	case "y_lim":

		return w.starlarkYLim, nil

	case "fill":

		return starlark.Bool(w.Fill), nil

	case "chart_type":

		return starlark.String(w.ChartType), nil

	case "fill_color":

		return w.starlarkFillColor, nil

	case "fill_color_inverted":

		return w.starlarkFillColorInverted, nil

	default:
		return nil, nil
	}
}

func (w *Plot) String() string       { return "Plot(...)" }
func (w *Plot) Type() string         { return "Plot" }
func (w *Plot) Freeze()              {}
func (w *Plot) Truth() starlark.Bool { return true }

func (w *Plot) Hash() (uint32, error) {
	sum, err := hashstructure.Hash(w, hashstructure.FormatV2, nil)
	return uint32(sum), err
}

type Root struct {
	render.Root

	starlarkChild starlark.Value
}

func newRoot(
	thread *starlark.Thread,
	_ *starlark.Builtin,
	args starlark.Tuple,
	kwargs []starlark.Tuple,
) (starlark.Value, error) {

	var (
		child               starlark.Value
		delay               starlark.Int
		max_age             starlark.Int
		show_full_animation starlark.Bool
	)

	if err := starlark.UnpackArgs(
		"Root",
		args, kwargs,
		"child", &child,
		"delay?", &delay,
		"max_age?", &max_age,
		"show_full_animation?", &show_full_animation,
	); err != nil {
		return nil, fmt.Errorf("unpacking arguments for Root: %s", err)
	}

	w := &Root{}

	if child != nil {
		childWidget, ok := child.(Widget)
		if !ok {
			return nil, fmt.Errorf(
				"invalid type for child: %s (expected Widget)",
				child.Type(),
			)
		}
		w.Child = childWidget.AsRenderWidget()
		w.starlarkChild = child
	}

	if val, err := starlark.AsInt32(delay); err == nil {
		w.Delay = int32(val)
	} else {
		return nil, err
	}

	if val, err := starlark.AsInt32(max_age); err == nil {
		w.MaxAge = int32(val)
	} else {
		return nil, err
	}

	w.ShowFullAnimation = bool(show_full_animation)

	return w, nil
}

func (w *Root) AsRenderRoot() render.Root {
	return w.Root
}

func (w *Root) AttrNames() []string {
	return []string{
		"child", "delay", "max_age", "show_full_animation",
	}
}

func (w *Root) Attr(name string) (starlark.Value, error) {
	switch name {

	case "child":

		return w.starlarkChild, nil

	case "delay":

		return starlark.MakeInt(int(w.Delay)), nil

	case "max_age":

		return starlark.MakeInt(int(w.MaxAge)), nil

	case "show_full_animation":

		return starlark.Bool(w.ShowFullAnimation), nil

	default:
		return nil, nil
	}
}

func (w *Root) String() string       { return "Root(...)" }
func (w *Root) Type() string         { return "Root" }
func (w *Root) Freeze()              {}
func (w *Root) Truth() starlark.Bool { return true }

func (w *Root) Hash() (uint32, error) {
	sum, err := hashstructure.Hash(w, hashstructure.FormatV2, nil)
	return uint32(sum), err
}

type Row struct {
	Widget

	render.Row

	starlarkChildren *starlark.List
}

func newRow(
	thread *starlark.Thread,
	_ *starlark.Builtin,
	args starlark.Tuple,
	kwargs []starlark.Tuple,
) (starlark.Value, error) {

	var (
		children    *starlark.List
		main_align  starlark.String
		cross_align starlark.String
		expanded    starlark.Bool
	)

	if err := starlark.UnpackArgs(
		"Row",
		args, kwargs,
		"children", &children,
		"main_align?", &main_align,
		"cross_align?", &cross_align,
		"expanded?", &expanded,
	); err != nil {
		return nil, fmt.Errorf("unpacking arguments for Row: %s", err)
	}

	w := &Row{}

	var childrenVal starlark.Value
	childrenIter := children.Iterate()
	defer childrenIter.Done()
	for i := 0; childrenIter.Next(&childrenVal); {
		if _, isNone := childrenVal.(starlark.NoneType); isNone {
			continue
		}

		childrenChild, ok := childrenVal.(Widget)
		if !ok {
			return nil, fmt.Errorf(
				"expected children to be a list of Widget but found: %s (at index %d)",
				childrenVal.Type(),
				i,
			)
		}

		w.Children = append(w.Children, childrenChild.AsRenderWidget())
	}
	w.starlarkChildren = children

	w.MainAlign = main_align.GoString()

	w.CrossAlign = cross_align.GoString()

	w.Expanded = bool(expanded)

	return w, nil
}

func (w *Row) AsRenderWidget() render.Widget {
	return &w.Row
}

func (w *Row) AttrNames() []string {
	return []string{
		"children", "main_align", "cross_align", "expanded",
	}
}

func (w *Row) Attr(name string) (starlark.Value, error) {
	switch name {

	case "children":

		return w.starlarkChildren, nil

	case "main_align":

		return starlark.String(w.MainAlign), nil

	case "cross_align":

		return starlark.String(w.CrossAlign), nil

	case "expanded":

		return starlark.Bool(w.Expanded), nil

	default:
		return nil, nil
	}
}

func (w *Row) String() string       { return "Row(...)" }
func (w *Row) Type() string         { return "Row" }
func (w *Row) Freeze()              {}
func (w *Row) Truth() starlark.Bool { return true }

func (w *Row) Hash() (uint32, error) {
	sum, err := hashstructure.Hash(w, hashstructure.FormatV2, nil)
	return uint32(sum), err
}

type Sequence struct {
	Widget

	render.Sequence

	starlarkChildren *starlark.List
}

func newSequence(
	thread *starlark.Thread,
	_ *starlark.Builtin,
	args starlark.Tuple,
	kwargs []starlark.Tuple,
) (starlark.Value, error) {

	var (
		children *starlark.List
	)

	if err := starlark.UnpackArgs(
		"Sequence",
		args, kwargs,
		"children", &children,
	); err != nil {
		return nil, fmt.Errorf("unpacking arguments for Sequence: %s", err)
	}

	w := &Sequence{}

	var childrenVal starlark.Value
	childrenIter := children.Iterate()
	defer childrenIter.Done()
	for i := 0; childrenIter.Next(&childrenVal); {
		if _, isNone := childrenVal.(starlark.NoneType); isNone {
			continue
		}

		childrenChild, ok := childrenVal.(Widget)
		if !ok {
			return nil, fmt.Errorf(
				"expected children to be a list of Widget but found: %s (at index %d)",
				childrenVal.Type(),
				i,
			)
		}

		w.Children = append(w.Children, childrenChild.AsRenderWidget())
	}
	w.starlarkChildren = children

	return w, nil
}

func (w *Sequence) AsRenderWidget() render.Widget {
	return &w.Sequence
}

func (w *Sequence) AttrNames() []string {
	return []string{
		"children",
	}
}

func (w *Sequence) Attr(name string) (starlark.Value, error) {
	switch name {

	case "children":

		return w.starlarkChildren, nil

	default:
		return nil, nil
	}
}

func (w *Sequence) String() string       { return "Sequence(...)" }
func (w *Sequence) Type() string         { return "Sequence" }
func (w *Sequence) Freeze()              {}
func (w *Sequence) Truth() starlark.Bool { return true }

func (w *Sequence) Hash() (uint32, error) {
	sum, err := hashstructure.Hash(w, hashstructure.FormatV2, nil)
	return uint32(sum), err
}

type Stack struct {
	Widget

	render.Stack

	starlarkChildren *starlark.List
}

func newStack(
	thread *starlark.Thread,
	_ *starlark.Builtin,
	args starlark.Tuple,
	kwargs []starlark.Tuple,
) (starlark.Value, error) {

	var (
		children *starlark.List
	)

	if err := starlark.UnpackArgs(
		"Stack",
		args, kwargs,
		"children", &children,
	); err != nil {
		return nil, fmt.Errorf("unpacking arguments for Stack: %s", err)
	}

	w := &Stack{}

	var childrenVal starlark.Value
	childrenIter := children.Iterate()
	defer childrenIter.Done()
	for i := 0; childrenIter.Next(&childrenVal); {
		if _, isNone := childrenVal.(starlark.NoneType); isNone {
			continue
		}

		childrenChild, ok := childrenVal.(Widget)
		if !ok {
			return nil, fmt.Errorf(
				"expected children to be a list of Widget but found: %s (at index %d)",
				childrenVal.Type(),
				i,
			)
		}

		w.Children = append(w.Children, childrenChild.AsRenderWidget())
	}
	w.starlarkChildren = children

	return w, nil
}

func (w *Stack) AsRenderWidget() render.Widget {
	return &w.Stack
}

func (w *Stack) AttrNames() []string {
	return []string{
		"children",
	}
}

func (w *Stack) Attr(name string) (starlark.Value, error) {
	switch name {

	case "children":

		return w.starlarkChildren, nil

	default:
		return nil, nil
	}
}

func (w *Stack) String() string       { return "Stack(...)" }
func (w *Stack) Type() string         { return "Stack" }
func (w *Stack) Freeze()              {}
func (w *Stack) Truth() starlark.Bool { return true }

func (w *Stack) Hash() (uint32, error) {
	sum, err := hashstructure.Hash(w, hashstructure.FormatV2, nil)
	return uint32(sum), err
}

type Text struct {
	Widget

	render.Text

	starlarkColor starlark.String

	size *starlark.Builtin
}

func newText(
	thread *starlark.Thread,
	_ *starlark.Builtin,
	args starlark.Tuple,
	kwargs []starlark.Tuple,
) (starlark.Value, error) {

	var (
		content starlark.String
		font    starlark.String
		height  starlark.Int
		offset  starlark.Int
		color   starlark.String
	)

	if err := starlark.UnpackArgs(
		"Text",
		args, kwargs,
		"content", &content,
		"font?", &font,
		"height?", &height,
		"offset?", &offset,
		"color?", &color,
	); err != nil {
		return nil, fmt.Errorf("unpacking arguments for Text: %s", err)
	}

	w := &Text{}

	w.Content = content.GoString()

	w.Font = font.GoString()

	w.Height = int(height.BigInt().Int64())

	w.Offset = int(offset.BigInt().Int64())

	w.starlarkColor = color
	if color.Len() > 0 {
		c, err := render.ParseColor(color.GoString())
		if err != nil {
			return nil, fmt.Errorf("color is not a valid hex string: %s", color.String())
		}
		w.Color = c
	}

	w.size = starlark.NewBuiltin("size", textSize)

	if err := w.Init(); err != nil {
		return nil, err
	}

	return w, nil
}

func (w *Text) AsRenderWidget() render.Widget {
	return &w.Text
}

func (w *Text) AttrNames() []string {
	return []string{
		"content", "font", "height", "offset", "color",
	}
}

func (w *Text) Attr(name string) (starlark.Value, error) {
	switch name {

	case "content":

		return starlark.String(w.Content), nil

	case "font":

		return starlark.String(w.Font), nil

	case "height":

		return starlark.MakeInt(int(w.Height)), nil

	case "offset":

		return starlark.MakeInt(int(w.Offset)), nil

	case "color":

		return w.starlarkColor, nil

	case "size":
		return w.size.BindReceiver(w), nil

	default:
		return nil, nil
	}
}

func (w *Text) String() string       { return "Text(...)" }
func (w *Text) Type() string         { return "Text" }
func (w *Text) Freeze()              {}
func (w *Text) Truth() starlark.Bool { return true }

func (w *Text) Hash() (uint32, error) {
	sum, err := hashstructure.Hash(w, hashstructure.FormatV2, nil)
	return uint32(sum), err
}

func textSize(
	thread *starlark.Thread,
	b *starlark.Builtin,
	args starlark.Tuple,
	kwargs []starlark.Tuple) (starlark.Value, error) {

	w := b.Receiver().(*Text)
	width, height := w.Size()

	return starlark.Tuple([]starlark.Value{
		starlark.MakeInt(width),
		starlark.MakeInt(height),
	}), nil
}

type WrappedText struct {
	Widget

	render.WrappedText

	starlarkColor starlark.String
}

func newWrappedText(
	thread *starlark.Thread,
	_ *starlark.Builtin,
	args starlark.Tuple,
	kwargs []starlark.Tuple,
) (starlark.Value, error) {

	var (
		content     starlark.String
		font        starlark.String
		height      starlark.Int
		width       starlark.Int
		linespacing starlark.Int
		color       starlark.String
		align       starlark.String
	)

	if err := starlark.UnpackArgs(
		"WrappedText",
		args, kwargs,
		"content", &content,
		"font?", &font,
		"height?", &height,
		"width?", &width,
		"linespacing?", &linespacing,
		"color?", &color,
		"align?", &align,
	); err != nil {
		return nil, fmt.Errorf("unpacking arguments for WrappedText: %s", err)
	}

	w := &WrappedText{}

	w.Content = content.GoString()

	w.Font = font.GoString()

	w.Height = int(height.BigInt().Int64())

	w.Width = int(width.BigInt().Int64())

	w.LineSpacing = int(linespacing.BigInt().Int64())

	w.starlarkColor = color
	if color.Len() > 0 {
		c, err := render.ParseColor(color.GoString())
		if err != nil {
			return nil, fmt.Errorf("color is not a valid hex string: %s", color.String())
		}
		w.Color = c
	}

	w.Align = align.GoString()

	return w, nil
}

func (w *WrappedText) AsRenderWidget() render.Widget {
	return &w.WrappedText
}

func (w *WrappedText) AttrNames() []string {
	return []string{
		"content", "font", "height", "width", "linespacing", "color", "align",
	}
}

func (w *WrappedText) Attr(name string) (starlark.Value, error) {
	switch name {

	case "content":

		return starlark.String(w.Content), nil

	case "font":

		return starlark.String(w.Font), nil

	case "height":

		return starlark.MakeInt(int(w.Height)), nil

	case "width":

		return starlark.MakeInt(int(w.Width)), nil

	case "linespacing":

		return starlark.MakeInt(int(w.LineSpacing)), nil

	case "color":

		return w.starlarkColor, nil

	case "align":

		return starlark.String(w.Align), nil

	default:
		return nil, nil
	}
}

func (w *WrappedText) String() string       { return "WrappedText(...)" }
func (w *WrappedText) Type() string         { return "WrappedText" }
func (w *WrappedText) Freeze()              {}
func (w *WrappedText) Truth() starlark.Bool { return true }

func (w *WrappedText) Hash() (uint32, error) {
	sum, err := hashstructure.Hash(w, hashstructure.FormatV2, nil)
	return uint32(sum), err
}
