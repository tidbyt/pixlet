package canvas

/*
#cgo LDFLAGS: /Users/rohan/code/tidbyt/bazel-bin/canvas-server/libcanvas.a
#
#cgo LDFLAGS: -L/Users/rohan/code/tidbyt/skia/out/Release-macos-arm64
#cgo LDFLAGS: -ldng_sdk
#cgo LDFLAGS: -lexpat
#cgo LDFLAGS: -lharfbuzz
#cgo LDFLAGS: -licu
#cgo LDFLAGS: -ljpeg
#cgo LDFLAGS: -lparticles
#cgo LDFLAGS: -lpiex
#cgo LDFLAGS: -lpng
#cgo LDFLAGS: -lskcms
#cgo LDFLAGS: -lskia
#cgo LDFLAGS: -lskottie
#cgo LDFLAGS: -lskparagraph
#cgo LDFLAGS: -lskresources
#cgo LDFLAGS: -lsksg
#cgo LDFLAGS: -lskshaper
#cgo LDFLAGS: -lskunicode
#cgo LDFLAGS: -lsvg
#cgo LDFLAGS: -lwebp_sse41
#cgo LDFLAGS: -lwuffs
#cgo LDFLAGS: -lzlib
#
#cgo LDFLAGS: -framework Foundation
#cgo LDFLAGS: -framework CoreFoundation
#cgo LDFLAGS: -framework CoreGraphics
#cgo LDFLAGS: -framework CoreText
#cgo LDFLAGS: -framework ImageIO
#cgo LDFLAGS: -framework Carbon
#cgo LDFLAGS: -framework Metal
#
#cgo LDFLAGS: -lstdc++
#cgo LDFLAGS: -mmacosx-version-min=14.0
#
#include <stdlib.h>
#include "/Users/rohan/code/tidbyt/canvas-server/canvas.h"
#
*/
import "C"

import (
	"bytes"
	"image"
	"image/color"
	"image/png"
	"sync"
	"unsafe"

	flatbuffers "github.com/google/flatbuffers/go"
	"github.com/tidbyt/go-libwebp/webp"

	"tidbyt.dev/pixlet/fonts"
	"tidbyt.dev/pixlet/render/canvas/fb"
)

var (
	onceFonts sync.Once
)

type SkiaCanvas struct {
	builder    *flatbuffers.Builder
	operations []flatbuffers.UOffsetT
}

func init() {
	register("skia", NewSkiaCanvas)
}

func loadFonts() {
	names := fonts.Names()

	for _, name := range names {
		font := fonts.GetFont(name)
		if len(font.TTF) == 0 {
			continue
		}

		cFont := C.CBytes(font.TTF)
		defer C.free(cFont)

		C.canvas_register_typeface((*C.uint8_t)(cFont), C.size_t(len(font.TTF)), C.CString(name))
	}
}

func NewSkiaCanvas(width, height int) Canvas {
	onceFonts.Do(loadFonts)

	return &SkiaCanvas{
		builder: flatbuffers.NewBuilder(1024),
	}
}

func (c *SkiaCanvas) AddArc(x, y, r, startAngle, endAngle float64) {
	// Start constructing the AddArc object.
	fb.AddArcStart(c.builder)
	fb.AddArcAddX(c.builder, x)
	fb.AddArcAddY(c.builder, y)
	fb.AddArcAddRadius(c.builder, r)
	fb.AddArcAddStartAngle(c.builder, startAngle)
	fb.AddArcAddEndAngle(c.builder, endAngle)
	addArc := fb.AddArcEnd(c.builder)

	fb.CanvasOperationStart(c.builder)
	fb.CanvasOperationAddAddArc(c.builder, addArc)
	canvasOperationOffset := fb.CanvasOperationEnd(c.builder)

	// Append the operation to the list.
	c.operations = append(c.operations, canvasOperationOffset)
}

func (c *SkiaCanvas) AddCircle(x, y, r float64) {
	fb.AddCircleStart(c.builder)
	fb.AddCircleAddX(c.builder, x)
	fb.AddCircleAddY(c.builder, y)
	fb.AddCircleAddRadius(c.builder, r)
	addCircle := fb.AddCircleEnd(c.builder)

	fb.CanvasOperationStart(c.builder)
	fb.CanvasOperationAddAddCircle(c.builder, addCircle)
	operation := fb.CanvasOperationEnd(c.builder)

	c.operations = append(c.operations, operation)
}

