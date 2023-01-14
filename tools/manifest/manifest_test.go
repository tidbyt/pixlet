package manifest_test

import (
	_ "embed"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"tidbyt.dev/pixlet/tools/manifest"
)

//go:embed testdata/source.star
var source []byte

func TestManifest(t *testing.T) {
	m := manifest.Manifest{
		ID:          "foo-tracker",
		Name:        "Foo Tracker",
		Summary:     "Track realtime foo",
		Desc:        "The foo tracker provides realtime feeds for foo.",
		Author:      "Tidbyt",
		FileName:    "foo_tracker.star",
		PackageName: "footracker",
		Source:      source,
	}

	expected, err := os.ReadFile("testdata/source.star")
	assert.NoError(t, err)
	assert.Equal(t, m.Source, expected)
}

func TestGeneratePackageName(t *testing.T) {
	type test struct {
		input string
		want  string
	}

	tests := []test{
		{input: "Cool App", want: "coolapp"},
		{input: "CoolApp", want: "coolapp"},
		{input: "cool-app", want: "coolapp"},
		{input: "cool_app", want: "coolapp"},
	}

	for _, tc := range tests {
		got := manifest.GeneratePackageName(tc.input)
		assert.Equal(t, tc.want, got)
	}
}

func TestGenerateFileName(t *testing.T) {
	type test struct {
		input string
		want  string
	}

	tests := []test{
		{input: "Cool App", want: "cool_app.star"},
		{input: "CoolApp", want: "coolapp.star"},
		{input: "cool-app", want: "cool_app.star"},
		{input: "cool_app", want: "cool_app.star"},
	}

	for _, tc := range tests {
		got := manifest.GenerateFileName(tc.input)
		assert.Equal(t, tc.want, got)
	}
}

func TestGenerateID(t *testing.T) {
	type test struct {
		input string
		want  string
	}

	tests := []test{
		{input: "Cool App", want: "cool-app"},
		{input: "CoolApp", want: "coolapp"},
		{input: "cool-app", want: "cool-app"},
		{input: "cool_app", want: "cool-app"},
	}

	for _, tc := range tests {
		got := manifest.GenerateID(tc.input)
		assert.Equal(t, tc.want, got)
	}
}
