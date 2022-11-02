load("render.star", "render")
load("schema.star", "schema")
load("math.star", "math")

def pie(i):
    if i % 2 == 0:
        color = "#fff"
    else:
        color = "#0f0"
    return render.PieChart(
            weights = [1, 2, 3],
            colors=["#f00", "#0f0", "#00f"],
            diameter=32,
        )

def main(config):
    msg = config.get("msg", "Hello")
    return render.Root(
        child = render.Stack(
            [pie(i) for i in range(1)]
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
