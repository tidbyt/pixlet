package hmac_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"tidbyt.dev/pixlet/runtime"
)

var hmacSource = `
load("hmac.star", "hmac")

def assert(success, message=None):
    if not success:
        fail(message or "assertion failed")

# Assert.

assert(hmac.md5("secret", "helloworld") == "8bd4df4530c3c2cafabf6986740e44bd")
assert(hmac.sha1("secret", "helloworld") == "e92eb69939a8b8c9843a75296714af611c73fb53")
assert(hmac.sha256("secret", "helloworld") == "7a7c2bf41973489be3b318ad2f16c75fc875c340deecb12a3f79b28bb7135c97")

def main():
	return []
`

func TestHmac(t *testing.T) {
	app := &runtime.Applet{}
	err := app.Load("hmac_test.star", []byte(hmacSource), nil)
	assert.NoError(t, err)

	screens, err := app.Run(map[string]string{})
	assert.NoError(t, err)
	assert.NotNil(t, screens)
}
