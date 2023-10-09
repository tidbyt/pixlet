#ifndef CANVAS_DRAW_H
#define CANVAS_DRAW_H

#include <stdint.h>

#ifdef __cplusplus

#include "canvas_generated.h"

/**
 * Renders a list of canvas operations to a buffer.
 *
 * @param operations The list of canvas operations to render.
 * @param size A pointer to a size_t variable that will be set to the size of
 * the rendered buffer.
 * @return A pointer to the rendered buffer. The caller is responsible for
 *         freeing this buffer.
 */
uint8_t *CanvasDraw(const fb::CanvasOperationList *const operations,
                    size_t *outSize);

extern "C"
{
#endif

    uint8_t *canvas_draw(const uint8_t *operations, size_t *outSize);

    void canvas_register_typeface(const uint8_t *data, size_t size,
                                  const char *alias);

#ifdef __cplusplus
}
#endif

#endif // CANVAS_DRAW_H