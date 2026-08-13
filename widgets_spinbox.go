package lvgl

/*
#include <lvgl.h>
*/
import "C"

// SpinBox wraps an lv_spinbox widget (a numeric input with +/- stepping).
type SpinBox struct{ *Obj }

// NewSpinBox creates a spinbox as a child of parent.
func NewSpinBox(parent *Obj) *SpinBox {
	return &SpinBox{wrapObj(C.lv_spinbox_create(parent.c))}
}

// SetValue sets the spinbox's numeric value.
func (s *SpinBox) SetValue(v int32) { C.lv_spinbox_set_value(s.c, C.int32_t(v)) }

// SetRange sets the min/max value range.
func (s *SpinBox) SetRange(min, max int32) {
	C.lv_spinbox_set_range(s.c, C.int32_t(min), C.int32_t(max))
}

// SetDigitFormat sets the total digit count and the decimal separator's
// position from the right (0 = no decimal point).
func (s *SpinBox) SetDigitFormat(digitCount, sepPos uint32) {
	C.lv_spinbox_set_digit_format(s.c, C.uint32_t(digitCount), C.uint32_t(sepPos))
}

// SetStep sets how much StepNext/StepPrev change the value by.
func (s *SpinBox) SetStep(step uint32) { C.lv_spinbox_set_step(s.c, C.uint32_t(step)) }

// StepNext increments the value by one step.
func (s *SpinBox) StepNext() { C.lv_spinbox_step_next(s.c) }

// StepPrev decrements the value by one step.
func (s *SpinBox) StepPrev() { C.lv_spinbox_step_prev(s.c) }

// SetRollover enables/disables wrapping from max back to min (and vice
// versa) when stepping past the range.
func (s *SpinBox) SetRollover(enable bool) { C.lv_spinbox_set_rollover(s.c, C.bool(enable)) }

// SetDigitCount sets the total number of digits shown.
func (s *SpinBox) SetDigitCount(count uint32) { C.lv_spinbox_set_digit_count(s.c, C.uint32_t(count)) }

// SetDecPointPos sets the decimal separator's position from the right
// (0 = no decimal point).
func (s *SpinBox) SetDecPointPos(pos uint32) { C.lv_spinbox_set_dec_point_pos(s.c, C.uint32_t(pos)) }

// SetMinValue sets the minimum value.
func (s *SpinBox) SetMinValue(min int32) { C.lv_spinbox_set_min_value(s.c, C.int32_t(min)) }

// SetMaxValue sets the maximum value.
func (s *SpinBox) SetMaxValue(max int32) { C.lv_spinbox_set_max_value(s.c, C.int32_t(max)) }

// SetCursorPos moves the cursor to the given digit position.
func (s *SpinBox) SetCursorPos(pos uint32) { C.lv_spinbox_set_cursor_pos(s.c, C.uint32_t(pos)) }

// SetDigitStepDirection sets which direction (left/right) the cursor
// moves after a StepNext/StepPrev.
func (s *SpinBox) SetDigitStepDirection(dir Dir) {
	C.lv_spinbox_set_digit_step_direction(s.c, C.lv_dir_t(dir))
}

// Value returns the spinbox's current numeric value.
func (s *SpinBox) Value() int32 { return int32(C.lv_spinbox_get_value(s.c)) }

// Rollover reports whether value rollover is enabled.
func (s *SpinBox) Rollover() bool { return bool(C.lv_spinbox_get_rollover(s.c)) }

// Step returns the current step amount.
func (s *SpinBox) Step() int32 { return int32(C.lv_spinbox_get_step(s.c)) }

// DigitCount returns the total number of digits shown.
func (s *SpinBox) DigitCount() uint32 { return uint32(C.lv_spinbox_get_digit_count(s.c)) }

// DecPointPos returns the decimal separator's position from the right.
func (s *SpinBox) DecPointPos() uint32 { return uint32(C.lv_spinbox_get_dec_point_pos(s.c)) }

// MinValue returns the minimum value.
func (s *SpinBox) MinValue() int32 { return int32(C.lv_spinbox_get_min_value(s.c)) }

// MaxValue returns the maximum value.
func (s *SpinBox) MaxValue() int32 { return int32(C.lv_spinbox_get_max_value(s.c)) }

// BindValue links the spinbox's value to an int Subject.
func (s *SpinBox) BindValue(subject *Subject) *Observer {
	return &Observer{c: C.lv_spinbox_bind_value(s.c, subject.c)}
}
