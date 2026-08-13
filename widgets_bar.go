package lvgl

/*
#include <lvgl.h>
*/
import "C"

// BarMode mirrors lv_bar_mode_t.
type BarMode uint32

var (
	BarModeNormal      = BarMode(C.LV_BAR_MODE_NORMAL)
	BarModeSymmetrical = BarMode(C.LV_BAR_MODE_SYMMETRICAL)
	BarModeRange       = BarMode(C.LV_BAR_MODE_RANGE)
)

// BarOrientation mirrors lv_bar_orientation_t.
type BarOrientation uint32

var (
	BarOrientationAuto       = BarOrientation(C.LV_BAR_ORIENTATION_AUTO)
	BarOrientationHorizontal = BarOrientation(C.LV_BAR_ORIENTATION_HORIZONTAL)
	BarOrientationVertical   = BarOrientation(C.LV_BAR_ORIENTATION_VERTICAL)
)

// Bar wraps an lv_bar widget.
type Bar struct{ *Obj }

// NewBar creates a bar as a child of parent.
func NewBar(parent *Obj) *Bar {
	return &Bar{wrapObj(C.lv_bar_create(parent.c))}
}

// SetValue sets the bar's value, optionally animating to it.
func (b *Bar) SetValue(v int32, animate bool) {
	C.lv_bar_set_value(b.c, C.int32_t(v), C.lv_anim_enable_t(animate))
}

// SetStartValue sets the lower value in BarModeRange mode.
func (b *Bar) SetStartValue(v int32, animate bool) {
	C.lv_bar_set_start_value(b.c, C.int32_t(v), C.lv_anim_enable_t(animate))
}

// SetRange sets the bar's min/max range.
func (b *Bar) SetRange(min, max int32) { C.lv_bar_set_range(b.c, C.int32_t(min), C.int32_t(max)) }

// SetMinValue sets the bar's minimum value.
func (b *Bar) SetMinValue(min int32) { C.lv_bar_set_min_value(b.c, C.int32_t(min)) }

// SetMaxValue sets the bar's maximum value.
func (b *Bar) SetMaxValue(max int32) { C.lv_bar_set_max_value(b.c, C.int32_t(max)) }

// SetMode sets the bar's mode (normal, symmetrical, or range).
func (b *Bar) SetMode(mode BarMode) { C.lv_bar_set_mode(b.c, C.lv_bar_mode_t(mode)) }

// SetOrientation sets the bar's orientation.
func (b *Bar) SetOrientation(o BarOrientation) {
	C.lv_bar_set_orientation(b.c, C.lv_bar_orientation_t(o))
}

// Value returns the bar's current value.
func (b *Bar) Value() int32 { return int32(C.lv_bar_get_value(b.c)) }

// StartValue returns the lower value in BarModeRange mode.
func (b *Bar) StartValue() int32 { return int32(C.lv_bar_get_start_value(b.c)) }

// MinValue returns the bar's minimum value.
func (b *Bar) MinValue() int32 { return int32(C.lv_bar_get_min_value(b.c)) }

// MaxValue returns the bar's maximum value.
func (b *Bar) MaxValue() int32 { return int32(C.lv_bar_get_max_value(b.c)) }

// Mode returns the bar's current mode.
func (b *Bar) Mode() BarMode { return BarMode(C.lv_bar_get_mode(b.c)) }

// Orientation returns the bar's current orientation.
func (b *Bar) Orientation() BarOrientation { return BarOrientation(C.lv_bar_get_orientation(b.c)) }

// BindValue links the bar's value to an int Subject.
func (b *Bar) BindValue(subject *Subject) *Observer {
	return &Observer{c: C.lv_bar_bind_value(b.c, subject.c)}
}
