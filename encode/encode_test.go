package encode

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"tidbyt.dev/pixlet/runtime"
)

var TestDotStar = `
load("render.star", "render")
load("encoding/base64.star", "base64")

def assert(success, message=None):
    if not success:
        fail(message or "assertion failed")

# Font tests
assert(render.fonts["6x13"] == "6x13")
assert(render.fonts["Dina_r400-6"] == "Dina_r400-6")

# Box tests
b1 = render.Box(
    width = 64,
    height = 32,
    color = "#000",
)

assert(b1.width == 64)
assert(b1.height == 32)
assert(b1.color == "#000000")

b2 = render.Box(
    child = b1,
)

assert(b2.child == b1)

# Text tests
t1 = render.Text(
    height = 10,
    font = render.fonts["6x13"],
    color = "#fff",
    content = "foo",
)
assert(t1.height == 10)
assert(t1.font == "6x13")
assert(t1.color == "#ffffff")
assert(0 < t1.size()[0])
assert(0 < t1.size()[1])

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

assert(f.child.width == 123)
assert(f.child.child.content == "hello")

# Padding
p = render.Padding(pad=3, child=render.Box(width=1, height=2))
p2 = render.Padding(pad=(1,2,3,4), child=render.Box(width=1, height=2))
p3 = render.Padding(pad=1, child=render.Box(width=1, height=2), expanded=True)

# Image tests
png_src = base64.decode("iVBORw0KGgoAAAANSUhEUgAAAAEAAAABAQMAAAAl21bKAAAAA1BMVEX/AAAZ4gk3AAAACklEQVR4nGNiAAAABgADNjd8qAAAAABJRU5ErkJggg==")
img = render.Image(src = png_src)
assert(img.src == png_src)
assert(0 < img.size()[0])
assert(0 < img.size()[1])

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
            ],
        ),
    ],
)

assert(r1.main_align == "space_evenly")
assert(r1.cross_align == "center")
assert(r1.children[1].main_align == "start")
assert(r1.children[1].cross_align == "end")
assert(len(r1.children) == 2)
assert(len(r1.children[1].children) == 2)

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
