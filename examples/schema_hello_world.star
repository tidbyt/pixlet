load("render.star", "render")
load("schema.star", "schema")

DEFAULT = "false"

def main(config):
    small = config.get("small", DEFAULT)
    msg = render.Text("Hello, World!")
    if small == "false":
        msg = render.Text("Hello, World!", font = "CG-pixel-3x5-mono")

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
                default = DEFAULT,
            ),
        ],
    )
