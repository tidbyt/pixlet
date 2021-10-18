package render

//go:generate go run gen/embedfonts.go

import (
	"encoding/base64"
	"log"

	"github.com/zachomedia/go-bdf"
	"golang.org/x/image/font"
)

var Font = map[string]font.Face{}

func init() {
	for name, dataB64 := range FontDataRaw {
		data, err := base64.StdEncoding.DecodeString(dataB64)
		if err != nil {
			log.Printf("couldn't decode %s: %s", name, err)
			continue
		}

		f, err := bdf.Parse(data)
		if err != nil {
			log.Printf("couldn't parse %s: %s", name, err)
			continue
		}

		Font[name] = f.NewFace()
	}
}
