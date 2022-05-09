package animation

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func AssertInterpolate(
	t *testing.T,
	expected Transform,
	from Transform,
	to Transform,
	progress float64,
) {
	actual, ok := from.Interpolate(to, progress)
	assert.True(t, ok)
	assert.Equal(t, expected, actual)
}
