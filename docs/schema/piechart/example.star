load("render.star", "render")
load("schema.star", "schema")

def main(config):
    return render.Root(
        render.PieChart(
            weights = [1, 2, 3],
            colors=["#f00", "#0f0", "#00f"],
            diameter=32,
        )
    )

def get_schema():
    return schema.Schema(
        version = "1",
        fields = [],
    )
