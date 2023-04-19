package encode

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"tidbyt.dev/pixlet/render"
)

func TestScreensFromRoots(t *testing.T) {
	// check that widget trees and params are copied correctly
	s := ScreensFromRoots([]render.Root{
		{Child: &render.Text{Content: "tree 0"}},
		{Child: &render.Text{Content: "tree 1"}},
	})
	assert.Equal(t, 2, len(s.roots))
	assert.Equal(t, "tree 0", s.roots[0].Child.(*render.Text).Content)
	assert.Equal(t, "tree 1", s.roots[1].Child.(*render.Text).Content)
	assert.Equal(t, 0, len(s.images))
	assert.Equal(t, int32(50), s.delay)
	assert.Equal(t, int32(0), s.MaxAge)

	// check that delay and maxAge are copied from first root only
	s = ScreensFromRoots([]render.Root{
		{Child: &render.Text{Content: "tree 0"}, Delay: 4711, MaxAge: 42},
		{Child: &render.Text{Content: "tree 1"}, Delay: 31415, MaxAge: 926535},
	})
	assert.Equal(t, 2, len(s.roots))
	assert.Equal(t, int32(4711), s.delay)
	assert.Equal(t, int32(42), s.MaxAge)
}
