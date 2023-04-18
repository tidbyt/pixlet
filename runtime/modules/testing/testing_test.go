package testing_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"go.starlark.net/starlark"
	"tidbyt.dev/pixlet/runtime"
	runtime_testing "tidbyt.dev/pixlet/runtime/modules/testing"
)

var TestApp = `
load("http.star", "http")
load("render.star", "render")

def main(config):
    response = http.get("https://ifconfig.me")

    if response.status_code != 200:
        ip_address = "ERROR"
    else:
        ip_address = response.body()

    return render.Root(
        child = render.Column(
            main_align = "center",
            cross_align = "center",
            expanded = True,
            children = [
                render.Row(
                    main_align = "center",
                    cross_align = "center",
                    expanded = True,
                    children = [
                        render.Text(
                            content = ip_address,
                            font = "6x13"
                        )
                    ]
                )
            ],
        )
    )
`

var TestAppTests = `
load("testing.star", "testing")
load("assert.star", "assert")

def test_ok(config):
    assert.true(True)
    assert.eq(1, 1)

def test_fail(config):
    assert.true(False)
    assert.eq(1, 2)

def test_unstubbed_http(config):
    root = main(config)

def test_http_success(config):
    testing.stub_http(
        lambda method, url, **kwargs: method == "GET" and url == "https://ifconfig.me",
        lambda method, url, **kwargs: { "status_code": 200, "body": "8.8.8.8" }
    )
    root = main(config)
    testing.assert_render_webp(root, "test_app.test.stub.webp")
`

func TestTesting(t *testing.T) {
	initializers := []runtime.ThreadInitializer{}
	stubbedHttpTransport := runtime_testing.InitHttpStub(false)
	initializers = append(initializers, func(thread *starlark.Thread) *starlark.Thread {
		stubbedHttpTransport.ClearStubs()
		return thread
	})
	app := &runtime.Applet{}
	err := app.Load("test_app.star", []byte(TestApp), runtime.TestingLoader)
	assert.NoError(t, err)

	err = app.LoadTests("test_app.test.star", []byte(TestAppTests), nil)
	assert.NoError(t, err)

	testResult, err := app.RunTests(map[string]string{}, initializers...)
	assert.NoError(t, err)
	assert.NotNil(t, testResult)
	assert.Equal(t, 4, len(testResult))

	for _, res := range testResult {
		if res.FunctionName == "test_ok" {
			assert.True(t, res.Success())
			assert.Nil(t, res.Errors)
		} else if res.FunctionName == "test_fail" {
			assert.False(t, res.Success())
			assert.Equal(t, 2, len(res.Errors))
			assert.Regexp(t, "assertion failed", res.Errors[0].String())
			assert.Regexp(t, "1 != 2", res.Errors[1].String())
		} else if res.FunctionName == "test_unstubbed_http" {
			assert.False(t, res.Success())
			assert.Equal(t, 1, len(res.Errors))
			assert.Regexp(t, "Unstubbed HTTP requests not allowed.", res.Errors[0].String())
		} else if res.FunctionName == "test_http_success" {
			assert.True(t, res.Success())
			assert.Nil(t, res.Errors)
		} else {
			t.Errorf("Unexpected test: %s", res.FunctionName)
		}
	}
}