func (c *SkiaCanvas) AddLineTo(x, y float64) {
	fb.AddLineToStart(c.builder)
	fb.AddLineToAddX(c.builder, x)
	fb.AddLineToAddY(c.builder, y)
	addLineTo := fb.AddLineToEnd(c.builder)

	fb.CanvasOperationStart(c.builder)
	fb.CanvasOperationAddAddLineTo(c.builder, addLineTo)
	operation := fb.CanvasOperationEnd(c.builder)

	c.operations = append(c.operations, operation)
}

func (c *SkiaCanvas) AddRectangle(x, y, w, h float64) {
	fb.AddRectangleStart(c.builder)
	fb.AddRectangleAddX(c.builder, x)
	fb.AddRectangleAddY(c.builder, y)
	fb.AddRectangleAddWidth(c.builder, w)
	fb.AddRectangleAddHeight(c.builder, h)
	addRectangle := fb.AddRectangleEnd(c.builder)

	fb.CanvasOperationStart(c.builder)
	fb.CanvasOperationAddAddRectangle(c.builder, addRectangle)
	operation := fb.CanvasOperationEnd(c.builder)

	c.operations = append(c.operations, operation)
}

func (c *SkiaCanvas) Clear() {
	fb.ClearStart(c.builder)
	clear := fb.ClearEnd(c.builder)

	fb.CanvasOperationStart(c.builder)
	fb.CanvasOperationAddClear(c.builder, clear)
	operation := fb.CanvasOperationEnd(c.builder)

	c.operations = append(c.operations, operation)
}

func (c *SkiaCanvas) ClipRectangle(x, y, w, h float64) {
	fb.ClipRectangleStart(c.builder)
	fb.ClipRectangleAddX(c.builder, x)
	fb.ClipRectangleAddY(c.builder, y)
	fb.ClipRectangleAddWidth(c.builder, w)
	fb.ClipRectangleAddHeight(c.builder, h)
	clipRectangle := fb.ClipRectangleEnd(c.builder)

	fb.CanvasOperationStart(c.builder)
	fb.CanvasOperationAddClipRectangle(c.builder, clipRectangle)
	operation := fb.CanvasOperationEnd(c.builder)

	c.operations = append(c.operations, operation)
}

func (c *SkiaCanvas) DrawGoImage(x, y float64, img image.Image) {
	var buf bytes.Buffer
	png.Encode(&buf, img)

	imgBytes := c.builder.CreateByteVector(buf.Bytes())
	fb.ImageStart(c.builder)
	fb.ImageAddImageData(c.builder, imgBytes)
	image := fb.ImageEnd(c.builder)

	fb.DrawImageStart(c.builder)
	fb.DrawImageAddImage(c.builder, image)
	fb.DrawImageAddX(c.builder, x)
	fb.DrawImageAddY(c.builder, y)
	fb.DrawImageAddWidth(c.builder, float64(img.Bounds().Dx()))
	fb.DrawImageAddHeight(c.builder, float64(img.Bounds().Dy()))
	drawImage := fb.DrawImageEnd(c.builder)

	fb.CanvasOperationStart(c.builder)
	fb.CanvasOperationAddDrawImage(c.builder, drawImage)
	operation := fb.CanvasOperationEnd(c.builder)

	c.operations = append(c.operations, operation)
}

func (c *SkiaCanvas) DrawImageFromBuffer(x, y, w, h float64, img []byte) {
	imgBytes := c.builder.CreateByteVector(img)
	fb.ImageStart(c.builder)
	fb.ImageAddImageData(c.builder, imgBytes)
	image := fb.ImageEnd(c.builder)

	fb.DrawImageStart(c.builder)
	fb.DrawImageAddImage(c.builder, image)
	fb.DrawImageAddX(c.builder, x)
	fb.DrawImageAddY(c.builder, y)
	fb.DrawImageAddWidth(c.builder, w)
	fb.DrawImageAddHeight(c.builder, h)
	drawImage := fb.DrawImageEnd(c.builder)

	fb.CanvasOperationStart(c.builder)
	fb.CanvasOperationAddDrawImage(c.builder, drawImage)
	operation := fb.CanvasOperationEnd(c.builder)

	c.operations = append(c.operations, operation)
}

