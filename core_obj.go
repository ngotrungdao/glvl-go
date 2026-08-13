package lvgl

/*
#include <stdlib.h>
#include <lvgl.h>
*/
import "C"
import (
	"runtime/cgo"
	"unsafe"
)

// Align mirrors lv_align_t.
type Align uint32

const (
	AlignDefault Align = iota
	AlignTopLeft
	AlignTopMid
	AlignTopRight
	AlignBottomLeft
	AlignBottomMid
	AlignBottomRight
	AlignLeftMid
	AlignRightMid
	AlignCenter
	AlignOutTopLeft
	AlignOutTopMid
	AlignOutTopRight
	AlignOutBottomLeft
	AlignOutBottomMid
	AlignOutBottomRight
	AlignOutLeftTop
	AlignOutLeftMid
	AlignOutLeftBottom
	AlignOutRightTop
	AlignOutRightMid
	AlignOutRightBottom
)

// State mirrors lv_state_t bit flags.
type State uint16

const (
	StateDefault  State = 0
	StateAlt      State = 1 << 0
	StateChecked  State = 1 << 2
	StateFocused  State = 1 << 3
	StateFocusKey State = 1 << 4
	StateEdited   State = 1 << 5
	StateHovered  State = 1 << 6
	StatePressed  State = 1 << 7
	StateScrolled State = 1 << 8
	StateDisabled State = 1 << 9
	StateUser1    State = 1 << 12
	StateUser2    State = 1 << 13
	StateUser3    State = 1 << 14
	StateUser4    State = 1 << 15
	StateAny      State = 0xFFFF
)

// Coordinate helpers. LVGL encodes LV_SIZE_CONTENT and percentage values by
// tagging the top bits of an otherwise plain int32 coordinate; these tags
// are set via preprocessor macros in C (LV_PCT, LV_SIZE_CONTENT) that cgo
// cannot call directly, so they are reimplemented here from
// core/lv_area.h's definitions.
const (
	coordTypeShift = 29
	coordTypeSpec  = int32(1) << coordTypeShift

	// CoordMax is the largest plain coordinate LVGL accepts (LV_COORD_MAX).
	CoordMax = (int32(1) << coordTypeShift) - 1

	// SizeContent tells a widget to size itself to fit its content
	// (LV_SIZE_CONTENT).
	SizeContent = CoordMax | coordTypeSpec

	pctStoredMax = CoordMax - 1
	pctPosMax    = pctStoredMax / 2
)

// Pct returns a coordinate value meaning "x percent of the parent size",
// equivalent to the LV_PCT(x) macro.
func Pct(x int32) int32 {
	v := x
	if v < 0 {
		if v < -pctPosMax {
			v = -pctPosMax
		}
		v = pctPosMax - (-v)
	} else if v > pctPosMax {
		v = pctPosMax
	}
	return v | coordTypeSpec
}

// Obj wraps an lv_obj_t. LVGL owns the underlying memory; call Delete to
// free it (which also deletes all children), mirroring lv_obj_delete.
type Obj struct {
	c *C.lv_obj_t

	handles     []cgo.Handle // event/anim callback handles registered on this obj, freed on LV_EVENT_DELETE
	autoCleanup bool         // whether the LV_EVENT_DELETE cleanup hook has been registered

	gridCols, gridRows *GridDsc // kept alive only; caller still owns Delete (see SetGridDscArray)

	bgImageSrc, arcImageSrc, bitmapMaskSrc *C.char // pinned strings set via SetStyleBgImageSrc etc.
}

func wrapObj(c *C.lv_obj_t) *Obj {
	if c == nil {
		return nil
	}
	return &Obj{c: c}
}

// NewObj creates a plain container object. If parent is nil, a new screen
// is created (lv_obj_create(NULL)).
func NewObj(parent *Obj) *Obj {
	var p *C.lv_obj_t
	if parent != nil {
		p = parent.c
	}
	return wrapObj(C.lv_obj_create(p))
}

