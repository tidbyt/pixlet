package render

import (
	"image"
)

// Column lays out and draws its children vertically (in a column).
//
// By default, a Column is as small as possible, while still holding
// all its children. However, if `expanded` is set, the Column will
// fill all available space vertically. The width of a Column is
// always that of its widest child.
//
// Alignment along the vertical main axis is controlled by passing
// one of the following `main_align` values:
// - `"start"`: place children at the beginning of the column
// - `"end"`: place children at the end of the column
// - `"center"`: place children in the middle of the column
// - `"space_between"`: place equal space between children
// - `"space_evenly"`: equal space between children and before/after first/last child
// - `"space_around"`: equal space between children, and half of that before/after first/last child
//
// Alignment along the horizontal cross axis is controlled by passing
// one of the following `cross_align` values:
// - `"start"`: place children at the left
// - `"end"`: place children at the right
// - `"center"`: place children in the center
//
// DOC(Children): Child widgets to lay out
// DOC(Expanded): Column should expand to fill all available vertical space
// DOC(MainAlign): Alignment along vertical main axis
// DOC(CrossAlign): Alignment along horizontal cross axis
//
// EXAMPLE BEGIN
// render.Column(
//      children=[
//           render.Box(width=10, height=8, color="#a00"),
//           render.Box(width=14, height=6, color="#0a0"),
//           render.Box(width=16, height=4, color="#00a"),
//      ],
// )
// EXAMPLE END
//
// EXAMPLE BEGIN
// render.Column(
//      expanded=True,
//      main_align="space_around",
//      cross_align="center",
//      children=[
//           render.Box(width=10, height=8, color="#a00"),
//           render.Box(width=14, height=6, color="#0a0"),
//           render.Box(width=16, height=4, color="#00a"),
//      ],
// )
// EXAMPLE END
type Column struct {
	Widget

	Children   []Widget `starlark:"children,required"`
	MainAlign  string   `starlark:"main_align"`
	CrossAlign string   `starlark:"cross_align"`
	Expanded   bool
}

func (c Column) Paint(bounds image.Rectangle, frameIdx int) image.Image {
	v := Vector{
		Vertical:   true,
		Children:   c.Children,
		MainAlign:  c.MainAlign,
		CrossAlign: c.CrossAlign,
		Expanded:   c.Expanded,
	}
	return v.Paint(bounds, frameIdx)
}

func (c Column) FrameCount() int {
	return MaxFrameCount(c.Children)
}
