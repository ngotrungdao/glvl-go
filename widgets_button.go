package lvgl

/*
#include <lvgl.h>
*/
import "C"

// Button wraps an lv_button widget. It has no dedicated API beyond the
// base Obj; typical use is creating a child (e.g. a Label) inside it and
// listening for EventClicked.
type Button struct{ *Obj }

// NewButton creates a button as a child of parent.
func NewButton(parent *Obj) *Button {
	return &Button{wrapObj(C.lv_button_create(parent.c))}
}
