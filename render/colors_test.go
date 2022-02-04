package render

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func ParseAndAssertColor(
	t *testing.T,
	scol string,
	expectedR uint32,
	expectedG uint32,
	expectedB uint32,
	expectedA uint32,
) {
	c, err := ParseColor(scol)
	assert.Nil(t, err)
	r, g, b, a := c.RGBA()
	assert.Equal(t, expectedR, r)
	assert.Equal(t, expectedG, g)
	assert.Equal(t, expectedB, b)
	assert.Equal(t, expectedA, a)
}

func TestParseColorRGB(t *testing.T) {
	ParseAndAssertColor(t, "#5ad", 0x5555, 0xaaaa, 0xdddd, 0xffff)
}

func TestParseColorRGBA(t *testing.T) {
	ParseAndAssertColor(t, "#5ad8", 0x2d82, 0x5b05, 0x7653, 0x8888)
}

func TestParseColorRRGGBB(t *testing.T) {
	ParseAndAssertColor(t, "#257adb", 0x2525, 0x7a7a, 0xdbdb, 0xffff)
}

func TestParseColorRRGGBBAA(t *testing.T) {
	ParseAndAssertColor(t, "#257adb75", 0x110a, 0x3831, 0x64df, 0x7575)
}

func TestParseColorBadValue(t *testing.T) {
	_, err := ParseColor("5ad")
	assert.NotNil(t, err)

	_, err = ParseColor("#5a")
	assert.NotNil(t, err)

	_, err = ParseColor("#5ad8f")
	assert.NotNil(t, err)

	_, err = ParseColor("#5ad8f33da")
	assert.NotNil(t, err)

	_, err = ParseColor("#xyz")
	assert.NotNil(t, err)
}
