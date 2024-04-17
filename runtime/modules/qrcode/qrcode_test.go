package qrcode_test

import (
	"context"
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
	app, err := runtime.NewApplet("test.star", []byte(qrCodeSource))
	assert.NoError(t, err)

	screens, err := app.Run(context.Background())
	assert.NoError(t, err)
	assert.NotNil(t, screens)
}