func (c *SkiaCanvas) DrawPixel(x, y int) {
	fb.DrawPixelStart(c.builder)
	fb.DrawPixelAddX(c.builder, int32(x))
	fb.DrawPixelAddY(c.builder, int32(y))
	drawPixel := fb.DrawPixelEnd(c.builder)

	fb.CanvasOperationStart(c.builder)
	fb.CanvasOperationAddDrawPixel(c.builder, drawPixel)
	operation := fb.CanvasOperationEnd(c.builder)

	c.operations = append(c.operations, operation)
}

func (c *SkiaCanvas) DrawString(x, y float64, text string) {
	textStr := c.builder.CreateString(text)

	fb.DrawStringStart(c.builder)
	fb.DrawStringAddText(c.builder, textStr)
	fb.DrawStringAddX(c.builder, x)
	fb.DrawStringAddY(c.builder, y)
	drawString := fb.DrawStringEnd(c.builder)

	fb.CanvasOperationStart(c.builder)
	fb.CanvasOperationAddDrawString(c.builder, drawString)
	operation := fb.CanvasOperationEnd(c.builder)

	c.operations = append(c.operations, operation)
}

func (c *SkiaCanvas) DrawStringWrapped(x, y, w, spacing float64, text string, align TextAlign) {
	var alignFlag fb.Align
	switch align {
	case AlignLeft:
		alignFlag = fb.AlignLEFT
	case AlignCenter:
		alignFlag = fb.AlignCENTER
	case AlignRight:
		alignFlag = fb.AlignRIGHT
	}

	textStr := c.builder.CreateString(text)

	fb.DrawStringWrappedStart(c.builder)
	fb.DrawStringWrappedAddText(c.builder, textStr)
	fb.DrawStringWrappedAddX(c.builder, x)
	fb.DrawStringWrappedAddY(c.builder, y)
	fb.DrawStringWrappedAddWidth(c.builder, w)
	fb.DrawStringWrappedAddSpacing(c.builder, spacing)
	fb.DrawStringWrappedAddAlign(c.builder, alignFlag)
	drawStringWrapped := fb.DrawStringWrappedEnd(c.builder)

	fb.CanvasOperationStart(c.builder)
	fb.CanvasOperationAddDrawStringWrapped(c.builder, drawStringWrapped)
	operation := fb.CanvasOperationEnd(c.builder)

	c.operations = append(c.operations, operation)
}

func (c *SkiaCanvas) FillPath() {
	fb.FillPathStart(c.builder)
	fillPath := fb.FillPathEnd(c.builder)

	fb.CanvasOperationStart(c.builder)
	fb.CanvasOperationAddFillPath(c.builder, fillPath)
	operation := fb.CanvasOperationEnd(c.builder)

	c.operations = append(c.operations, operation)
}

func (c *SkiaCanvas) FontHeight() float64 {
	return 0
}

func (c *SkiaCanvas) Image() image.Image {
	fb.CanvasOperationListStartOperationsVector(c.builder, len(c.operations))
	for i := len(c.operations) - 1; i >= 0; i-- {
		c.builder.PrependUOffsetT(c.operations[i])
	}
	operationsVec := c.builder.EndVector(len(c.operations))

	fb.CanvasOperationListStart(c.builder)
	fb.CanvasOperationListAddOperations(c.builder, operationsVec)
	canvasOpList := fb.CanvasOperationListEnd(c.builder)
	c.builder.Finish(canvasOpList)

	// Call the C function.
	var outSize C.size_t

	cOperations := C.CBytes(c.builder.FinishedBytes())
	defer C.free(cOperations)
	cImageData := C.canvas_draw((*C.uint8_t)(cOperations), &outSize)
	defer C.free(unsafe.Pointer(cImageData))

	// Convert the C image data to a Go byte slice.
	imageData := C.GoBytes(unsafe.Pointer(cImageData), C.int(outSize))

	// Decode the image data.
	image, err := webp.DecodeRGBA(imageData, &webp.DecoderOptions{})
	if err != nil {
		panic(err)
	}

	return image
}

