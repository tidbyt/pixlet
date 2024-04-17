package xpath_test

import (
	"context"
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
   <baz>
      <qux>999</qux>
      <qux>888</qux>
   </baz>
   <baz>
      <qux>777</qux>
   </baz>
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

    n = d.query_node("/foo/baz")
    t = n.query("/qux")
    if t != "999":
        fail(t)

    n = d.query_all_nodes("/foo/baz")
    if len(n) != 2:
        fail(len(n))
    t = n[0].query_all("/qux")
    if len(t) != 2:
        fail(len(t))
    if t[0] != "999":
        fail(t[0])
    if t[1] != "888":
        fail(t[1])
    t = n[1].query_all("/qux")
    if len(t) != 1:
        fail(len(t))
    if t[0] != "777":
        fail(t[0])

    return [r.Root(child=r.Text("1337"))]
`
	app, err := runtime.NewApplet("test.star", []byte(src))
	require.NoError(t, err)
	screens, err := app.Run(context.Background())
	require.NoError(t, err)
	assert.NotNil(t, screens)
}
