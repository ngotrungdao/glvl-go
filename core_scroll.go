package lvgl

/*
#include <lvgl.h>
*/
import "C"

// ScrollbarMode mirrors lv_scrollbar_mode_t.
type ScrollbarMode uint32

var (
	ScrollbarModeOff    = ScrollbarMode(C.LV_SCROLLBAR_MODE_OFF)
	ScrollbarModeOn     = ScrollbarMode(C.LV_SCROLLBAR_MODE_ON)
	ScrollbarModeActive = ScrollbarMode(C.LV_SCROLLBAR_MODE_ACTIVE)
	ScrollbarModeAuto   = ScrollbarMode(C.LV_SCROLLBAR_MODE_AUTO)
)

// ScrollSnap mirrors lv_scroll_snap_t.
type ScrollSnap uint32

var (
	ScrollSnapNone   = ScrollSnap(C.LV_SCROLL_SNAP_NONE)
	ScrollSnapStart  = ScrollSnap(C.LV_SCROLL_SNAP_START)
	ScrollSnapEnd    = ScrollSnap(C.LV_SCROLL_SNAP_END)
	ScrollSnapCenter = ScrollSnap(C.LV_SCROLL_SNAP_CENTER)
)

// SetScrollbarMode sets when the object's scrollbar(s) are shown.
func (o *Obj) SetScrollbarMode(mode ScrollbarMode) {
	C.lv_obj_set_scrollbar_mode(o.c, C.lv_scrollbar_mode_t(mode))
}

// SetScrollDir sets which direction(s) the object can be scrolled in.
func (o *Obj) SetScrollDir(dir Dir) { C.lv_obj_set_scroll_dir(o.c, C.lv_dir_t(dir)) }

// SetScrollSnapX sets horizontal scroll-snap alignment for children.
func (o *Obj) SetScrollSnapX(align ScrollSnap) {
	C.lv_obj_set_scroll_snap_x(o.c, C.lv_scroll_snap_t(align))
}

// SetScrollSnapY sets vertical scroll-snap alignment for children.
func (o *Obj) SetScrollSnapY(align ScrollSnap) {
	C.lv_obj_set_scroll_snap_y(o.c, C.lv_scroll_snap_t(align))
}

// ScrollBy scrolls the object by a relative amount, optionally animating.
func (o *Obj) ScrollBy(dx, dy int32, animate bool) {
	C.lv_obj_scroll_by(o.c, C.int32_t(dx), C.int32_t(dy), C.lv_anim_enable_t(animate))
}

// ScrollByBounded is like ScrollBy but clamps to the scrollable range.
func (o *Obj) ScrollByBounded(dx, dy int32, animate bool) {
	C.lv_obj_scroll_by_bounded(o.c, C.int32_t(dx), C.int32_t(dy), C.lv_anim_enable_t(animate))
}

// ScrollTo scrolls to an absolute position, optionally animating.
func (o *Obj) ScrollTo(x, y int32, animate bool) {
	C.lv_obj_scroll_to(o.c, C.int32_t(x), C.int32_t(y), C.lv_anim_enable_t(animate))
}

// ScrollToX scrolls to an absolute horizontal position, optionally animating.
func (o *Obj) ScrollToX(x int32, animate bool) {
	C.lv_obj_scroll_to_x(o.c, C.int32_t(x), C.lv_anim_enable_t(animate))
}

// ScrollToY scrolls to an absolute vertical position, optionally animating.
func (o *Obj) ScrollToY(y int32, animate bool) {
	C.lv_obj_scroll_to_y(o.c, C.int32_t(y), C.lv_anim_enable_t(animate))
}

// ScrollToView scrolls the object's ancestors so it becomes visible.
func (o *Obj) ScrollToView(animate bool) {
	C.lv_obj_scroll_to_view(o.c, C.lv_anim_enable_t(animate))
}

// ScrollToViewRecursive is like ScrollToView but scrolls every scrollable
// ancestor, not just the immediate parent.
func (o *Obj) ScrollToViewRecursive(animate bool) {
	C.lv_obj_scroll_to_view_recursive(o.c, C.lv_anim_enable_t(animate))
}

// IsScrolling reports whether the object is currently being scrolled.
func (o *Obj) IsScrolling() bool { return bool(C.lv_obj_is_scrolling(o.c)) }

// ScrollbarInvalidate marks the object's scrollbar(s) for redraw.
func (o *Obj) ScrollbarInvalidate() { C.lv_obj_scrollbar_invalidate(o.c) }

// ScrollbarMode returns the current scrollbar visibility mode.
func (o *Obj) ScrollbarMode() ScrollbarMode {
	return ScrollbarMode(C.lv_obj_get_scrollbar_mode(o.c))
}

// ScrollDir returns the direction(s) the object can be scrolled in.
func (o *Obj) ScrollDir() Dir { return Dir(C.lv_obj_get_scroll_dir(o.c)) }

// ScrollSnapX returns the horizontal scroll-snap alignment.
func (o *Obj) ScrollSnapX() ScrollSnap { return ScrollSnap(C.lv_obj_get_scroll_snap_x(o.c)) }

// ScrollSnapY returns the vertical scroll-snap alignment.
func (o *Obj) ScrollSnapY() ScrollSnap { return ScrollSnap(C.lv_obj_get_scroll_snap_y(o.c)) }

// ScrollX returns the current horizontal scroll position.
func (o *Obj) ScrollX() int32 { return int32(C.lv_obj_get_scroll_x(o.c)) }

// ScrollY returns the current vertical scroll position.
func (o *Obj) ScrollY() int32 { return int32(C.lv_obj_get_scroll_y(o.c)) }

// ScrollTop returns how far the content can still scroll upward.
func (o *Obj) ScrollTop() int32 { return int32(C.lv_obj_get_scroll_top(o.c)) }

// ScrollBottom returns how far the content can still scroll downward.
func (o *Obj) ScrollBottom() int32 { return int32(C.lv_obj_get_scroll_bottom(o.c)) }

// ScrollLeft returns how far the content can still scroll left.
func (o *Obj) ScrollLeft() int32 { return int32(C.lv_obj_get_scroll_left(o.c)) }

// ScrollRight returns how far the content can still scroll right.
func (o *Obj) ScrollRight() int32 { return int32(C.lv_obj_get_scroll_right(o.c)) }
