# Exmaple of how to use Stack with Rows and Columns to control quadrant position of text on the tidbyt.
# by tavdog

load("render.star", "render")

def main():
    return render.Root(
        child = render.Stack(
            children = [

                # column to hold stuff at the top of the screen
                render.Column(
                    main_align = "start",  # this controls position of children, start = top
                    expanded = True,
                    children = [

                        # row to hold text with equal space between them
                        render.Row(
                            main_align = "space_between",
                            expanded = True,
                            children = [
                                render.Text("TL", color = "#044"),
                                render.Text("TR", color = "#440"),
                            ],
                        ),
                    ],
                ),

                #column to hold the stuff at the bottom of the screen
                render.Column(
                    main_align = "end",  # bottom
                    expanded = True,
                    children = [

                        # row to hold text evenly distrubuted accross the row
                        render.Row(
                            main_align = "space_evenly",
                            expanded = True,
                            children = [
                                render.Text("BL", color = "#911"),
                                render.Text("BR", color = "#119"),
                            ],
                        ),
                    ],
                ),

                # column to hold the stuff in the middle (center) of the screen
                render.Column(
                    main_align = "center",
                    expanded = True,
                    children = [
                        # row to hold text with equal space around each item
                        render.Row(
                            main_align = "space_around",
                            expanded = True,
                            children = [
                                render.Text("ML", color = "#235"),
                                render.Text("MR", color = "#532"),
                            ],
                        ),
                    ],
                ),
            ],
        ),
    )
