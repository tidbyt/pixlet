package qrcode_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"tidbyt.dev/pixlet/runtime"
)

var qrCodeSource = `
load("qrcode.star", "qrcode")

def assert(success, message=None):
    if not success:
        fail(message or "assertion failed")


url = "https://tidbyt.com?utm_source=pixlet_example"
code = qrcode.generate(
    url = url,
    size = "large",
    color = "#fff",
    background = "#000",
)

def main():
	return []
`

func TestQRCode(t *testing.T) {
	app := &runtime.Applet{}
	err := app.Load("testid", "test.star", []byte(qrCodeSource), nil)
	assert.NoError(t, err)

	screens, err := app.Run(map[string]string{})
	assert.NoError(t, err)
	assert.NotNil(t, screens)
}
