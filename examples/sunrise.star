load("render.star", "render")
load("schema.star", "schema")
load("time.star", "time")
load("encoding/json.star", "json")
load("sunrise.star", "sunrise")

DEFAULT_LOCATION = """
{
	"lat": "40.6781784",
	"lng": "-73.9441579",
	"description": "Brooklyn, NY, USA",
	"locality": "Brooklyn",
	"place_id": "ChIJCSF8lBZEwokRhngABHRcdoI",
	"timezone": "America/New_York"
}
"""

def main(config):
    location = config.get("location", DEFAULT_LOCATION)
    loc = json.decode(location)
    lat, lng = float(loc["lat"]), float(loc["lng"])

    now = time.now()
    rise = sunrise.sunrise(lat, lng, now)
    set = sunrise.sunset(lat, lng, now)

    # Check if the sun does not rise or set today. This would happen if the
    # location of the deivce is close to the north or south pole where there are
    # many days of light or darkness. Maybe someone brought their Tidbyt to the
    # Amundsen-Scott South Pole Station! How cool would that be?
    if rise == None or set == None:
        return render.Root(
            child = render.Column(
                children = [
                    render.Text("Now: %s" % now.in_location(loc["timezone"]).format("3:04 PM")),
                    render.Marquee(
                        width = 64,
                        child = render.Text("Sun doesn't rise or set today."),
                    ),
                ],
            ),
        )

    return render.Root(
        child = render.Column(
            children = [
                render.Text("Now: %s" % now.in_location(loc["timezone"]).format("3:04 PM")),
                render.Text("Rise: %s" % rise.in_location(loc["timezone"]).format("3:04 PM")),
                render.Text("Set: %s" % set.in_location(loc["timezone"]).format("3:04 PM")),
            ],
        ),
    )

def get_schema():
    return schema.Schema(
        version = "1",
        fields = [
            schema.Location(
                id = "location",
                name = "Location",
                desc = "Location for which to display time.",
                icon = "locationDot",
            ),
        ],
    )
