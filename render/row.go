package render

import (
	"encoding/json"
	"image"

	"github.com/tidbyt/gg"
)

// Row lays out and draws its children horizontally (in a row).
//
// By default, a Row is as small as possible, while still holding all
// its children. However, if `expanded` is set, the Row will fill all
// available space horizontally. The height of a Row is always that of
// its tallest child.
//
// Alignment along the horizontal main axis is controlled by passing
// one of the following `main_align` values:
// - `"start"`: place children at the beginning of the row
// - `"end"`: place children at the end of the row
// - `"center"`: place children in the middle of the row
// - `"space_between"`: place equal space between children
// - `"space_evenly"`: equal space between children and before/after first/last child
// - `"space_around"`: equal space between children, and half of that before/after first/last child
//
// Alignment along the vertical cross axis is controlled by passing
// one of the following `cross_align` values:
// - `"start"`: place children at the top
// - `"end"`: place children at the bottom
// - `"center"`: place children at the center
//
// DOC(Children): Child widgets to lay out
// DOC(Expanded): Row should expand to fill all available horizontal space
// DOC(MainAlign): Alignment along horizontal main axis
// DOC(CrossAlign): Alignment along vertical cross axis
//
// EXAMPLE BEGIN
// render.Row(
//      children=[
//           render.Box(width=10, height=8, color="#a00"),
//           render.Box(width=14, height=6, color="#0a0"),
//           render.Box(width=16, height=4, color="#00a"),
//      ],
// )
// EXAMPLE END
//
// EXAMPLE BEGIN
// render.Row(
//      expanded=True,
//      main_align="space_between",
//      cross_align="end",
//      children=[
//           render.Box(width=10, height=8, color="#a00"),
//           render.Box(width=14, height=6, color="#0a0"),
//           render.Box(width=16, height=4, color="#00a"),
//      ],
// )
// EXAMPLE END
type Row struct {
	Type string `starlark:"-"`

	Children   []Widget `starlark:"children,required"`
	MainAlign  string   `starlark:"main_align"`
	CrossAlign string   `starlark:"cross_align"`
	Expanded   bool
}

func (r Row) PaintBounds(bounds image.Rectangle, frameIdx int) image.Rectangle {
	v := Vector{
		Vertical:   false,
		Children:   r.Children,
		MainAlign:  r.MainAlign,
		CrossAlign: r.CrossAlign,
		Expanded:   r.Expanded,
	}
	return v.PaintBounds(bounds, frameIdx)
}

func (r Row) Paint(dc *gg.Context, bounds image.Rectangle, frameIdx int) {
	v := Vector{
		Vertical:   false,
		Children:   r.Children,
		MainAlign:  r.MainAlign,
		CrossAlign: r.CrossAlign,
		Expanded:   r.Expanded,
	}
	v.Paint(dc, bounds, frameIdx)
}

func (r Row) FrameCount() int {
	return MaxFrameCount(r.Children)
}

func (r *Row) UnmarshalJSON(data []byte) error {
	type Alias Row
	aux := &struct {
		Children []json.RawMessage
		*Alias
	}{
		Alias: (*Alias)(r),
	}
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}
	children := []Widget{}
	for _, childData := range aux.Children {
		child, err := UnmarshalWidgetJSON(childData)
		if err != nil {
			return err
		}
		children = append(children, child)
	}
	r.Children = children
	return nil
}
