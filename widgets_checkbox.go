package lvgl

/*
#include <stdlib.h>
#include <lvgl.h>
*/
import "C"
import "unsafe"

// Checkbox wraps an lv_checkbox widget.
type Checkbox struct{ *Obj }

// NewCheckbox creates a checkbox as a child of parent.
func NewCheckbox(parent *Obj) *Checkbox {
	return &Checkbox{wrapObj(C.lv_checkbox_create(parent.c))}
}

// SetText sets the checkbox's label text.
func (c *Checkbox) SetText(text string) {
	cText := C.CString(text)
	defer C.free(unsafe.Pointer(cText))
	C.lv_checkbox_set_text(c.c, cText)
}

// Text returns the checkbox's label text.
func (c *Checkbox) Text() string {
	return C.GoString(C.lv_checkbox_get_text(c.c))
}

// Checked reports whether the checkbox is currently checked.
func (c *Checkbox) Checked() bool {
	return c.HasState(StateChecked)
}
