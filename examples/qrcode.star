load("cache.star", "cache")
load("encoding/base64.star", "base64")
load("render.star", "render")
load("qrcode.star", "qrcode")

def main(config):
    url = "https://tidbyt.com?utm_source=pixlet_example"

    data = cache.get(url)
    if data == None:
        code = qrcode.generate(
            url = url,
            size = "large",
            color = "#fff",
            background = "#000",
        )
        cache.set(url, base64.encode(code), ttl_seconds = 3600)
    else:
        code = base64.decode(data)

    return render.Root(
        child = render.Padding(
            child = render.Image(src = code),
            pad = 1,
        ),
    )
