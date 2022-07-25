load("encoding/json.star", "json")
load("render.star", "render")
load("schema.star", "schema")

def main(config):
    option = config.get("search", '{"display": "Blueberry", "value": "blueberry"}')
    fruit = json.decode(option)

    return render.Root(
        child = render.Marquee(
            width = 64,
            child = render.Text(fruit["display"]),
        ),
    )

def search(pattern):
    if pattern.startswith("a"):
        return [
            schema.Option(
                display = "Apple",
                value = "apple",
            ),
            schema.Option(
                display = "Apricot",
                value = "apricot",
            ),
        ]
    else:
        return []

def get_schema():
    return schema.Schema(
        version = "1",
        fields = [
            schema.Typeahead(
                id = "search",
                name = "Search",
                desc = "A list of items that match search.",
                icon = "gear",
                handler = search,
            ),
        ],
    )
