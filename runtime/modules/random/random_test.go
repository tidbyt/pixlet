package random_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"tidbyt.dev/pixlet/runtime"
)

var randomSrc = `
load("random.star", "random")

min = 100
max = 200

def run_test():
	for x in range(0, 100):
		num = random.number(min, max)
		if num < min:
			fail("random number less then min")
		if num > max:
			fail("random number greater then max")

run_test()

def main():
	return []
`

func TestRandom(t *testing.T) {
	app := &runtime.Applet{}
	err := app.Load("random_test.star", []byte(randomSrc), nil)
	assert.NoError(t, err)

	screens, err := app.Run(map[string]string{})
	assert.NoError(t, err)
	assert.NotNil(t, screens)
}
