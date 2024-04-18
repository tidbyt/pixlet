package runtime

import (
	"testing"
	"testing/fstest"

	"github.com/stretchr/testify/require"
)

func TestReadFile(t *testing.T) {
	src := `
load("hello.txt", hello = "file")

def assert_eq(message, actual, expected):
	if not expected == actual:
		fail(message, "-", "expected", expected, "actual", actual)

def test_read():
	f = hello.open()
	assert_eq("read", f.read(), "hello world")

def test_read_binary():
	f = hello.open(mode="rb")
	assert_eq("read_binary", f.read(), b"hello world")

def test_readall():
	assert_eq("readall", hello.readall(), "hello world")

def test_readall_binary():
	assert_eq("readall_binary", hello.readall("rb"), b"hello world")

def main():
	pass

`

	helloTxt := `hello world`

	vfs := &fstest.MapFS{
		"main.star": {Data: []byte(src)},
		"hello.txt": {Data: []byte(helloTxt)},
	}

	app, err := NewAppletFromFS("test_read_file", vfs)
	require.NoError(t, err)
	app.RunTests(t)
}
