"""
Applet:
Summary:
Description:
Author:
"""

load("render.star", "render")

def main():
    return render.Root(
        child = render.Text("Hello, World!"),
    )
