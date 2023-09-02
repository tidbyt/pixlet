package encode

import (
	"testing"

	"tidbyt.dev/pixlet/runtime"
)

var BenchmarkDotStar = `
"""
Applet: Gradient
Summary: Displays dynamic gradients
Description: Customize gradient fills for your Tidbyt.
Author: Jeffrey Lancaster
"""

load("math.star", "math")
load("random.star", "random")
load("render.star", "render")
load("schema.star", "schema")
load("time.star", "time")

PIXLET_W = 64
PIXLET_H = 32
GLOBAL_FONT = "tom-thumb"  # or "CG-pixel-3x5-mono"

def median(val1, val2):
    return math.floor((val1 + val2) / 2)

def makeRange(minValue, maxValue, numValues):
    rangeArray = []
    for i in range(0, numValues):
        step = (maxValue - minValue) / numValues
        calcValue = math.round(minValue + (i * step))
        rangeArray.append(calcValue)
    return rangeArray

def rgbRange(start, end, steps):
    rRange = makeRange(start[0], end[0], steps)
    gRange = makeRange(start[1], end[1], steps)
    bRange = makeRange(start[2], end[2], steps)
    returnRange = []
    for n in range(0, steps):
        returnRange.append([rRange[n], gRange[n], bRange[n]])
    return returnRange

# from: https://www.educative.io/answers/how-to-convert-hex-to-rgb-and-rgb-to-hex-in-python
def hex_to_rgb(hex):
    hex = hex.replace("#", "")
    rgb = []
    for i in (0, 2, 4):
        decimal = int(hex[i:i + 2], 16)
        rgb.append(decimal)
    return tuple(rgb)

def rgb_to_hex(r, g, b):
    rgbArr = ["0", "1", "2", "3", "4", "5", "6", "7", "8", "9", "A", "B", "C", "D", "E", "F"]
    hex = "#"
    r = math.floor(r)
    g = math.floor(g)
    b = math.floor(b)
    for i in (r, g, b):
        secondNum = i % 16
        firstNum = math.floor((i - secondNum) / 16)
        hex += rgbArr[firstNum] + rgbArr[secondNum]
    return hex

def randomColor():
    randomRed = random.number(0, 255)
    randomGreen = random.number(0, 255)
    randomBlue = random.number(0, 255)
    return rgb_to_hex(randomRed, randomGreen, randomBlue)

def shiftLeft(thisArray):
    newThisArray = []
    for i in thisArray:
        newThisArray.append(i[1:] + i[:1])
    return newThisArray

def shiftRight(thisArray):
    newThisArray = []
    for i in thisArray:
        newThisArray.append(i[-1:] + i[:-1])
    return newThisArray

# from: https://stackoverflow.com/questions/2150108/efficient-way-to-rotate-a-list-in-python
def shiftUp(thisArray):
    return thisArray[1:] + thisArray[:1]

def shiftDown(thisArray):
    return thisArray[-1:] + thisArray[:-1]

def four_color_gradient(topL, topR, botL, botR, config):
    # convert inputs to rgb
    topLrgb = hex_to_rgb(topL)
    topRrgb = hex_to_rgb(topR)
    botLrgb = hex_to_rgb(botL)
    botRrgb = hex_to_rgb(botR)

    # determine left column and right column ranges: PIXLET_H
    leftCol = rgbRange(topLrgb, botLrgb, PIXLET_H)
    rightCol = rgbRange(topRrgb, botRrgb, PIXLET_H)

    # for each row, determine range: PIXLET_W
    gradientArray = []
    animatedArray = []

    # make basic gradient array
    for n in range(0, PIXLET_H):
        rowGradient = rgbRange(leftCol[n], rightCol[n], PIXLET_W)
        gradientArray.append(rowGradient)

    # convert each value in gradientArray from RGB to hex
    for n, i in enumerate(gradientArray):
        for m, j in enumerate(i):
            gradientArray[n][m] = rgb_to_hex(j[0], j[1], j[2])

    # if animated, expand gradientArray in animatedArray
    if config.bool("animation", False) == True:
        if config.get("direction") == "up" or config.get("direction") == "down":
            # append the mirror image (adds more rows down)
            mirrorArray = gradientArray[::-1]
            gradientArray += mirrorArray
        elif config.get("direction") == "left" or config.get("direction") == "right":
            # append the mirror image (adds more pixels across)
            for n, i in enumerate(gradientArray):
                mirrorRow = i[::-1]
                gradientArray[n] += mirrorRow

        # add to animatedArray
        if config.get("direction") == "left" or config.get("direction") == "right":
            numFrames = len(gradientArray[0])  # PIXLET_W * 2
        else:
            numFrames = len(gradientArray)  # PIXLET_H * 2

        for i in range(0, numFrames):
            animatedArray.append(gradientArray)

            # shift
            if config.get("direction") == "up":
                gradientArray = shiftUp(gradientArray)
            elif config.get("direction") == "down":
                gradientArray = shiftDown(gradientArray)
            elif config.get("direction") == "left":
                gradientArray = shiftLeft(gradientArray)
            elif config.get("direction") == "right":
                gradientArray = shiftRight(gradientArray)
    else:
        animatedArray = [gradientArray]

    return animatedArray

def two_color_gradient(topL, botR, config):
    topLrgb = hex_to_rgb(topL)
    botRrgb = hex_to_rgb(botR)
    medianR = median(topLrgb[0], botRrgb[0])
    medianG = median(topLrgb[1], botRrgb[1])
    medianB = median(topLrgb[2], botRrgb[2])

    # average r, g, b for other two corners
    medianRGB = rgb_to_hex(medianR, medianG, medianB)
    return four_color_gradient(topL, medianRGB, medianRGB, botR, config)

def displayArray(array, labelCount, config):
    animationChildren = []
    for n in array:  # frames
        columnChildren = []
        stackChildren = []
        for m in n:  # column of rows
            rowChildren = []
            for p in m:  # cells in row
                rowChildren.append(
                    render.Box(width = 1, height = 1, color = p),
                )
            columnChildren.append(
                render.Row(
                    children = rowChildren,
                ),
            )

        # add the gradient to the stack
        stackChildren.append(
            render.Column(children = columnChildren),
        )

        # add labels to the stack
        if config.bool("labels"):
            if labelCount == 4:
                topL = n[0][0]
                topR = n[0][PIXLET_W - 1]
                botL = n[PIXLET_H - 1][0]
                botR = n[PIXLET_H - 1][PIXLET_W - 1]
                if config.bool("animation") == False:
                    if config.get("gradient_type") == "4color":
                        topL = config.get("color1")
                        topR = config.get("color2")
                        botL = config.get("color3")
                        botR = config.get("color4")
                    elif config.get("gradient_type") == "default":
                        topL = "#FF0000"
                        topR = "#FFFF00"
                        botL = "#0000FF"
                        botR = "#FFFFFF"
                stackChildren.extend([
                    render.Padding(
                        child = render.Text(content = topL.upper().replace("#", ""), color = "#000", font = GLOBAL_FONT),
                        pad = (1, 1, 1, 1),
                    ),
                    render.Padding(
                        child = render.Text(content = topR.upper().replace("#", ""), color = "#000", font = GLOBAL_FONT),
                        pad = (40, 1, 1, 1),
                    ),
                    render.Padding(
                        child = render.Text(content = botL.upper().replace("#", ""), color = "#000", font = GLOBAL_FONT),
                        pad = (1, 26, 1, 1),
                    ),
                    render.Padding(
                        child = render.Text(content = botR.upper().replace("#", ""), color = "#000", font = GLOBAL_FONT),
                        pad = (40, 26, 1, 1),
                    ),
                ])
            elif labelCount == 2:
                topL = n[0][0]
                botR = n[PIXLET_H - 1][PIXLET_W - 1]
                if config.bool("animation") == False:
                    topL = config.get("color1")
                    botR = config.get("color2")
                stackChildren.extend([
                    render.Padding(
                        child = render.Text(content = topL.upper().replace("#", ""), color = "#000", font = GLOBAL_FONT),
                        pad = (1, 1, 1, 1),
                    ),
                    render.Padding(
                        child = render.Text(content = botR.upper().replace("#", ""), color = "#000", font = GLOBAL_FONT),
                        pad = (40, 26, 1, 1),
                    ),
                ])

        # add the stack (frame) to the animation
        animationChildren.append(
            render.Stack(children = stackChildren),
        )

    return animationChildren

def main(config):
    random.seed(time.now().unix // 15)

    # define gradientArray and labels
    animatedArray = []
    labelCount = 4
    if config.get("gradient_type") == "random":
        color1 = randomColor()
        color2 = randomColor()
        color3 = randomColor()
        color4 = randomColor()
        animatedArray = four_color_gradient(color1, color2, color3, color4, config)
    elif config.get("gradient_type") == "4color":
        color1 = config.get("color1")
        color2 = config.get("color2")
        color3 = config.get("color3")
        color4 = config.get("color4")
        animatedArray = four_color_gradient(color1, color2, color3, color4, config)
    elif config.get("gradient_type") == "2color":
        color1 = config.get("color1")
        color2 = config.get("color2")
        animatedArray = two_color_gradient(color1, color2, config)
        labelCount = 2
    else:
        animatedArray = four_color_gradient("#FF0000", "#FFFF00", "#0000FF", "#FFFFFF", config)

    # show animatedArray with labels
    animationChildren = displayArray(animatedArray, labelCount, config)

    # get the delay preference
    if config.get("speed") == "fast":
        animation_delay = 10
    else:
        animation_delay = 500

    # show the animation
    return render.Root(
        child = render.Animation(
            children = animationChildren,
        ),
        delay = animation_delay,
    )

def more_gradient_options(gradient_type):
    if gradient_type == "2color":
        return [
            schema.Color(
                id = "color1",
                name = "Color #1",
                desc = "Top left corner",
                icon = "brush",
                default = "#FF0000",
            ),
            schema.Color(
                id = "color2",
                name = "Color #2",
                desc = "Bottom right corner",
                icon = "brush",
                default = "#0000FF",
            ),
        ]
    elif gradient_type == "4color":
        return [
            schema.Color(
                id = "color1",
                name = "Color #1",
                desc = "Top left corner",
                icon = "brush",
                default = "#FF0000",
            ),
            schema.Color(
                id = "color2",
                name = "Color #2",
                desc = "Top right corner",
                icon = "brush",
                default = "#FFFF00",
            ),
            schema.Color(
                id = "color3",
                name = "Color #3",
                desc = "Bottom left corner",
                icon = "brush",
                default = "#0000FF",
            ),
            schema.Color(
                id = "color4",
                name = "Color #4",
                desc = "Bottom right corner",
                icon = "brush",
                default = "#FFFFFF",
            ),
        ]
    else:
        return []

def get_schema():
    gradientOptions = [
        schema.Option(
            display = "Default",
            value = "default",
        ),
        schema.Option(
            display = "Random",
            value = "random",
        ),
        schema.Option(
            display = "Pick 2",
            value = "2color",
        ),
        schema.Option(
            display = "Pick 4",
            value = "4color",
        ),
    ]

    animationSpeedOptions = [
        schema.Option(
            display = "Fast",
            value = "fast",
        ),
        schema.Option(
            display = "Slow",
            value = "slow",
        ),
    ]

    animationDirectionOptions = [
        schema.Option(
            display = "Scroll up",
            value = "up",
        ),
        schema.Option(
            display = "Scroll down",
            value = "down",
        ),
        schema.Option(
            display = "Scroll left",
            value = "left",
        ),
        schema.Option(
            display = "Scroll right",
            value = "right",
        ),
    ]

    # icons from: https://fontawesome.com/
    return schema.Schema(
        version = "1",
        fields = [
            schema.Dropdown(
                id = "gradient_type",
                name = "Gradient Type",
                icon = "circleHalfStroke",
                desc = "Which gradient to show",
                default = gradientOptions[0].value,
                options = gradientOptions,
            ),
            schema.Toggle(
                id = "labels",
                name = "Text",
                desc = "Show hex values?",
                icon = "font",
                default = False,
            ),
            schema.Toggle(
                id = "animation",
                name = "Animation",
                desc = "Animate the gradient?",
                icon = "play",
                default = False,
            ),
            schema.Generated(
                id = "gradient_generated",
                source = "gradient_type",
                handler = more_gradient_options,
            ),
            # schema.Generated(
            #     id = "animation_generated",
            #     source = "animation",
            #     handler = more_animation_options,
            # ),
            schema.Dropdown(
                id = "speed",
                name = "Animation Speed",
                icon = "forward",
                desc = "How fast to scroll",
                default = "slow",
                options = animationSpeedOptions,
            ),
            schema.Dropdown(
                id = "direction",
                name = "Direction",
                icon = "arrowsUpDownLeftRight",
                desc = "Which way to scroll",
                default = animationDirectionOptions[0].value,
                options = animationDirectionOptions,
            ),
        ],
    )
`

func BenchmarkRunAndRender(b *testing.B) {
	app := &runtime.Applet{}
	err := app.Load("benchmark.star", []byte(BenchmarkDotStar), nil)
	if err != nil {
		b.Error(err)
	}

	config := map[string]string{}
	roots, err := app.Run(config)
	if err != nil {
		b.Error(err)
	}

	for i := 0; i < b.N; i++ {
		webp, err := ScreensFromRoots(roots).EncodeWebP(15000)
		if err != nil {
			b.Error(err)
		}

		if len(webp) == 0 {
			b.Error()
		}
	}
}
