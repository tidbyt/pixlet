"""
Applet: CMS
Summary: CMS Backend
Description: An app that supports creating a CMS display.
Author: Tidbyt
"""

load("encoding/base64.star", "base64")
load("render.star", "render")
load("schema.star", "schema")

DEFAULT_DISPLAY_TYPE = "horizontal"

DEFAULT_BACKGROUND_COLOR = "#000000"

DEFAULT_HEADING = "Example"
DEFAULT_HEADING_COLOR = "#7AB0FF"
DEFAULT_HIDE_HEADING = False

DEFAULT_BODY = "Example body text that is a bit longer than the heading text."
DEFAULT_BODY_COLOR = "#FFFFFF"
DEFAULT_HIDE_BODY = False

DEFAULT_HIDE_LOGO = False
DEFAULT_LOGO = """
iVBORw0KGgoAAAANSUhEUgAAABAAAAASCAYAAABSO15qAAAACXBIWXMAAC4jAAAuIwF4pT92AAABFklE
QVQ4jZXTzU7CQBQF4K+kLBCRRCNCqkRf06fyPfxjQ2JcGk1cqRElYaF10ZlQYGjwJE2nt+eeOXPn3uzy
qtSAz/A+2EbIt8Q76KMXvgu8Y75ObCWSj/AQnogpHjFoEshwjAnGwUFEP7iY4CRwNwRK3OA04SpihNvA
XRHIgvK4ITmiCEJZXaCD67pyA8rA7bK8hRbOdkiObou4eXRQ4mNHAar+KOsC33azH1FiVhc4lu6JbWhh
GBcD3Gto1wT2cYdhjme0/5FMVcgRnnL8YIHf8LPEm82aZDi07MIM7Vx1/myN3EvEWE5nRJnjK0GcBkfR
RaY6ZrFOTI1zVzWRe2vxharQKy5SVzfHBV5Udz3DK87DegV/vJUy/l4a2gQAAAAASUVORK5CYII=
"""

def main(config):
    display_type = config.str("display_type", DEFAULT_DISPLAY_TYPE)

    if display_type == "horizontal":
        return render_horizontal(config)

    return render_vertical(config)

def render_horizontal(config):
    background_color = config.str("background_color", DEFAULT_BACKGROUND_COLOR)

    heading = config.str("heading", DEFAULT_HEADING)
    heading_color = config.str("heading_color", DEFAULT_HEADING_COLOR)
    hide_heading = config.bool("hide_heading", DEFAULT_HIDE_HEADING)

    body = config.str("body", DEFAULT_BODY)
    body_color = config.str("body_color", DEFAULT_BODY_COLOR)
    hide_body = config.bool("hide_body", DEFAULT_HIDE_BODY)

    logo = config.str("logo", DEFAULT_LOGO)
    hide_logo = config.bool("hide_logo", DEFAULT_HIDE_LOGO)

    row_one = []
    if not hide_logo:
        row_one.append(
            render.Image(
                src = base64.decode(logo),
            ),
        )

    if not hide_heading:
        row_one.append(
            render.Text(
                content = heading,
                color = heading_color,
            ),
        )

    column_one = []
    if len(row_one) != 0:
        column_one.append(
            render.Row(
                expanded = True,
                main_align = "space_evenly",
                cross_align = "center",
                children = row_one,
            ),
        )

    if not hide_body:
        column_one.append(
            render.Marquee(
                width = 60,
                child = render.Text(
                    content = body,
                    color = body_color,
                ),
            ),
        )

    if len(column_one) == 0:
        return []

    return render.Root(
        show_full_animation = True,
        child = render.Box(
            padding = 2,
            color = background_color,
            child = render.Column(
                expanded = True,
                main_align = "space_around",
                cross_align = "center",
                children = column_one,
            ),
        ),
    )

