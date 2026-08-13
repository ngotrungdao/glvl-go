package lvgl

/*
#include <stdlib.h>
#include <lvgl.h>
*/
import "C"
import "unsafe"

// ArcLabelDir mirrors lv_arclabel_dir_t.
type ArcLabelDir uint32

var (
	ArcLabelDirClockwise        = ArcLabelDir(C.LV_ARCLABEL_DIR_CLOCKWISE)
	ArcLabelDirCounterClockwise = ArcLabelDir(C.LV_ARCLABEL_DIR_COUNTER_CLOCKWISE)
)

// ArcLabelTextAlign mirrors lv_arclabel_text_align_t.
type ArcLabelTextAlign uint32

var (
	ArcLabelTextAlignDefault  = ArcLabelTextAlign(C.LV_ARCLABEL_TEXT_ALIGN_DEFAULT)
	ArcLabelTextAlignLeading  = ArcLabelTextAlign(C.LV_ARCLABEL_TEXT_ALIGN_LEADING)
	ArcLabelTextAlignCenter   = ArcLabelTextAlign(C.LV_ARCLABEL_TEXT_ALIGN_CENTER)
	ArcLabelTextAlignTrailing = ArcLabelTextAlign(C.LV_ARCLABEL_TEXT_ALIGN_TRAILING)
)

// ArcLabelOverflow mirrors lv_arclabel_overflow_t.
type ArcLabelOverflow uint32

var (
	ArcLabelOverflowVisible  = ArcLabelOverflow(C.LV_ARCLABEL_OVERFLOW_VISIBLE)
	ArcLabelOverflowEllipsis = ArcLabelOverflow(C.LV_ARCLABEL_OVERFLOW_ELLIPSIS)
	ArcLabelOverflowClip     = ArcLabelOverflow(C.LV_ARCLABEL_OVERFLOW_CLIP)
)

// ArcLabel wraps an lv_arclabel widget: text flowed along an arc.
type ArcLabel struct{ *Obj }

// NewArcLabel creates an arc label as a child of parent.
func NewArcLabel(parent *Obj) *ArcLabel {
	return &ArcLabel{wrapObj(C.lv_arclabel_create(parent.c))}
}

// SetText sets the label's text.
func (a *ArcLabel) SetText(text string) {
	cText := C.CString(text)
	defer C.free(unsafe.Pointer(cText))
	C.lv_arclabel_set_text(a.c, cText)
}

// SetAngleStart sets the starting angle, in degrees, of the text arc.
func (a *ArcLabel) SetAngleStart(start float32) {
	C.lv_arclabel_set_angle_start(a.c, C.lv_value_precise_t(start))
}

// SetAngleSize sets the angular length, in degrees, the text can span.
func (a *ArcLabel) SetAngleSize(size float32) {
	C.lv_arclabel_set_angle_size(a.c, C.lv_value_precise_t(size))
}

// SetRadius sets the radius, in pixels, of the text arc.
func (a *ArcLabel) SetRadius(radius uint32) { C.lv_arclabel_set_radius(a.c, C.uint32_t(radius)) }

// SetDir sets whether text flows clockwise or counter-clockwise.
func (a *ArcLabel) SetDir(dir ArcLabelDir) { C.lv_arclabel_set_dir(a.c, C.lv_arclabel_dir_t(dir)) }

// SetTextVerticalAlign sets the text's alignment relative to the arc radius.
func (a *ArcLabel) SetTextVerticalAlign(align ArcLabelTextAlign) {
	C.lv_arclabel_set_text_vertical_align(a.c, C.lv_arclabel_text_align_t(align))
}

// SetOverflow sets how text longer than the arc is handled.
func (a *ArcLabel) SetOverflow(overflow ArcLabelOverflow) {
	C.lv_arclabel_set_overflow(a.c, C.lv_arclabel_overflow_t(overflow))
}

// SetOffset offsets the text's starting position along the arc, in pixels.
func (a *ArcLabel) SetOffset(offset int32) { C.lv_arclabel_set_offset(a.c, C.int32_t(offset)) }

// SetRecolor enables/disables inline `#color# text#` recoloring.
func (a *ArcLabel) SetRecolor(enable bool) { C.lv_arclabel_set_recolor(a.c, C.bool(enable)) }

// SetCenterOffset sets the arc's center point offset from the widget's
// own center, in pixels.
func (a *ArcLabel) SetCenterOffset(x, y uint32) {
	C.lv_arclabel_set_center_offset_x(a.c, C.uint32_t(x))
	C.lv_arclabel_set_center_offset_y(a.c, C.uint32_t(y))
}

// SetTextHorizontalAlign sets the text's alignment along the arc's length.
func (a *ArcLabel) SetTextHorizontalAlign(align ArcLabelTextAlign) {
	C.lv_arclabel_set_text_horizontal_align(a.c, C.lv_arclabel_text_align_t(align))
}

// SetEndOverlap enables/disables allowing text to overlap itself if it
// wraps all the way around the arc.
func (a *ArcLabel) SetEndOverlap(overlap bool) { C.lv_arclabel_set_end_overlap(a.c, C.bool(overlap)) }

// Recolor reports whether inline recoloring is enabled.
func (a *ArcLabel) Recolor() bool { return bool(C.lv_arclabel_get_recolor(a.c)) }

// Radius returns the text arc's radius, in pixels.
func (a *ArcLabel) Radius() uint32 { return uint32(C.lv_arclabel_get_radius(a.c)) }

// CenterOffset returns the arc's center point offset, in pixels.
func (a *ArcLabel) CenterOffset() (x, y uint32) {
	return uint32(C.lv_arclabel_get_center_offset_x(a.c)), uint32(C.lv_arclabel_get_center_offset_y(a.c))
}

// AngleStart returns the text arc's starting angle, in degrees.
func (a *ArcLabel) AngleStart() float32 { return float32(C.lv_arclabel_get_angle_start(a.c)) }

// AngleSize returns the text arc's angular span, in degrees.
func (a *ArcLabel) AngleSize() float32 { return float32(C.lv_arclabel_get_angle_size(a.c)) }

// TextAngle returns the angular length, in degrees, the current text
// actually occupies.
func (a *ArcLabel) TextAngle() float32 { return float32(C.lv_arclabel_get_text_angle(a.c)) }

// Dir returns the text flow direction.
func (a *ArcLabel) Dir() ArcLabelDir { return ArcLabelDir(C.lv_arclabel_get_dir(a.c)) }

// TextVerticalAlign returns the text's alignment relative to the arc radius.
func (a *ArcLabel) TextVerticalAlign() ArcLabelTextAlign {
	return ArcLabelTextAlign(C.lv_arclabel_get_text_vertical_align(a.c))
}

// TextHorizontalAlign returns the text's alignment along the arc's length.
func (a *ArcLabel) TextHorizontalAlign() ArcLabelTextAlign {
	return ArcLabelTextAlign(C.lv_arclabel_get_text_horizontal_align(a.c))
}

// Overflow returns how text longer than the arc is handled.
func (a *ArcLabel) Overflow() ArcLabelOverflow {
	return ArcLabelOverflow(C.lv_arclabel_get_overflow(a.c))
}

// EndOverlap reports whether text is allowed to overlap itself.
func (a *ArcLabel) EndOverlap() bool { return bool(C.lv_arclabel_get_end_overlap(a.c)) }
