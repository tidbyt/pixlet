load("render.star", "render")
load("schema.star", "schema")
load("encoding/json.star", "json")

EXAMPLE_LOCATION = """
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
    option = config.get("station", '{"display": "Back Bay", "value": "back_bay"}')
    station = json.decode(option)

    return render.Root(
        child = render.Marquee(
            width = 64,
            child = render.Text(station["display"]),
        ),
    )

def get_stations(location):
    loc = json.decode(location)  # See example location above.
    locality = loc["locality"]

    if locality == "New York":
        return [
            schema.Option(
                display = "Grand Central",
                value = "grand_central",
            ),
            schema.Option(
                display = "34th Street Penn Station",
                value = "34th_street_penn_station",
            ),
        ]
    elif locality == "Philadelphia":
        return [
            schema.Option(
                display = "30th Street Station",
                value = "30th_street_station",
            ),
        ]
    else:
        return [
            schema.Option(
                display = "Back Bay",
                value = "back_bay",
            ),
        ]

def get_schema():
    return schema.Schema(
        version = "1",
        fields = [
            schema.LocationBased(
                id = "station",
                name = "Train Station",
                desc = "A list of train stations based on a location.",
                icon = "train",
                handler = get_stations,
            ),
        ],
    )
