load("render.star", "render")
load("schema.star", "schema")

def main(config):
    color = config.get("color", "#BFEDC4")

    return render.Root(
        child = render.Marquee(
            width = 64,
            child = render.Text("Text color", color = color),
        ),
    )

def get_schema():
    options = [
        schema.Option(
            display = "Pink",
            value = "#FF94FF",
        ),
        schema.Option(
            display = "Mustard",
            value = "#FFD10D",
        ),
    ]

    return schema.Schema(
        version = "1",
        fields = [
            schema.Dropdown(
                id = "color",
                name = "Text Color",
                desc = "The color of text to be displayed.",
                icon = "brush",
                default = options[0].value,
                options = options,
            ),
        ],
    )
