load("render.star", "render")
load("schema.star", "schema")
load("encoding/base64.star", "base64")

def main(config):
    encoded = config.get("photo", DEFAULT_PHOTO)
    photo = base64.decode(encoded)
    return render.Root(
        child = render.Image(photo),
    )

def get_schema():
    return schema.Schema(
        version = "1",
        fields = [
            schema.PhotoSelect(
                id = "photo",
                name = "Add Photo",
                desc = "A photo to display.",
                icon = "gear",
            ),
        ],
    )

DEFAULT_PHOTO = """
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
RU5ErkJggg==
"""
