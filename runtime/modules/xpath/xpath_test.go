package xpath_test

import (
    "testing"

    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/require"
    "tidbyt.dev/pixlet/runtime"
)

func TestXPath(t *testing.T) {
	src := `
load("render.star", r="render")
load("xpath.star", "xpath")

def main():
    xml = """
<foo>
   <bar>1337</bar>
   <bar>4711</bar>
</foo>
"""

    d = xpath.loads(xml)

    t = d.query("/foo/bar")
    if t != "1337":
        fail(t)

    t = d.query_all("/foo/bar")
    if len(t) != 2:
        fail(len(t))
    if t[0] != "1337":
        fail(t[0])
    if t[1] != "4711":
        fail(t[1])

    t = d.query("/foo/doesntexist")
    if t != None:
        fail(t)

    t = d.query_all("/foo/doesntexist")
    if len(t) != 0:
        fail(t)

    return [r.Root(child=r.Text("1337"))]
`
	app := &runtime.Applet{}
	err := app.Load("test.star", []byte(src), nil)
	require.NoError(t, err)
	screens, err := app.Run(map[string]string{})
	require.NoError(t, err)
	assert.NotNil(t, screens)
}