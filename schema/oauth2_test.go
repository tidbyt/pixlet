package schema_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"tidbyt.dev/pixlet/runtime"
)

var oauth2Source = `
load("encoding/json.star", "json")
load("schema.star", "schema")

def assert(success, message = None):
    if not success:
        fail(message or "assertion failed")

def oauth_handler(params):
    params = json.decode(params)
    return "foobar123"

t = schema.OAuth2(
    id = "auth",
    name = "GitHub",
    desc = "Connect your GitHub account.",
    icon = "github",
    handler = oauth_handler,
    client_id = "the-oauth2-client-id",
    authorization_endpoint = "https://example.com/",
    scopes = [
        "read:user",
    ],
)

assert(t.id == "auth")
assert(t.name == "GitHub")
assert(t.desc == "Connect your GitHub account.")
assert(t.icon == "github")
assert(t.handler("{}") == "foobar123")
assert(t.client_id == "the-oauth2-client-id")
assert(t.authorization_endpoint == "https://example.com/")
assert(t.scopes == ["read:user"])

def main():
    return []

`

func TestOAuth2(t *testing.T) {
	app := &runtime.Applet{}
	err := app.Load("oauth2.star", []byte(oauth2Source), nil)
	assert.NoError(t, err)

	screens, err := app.Run(map[string]string{})
	assert.NoError(t, err)
	assert.NotNil(t, screens)
}
