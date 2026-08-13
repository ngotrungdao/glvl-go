package lvgl

/*
#include <stdlib.h>
#include <lvgl.h>
*/
import "C"
import "unsafe"

// Canvas wraps an lv_canvas widget: a raster surface you can draw into
// pixel-by-pixel.
//
// Like Style/GridDsc, LVGL stores a pointer to the pixel buffer rather
// than copying it, so SetBuffer allocates and pins it on the C heap for
// the widget's lifetime.
type Canvas struct {
	*Obj
	buf unsafe.Pointer
}

// NewCanvas creates a canvas as a child of parent.
func NewCanvas(parent *Obj) *Canvas {
	return &Canvas{Obj: wrapObj(C.lv_canvas_create(parent.c))}
}

// SetBuffer allocates an ARGB8888 (4 bytes/pixel) buffer of the given
// pixel size and assigns it as the canvas's draw target.
func (cv *Canvas) SetBuffer(w, h int32) {
	size := C.size_t(w) * C.size_t(h) * 4
	buf := C.malloc(size)
	C.lv_canvas_set_buffer(cv.c, buf, C.int32_t(w), C.int32_t(h), C.LV_COLOR_FORMAT_ARGB8888)

	if cv.buf != nil {
		C.free(cv.buf)
	}
	cv.buf = buf
}

// SetPixel sets a single pixel's color and opacity. The canvas must have
// a buffer assigned via SetBuffer first.
func (cv *Canvas) SetPixel(x, y int32, color Color, opa Opa) {
	C.lv_canvas_set_px(cv.c, C.int32_t(x), C.int32_t(y), color.toC(), C.lv_opa_t(opa))
}

// FillBg fills the entire canvas with a solid color.
func (cv *Canvas) FillBg(color Color, opa Opa) {
	C.lv_canvas_fill_bg(cv.c, color.toC(), C.lv_opa_t(opa))
}
