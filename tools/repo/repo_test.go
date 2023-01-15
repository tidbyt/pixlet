package repo_test

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"tidbyt.dev/pixlet/tools/repo"
)

func TestIsInRepo(t *testing.T) {
	root, err := os.Getwd()
	assert.NoError(t, err)

	tests := map[string]struct {
		repo string
		want bool
	}{
		"Pixlet repo should always be true": {
			repo: "pixlet",
			want: true,
		},
		"Any other repo should always be false": {
			repo: "foo",
			want: false,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got := repo.IsInRepo(root, tc.repo)
			assert.Equal(t, tc.want, got)
		})
	}
}
