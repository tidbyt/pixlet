load("render.star", "render")
load("schema.star", "schema")

def main(config):
    pet = config.get("pet", "turtle")
    children = [
        render.Marquee(
            width = 64,
            child = render.Text("pet: %s" % pet),
        ),
    ]

    if pet == "dog":
        has_leash = config.bool("leash", False)
        children.append(
            render.Marquee(
                width = 64,
                child = render.Text("has leash: %s" % has_leash),
            ),
        )

    if pet == "cat":
        has_litter_box = config.bool("litter-box", False)
        children.append(
            render.Marquee(
                width = 64,
                child = render.Text("has litter box: %s" % has_litter_box),
            ),
        )

    return render.Root(
        child = render.Column(
            children = children,
        ),
    )

def more_options(pet):
    if pet == "dog":
        return [
            schema.Toggle(
                id = "leash",
                name = "Leash",
                desc = "A toggle to enable a dog leash.",
                icon = "gear",
                default = False,
            ),
        ]
    elif pet == "cat":
        return [
            schema.Toggle(
                id = "litter-box",
                name = "Litter Box",
                desc = "A toggle to enable a litter box.",
                icon = "gear",
                default = False,
            ),
        ]
    else:
        return []

def get_schema():
    return schema.Schema(
        version = "1",
        fields = [
            schema.Text(
                id = "pet",
                name = "Pet Type",
                desc = "What type of pet do you have?",
                icon = "gear",
            ),
            schema.Generated(
                id = "generated",
                source = "pet",
                handler = more_options,
            ),
        ],
    )
