package manifest_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"tidbyt.dev/pixlet/tools/manifest"
)

func TestValidateName(t *testing.T) {
	type test struct {
		input     string
		shouldErr bool
	}

	tests := []test{
		{input: "Cool App", shouldErr: false},
		{input: "Cool app", shouldErr: true},
		{input: "cool app", shouldErr: true},
		{input: "coolApp", shouldErr: true},
		{input: "Really Really Long App Name", shouldErr: true},
		{input: "", shouldErr: true},
		{input: "Clark's App", shouldErr: false},
	}

	for _, tc := range tests {
		err := manifest.ValidateName(tc.input)

		if tc.shouldErr {
			assert.Error(t, err)
		} else {
			assert.NoError(t, err)
		}
	}
}

func TestValidateSummary(t *testing.T) {
	type test struct {
		input     string
		shouldErr bool
	}

	tests := []test{
		{input: "A cool app", shouldErr: false},
		{input: "A really really really cool app", shouldErr: true},
		{input: "A cool app.", shouldErr: true},
		{input: "A cool app!", shouldErr: true},
		{input: "A cool app?", shouldErr: true},
		{input: "a cool app", shouldErr: true},
		{input: "NYC Subway departures", shouldErr: false},
		{input: "", shouldErr: true},
	}

	for _, tc := range tests {
		err := manifest.ValidateSummary(tc.input)

		if tc.shouldErr {
			assert.Error(t, err)
		} else {
			assert.NoError(t, err)
		}
	}
}

func TestValidateDesc(t *testing.T) {
	type test struct {
		input     string
		shouldErr bool
	}

	tests := []test{
		{input: "A really cool app that does really cool app things.", shouldErr: false},
		{input: "a really cool app that does really cool app things.", shouldErr: true},
		{input: "A really cool app that does really cool app things", shouldErr: true},
		{input: "", shouldErr: true},
	}

	for _, tc := range tests {
		err := manifest.ValidateDesc(tc.input)

		if tc.shouldErr {
			assert.Error(t, err)
		} else {
			assert.NoError(t, err)
		}
	}
}

func TestValidateID(t *testing.T) {
	type test struct {
		input     string
		shouldErr bool
	}

	tests := []test{
		{input: "foo-bar", shouldErr: false},
		{input: "foobar", shouldErr: false},
		{input: "FooBar", shouldErr: true},
		{input: "foo$", shouldErr: true},
		{input: "", shouldErr: true},
	}

	for _, tc := range tests {
		err := manifest.ValidateID(tc.input)

		if tc.shouldErr {
			assert.Error(t, err)
		} else {
			assert.NoError(t, err)
		}
	}

}

func TestValidateFileName(t *testing.T) {
	type test struct {
		input     string
		shouldErr bool
	}

	tests := []test{
		{input: "foo_bar.star", shouldErr: false},
		{input: "foo_bar", shouldErr: true},
		{input: "FooBar.star", shouldErr: true},
		{input: "foo$.star", shouldErr: true},
		{input: "", shouldErr: true},
	}

	for _, tc := range tests {
		err := manifest.ValidateFileName(tc.input)

		if tc.shouldErr {
			assert.Error(t, err)
		} else {
			assert.NoError(t, err)
		}
	}

}

func TestValidatePackageName(t *testing.T) {
	type test struct {
		input     string
		shouldErr bool
	}

	tests := []test{
		{input: "foobar", shouldErr: false},
		{input: "foo_bar", shouldErr: true},
		{input: "FooBar", shouldErr: true},
		{input: "foo$", shouldErr: true},
		{input: "", shouldErr: true},
	}

	for _, tc := range tests {
		err := manifest.ValidatePackageName(tc.input)

		if tc.shouldErr {
			assert.Error(t, err)
		} else {
			assert.NoError(t, err)
		}
	}

}
