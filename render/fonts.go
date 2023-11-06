package render

//go:generate go run gen/embedfonts.go

import (
	"encoding/base64"
	"fmt"
	"sync"

	"github.com/zachomedia/go-bdf"
	"golang.org/x/image/font"
)

var fontCache = map[string]font.Face{}
var fontMutex = &sync.Mutex{}

func GetFontList() []string {
	fontNames := []string{}
	for key := range fontDataRaw {
		fontNames = append(fontNames, key)
	}
	return fontNames
}

func GetFont(name string) (font.Face, error) {
	fontMutex.Lock()
	defer fontMutex.Unlock()

	if font, ok := fontCache[name]; ok {
		return font, nil
	}

	dataB64, ok := fontDataRaw[name]
	if !ok {
		return nil, fmt.Errorf("unknown font '%s'", name)
	}

	data, err := base64.StdEncoding.DecodeString(dataB64)
	if err != nil {
		return nil, fmt.Errorf("decoding font '%s': %w", name, err)
	}

	f, err := bdf.Parse(data)
	if err != nil {
		return nil, fmt.Errorf("parsing font '%s': %w", name, err)
	}

	fontCache[name] = f.NewFace()
	return fontCache[name], nil
}
