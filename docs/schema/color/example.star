load("render.star", "render")
load("schema.star", "schema")

DEFAULT_COLOR = "#FF59FF"

def main(config):
    color = config.str("color", DEFAULT_COLOR)

    return render.Root(
        child = render.Box(
            width = 64,
            height = 32,
            color = color,
        ),
    )

def get_schema():
    return schema.Schema(
        version = "1",
        fields = [
            schema.Color(
                id = "color",
                name = "Color",
                desc = "Color of the screen.",
                icon = "brush",
                default = DEFAULT_COLOR,
                palette = [
                    DEFAULT_COLOR,
                    "#7AB0FF",
                    "#BFEDC4",
                    "#78DECC",
                    "#DBB5FF",
                ],
            ),
        ],
    )
