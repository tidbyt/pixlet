load("render.star", "render")
load("assert.star", "assert")
load("encoding/base64.star", "base64")

# Font tests
assert.eq(render.fonts["6x13"], "6x13")
assert.eq(render.fonts["Dina_r400-6"], "Dina_r400-6")

# Box tests
b1 = render.Box(
    width = 64,
    height = 32,
    color = "#000",
)

assert.eq(b1.width, 64)
assert.eq(b1.height, 32)
assert.eq(b1.color, "#000000")

b2 = render.Box(
    child = b1,
)

assert.eq(b2.child, b1)

# Text tests
t1 = render.Text(
    height = 10,
    font = render.fonts["6x13"],
    color = "#fff",
    content = "foo",
)
assert.eq(t1.height, 10)
assert.eq(t1.font, "6x13")
assert.eq(t1.color, "#ffffff")
assert.lt(0, t1.size()[0])
assert.lt(0, t1.size()[1])

# WrappedText
tw = render.WrappedText(
    height = 16,
    width = 64,
    font = render.fonts["6x13"],
    color = "#f00",
    content = "hey ho foo bar wrap this line it's very long wrap it please",
)

# Frame tests
f = render.Frame(
    root = render.Box(
        width = 123,
        child = render.Text(
            content = "hello",
        ),
    ),
)

assert.eq(f.root.width, 123)
assert.eq(f.root.child.content, "hello")

# Padding
p = render.Padding(pad = 3, child = render.Box(width = 1, height = 2))
p2 = render.Padding(pad = (1, 2, 3, 4), child = render.Box(width = 1, height = 2))
p3 = render.Padding(pad = 1, child = render.Box(width = 1, height = 2), expanded = True)

# PNG tests
png_src = base64.decode("iVBORw0KGgoAAAANSUhEUgAAAAEAAAABAQMAAAAl21bKAAAAA1BMVEX/AAAZ4gk3AAAACklEQVR4nGNiAAAABgADNjd8qAAAAABJRU5ErkJggg==")
png = render.PNG(src = png_src)
assert.eq(png.src, png_src)
assert.lt(0, png.size()[0])
assert.lt(0, png.size()[1])

# Row and Column
r1 = render.Row(
    expanded = True,
    main_align = "space_evenly",
    cross_align = "center",
    children = [
        render.Box(width = 12, height = 14),
        render.Column(
            expanded = True,
            main_align = "start",
            cross_align = "end",
            children = [
                render.Box(width = 6, height = 7),
                render.Box(width = 4, height = 5),
            ],
        ),
    ],
)

assert.eq(r1.main_align, "space_evenly")
assert.eq(r1.cross_align, "center")
assert.eq(r1.children[1].main_align, "start")
assert.eq(r1.children[1].cross_align, "end")
assert.eq(len(r1.children), 2)
assert.eq(len(r1.children[1].children), 2)
