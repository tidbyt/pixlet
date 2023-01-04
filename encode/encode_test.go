package encode

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"tidbyt.dev/pixlet/render"
	"tidbyt.dev/pixlet/runtime"
)

var TestDotStar = `
load("render.star", "render")
load("encoding/base64.star", "base64")

def assert(success, message=None):
    if not success:
        fail(message or "assertion failed")

# Font tests
assert(render.fonts["6x13"] == "6x13", 'render.fonts["6x13"] == "6x13"')
assert(render.fonts["Dina_r400-6"] == "Dina_r400-6", 'render.fonts["Dina_r400-6"] == "Dina_r400-6"')

# Box tests
b1 = render.Box(
    width = 64,
    height = 32,
    color = "#000",
)

assert(b1.width == 64, "b1.width == 64")
assert(b1.height == 32, "b1.height == 32")
assert(b1.color == "#000", 'b1.color == "#000"')

b2 = render.Box(
    child = b1,
    color = "#0f0d",
)

assert(b2.child == b1, "b2.child == b1")
assert(b2.color == "#0f0d", 'b2.color == "#0f0d"')

# Text tests
t1 = render.Text(
    height = 10,
    font = render.fonts["6x13"],
    color = "#fff",
    content = "foo",
)
assert(t1.height == 10, "t1.height == 10")
assert(t1.font == "6x13", 't1.font == "6x13"')
assert(t1.color == "#fff", 't1.color == "#fff"')
assert(0 < t1.size()[0], "0 < t1.size()[0]")
assert(0 < t1.size()[1], "0 < t1.size()[1]")

# WrappedText
tw = render.WrappedText(
    height = 16,
    width = 64,
    font = render.fonts["6x13"],
    color = "#f00",
    content = "hey ho foo bar wrap this line it's very long wrap it please",
)

# Root tests
f = render.Root(
    child = render.Box(
        width = 123,
        child = render.Text(
            content = "hello",
        ),
    ),
)

assert(f.child.width == 123, "f.child.width == 123")
assert(f.child.child.content == "hello", 'f.child.child.content == "hello"')

# Padding
p = render.Padding(pad=3, child=render.Box(width=1, height=2))
p2 = render.Padding(pad=(1,2,3,4), child=render.Box(width=1, height=2))
p3 = render.Padding(pad=1, child=render.Box(width=1, height=2), expanded=True)

# Image tests
png_src = base64.decode("iVBORw0KGgoAAAANSUhEUgAAAAEAAAABAQMAAAAl21bKAAAAA1BMVEX/AAAZ4gk3AAAACklEQVR4nGNiAAAABgADNjd8qAAAAABJRU5ErkJggg==")
img = render.Image(src = png_src)
assert(img.src == png_src, "img.src == png_src")
assert(0 < img.size()[0], "0 < img.size()[0]")
assert(0 < img.size()[1], "0 < img.size()[1]")

# Row and Column
r1 = render.Row(
    expanded = True,
    main_align = "space_evenly",
    cross_align = "center",
    children = [
        render.Box(width=12, height=14),
        render.Column(
            expanded = True,
            main_align = "start",
            cross_align = "end",
            children = [
                render.Box(width=6, height=7),
                render.Box(width=4, height=5),
                tw,
            ],
        ),
        render.Plot(
            data = [
                (0, 3.35),
                (1, 2.15),
                (2, 2.37),
                (3, -0.31),
                (4, -3.53),
                (5, 1.31),
                (6, -1.3),
                (7, 4.60),
                (8, 3.33),
                (9, 5.92),
            ],
            width = 64,
            height = 17,
            x_lim = (None, 10),
            y_lim = (0, 10),
            color = "#0f0",
            color_inverted = "#f00",
            fill = True,
        ),
    ],
)

assert(r1.main_align == "space_evenly", 'r1.main_align == "space_evenly"')
assert(r1.cross_align == "center", 'r1.cross_align == "center"')
assert(r1.children[1].main_align == "start", 'r1.children[1].main_align == "start"')
assert(r1.children[1].cross_align == "end", 'r1.children[1].cross_align == "end"')
assert(len(r1.children) == 3, "len(r1.children) == 3")
assert(len(r1.children[1].children) == 3, "len(r1.children[1].children) == 3")

def main():
    return render.Root(child=r1)
`

func TestFile(t *testing.T) {
	app := runtime.Applet{}
	err := app.Load("test.star", []byte(TestDotStar), nil)
	assert.NoError(t, err)

	roots, err := app.Run(map[string]string{})
	assert.NoError(t, err)

	webp, err := ScreensFromRoots(roots).EncodeWebP()
	assert.NoError(t, err)
	assert.True(t, len(webp) > 0)
}

func TestHash(t *testing.T) {
	app := runtime.Applet{}
	err := app.Load("test.star", []byte(TestDotStar), nil)
	require.NoError(t, err)

	roots, err := app.Run(map[string]string{})
	require.NoError(t, err)

	// ensure we can calculate a hash
	hash, err := ScreensFromRoots(roots).Hash()
	require.NoError(t, err)
	require.True(t, len(hash) > 0)

	// ensure the hash doesn't change
	for i := 0; i < 20; i++ {
		h, err := ScreensFromRoots(roots).Hash()
		assert.NoError(t, err)
		assert.Equal(t, hash, h)
	}

	// change the app slightly
	modifiedSource := strings.Replace(TestDotStar, "foo bar", "bar foo", 1)
	app2 := runtime.Applet{}
	err = app2.Load("test.star", []byte(modifiedSource), nil)
	require.NoError(t, err)

	roots2, err := app2.Run(map[string]string{})
	require.NoError(t, err)

	// ensure we can calculate a hash on the new app
	hash2, err := ScreensFromRoots(roots2).Hash()

	// ensure hashes are different
	require.NotEqual(t, hash, hash2)
}

func TestHashDelayAndMaxAge(t *testing.T) {
	r := []render.Root{{Child: &render.Text{Content: "derp"}}}

	h1, err := ScreensFromRoots(r).Hash()
	assert.NoError(t, err)
	r[0].MaxAge = 12
	h2, err := ScreensFromRoots(r).Hash()
	assert.NoError(t, err)
	r[0].Delay = 1
	h3, err := ScreensFromRoots(r).Hash()
	assert.NoError(t, err)

	assert.NotEqual(t, h1, h2)
	assert.NotEqual(t, h2, h3)
}

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
