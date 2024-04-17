package schema_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"tidbyt.dev/pixlet/runtime"
)

func TestColorSuccess(t *testing.T) {
	src := `
load("schema.star", "schema")

def assert(success, message=None):
    if not success:
        fail(message or "assertion failed")

# no palette, 3 char default
s1 = schema.Color(
    id = "colors",
    name = "Colors",
    desc = "The color to display",
    icon = "brush",
    default = "#fff",
)

assert(s1.id == "colors")
assert(s1.name == "Colors")
assert(s1.desc == "The color to display")
assert(s1.icon == "brush")
assert(s1.default == "#fff")

# with palette
s2 = schema.Color(
    id = "colors",
    name = "Colors",
    desc = "The color to display",
    icon = "brush",
    default = "123456",
    palette = ["#f0f", "#aabbcd", "103", "323334"],
)

assert(s2.id == "colors")
assert(s2.name == "Colors")
assert(s2.desc == "The color to display")
assert(s2.icon == "brush")
assert(s2.default == "#123456")
print(s2.palette)
assert(len(s2.palette) == 4)
assert(s2.palette[0] == "#f0f")
assert(s2.palette[1] == "#aabbcd")
assert(s2.palette[2] == "#103")
assert(s2.palette[3] == "#323334")

def main():
    return []
`
	app, err := runtime.NewApplet("colors.star", []byte(src))
	assert.NoError(t, err)

	screens, err := app.Run(context.Background())
	assert.NoError(t, err)
	assert.NotNil(t, screens)
}

func TestColorMalformedColors(t *testing.T) {
	src := `
load("schema.star", "schema")
load("encoding/json.star", "json")

def assert(success, message=None):
    if not success:
        fail(message or "assertion failed")

def main(config):
    s = schema.Color(
        id = "colors",
        name = "Colors",
        desc = "The color to display",
        icon = "brush",
        default = config["default"],
        palette = json.decode(config["palette"]),
    )

    assert(s.id == "colors")
    assert(s.name == "Colors")
    assert(s.desc == "The color to display")
    assert(s.icon == "brush")

    return []
`
	app, err := runtime.NewApplet("colors.star", []byte(src))
	assert.NoError(t, err)

	// Well formed input -> success
	screens, err := app.RunWithConfig(context.Background(), map[string]string{"default": "#ffaa77", "palette": "[]"})
	assert.NoError(t, err)
	assert.NotNil(t, screens)

	// Bad default
	_, err = app.RunWithConfig(context.Background(), map[string]string{"default": "#nothex", "palette": "[]"})
	assert.Error(t, err)
	_, err = app.RunWithConfig(context.Background(), map[string]string{"default": "", "palette": "[]"})
	assert.Error(t, err)
	_, err = app.RunWithConfig(context.Background(), map[string]string{"default": "0", "palette": "[]"})
	assert.Error(t, err)
	_, err = app.RunWithConfig(context.Background(), map[string]string{"default": "01", "palette": "[]"})
	assert.Error(t, err)
	_, err = app.RunWithConfig(context.Background(), map[string]string{"default": "#01", "palette": "[]"})
	assert.Error(t, err)
	_, err = app.RunWithConfig(context.Background(), map[string]string{"default": "0123", "palette": "[]"})
	assert.Error(t, err)
	_, err = app.RunWithConfig(context.Background(), map[string]string{"default": "#0123", "palette": "[]"})
	assert.Error(t, err)
	_, err = app.RunWithConfig(context.Background(), map[string]string{"default": "01234", "palette": "[]"})
	assert.Error(t, err)
	_, err = app.RunWithConfig(context.Background(), map[string]string{"default": "#01234", "palette": "[]"})
	assert.Error(t, err)
	_, err = app.RunWithConfig(context.Background(), map[string]string{"default": "0123456", "palette": "[]"})
	assert.Error(t, err)
	_, err = app.RunWithConfig(context.Background(), map[string]string{"default": "#0123456", "palette": "[]"})
	assert.Error(t, err)

	// Bad palette
	_, err = app.RunWithConfig(context.Background(), map[string]string{"default": "#ffaa77", "palette": `["nothex"]`})
	assert.Error(t, err)
	_, err = app.RunWithConfig(context.Background(), map[string]string{"default": "#ffaa77", "palette": `["fff", "ffaabb", "0"]`})
	assert.Error(t, err)
	_, err = app.RunWithConfig(context.Background(), map[string]string{"default": "#ffaa77", "palette": `["fff", "ffaabb", "#0f"]`})
	assert.Error(t, err)
	_, err = app.RunWithConfig(context.Background(), map[string]string{"default": "#ffaa77", "palette": `["fff", "ffaabb", "0123"]`})
	assert.Error(t, err)
	_, err = app.RunWithConfig(context.Background(), map[string]string{"default": "#ffaa77", "palette": `["fff", "ffaabb", "0123456"]`})
	assert.Error(t, err)
}
