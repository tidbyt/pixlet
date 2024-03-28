package fonts

import (
	"embed"
	"fmt"
	"log"
	"strings"
	"sync"

	"github.com/zachomedia/go-bdf"
)

var (
	//go:embed *.bdf *.ttf
	files embed.FS

	once      sync.Once
	fontCache = map[string]*Font{}
)

type Font struct {
	Name string
	Font *bdf.Font

	// BDF is the raw BDF font data.
	BDF []byte

	// TTF is the raw TTF font data.
	TTF []byte
}

func loadFonts() {
	fontFileInfos, err := files.ReadDir(".")
	if err != nil {
		fmt.Println("files.ReadDir()", err)
		return
	}

	for _, ffi := range fontFileInfos {
		if !strings.HasSuffix(ffi.Name(), ".bdf") {
			continue
		}

		name := strings.TrimSuffix(ffi.Name(), ".bdf")

		bdfBuf, err := files.ReadFile(name + ".bdf")
		if err != nil {
			fmt.Printf("files.ReadFile(): %s\n", err)
			continue
		}

		ttfBuf, err := files.ReadFile(name + ".ttf")
		if err != nil {
			fmt.Printf("files.ReadFile(): %s\n", err)
			continue
		}

		fnt, err := bdf.Parse(bdfBuf)
		if err != nil {
			fmt.Printf("bdf.Parse(%s): %s\n", ffi.Name(), err)
		}

		fontCache[name] = &Font{
			Name: name,
			Font: fnt,
			BDF:  bdfBuf,
			TTF:  ttfBuf,
		}
	}
}

func Names() []string {
	once.Do(loadFonts)

	fontNames := []string{}
	for key := range fontCache {
		fontNames = append(fontNames, key)
	}
	return fontNames
}

func GetFont(name string) *Font {
	once.Do(loadFonts)

	font, ok := fontCache[name]
	if !ok {
		log.Panicf("Unknown font '%s', the available fonts are: %v", name, Names())
	}

	return font
}
