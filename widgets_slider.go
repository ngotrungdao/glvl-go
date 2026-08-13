package lvgl

/*
#include <lvgl.h>
*/
import "C"

// SliderMode mirrors lv_slider_mode_t.
type SliderMode uint32

var (
	SliderModeNormal      = SliderMode(C.LV_SLIDER_MODE_NORMAL)
	SliderModeSymmetrical = SliderMode(C.LV_SLIDER_MODE_SYMMETRICAL)
	SliderModeRange       = SliderMode(C.LV_SLIDER_MODE_RANGE)
)

// SliderOrientation mirrors lv_slider_orientation_t.
type SliderOrientation uint32

var (
	SliderOrientationAuto       = SliderOrientation(C.LV_SLIDER_ORIENTATION_AUTO)
	SliderOrientationHorizontal = SliderOrientation(C.LV_SLIDER_ORIENTATION_HORIZONTAL)
	SliderOrientationVertical   = SliderOrientation(C.LV_SLIDER_ORIENTATION_VERTICAL)
)

// Slider wraps an lv_slider widget.
type Slider struct{ *Obj }

// NewSlider creates a slider as a child of parent.
func NewSlider(parent *Obj) *Slider {
	return &Slider{wrapObj(C.lv_slider_create(parent.c))}
}

// SetValue sets the slider's value, optionally animating to it.
func (s *Slider) SetValue(value int32, animate bool) {
	C.lv_slider_set_value(s.c, C.int32_t(value), C.lv_anim_enable_t(animate))
}

// SetStartValue sets the lower value in SliderModeRange mode.
func (s *Slider) SetStartValue(value int32, animate bool) {
	C.lv_slider_set_start_value(s.c, C.int32_t(value), C.lv_anim_enable_t(animate))
}

// SetRange sets the slider's min/max range.
func (s *Slider) SetRange(min, max int32) {
	C.lv_slider_set_range(s.c, C.int32_t(min), C.int32_t(max))
}

// SetMinValue sets the slider's minimum value.
func (s *Slider) SetMinValue(min int32) { C.lv_slider_set_min_value(s.c, C.int32_t(min)) }

// SetMaxValue sets the slider's maximum value.
func (s *Slider) SetMaxValue(max int32) { C.lv_slider_set_max_value(s.c, C.int32_t(max)) }

// SetMode sets the slider's mode (normal, symmetrical, or range).
func (s *Slider) SetMode(mode SliderMode) { C.lv_slider_set_mode(s.c, C.lv_slider_mode_t(mode)) }

// SetOrientation sets the slider's orientation.
func (s *Slider) SetOrientation(o SliderOrientation) {
	C.lv_slider_set_orientation(s.c, C.lv_slider_orientation_t(o))
}

// Value returns the slider's current value.
func (s *Slider) Value() int32 {
	return int32(C.lv_slider_get_value(s.c))
}

// LeftValue returns the lower value in SliderModeRange mode.
func (s *Slider) LeftValue() int32 { return int32(C.lv_slider_get_left_value(s.c)) }

// MinValue returns the slider's minimum value.
func (s *Slider) MinValue() int32 { return int32(C.lv_slider_get_min_value(s.c)) }

// MaxValue returns the slider's maximum value.
func (s *Slider) MaxValue() int32 { return int32(C.lv_slider_get_max_value(s.c)) }

// IsDragged reports whether the slider is currently being dragged.
func (s *Slider) IsDragged() bool { return bool(C.lv_slider_is_dragged(s.c)) }

// Mode returns the slider's current mode.
func (s *Slider) Mode() SliderMode { return SliderMode(C.lv_slider_get_mode(s.c)) }

// Orientation returns the slider's current orientation.
func (s *Slider) Orientation() SliderOrientation {
	return SliderOrientation(C.lv_slider_get_orientation(s.c))
}

// BindValue links the slider's value to an int Subject: moving the
// slider updates the subject, and changing the subject moves the slider.
func (s *Slider) BindValue(subject *Subject) *Observer {
	return &Observer{c: C.lv_slider_bind_value(s.c, subject.c)}
}
