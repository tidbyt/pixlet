package runtime

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"tidbyt.dev/pixlet/render"
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
imgPng = render.Image(src = png_src)
assert(imgPng.src == png_src, "imgPng.src == png_src")
assert(0 < imgPng.size()[0], "0 < imgPng.size()[0]")
assert(0 < imgPng.size()[1], "0 < imgPng.size()[1]")

gif_src = base64.decode("R0lGODlhBQAEAPAAAAAAAAAAACH5BAF7AAAAIf8LTkVUU0NBUEUyLjADAQAAACwAAAAABQAEAAACBgRiaLmLBQAh+QQBewAAACwAAAAABQAEAAACBYRzpqhXACH5BAF7AAAALAAAAAAFAAQAAAIGDG6Qp8wFACH5BAF7AAAALAAAAAAFAAQAAAIGRIBnyMoFADs=")
imgGif = render.Image(src = gif_src)
assert(5 == imgGif.size()[0], "5 == imgGif.size()[0]")
assert(4 == imgGif.size()[1], "4 == imgGif.size()[1]")
assert(1230 == imgGif.delay, "1230 == imgGif.delay")

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

assert(r1.main_align == "space_evenly", 'r1.main_align == "space_evenly"')
assert(r1.cross_align == "center", 'r1.cross_align == "center"')
assert(r1.children[1].main_align == "start", 'r1.children[1].main_align == "start"')
assert(r1.children[1].cross_align == "end", 'r1.children[1].cross_align == "end"')
assert(len(r1.children) == 2, "len(r1.children) == 2")
assert(len(r1.children[1].children) == 2, "len(r1.children[1].children) == 2")

def main():
    return render.Root(child=r1)
`

func TestBigDotStar(t *testing.T) {
	app := &Applet{}
	err := app.Load("big.star", []byte(TestDotStar), nil)
	assert.NoError(t, err)
	screens, err := app.Run(map[string]string{})
	assert.NoError(t, err)
	assert.NotNil(t, screens)
}

func TestBox(t *testing.T) {
	const (
		filename = "test_box.star"
		src      = `
load("render.star", "render")
b = render.Box(
	width = 2,
	height = 1,
	child = render.Box(height=2),
)
def main():
    return render.Root(child=b)
`
	)

	app := &Applet{}
	err := app.Load(filename, []byte(src), nil)
	assert.NoError(t, err)

	b := app.Globals["b"]
	assert.IsType(t, &Box{}, b)

	widget := b.(*Box).AsRenderWidget()
	assert.IsType(t, &render.Box{}, widget)

	box := widget.(*render.Box)
	assert.Equal(t, 2, box.Width)
	assert.Equal(t, 1, box.Height)

	assert.IsType(t, &render.Box{}, box.Child)
	assert.Equal(t, box.Child.(*render.Box).Height, 2)

	assert.Equal(t, image.Rect(0, 0, 2, 1), widget.Paint(image.Rect(0, 0, 64, 32), 0).Bounds())
}

func TestText(t *testing.T) {
	const (
		filename = "test_text.star"
		src      = `
load("render.star", "render")
t = render.Text(
	height = 10,
	content = "hello",
	font = render.fonts["6x13"],
	color = "#ffffff",
)
def main():
    return render.Root(child=t)
`
	)

	app := &Applet{}
	err := app.Load(filename, []byte(src), nil)
	assert.NoError(t, err)

	txt := app.Globals["t"]
	assert.IsType(t, &Text{}, txt)

	widget := txt.(*Text).AsRenderWidget()
	assert.IsType(t, &render.Text{}, widget)

	text := widget.(*render.Text)
	assert.Equal(t, 10, text.Height)
	assert.Equal(t, "hello", text.Content)

	eR, eG, eB, eA := color.White.RGBA()
	r, g, b, a := text.Color.RGBA()
	assert.Equal(t, []uint32{eR, eG, eB, eA}, []uint32{r, g, b, a})

	rendered := widget.Paint(image.Rect(0, 0, 64, 32), 0)
	assert.Greater(t, rendered.Bounds().Dx(), 0)
	assert.Equal(t, text.Height, rendered.Bounds().Dy())
}

func TestImage(t *testing.T) {
	// create a new PNG with a single blue pixel
	bounds := image.Rect(0, 0, 64, 32)
	blue := color.NRGBA{0, 0, 255, 255}

	im := image.NewRGBA(bounds)
	im.Set(12, 12, blue)

	var p bytes.Buffer
	require.NoError(t, png.Encode(&p, im))

	const filename = "test_png.star"
	src := fmt.Sprintf(`
load("render.star", "render")
load("encoding/base64.star", "base64")

img = render.Image(src = base64.decode("%s"))
def main():
    return render.Root(child=img)

`, base64.StdEncoding.EncodeToString(p.Bytes()))

	app := &Applet{}
	err := app.Load(filename, []byte(src), nil)
	assert.NoError(t, err)

	starlarkP := app.Globals["img"]
	require.IsType(t, &Image{}, starlarkP)

	actualIm := starlarkP.(*Image).AsRenderWidget().Paint(image.Rect(0, 0, 64, 32), 0)
	assert.Equal(t, bounds, actualIm.Bounds())
	assert.Equal(t, blue, actualIm.At(12, 12))
}
