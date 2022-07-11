package qrcode

import (
	"fmt"
	"image/color"
	"math/rand"
	"strings"
	"sync"
	"time"

	goqrcode "github.com/skip2/go-qrcode"
	"go.starlark.net/starlark"
	"go.starlark.net/starlarkstruct"
	"tidbyt.dev/pixlet/render"
)

const (
	ModuleName = "qrcode"
)

var (
	once   sync.Once
	module starlark.StringDict
)

func LoadModule() (starlark.StringDict, error) {
	once.Do(func() {
		rand.Seed(time.Now().UnixNano())
		module = starlark.StringDict{
			ModuleName: &starlarkstruct.Module{
				Name: ModuleName,
				Members: starlark.StringDict{
					"generate": starlark.NewBuiltin("generate", generateQRCode),
				},
			},
		}
	})

	return module, nil
}

func generateQRCode(thread *starlark.Thread, _ *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
	var (
		starUrl        starlark.String
		starSize       starlark.String
		starColor      starlark.String
		starBackground starlark.String
	)

	if err := starlark.UnpackArgs(
		"generate",
		args, kwargs,
		"url", &starUrl,
		"size", &starSize,
		"color?", &starColor,
		"background?", &starBackground,
	); err != nil {
		return nil, fmt.Errorf("unpacking arguments for generate: %w", err)
	}

	// Validate size.
	if !contains([]string{"small", "medium", "large"}, starSize.GoString()) {
		return nil, fmt.Errorf("size must be small, medium, or large")
	}

	// Determine QRCode sizing information.
	version := 0
	imgSize := 0
	switch starSize.GoString() {
	case "small":
		version = 1
		imgSize = 21
	case "medium":
		version = 2
		imgSize = 25
	case "large":
		version = 3
		imgSize = 29
	}

	url := starUrl.String()
	code, err := goqrcode.NewWithForcedVersion(url, version, goqrcode.Low)
	if err != nil {
		return nil, err
	}

	// Set default styles.
	code.DisableBorder = true
	code.ForegroundColor = color.White
	code.BackgroundColor = color.Transparent

	// Override color if one is provided.
	if starColor.Len() > 0 {
		color, err := render.ParseColor(starColor.GoString())
		if err != nil {
			return nil, fmt.Errorf("color is not a valid hex string: %s", starColor.String())
		}
		code.ForegroundColor = color
	}

	// Override background if one is provided.
	if starBackground.Len() > 0 {
		background, err := render.ParseColor(starBackground.GoString())
		if err != nil {
			return nil, fmt.Errorf("color is not a valid hex string: %s", starColor.String())
		}
		code.BackgroundColor = background
	}

	png, err := code.PNG(imgSize)
	if err != nil {
		return nil, err
	}

	return starlark.String(string(png)), nil
}

func contains(slice []string, str string) bool {
	for _, item := range slice {
		if strings.Compare(item, str) == 0 {
			return true
		}
	}

	return false
}
