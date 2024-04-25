package manifest_test

import (
	"bytes"
	_ "embed"
	"io"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"tidbyt.dev/pixlet/manifest"
)

//go:embed testdata/source.star
var source []byte

var output string = `---
id: foo-tracker
name: Foo Tracker
summary: Track realtime foo
desc: The foo tracker provides realtime feeds for foo.
author: Tidbyt
`

func TestManifest(t *testing.T) {
	m := manifest.Manifest{
		ID:      "foo-tracker",
		Name:    "Foo Tracker",
		Summary: "Track realtime foo",
		Desc:    "The foo tracker provides realtime feeds for foo.",
		Author:  "Tidbyt",
		Source:  source,
	}

	expected, err := os.ReadFile("testdata/source.star")
	assert.NoError(t, err)
	assert.Equal(t, m.Source, expected)
}

func TestLoadManifest(t *testing.T) {
	p := filepath.Join("testdata", "manifest.yaml")
	f, err := os.Open(p)
	assert.NoError(t, err)
	defer f.Close()

	m, err := manifest.LoadManifest(f)
	assert.NoError(t, err)

	assert.Equal(t, m.ID, "fuzzy-clock")
	assert.Equal(t, m.Name, "Fuzzy Clock")
	assert.Equal(t, m.Author, "Max Timkovich")
	assert.Equal(t, m.Summary, "Human readable time")
	assert.Equal(t, m.Desc, "Display the time in a groovy, human-readable way.")
}

func TestWriteManifest(t *testing.T) {
	m := manifest.Manifest{
		ID:      "foo-tracker",
		Name:    "Foo Tracker",
		Summary: "Track realtime foo",
		Desc:    "The foo tracker provides realtime feeds for foo.",
		Author:  "Tidbyt",
		Source:  source,
	}

	buff := bytes.Buffer{}
	err := m.WriteManifest(&buff)
	assert.NoError(t, err)

	b, err := io.ReadAll(&buff)
	assert.NoError(t, err)

	assert.Equal(t, output, string(b))
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
		got := manifest.GenerateDirName(tc.input)
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
