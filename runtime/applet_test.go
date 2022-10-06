package runtime

import (
	"fmt"
	"testing"

	starlibbase64 "github.com/qri-io/starlib/encoding/base64"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.starlark.net/starlark"
)

func TestLoadEmptySrc(t *testing.T) {
	app := &Applet{}
	err := app.Load("test.star", []byte{}, nil)
	assert.Error(t, err)
}

func TestLoadMalformed(t *testing.T) {
	src := "this is not valid starlark"
	app := &Applet{}
	err := app.Load("test.star", []byte(src), nil)
	assert.Error(t, err)
}

func TestLoadMainMustBeFunction(t *testing.T) {
	// This is legal
	src := `
load("render.star", "render")
def main():
    return render.Root(child=render.Box())
`
	app := &Applet{}
	err := app.Load("test.star", []byte(src), nil)
	assert.NoError(t, err)

	// As is this
	src = `
load("render.star", "render")
def main2():
    return render.Root(child=render.Box())

main = main2
`
	app = &Applet{}
	err = app.Load("test.star", []byte(src), nil)
	assert.NoError(t, err)

	// And this (a lambda is a function)
	src = `
load("render.star", "render")
def main2():
    return render.Root(child=render.Box())

main = lambda: main2()
`
	app = &Applet{}
	err = app.Load("test.star", []byte(src), nil)
	assert.NoError(t, err)

	// But not this, because a string is not a function
	src = `
load("render.star", "render")
def main2():
    return render.Root(child=render.Box())

main = "main2"
`
	app = &Applet{}
	err = app.Load("test.star", []byte(src), nil)
	assert.Error(t, err)

	// And not this either, because here main is gone
	src = `
load("render.star", "render")
def main2():
    return render.Root(child=render.Box())
`
	app = &Applet{}
	err = app.Load("test.star", []byte(src), nil)
	assert.Error(t, err)

}

func TestRunMainReturnsFrames(t *testing.T) {
	// This fails when run, because a box is not a frame
	src := `
load("render.star", "render")
def main():
    return [render.Box()]
`
	app := &Applet{}
	err := app.Load("test.star", []byte(src), nil)
	assert.NoError(t, err)
	screens, err := app.Run(map[string]string{})
	assert.Error(t, err)
	assert.Nil(t, screens)

	// But a single frame is ok
	src = `
load("render.star", "render")
def main():
    return render.Root(child=render.Box())
`

	app = &Applet{}
	err = app.Load("test.star", []byte(src), nil)
	assert.NoError(t, err)
	screens, err = app.Run(map[string]string{})
	assert.NoError(t, err)
	assert.NotNil(t, screens)

	// And a list of frames is ok
	src = `
load("render.star", "render")
def main():
    return [render.Root(child=render.Box()), render.Root(child=render.Text("hi"))]
`
	app = &Applet{}
	err = app.Load("test.star", []byte(src), nil)
	assert.NoError(t, err)
	screens, err = app.Run(map[string]string{})
	assert.NoError(t, err)
	assert.NotNil(t, screens)
}

func TestRunMainAcceptsConfig(t *testing.T) {
	config := map[string]string{
		"one":     "1",
		"two":     "2",
		"toggle1": "true",
		"toggle2": "false",
	}

	// It's ok for main() to accept no args at all
	src := `
load("render.star", "render")
def main():
    return render.Root(child=render.Box())
`
	app := &Applet{}
	err := app.Load("test.star", []byte(src), nil)
	assert.NoError(t, err)
	roots, err := app.Run(config)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(roots))

	// And it can accept a (the) config dict
	src = `
load("render.star", "render")

def assert_eq(message, actual, expected):
	if not expected == actual:
		fail(message, "-", "expected", expected, "actual", actual)

def main(config):
	assert_eq("config.get with fallback", config.get("doesnt_exist", "foo"), "foo")

	assert_eq("config.str with fallback", config.str("doesnt_exist", "foo"), "foo")
	assert_eq("config.str non-existent value", config.str("doesnt_exist"), None)

	assert_eq("config.bool with fallback", config.bool("doesnt_exist", True), True)
	assert_eq("config.bool non-existent value", config.bool("doesnt_exist"), None)

	assert_eq("config.bool('toggle1')", config.bool("toggle1"), True)
	assert_eq("config.bool('toggle2')", config.bool("toggle2"), False)

	return [render.Root(child=render.Box()) for _ in range(int(config["one"]) + int(config["two"]))]
`
	app = &Applet{}
	err = app.Load("test.star", []byte(src), nil)
	require.NoError(t, err)
	roots, err = app.Run(config)
	require.NoError(t, err)
	assert.Equal(t, 3, len(roots))
}

func TestModuleLoading(t *testing.T) {
	// Our basic set of modules can be imported
	src := `
load("render.star", "render")
load("encoding/base64.star", "base64")
load("encoding/json.star", "json")
load("http.star", "http")
load("math.star", "math")
load("re.star", "re")
load("time.star", "time")

hello_b64 = "aGVsbG8gdGhlcmU="
hello_json = '{"hello": "there"}'
hello_re = 'he[l]{2}o\\sthere'

def main():
    if base64.decode(hello_b64) != "hello there":
        fail("base64 broken")
    if json.decode(hello_json)["hello"] != "there":
        fail("json broken")
    if http.get == None:
        fail("http broken")
    if math.ceil(3.14159265358979) != 4:
        fail("math broken")
    if re.findall(hello_re, "well hello there friend") != ("hello there",):
        fail("re broken")
    if time.parse_duration("10s").seconds != 10:
        fail("time broken")
    return render.Root(child=render.Box())
`
	app := &Applet{}
	err := app.Load("test.star", []byte(src), nil)
	assert.NoError(t, err)
	roots, err := app.Run(map[string]string{})
	assert.NoError(t, err)
	assert.Equal(t, 1, len(roots))

	// An additional module loader can be added
	loader := func(thread *starlark.Thread, module string) (starlark.StringDict, error) {
		if module == "hello.star" {
			return starlibbase64.LoadModule()
		}
		return nil, fmt.Errorf("invalid module: %s", module)
	}
	src = `
load("render.star", "render")
load("hello.star", "base64")
def main():
    if int(base64.decode("NDI=")) != 42:
        fail("something went wrong")
    return render.Root(child=render.Box())
`
	app = &Applet{}
	err = app.Load("test.star", []byte(src), loader)
	assert.NoError(t, err)
	roots, err = app.Run(map[string]string{})
	assert.NoError(t, err)
	assert.Equal(t, 1, len(roots))

}

func TestThreadInitializer(t *testing.T) {
	src := `
load("render.star", "render")
def main():
	print('foobar')
	return render.Root(child=render.Box())
`
	// override the print function of the thread
	var printedText string
	initializer := func(thread *starlark.Thread) *starlark.Thread {
		thread.Print = func(thread *starlark.Thread, msg string) {
			printedText += msg
		}
		return thread
	}

	app := &Applet{}
	err := app.Load("test.star", []byte(src), nil)
	assert.NoError(t, err)
	_, err = app.Run(map[string]string{}, initializer)
	assert.NoError(t, err)

	// our print function should have been called
	assert.Equal(t, "foobar", printedText)
}

// TODO: test Screens, especially Screens.Render()
