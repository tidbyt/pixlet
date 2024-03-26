#include <iostream>
#include <memory>
#include <string>

#include "canvas_generated.h"
#include "include/codec/SkCodec.h"
#include "include/core/SkBitmap.h"
#include "include/core/SkCanvas.h"
#include "include/core/SkColor.h"
#include "include/core/SkData.h"
#include "include/core/SkEncodedImageFormat.h"
#include "include/core/SkFont.h"
#include "include/core/SkFontTypes.h"
#include "include/core/SkGraphics.h"
#include "include/core/SkImage.h"
#include "include/core/SkPaint.h"
#include "include/core/SkPath.h"
#include "include/core/SkPoint.h"
#include "include/core/SkRect.h"
#include "include/core/SkShader.h"
#include "include/core/SkStream.h"
#include "include/core/SkString.h"
#include "include/core/SkSurface.h"
#include "include/core/SkTextBlob.h"
#include "include/core/SkTileMode.h"
#include "include/encode/SkWebpEncoder.h"
#include "modules/skparagraph/include/TypefaceFontProvider.h"

sk_sp<skia::textlayout::TypefaceFontProvider> fontProvider =
    sk_make_sp<skia::textlayout::TypefaceFontProvider>();

uint8_t *CanvasDraw(const fb::CanvasOperationList *const operations,
                    size_t *size) {
  sk_sp<SkSurface> surface =
      SkSurfaces::Raster(SkImageInfo::MakeN32Premul(64, 32));
  SkCanvas *canvas = surface->getCanvas();

  SkPaint paint;
  paint.setStrokeCap(SkPaint::kSquare_Cap);
  paint.setStyle(SkPaint::kFill_Style);
  paint.setAntiAlias(false);

  SkFont font;
  SkPath path;

  for (const auto &operation : *operations->operations()) {
    if (operation->add_arc() != nullptr) {
      const auto &arc = operation->add_arc();
      SkRect oval =
          SkRect::MakeXYWH(arc->x() - arc->radius(), arc->y() - arc->radius(),
                           2 * arc->radius(), 2 * arc->radius());
      float startAngleDeg = arc->start_angle() * (180.0 / M_PI);
      float endAngleDeg = arc->end_angle() * (180.0 / M_PI);
      float sweepAngleDeg = endAngleDeg - startAngleDeg;
      path.addArc(oval, startAngleDeg, sweepAngleDeg);
    }

    if (operation->add_circle() != nullptr) {
      const auto &circle = operation->add_circle();
      path.addCircle(circle->x(), circle->y(), circle->radius());
    }

    if (operation->add_line_to() != nullptr) {
      const auto &line = operation->add_line_to();
      path.lineTo(line->x(), line->y());
    }

    if (operation->add_rectangle() != nullptr) {
      const auto &rect = operation->add_rectangle();
      path.addRect(SkRect::MakeXYWH(rect->x(), rect->y(), rect->width(),
                                    rect->height()));
    }

    if (operation->clear() != nullptr) {
      canvas->clear(paint.getColor());
    }

    if (operation->clip_rectangle() != nullptr) {
      const auto &rect = operation->clip_rectangle();
      SkRect skRect =
          SkRect::MakeXYWH(rect->x(), rect->y(), rect->width(), rect->height());
      canvas->clipRect(skRect);
    }

    if (operation->draw_image() != nullptr) {
      const auto &drawImage = operation->draw_image();
      const auto &imageData = drawImage->image()->image_data();
      sk_sp<SkData> data =
          SkData::MakeWithoutCopy(imageData->data(), imageData->size());
      std::unique_ptr<SkCodec> codec = SkCodec::MakeFromData(data);

      if (!codec) {
        return NULL;
      }

      SkBitmap bitmap;
      bitmap.allocPixels(codec->getInfo());
      if (codec->getPixels(codec->getInfo(), bitmap.getPixels(),
                           bitmap.rowBytes()) != SkCodec::kSuccess) {
        return NULL;
      }

      SkRect skRect = SkRect::MakeXYWH(drawImage->x(), drawImage->y(),
                                       drawImage->width(), drawImage->height());
      SkSamplingOptions sampling(SkFilterMode::kNearest);

      sk_sp<SkImage> image = SkImages::RasterFromBitmap(bitmap);
      canvas->drawImageRect(image, skRect, sampling, nullptr);
    }

    if (operation->draw_pixel() != nullptr) {
      const auto &pixel = operation->draw_pixel();
      canvas->drawPoint(pixel->x(), pixel->y(), paint);
    }

    if (operation->draw_string() != nullptr) {
      const auto &drawString = operation->draw_string();
      canvas->drawString(drawString->text()->c_str(), drawString->x(),
                         drawString->y(), font, paint);
    }

    if (operation->fill_path() != nullptr) {
      canvas->drawPath(path, paint);
      path.reset();
    }

    if (operation->pop() != nullptr) {
      canvas->restore();
    }

    if (operation->push() != nullptr) {
      canvas->save();
    }

    if (operation->rotate() != nullptr) {
      const auto &rotate = operation->rotate();
      canvas->rotate(rotate->angle() * (180.0 / M_PI));
    }

    if (operation->scale() != nullptr) {
      const auto &scale = operation->scale();
      canvas->scale(scale->x(), scale->y());
    }

    if (operation->set_color() != nullptr) {
      const auto &color = operation->set_color()->color();
      paint.setColor(
          SkColorSetARGB(color->alpha(), color->r(), color->g(), color->b()));
    }

    if (operation->set_font_face() != nullptr) {
      const auto &fontFace = operation->set_font_face()->font_face();

      sk_sp<SkTypeface> skFace(fontProvider.get()
                                   ->matchFamily(fontFace->name()->c_str())
                                   ->createTypeface(0));
      font = SkFont(skFace, fontFace->size());
      font.setEdging(SkFont::Edging::kAlias);  // disable anti-aliasing
    }

    if (operation->translate() != nullptr) {
      const auto &translate = operation->translate();
      canvas->translate(translate->dx(), translate->dy());
    }
  }

  auto image = surface->makeImageSnapshot();

  SkBitmap bm;
  if (!bm.tryAllocPixels(SkImageInfo::MakeN32Premul(64, 32))) {
    std::cerr << "Could not allocate bitmap.\n";
    return NULL;
  }

  image->readPixels(nullptr, bm.info(), bm.getPixels(), 64 * 4, 0, 0);

  // encode and return SkBitmap
  SkWebpEncoder::Options encodeOpts;
  encodeOpts.fCompression = SkWebpEncoder::Compression::kLossless;

  SkDynamicMemoryWStream memStream;
  if (!SkWebpEncoder::Encode(&memStream, bm.pixmap(), encodeOpts)) {
    std::cerr << "Could not encode bitmap.\n";
    return NULL;
  }

  // Retrieve the PNG data
  *size = memStream.bytesWritten();
  uint8_t *output = (uint8_t *)malloc(*size);
  memStream.copyTo(output);

  return output;
}

extern "C" uint8_t *canvas_draw(const uint8_t *operations, size_t *outSize) {
  const fb::CanvasOperationList *const ops =
      flatbuffers::GetRoot<fb::CanvasOperationList>(operations);
  return CanvasDraw(ops, outSize);
}

extern "C" void canvas_register_typeface(const uint8_t *data, size_t size,
                                         const char *alias) {
  sk_sp<SkData> skData = SkData::MakeWithoutCopy(data, size);
  sk_sp<SkTypeface> skFace = SkTypeface::MakeFromData(skData);
  fontProvider.get()->registerTypeface(skFace, SkString(alias));
}