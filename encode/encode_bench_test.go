package encode

import (
	"testing"

	"tidbyt.dev/pixlet/runtime"
)

var BenchmarkDotStar = `
"""benchmark

A mock app that returns a widget tree of what is meant to be average
size.
"""

load("render.star", "render")
load("encoding/base64.star", "base64")

SUNNY_PNG = """
iVBORw0KGgoAAAANSUhEUgAAAA0AAAANCAYAAABy6+R8AAABhGlDQ1BJQ0MgcHJvZmlsZQ
AAKJF9kT1Iw0AcxV/TiiIVBTuIOGSoThZERR21CkWoEGqFVh1MLv2CJg1Jiouj4Fpw8GOx
6uDirKuDqyAIfoA4OTopukiJ/0sKLWI8OO7Hu3uPu3eAUC8zzQqNAZpum6lEXMxkV8XOVw
QRQgTT6JOZZcxJUhK+4+seAb7exXiW/7k/R4+asxgQEIlnmWHaxBvEU5u2wXmfOMKKskp8
Tjxq0gWJH7muePzGueCywDMjZjo1TxwhFgttrLQxK5oa8SRxVNV0yhcyHquctzhr5Spr3p
O/MJzTV5a5TnMICSxiCRJEKKiihDJsxGjVSbGQov24j3/Q9UvkUshVAiPHAirQILt+8D/4
3a2Vnxj3ksJxoOPFcT6Ggc5doFFznO9jx2mcAMFn4Epv+St1YOaT9FpLix4BvdvAxXVLU/
aAyx1g4MmQTdmVgjSFfB54P6NvygL9t0D3mtdbcx+nD0CaukreAAeHwEiBstd93t3V3tu/
Z5r9/QBe53KfIhP12QAAAAZiS0dEAP8A/wD/oL2nkwAAAAlwSFlzAAALEwAACxMBAJqcGA
AAAAd0SU1FB+QDBBQ5HxFglVAAAAAZdEVYdENvbW1lbnQAQ3JlYXRlZCB3aXRoIEdJTVBX
gQ4XAAAAYUlEQVQoz2NgYGD4jw3/f239H5ccEwMDA8P/19YMxABkdQRNRpdnhDJwmsgoeh
TDRiZ8GnA5nYmQAmziLDABbM7ApZmFGMXo/mPCJohLMc6AQFeAyyCS44kBlwYCYqSnPQAb
5W9EvIXnIQAAAABJRU5ErkJggg==
"""

def main(config):
    return render.Root(
        delay = 60,
        child = render.Box(
            color = "#002b36",
            child = render.Marquee(
                height = 32,
                offset_start = 16,
                offset_end = 32,
                scroll_direction = "vertical",
                child = render.Padding(
                    pad = 1,
                    child = render.Box(
                        height = 250,
                        child = render.Column(
                            main_align = "start",
                            expanded = True,
                            children = [
                                render.Row(
                                    cross_align = "center",
                                    main_align = "space_around",
                                    expanded = True,
                                    children = [
                                        render.Image(base64.decode(SUNNY_PNG)),
                                        render.Text("Aug 15", color = "#cb4b16"),
                                    ],
                                ),
                                render.WrappedText("This is a cool title", color = "#b58900"),
                                render.WrappedText("Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat. Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur. Excepteur sint occaecat cupidatat non proident, sunt in culpa qui officia deserunt mollit anim id est laborum.", color = "#93a1a1"),
                            ],
                        ),
                    ),
                ),
            ),
        ),
    )
`

func BenchmarkRunAndRender(b *testing.B) {
	app := &runtime.Applet{}
	err := app.Load("benchmark.star", []byte(BenchmarkDotStar), nil)
	if err != nil {
		b.Error(err)
	}

	config := map[string]string{}
	for i := 0; i < b.N; i++ {
		roots, err := app.Run(config)
		if err != nil {
			b.Error(err)
		}

		webp, err := ScreensFromRoots(roots).EncodeWebP()
		if err != nil {
			b.Error(err)
		}

		if len(webp) == 0 {
			b.Error()
		}
	}
}
