load("render.star", "render")

def main():
    return render.Root(
        child = render.Box(
            color = "#ff0000",
            child = render.Box(
                height = 16,
                width = 32,
                color = "#00ff00",
            )
        ),
    )
