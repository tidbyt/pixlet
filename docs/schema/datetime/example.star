load("render.star", "render")
load("schema.star", "schema")
load("time.star", "time")

def main(config):
    user_configured = config.get("event_time", "2022-02-02T20:00:00Z")
    event_time = time.parse_time(user_configured).in_location("America/New_York")

    return render.Root(
        child = render.Marquee(
            width = 64,
            child = render.Text(event_time.format("3:04 PM")),
        ),
    )

def get_schema():
    return schema.Schema(
        version = "1",
        fields = [
            schema.DateTime(
                id = "event_time",
                name = "Event Time",
                desc = "The time of the event.",
                icon = "gear",
            ),
        ],
    )
