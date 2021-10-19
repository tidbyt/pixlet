load("render.star", "render")

def main(config):
    font = config.get("font", "tb-8")
    print("Using font: '{}'".format(font))
    return render.Root(
        child = render.Column(
            children = [
                render.Box(
                    width = 64,
                    height = 1,
                    color = "#78DECC",
                ),
                render.Marquee(
                    width = 64,
                    child = render.Text("The quick brown fox jumps over the lazy dog", font = font),
                ),
                render.Box(
                    width = 64,
                    height = 1,
                    color = "#78DECC",
                ),
                render.Marquee(
                    width = 64,
                    child = render.Text("THE QUICK BROWN FOX JUMPS OVER THE LAZY DOG", font = font),
                ),
                render.Box(
                    width = 64,
                    height = 1,
                    color = "#78DECC",
                ),
                render.Marquee(
                    width = 64,
                    child = render.Text("!@#$%^&*()_+:?><~`", font = font),
                ),
                render.Box(
                    width = 64,
                    height = 1,
                    color = "#78DECC",
                ),
            ],
        ),
    )
