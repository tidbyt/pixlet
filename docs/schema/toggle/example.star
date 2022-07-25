load("render.star", "render")
load("schema.star", "schema")

def main(config):
    party_mode = config.bool("party_mode", False)
    if party_mode:
        msg = "Party mode enabled"
    else:
        msg = "Party mode disabled"

    return render.Root(
        child = render.Marquee(
            width = 64,
            child = render.Text(msg),
        ),
    )

def get_schema():
    return schema.Schema(
        version = "1",
        fields = [
            schema.Toggle(
                id = "party_mode",
                name = "Party Mode",
                desc = "A toggle to enable party mode.",
                icon = "gear",
                default = False,
            ),
        ],
    )
