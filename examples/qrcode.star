load("cache.star", "cache")
load("render.star", "render")
load("qrcode.star", "qrcode")

def main(config):
    url = "https://tidbyt.com?utm_source=pixlet_example"
    code = cache.get(url)
    if code == None:
        code = qrcode.generate(
            url = url,
            size = "large",
            color = "#fff",
            background = "#000",
        )

    return render.Root(
        child = render.Padding(
            child = render.Image(src = code),
            pad = 1,
        ),
    )