// Delete deletes the object and all of its children.
func (o *Obj) Delete() {
	C.lv_obj_delete(o.c)
	o.c = nil
	for _, p := range []*C.char{o.bgImageSrc, o.arcImageSrc, o.bitmapMaskSrc} {
		if p != nil {
			C.free(unsafe.Pointer(p))
		}
	}
	o.bgImageSrc, o.arcImageSrc, o.bitmapMaskSrc = nil, nil, nil
}

// Clean deletes all children of the object, keeping the object itself.
func (o *Obj) Clean() {
	C.lv_obj_clean(o.c)
}

// SetPos sets the object's x, y position relative to its parent.
func (o *Obj) SetPos(x, y int32) {
	C.lv_obj_set_pos(o.c, C.int32_t(x), C.int32_t(y))
}

// SetSize sets the object's width and height.
func (o *Obj) SetSize(w, h int32) {
	C.lv_obj_set_size(o.c, C.int32_t(w), C.int32_t(h))
}

// SetWidth sets the object's width.
func (o *Obj) SetWidth(w int32) { C.lv_obj_set_width(o.c, C.int32_t(w)) }

// SetHeight sets the object's height.
func (o *Obj) SetHeight(h int32) { C.lv_obj_set_height(o.c, C.int32_t(h)) }

// SetAlign sets the object's alignment relative to its parent, without
// moving it yet (see also Align).
func (o *Obj) SetAlign(a Align) {
	C.lv_obj_set_align(o.c, C.lv_align_t(a))
}

// AlignTo aligns the object relative to base with the given pixel offset.
func (o *Obj) AlignTo(base *Obj, a Align, xOfs, yOfs int32) {
	C.lv_obj_align_to(o.c, base.c, C.lv_align_t(a), C.int32_t(xOfs), C.int32_t(yOfs))
}

// AlignObj aligns the object relative to its parent with the given pixel
// offset (lv_obj_align).
func (o *Obj) AlignObj(a Align, xOfs, yOfs int32) {
	C.lv_obj_align(o.c, C.lv_align_t(a), C.int32_t(xOfs), C.int32_t(yOfs))
}

// Center centers the object within its parent.
func (o *Obj) Center() {
	C.lv_obj_center(o.c)
}

// X returns the object's x coordinate.
func (o *Obj) X() int32 { return int32(C.lv_obj_get_x(o.c)) }

// Y returns the object's y coordinate.
func (o *Obj) Y() int32 { return int32(C.lv_obj_get_y(o.c)) }

// Width returns the object's width.
func (o *Obj) Width() int32 { return int32(C.lv_obj_get_width(o.c)) }

// Height returns the object's height.
func (o *Obj) Height() int32 { return int32(C.lv_obj_get_height(o.c)) }

// Parent returns the object's parent, or nil if it is a screen.
func (o *Obj) Parent() *Obj { return wrapObj(C.lv_obj_get_parent(o.c)) }

// Screen returns the screen the object belongs to.
func (o *Obj) Screen() *Obj { return wrapObj(C.lv_obj_get_screen(o.c)) }

// ChildCount returns the number of direct children.
func (o *Obj) ChildCount() uint32 { return uint32(C.lv_obj_get_child_count(o.c)) }

// Child returns the idx-th direct child.
func (o *Obj) Child(idx int32) *Obj { return wrapObj(C.lv_obj_get_child(o.c, C.int32_t(idx))) }

// AddState sets the given state flag(s) on the object (e.g. StateChecked).
func (o *Obj) AddState(s State) { C.lv_obj_add_state(o.c, C.lv_state_t(s)) }

// RemoveState clears the given state flag(s).
func (o *Obj) RemoveState(s State) { C.lv_obj_remove_state(o.c, C.lv_state_t(s)) }

// HasState reports whether the given state flag(s) are all set.
func (o *Obj) HasState(s State) bool { return bool(C.lv_obj_has_state(o.c, C.lv_state_t(s))) }

// Same reports whether o and other wrap the same underlying LVGL object.
func (o *Obj) Same(other *Obj) bool {
	return o != nil && other != nil && o.c == other.c
}

// HandleCount returns the number of event/anim callback handles currently
// tracked for this object (mainly useful for tests/debugging).
func (o *Obj) HandleCount() int { return len(o.handles) }