func (c *SkiaCanvas) MeasureString(text string) (w, h float64) {
	return 0, 0
}

func (c *SkiaCanvas) Pop() {
	fb.PopStart(c.builder)
	pop := fb.PopEnd(c.builder)

	fb.CanvasOperationStart(c.builder)
	fb.CanvasOperationAddPop(c.builder, pop)
	operation := fb.CanvasOperationEnd(c.builder)

	c.operations = append(c.operations, operation)
}

func (c *SkiaCanvas) Push() {
	fb.PushStart(c.builder)
	push := fb.PushEnd(c.builder)

	fb.CanvasOperationStart(c.builder)
	fb.CanvasOperationAddPush(c.builder, push)
	operation := fb.CanvasOperationEnd(c.builder)

	c.operations = append(c.operations, operation)
}

func (c *SkiaCanvas) Rotate(angle float64) {
	fb.RotateStart(c.builder)
	fb.RotateAddAngle(c.builder, angle)
	rotate := fb.RotateEnd(c.builder)

	fb.CanvasOperationStart(c.builder)
	fb.CanvasOperationAddRotate(c.builder, rotate)
	operation := fb.CanvasOperationEnd(c.builder)

	c.operations = append(c.operations, operation)
}

func (c *SkiaCanvas) Scale(x, y float64) {
	fb.ScaleStart(c.builder)
	fb.ScaleAddX(c.builder, x)
	fb.ScaleAddY(c.builder, y)
	scale := fb.ScaleEnd(c.builder)

	fb.CanvasOperationStart(c.builder)
	fb.CanvasOperationAddScale(c.builder, scale)
	operation := fb.CanvasOperationEnd(c.builder)

	c.operations = append(c.operations, operation)
}

func (c *SkiaCanvas) SetColor(col color.Color) {
	colNRGBA := color.NRGBAModel.Convert(col).(color.NRGBA)

	fb.ColorStart(c.builder)
	fb.ColorAddR(c.builder, uint32(colNRGBA.R))
	fb.ColorAddG(c.builder, uint32(colNRGBA.G))
	fb.ColorAddB(c.builder, uint32(colNRGBA.B))
	fb.ColorAddAlpha(c.builder, uint32(colNRGBA.A))
	fbCol := fb.ColorEnd(c.builder)

	fb.SetColorStart(c.builder)
	fb.SetColorAddColor(c.builder, fbCol)
	setColor := fb.SetColorEnd(c.builder)

	fb.CanvasOperationStart(c.builder)
	fb.CanvasOperationAddSetColor(c.builder, setColor)
	operation := fb.CanvasOperationEnd(c.builder)

	c.operations = append(c.operations, operation)
}

func (c *SkiaCanvas) SetFont(font *fonts.Font) {
	fontNameStr := c.builder.CreateString(font.Name)

	fb.FontFaceStart(c.builder)
	fb.FontFaceAddName(c.builder, fontNameStr)
	fb.FontFaceAddSize(c.builder, int32(font.Font.PixelSize))
	fc := fb.FontFaceEnd(c.builder)

	fb.SetFontFaceStart(c.builder)
	fb.SetFontFaceAddFontFace(c.builder, fc)
	setFontFace := fb.SetFontFaceEnd(c.builder)

	fb.CanvasOperationStart(c.builder)
	fb.CanvasOperationAddSetFontFace(c.builder, setFontFace)
	operation := fb.CanvasOperationEnd(c.builder)

	c.operations = append(c.operations, operation)
}

func (c *SkiaCanvas) TransformPoint(x, y float64) (ax, ay float64) {
	// Assuming this method doesn't modify the operations queue.
	return 0, 0
}

func (c *SkiaCanvas) Translate(dx, dy float64) {
	fb.TranslateStart(c.builder)
	fb.TranslateAddDx(c.builder, dx)
	fb.TranslateAddDy(c.builder, dy)
	translate := fb.TranslateEnd(c.builder)

	fb.CanvasOperationStart(c.builder)
	fb.CanvasOperationAddTranslate(c.builder, translate)
	operation := fb.CanvasOperationEnd(c.builder)

	c.operations = append(c.operations, operation)
}
