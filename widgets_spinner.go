package lvgl

/*
#include <lvgl.h>
*/
import "C"

// Spinner wraps an lv_spinner widget (an indeterminate loading indicator).
type Spinner struct{ *Obj }

// NewSpinner creates a spinner as a child of parent.
func NewSpinner(parent *Obj) *Spinner {
	return &Spinner{wrapObj(C.lv_spinner_create(parent.c))}
}

// SetAnimDuration sets how long, in milliseconds, one full rotation takes.
func (s *Spinner) SetAnimDuration(ms uint32) { C.lv_spinner_set_anim_duration(s.c, C.uint32_t(ms)) }

// SetArcSweep sets the angular length, in degrees, of the spinning arc.
func (s *Spinner) SetArcSweep(degrees uint32) {
	C.lv_spinner_set_arc_sweep(s.c, C.uint32_t(degrees))
}

// SetAnimParams sets the rotation duration (ms) and arc sweep (degrees)
// together in one call.
func (s *Spinner) SetAnimParams(durationMs, sweepDegrees uint32) {
	C.lv_spinner_set_anim_params(s.c, C.uint32_t(durationMs), C.uint32_t(sweepDegrees))
}

// AnimDuration returns the current rotation duration, in milliseconds.
func (s *Spinner) AnimDuration() uint32 { return uint32(C.lv_spinner_get_anim_duration(s.c)) }

// ArcSweep returns the current arc sweep, in degrees.
func (s *Spinner) ArcSweep() uint32 { return uint32(C.lv_spinner_get_arc_sweep(s.c)) }
