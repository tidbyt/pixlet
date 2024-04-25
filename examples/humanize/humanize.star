load("humanize.star", "humanize")
load("render.star", "render")
load("schema.star", "schema")
load("time.star", "time")

DEFAULT_COUNTER = "1337"
DEFAULT_APPS = "42"
DEFAULT_TIMEZONE = "America/New_York"

def main(config):
    tz = config.get("$tz", DEFAULT_TIMEZONE)
    num_apps = config.get("num_apps", DEFAULT_APPS)
    now = time.now()

    return render.Root(
        child = render.Column(
            children = [
                render.Text(" %s rated" % humanize.plural(int(num_apps), "app")),
                render.Text(" Comma: %s" % humanize.comma(int(config.get("count", "1337")))),
            ],
        ),
    )

def get_schema():
    return schema.Schema(
        version = "1",
        fields = [
            schema.Text(
                id = "count",
                name = "Count",
                desc = "A cool counter that has comma separators",
                icon = "number",
                default = DEFAULT_COUNTER,
            ),
            schema.Text(
                id = "num_apps",
                name = "How many apps do you want?",
                desc = "The number of apps",
                icon = "number",
                default = DEFAULT_APPS,
            ),
        ],
    )
