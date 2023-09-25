package canvas

import (
	"os"
)

type TextAlign string

const (
	AlignLeft   TextAlign = "left"
	AlignCenter TextAlign = "center"
	AlignRight  TextAlign = "right"

	DefaultCanvasProvider = "gg"
)

type provider func(width, height int) Canvas

var (
	providers = make(map[string]provider)
)

func register(name string, p provider) {
	providers[name] = p
}

func NewCanvas(width, height int) Canvas {
	pc := os.Getenv("PIXLET_CANVAS")

	if pc == "" {
		pc = DefaultCanvasProvider
	}

	if provider, ok := providers[pc]; ok {
		return provider(width, height)
	} else {
		panic("unknown canvas provider: " + pc)
	}
}
