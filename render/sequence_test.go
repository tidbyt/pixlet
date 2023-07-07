package render

import (
  "image"
  "image/color"
  "testing"

  "github.com/stretchr/testify/assert"
)

func TestSequenceOnlyOneFrameAtATime(t *testing.T) {
  seq := Sequence{
    Children: []Widget {
      Box{Width: 3, Height: 3, Color: color.RGBA{0xff, 0, 0, 0xff}},
      Box{Width: 6, Height: 3, Color: color.RGBA{0, 0xff, 0, 0xff}},
      Box{Width: 9, Height: 3, Color: color.RGBA{0, 0, 0xff, 0xff}},
    },
  }

  // Frame 0
  im := PaintWidget(seq, image.Rect(0, 0, 10, 3), 0)
  assert.Equal(t, nil, checkImage([]string{
    "rrr",
    "rrr",
    "rrr",
  }, im))

  // Frame 1
  im = PaintWidget(seq, image.Rect(0, 0, 10, 3), 1)
  assert.Equal(t, nil, checkImage([]string{
    "gggggg",
    "gggggg",
    "gggggg",
  }, im))

  // Frame 2
  im = PaintWidget(seq, image.Rect(0, 0, 10, 3), 2)
  assert.Equal(t, nil, checkImage([]string{
    "bbbbbbbbb",
    "bbbbbbbbb",
    "bbbbbbbbb",
  }, im))
}
