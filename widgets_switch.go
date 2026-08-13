package lvgl

/*
#include <lvgl.h>
*/
import "C"

// Switch wraps an lv_switch widget.
type Switch struct{ *Obj }

// NewSwitch creates a switch as a child of parent.
func NewSwitch(parent *Obj) *Switch {
	return &Switch{wrapObj(C.lv_switch_create(parent.c))}
}

// On reports whether the switch is currently on.
func (s *Switch) On() bool {
	return s.HasState(StateChecked)
}
