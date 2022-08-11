package render

import (
	"image"
	"image/color"
	"runtime"
	"sync"

	"github.com/tidbyt/gg"
)

const (
	// DefaultFrameWidth is the normal width for a frame.
	DefaultFrameWidth = 64

	// DefaultFrameHeight is the normal height for a frame.
	DefaultFrameHeight = 32

	// DefaultMaxFrameCount is the default maximum number of frames to render.
	DefaultMaxFrameCount = 2000
)

// Every Widget tree has a Root.
//
// The child widget, and all its descendants, will be drawn on a 64x32
// canvas. Root places its child in the upper left corner of the
// canvas.
//
// If the tree contains animated widgets, the resulting animation will
// run with _delay_ milliseconds per frame.
//
// If the tree holds time sensitive information which must never be
// displayed past a certain point in time, pass _MaxAge_ to specify
// an expiration time in seconds. Display devices use this to avoid
// displaying stale data in the event of e.g. connectivity issues.
//
// DOC(Child): Widget to render
// DOC(Delay): Frame delay in milliseconds
// DOC(MaxAge): Expiration time in seconds
//
type Root struct {
	Child  Widget `starlark:"child,required"`
	Delay  int32  `starlark:"delay"`
	MaxAge int32  `starlark:"max_age"`

	maxParallelFrames int
	maxFrameCount     int
}

type RootPaintOption func(*Root)

// WithMaxParallelFrames sets the maximum number of frames that will
// be painted in parallel.
//
// By default, only `runtime.NumCPU()` frames are painted in parallel.
// Higher parallelism consumes more memory, and doesn't usually make
// sense since painting is CPU-bouond.
func WithMaxParallelFrames(max int) RootPaintOption {
	return func(r *Root) {
		r.maxParallelFrames = max
	}
}

// WithMaxFrameCount sets the maximum number of frames that will be
// rendered when calling `Paint`.
func WithMaxFrameCount(max int) RootPaintOption {
	return func(r *Root) {
		r.maxFrameCount = max
	}
}

// Paint renders the child widget onto the frame. It doesn't do
// any resizing or alignment.
func (r Root) Paint(solidBackground bool, opts ...RootPaintOption) []image.Image {
	for _, opt := range opts {
		opt(&r)
	}

	if r.maxFrameCount <= 0 {
		r.maxFrameCount = DefaultMaxFrameCount
	}

	numFrames := r.Child.FrameCount()
	if numFrames > r.maxFrameCount {
		numFrames = r.maxFrameCount
	}

	frames := make([]image.Image, numFrames)

	parallelism := r.maxParallelFrames
	if parallelism <= 0 {
		parallelism = runtime.NumCPU()
	}

	var wg sync.WaitGroup
	sem := make(chan bool, parallelism)
	for i := 0; i < numFrames; i++ {
		wg.Add(1)
		sem <- true

		go func(i int) {
			defer func() {
				<-sem
				wg.Done()
			}()

			dc := gg.NewContext(DefaultFrameWidth, DefaultFrameHeight)
			if solidBackground {
				dc.SetColor(color.Black)
				dc.Clear()
			}

			dc.Push()
			r.Child.Paint(dc, image.Rect(0, 0, DefaultFrameWidth, DefaultFrameHeight), i)
			dc.Pop()
			frames[i] = dc.Image()
		}(i)
	}

	wg.Wait()
	return frames
}

// PaintRoots draws >=1 Roots which must all have the same dimensions.
func PaintRoots(solidBackground bool, roots ...Root) []image.Image {
	var images []image.Image
	for _, r := range roots {
		images = append(images, r.Paint(solidBackground)...)
	}

	return images
}
