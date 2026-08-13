package lvgl

/*
#include <stdlib.h>
#include <lvgl.h>
*/
import "C"
import "unsafe"

// RollerMode mirrors lv_roller_mode_t.
type RollerMode uint32

var (
	RollerModeNormal   = RollerMode(C.LV_ROLLER_MODE_NORMAL)
	RollerModeInfinite = RollerMode(C.LV_ROLLER_MODE_INFINITE)
)

// Roller wraps an lv_roller widget.
//
// Unlike lv_label_set_text or lv_dropdown_set_options (which have a
// documented "_static" non-copying counterpart, implying the default
// copies), lv_roller_set_options has no such counterpart to confirm
// either way, so SetOptions conservatively keeps its C string pinned on
// the C heap for as long as it's current, freeing the previous one (if
// any) when replaced.
type Roller struct {
	*Obj
	options *C.char
}

// NewRoller creates a roller as a child of parent.
func NewRoller(parent *Obj) *Roller {
	return &Roller{Obj: wrapObj(C.lv_roller_create(parent.c))}
}

// SetOptions sets the list of selectable options, newline-separated.
func (r *Roller) SetOptions(options string, mode RollerMode) {
	cOpt := C.CString(options)
	C.lv_roller_set_options(r.c, cOpt, C.lv_roller_mode_t(mode))
	if r.options != nil {
		C.free(unsafe.Pointer(r.options))
	}
	r.options = cOpt
}

// SetSelected selects the option at the given index.
func (r *Roller) SetSelected(idx uint32, animate bool) {
	C.lv_roller_set_selected(r.c, C.uint32_t(idx), C.lv_anim_enable_t(animate))
}

// SetVisibleRowCount sets how many rows are shown at once.
func (r *Roller) SetVisibleRowCount(count uint32) {
	C.lv_roller_set_visible_row_count(r.c, C.uint32_t(count))
}

// Selected returns the index of the currently selected option.
func (r *Roller) Selected() uint32 {
	return uint32(C.lv_roller_get_selected(r.c))
}

// Options returns the newline-separated option list.
func (r *Roller) Options() string {
	return C.GoString(C.lv_roller_get_options(r.c))
}

// SelectedStr returns the text of the currently selected option.
func (r *Roller) SelectedStr() string {
	const bufSize = 128
	buf := (*C.char)(C.malloc(bufSize))
	defer C.free(unsafe.Pointer(buf))
	C.lv_roller_get_selected_str(r.c, buf, bufSize)
	return C.GoString(buf)
}

// OptionCount returns how many options there are.
func (r *Roller) OptionCount() uint32 { return uint32(C.lv_roller_get_option_count(r.c)) }

// BindValue links the roller's selected index to an int Subject.
func (r *Roller) BindValue(subject *Subject) *Observer {
	return &Observer{c: C.lv_roller_bind_value(r.c, subject.c)}
}
