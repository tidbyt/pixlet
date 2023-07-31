package render

import (
	"image/color"
	"testing"

	"github.com/stretchr/testify/assert"
)

func ParseAndAssertColor(
	t *testing.T,
	scol string,
	expectedR uint8,
	expectedG uint8,
	expectedB uint8,
	expectedA uint8,
) {
	col, err := ParseColor(scol)
	assert.Nil(t, err)

	c, ok := col.(color.NRGBA)
	assert.True(t, ok)

	assert.Equal(t, expectedR, c.R)
	assert.Equal(t, expectedG, c.G)
	assert.Equal(t, expectedB, c.B)
	assert.Equal(t, expectedA, c.A)
}

func TestParseColorRGB(t *testing.T) {
	ParseAndAssertColor(t, "#5ad", 0x55, 0xaa, 0xdd, 0xff)
	ParseAndAssertColor(t, "5ad", 0x55, 0xaa, 0xdd, 0xff)
}

func TestParseColorRGBA(t *testing.T) {
	ParseAndAssertColor(t, "#5ad8", 0x55, 0xaa, 0xdd, 0x88)
	ParseAndAssertColor(t, "5ad8", 0x55, 0xaa, 0xdd, 0x88)
}

func TestParseColorRRGGBB(t *testing.T) {
	ParseAndAssertColor(t, "#257adb", 0x25, 0x7a, 0xdb, 0xff)
	ParseAndAssertColor(t, "257adb", 0x25, 0x7a, 0xdb, 0xff)
}

func TestParseColorRRGGBBAA(t *testing.T) {
	ParseAndAssertColor(t, "#257adb75", 0x25, 0x7a, 0xdb, 0x75)
	ParseAndAssertColor(t, "257adb75", 0x25, 0x7a, 0xdb, 0x75)
}

func TestParseColorBadValue(t *testing.T) {
	_, err := ParseColor("5a")
	assert.Error(t, err)

	_, err = ParseColor("#5a")
	assert.Error(t, err)

	_, err = ParseColor("#5ad8f")
	assert.Error(t, err)

	_, err = ParseColor("5ad8f")
	assert.Error(t, err)

	_, err = ParseColor("#5ad8f33da")
	assert.Error(t, err)

	_, err = ParseColor("#xyz")
	assert.Error(t, err)

	_, err = ParseColor("##abc")
	assert.Error(t, err)
}
