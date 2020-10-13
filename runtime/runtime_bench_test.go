package runtime

import (
	"testing"
)

var BenchmarkDotStar = `
"""benchmark

A mock app that returns a widget tree of what is meant to be average
size.
"""

load("render.star", r = "render")
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

def main():
    return [r.Frame(
        delay=80,
        root=r.Column(children=[
            r.Marquee(
                width=50,
                child = r.Text(
                    content = "This is a pretty long message to scroll",
                    color = "#fff",
                    font = r.fonts["5x8"],
                    height = 6,
                ),
            ),
            r.Row(
                expanded=True,
                main_align="space_evenly",
                children=[
                    r.Box(width=10, height=7, color="#f00"),
                    r.Box(width=10, height=7, color="#0f0"),
                ]),
            r.Row(children=[
                r.Column(children=[
                    r.Box(width=10, height=5, color="#00f"),
                    r.Box(width=10, height=5, color="#f0f"),
                ]),
                r.PNG(base64.decode(SUNNY_PNG)),
            ]),
        ]),
    )]
`

func BenchmarkRunAndRender(b *testing.B) {
	app := &Applet{}
	err := app.Load("benchmark.star", []byte(BenchmarkDotStar), nil)
	if err != nil {
		b.Error(err)
	}

	config := map[string]string{}
	for i := 0; i < b.N; i++ {
		screens, err := app.Run(config)
		if err != nil {
			b.Error(err)
		}

		webp, err := screens.Render()
		if err != nil {
			b.Error(err)
		}

		if len(webp) == 0 {
			b.Error()
		}
	}
}
