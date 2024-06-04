package render

import (
	"encoding/json"
	"fmt"
	"image"

	"github.com/tidbyt/gg"
)

// A Widget is a self-contained object that can render itself as an image.
type Widget interface {
	// PaintBounds Returns the bounds of the area that will actually be drawn to when Paint() is called
	PaintBounds(bounds image.Rectangle, frameIdx int) image.Rectangle
	Paint(dc *gg.Context, bounds image.Rectangle, frameIdx int)
	FrameCount() int
}

// Widgets can require initialization
type WidgetWithInit interface {
	Init() error
}

// WidgetStaticSize has inherent size and width known before painting.
type WidgetStaticSize interface {
	Size() (int, int)
}

// Computes a (mod m). Useful for handling frameIdx > num available
// frames in Widget.Paint()
func ModInt(a, m int) int {
	a = a % m
	if a < 0 {
		a += m
	}
	return a
}

// Computes the maximum frame count of a slice of widgets.
func MaxFrameCount(widgets []Widget) int {
	m := 1

	for _, w := range widgets {
		if c := w.FrameCount(); c > m {
			m = c
		}
	}

	return m
}

func UnmarshalWidgetJSON(data []byte) (Widget, error) {
	if string(data) == "null" {
		return nil, nil
	}

	var widgetType struct {
		Type string `json:"type"`
	}
	if err := json.Unmarshal(data, &widgetType); err != nil {
		return nil, err
	}

	var widget Widget
	switch widgetType.Type {
	case "Animation":
		widget = &Animation{}
	case "Box":
		widget = &Box{}
	case "Column":
		widget = &Column{}
	case "Image":
		widget = &Image{}
	case "Marquee":
		widget = &Marquee{}
	case "Row":
		widget = &Row{}
	case "Text":
		widget = &Text{}
	case "Vector":
		widget = &Vector{}

	default:
		return nil, fmt.Errorf("unknown widget type: %s", widgetType.Type)
	}

	if err := json.Unmarshal(data, widget); err != nil {
		return nil, err
	}

	return widget, nil
}
