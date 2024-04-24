load("render.star", "render")
load("schema.star", "schema")

def main(config):
    message = "Hello, %s!" % config.get("who", "World")

    if config.bool("small"):
        msg = render.Text(message, font = "CG-pixel-3x5-mono")
    else:
        msg = render.Text(message)

    return render.Root(
        child = msg,
    )

def get_schema():
    return schema.Schema(
        version = "1",
        fields = [
            schema.Text(
                id = "who",
                name = "Who?",
                desc = "Who to say hello to.",
                icon = "user",
            ),
            schema.Toggle(
                id = "small",
                name = "Display small text",
                desc = "A toggle to display smaller text.",
                icon = "compress",
                default = False,
            ),
        ],
    )
