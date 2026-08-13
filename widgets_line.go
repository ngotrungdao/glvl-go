package lvgl

/*
#include <stdlib.h>
#include <lvgl.h>
*/
import "C"
import "unsafe"

// Line wraps an lv_line widget: draws straight segments connecting a list
// of points.
//
// LVGL keeps a reference to the points array passed to SetPoints rather
// than copying it, so SetPoints builds and pins the array on the C heap,
// freeing the previous one on replacement.
type Line struct {
	*Obj
	points *C.lv_point_precise_t
}

// NewLine creates a line as a child of parent.
func NewLine(parent *Obj) *Line {
	return &Line{Obj: wrapObj(C.lv_line_create(parent.c))}
}

// SetPoints sets the sequence of points the line connects.
func (l *Line) SetPoints(points []Point) {
	n := len(points)
	c := (*C.lv_point_precise_t)(C.malloc(C.size_t(n) * C.size_t(unsafe.Sizeof(C.lv_point_precise_t{}))))
	buf := unsafe.Slice(c, n)
	for i, p := range points {
		buf[i].x = C.lv_value_precise_t(p.X)
		buf[i].y = C.lv_value_precise_t(p.Y)
	}
	C.lv_line_set_points(l.c, c, C.uint32_t(n))

	if l.points != nil {
		C.free(unsafe.Pointer(l.points))
	}
	l.points = c
}

// SetYInvert flips the line vertically within its own coordinate space.
func (l *Line) SetYInvert(enable bool) { C.lv_line_set_y_invert(l.c, C.bool(enable)) }

// PointCount returns how many points the line currently has.
func (l *Line) PointCount() uint32 { return uint32(C.lv_line_get_point_count(l.c)) }

// YInvert reports whether the line is vertically flipped.
func (l *Line) YInvert() bool { return bool(C.lv_line_get_y_invert(l.c)) }