def render_vertical(config):
    background_color = config.str("background_color", DEFAULT_BACKGROUND_COLOR)

    heading = config.str("heading", DEFAULT_HEADING)
    heading_color = config.str("heading_color", DEFAULT_HEADING_COLOR)
    hide_heading = config.bool("hide_heading", DEFAULT_HIDE_HEADING)

    body = config.str("body", DEFAULT_BODY)
    body_color = config.str("body_color", DEFAULT_BODY_COLOR)
    hide_body = config.bool("hide_body", DEFAULT_HIDE_BODY)

    logo = config.str("logo", DEFAULT_LOGO)
    hide_logo = config.bool("hide_logo", DEFAULT_HIDE_LOGO)

    children = []
    if not hide_logo:
        children.append(
            render.Row(
                cross_align = "center",
                main_align = "space_around",
                expanded = True,
                children = [
                    render.Image(
                        src = base64.decode(logo),
                    ),
                ],
            ),
        )
        children.append(render.Box(height = 4))

    if not hide_heading:
        children.append(
            render.Row(
                cross_align = "center",
                main_align = "space_around",
                expanded = True,
                children = [
                    render.Text(
                        content = heading,
                        color = heading_color,
                    ),
                ],
            ),
        )
        children.append(render.Box(height = 4))

    if not hide_body:
        children.append(
            render.WrappedText(content = body, color = body_color),
        )
        children.append(render.Box(height = 10, width = 1))

    if len(children) == 0:
        return []

    return render.Root(
        show_full_animation = True,
        delay = 60,
        child = render.Box(
            color = background_color,
            child = render.Marquee(
                height = 32,
                offset_start = 16,
                offset_end = 32,
                scroll_direction = "vertical",
                child = render.Padding(
                    pad = 1,
                    child = render.Column(
                        main_align = "start",
                        children = children,
                    ),
                ),
            ),
        ),
    )

def get_schema():
    options = [
        schema.Option(
            display = "Horizontal Scroll",
            value = DEFAULT_DISPLAY_TYPE,
        ),
        schema.Option(
            display = "Vertical Scroll",
            value = "vertical",
        ),
    ]

    return schema.Schema(
        version = "1",
        fields = [
            schema.Dropdown(
                id = "display_type",
                name = "Display Type",
                desc = "The type of display to use.",
                icon = "display",
                default = DEFAULT_DISPLAY_TYPE,
                options = options,
            ),
            schema.Color(
                id = "background_color",
                name = "Background Color",
                desc = "Background color for the display.",
                icon = "brush",
                default = DEFAULT_BACKGROUND_COLOR,
            ),
            schema.Text(
                id = "heading",
                name = "Heading Text",
                desc = "The heading text to display",
                icon = "textHeight",
                default = DEFAULT_HEADING,
            ),
            schema.Color(
                id = "heading_color",
                name = "Heading Color",
                desc = "Heading text color.",
                icon = "brush",
                default = DEFAULT_HEADING_COLOR,
            ),
            schema.Toggle(
                id = "hide_heading",
                name = "Hide Heading",
                desc = "A toggle to hide the heading text.",
                icon = "eyeSlash",
                default = DEFAULT_HIDE_HEADING,
            ),
            schema.Text(
                id = "body",
                name = "Body Text",
                desc = "The body text to display",
                icon = "textHeight",
                default = DEFAULT_BODY,
            ),
            schema.Color(
                id = "body_color",
                name = "Body Color",
                desc = "Body text color.",
                icon = "brush",
                default = DEFAULT_BODY_COLOR,
            ),
            schema.Toggle(
                id = "hide_body",
                name = "Hide Body",
                desc = "A toggle to hide the body text.",
                icon = "eyeSlash",
                default = DEFAULT_HIDE_BODY,
            ),
            schema.PhotoSelect(
                id = "logo",
                name = "Logo",
                desc = "A logo to display.",
                icon = "photoFilm",
            ),
            schema.Toggle(
                id = "hide_logo",
                name = "Hide Logo",
                desc = "A toggle to hide the logo.",
                icon = "eyeSlash",
                default = DEFAULT_HIDE_LOGO,
            ),
        ],
    )
