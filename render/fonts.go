package render

//go:generate go run gen/embedfonts.go

import (
	"encoding/base64"
	"log"

	"github.com/zachomedia/go-bdf"
	"golang.org/x/image/font"
)

var fontCache = map[string]font.Face{}

func GetFontList() []string {
	fontNames := []string{}
	for key := range fontDataRaw {
		fontNames = append(fontNames, key)
	}
	return fontNames
}

func GetFont(name string) font.Face {
	if font, ok := fontCache[name]; ok {
		return font
	}

	dataB64, ok := fontDataRaw[name]
	if !ok {
		log.Panicf("Unknown font '%s', the available fonts are: %v", name, GetFontList())
	}

	data, err := base64.StdEncoding.DecodeString(dataB64)
	if err != nil {
		log.Panicf("couldn't decode %s: %s", name, err)
	}

	f, err := bdf.Parse(data)
	if err != nil {
		log.Panicf("couldn't parse %s: %s", name, err)
	}

	fontCache[name] = f.NewFace()
	return fontCache[name]
}
