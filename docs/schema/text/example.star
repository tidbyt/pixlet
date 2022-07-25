load("render.star", "render")
load("schema.star", "schema")

def main(config):
    msg = config.get("msg", "Hello")
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
            schema.Text(
                id = "msg",
                name = "Message",
                desc = "A message to display.",
                icon = "gear",
                default = "Hello",
            ),
        ],
    )
