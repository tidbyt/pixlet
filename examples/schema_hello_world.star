load("render.star", "render")
load("schema.star", "schema")

def main(config):
    if config.get("small"):
        msg = render.Text("Hello, World!", font = "CG-pixel-3x5-mono")
    else:
        msg = render.Text("Hello, World!")

    return render.Root(
        child = msg,
    )

def get_schema():
    return schema.Schema(
        version = "1",
        fields = [
            schema.Toggle(
                id = "small",
                name = "Display small text",
                desc = "A toggle to display smaller text.",
                icon = "compress",
                default = False,
            ),
        ],
    )
