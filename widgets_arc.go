package lvgl

/*
#include <lvgl.h>
*/
import "C"

// ArcMode mirrors lv_arc_mode_t.
type ArcMode uint32

var (
	ArcModeNormal      = ArcMode(C.LV_ARC_MODE_NORMAL)
	ArcModeSymmetrical = ArcMode(C.LV_ARC_MODE_SYMMETRICAL)
	ArcModeReverse     = ArcMode(C.LV_ARC_MODE_REVERSE)
)

// Arc wraps an lv_arc widget.
type Arc struct{ *Obj }

// NewArc creates an arc as a child of parent.
func NewArc(parent *Obj) *Arc {
	return &Arc{wrapObj(C.lv_arc_create(parent.c))}
}

// SetValue sets the arc's value (position of the knob along the arc).
func (a *Arc) SetValue(v int32) { C.lv_arc_set_value(a.c, C.int32_t(v)) }

// SetRange sets the arc's min/max value range.
func (a *Arc) SetRange(min, max int32) { C.lv_arc_set_range(a.c, C.int32_t(min), C.int32_t(max)) }

// SetMinValue sets the arc's minimum value.
func (a *Arc) SetMinValue(min int32) { C.lv_arc_set_min_value(a.c, C.int32_t(min)) }

// SetMaxValue sets the arc's maximum value.
func (a *Arc) SetMaxValue(max int32) { C.lv_arc_set_max_value(a.c, C.int32_t(max)) }

// SetAngles sets the start/end angles (in degrees) of the arc itself.
func (a *Arc) SetAngles(start, end float32) {
	C.lv_arc_set_angles(a.c, C.lv_value_precise_t(start), C.lv_value_precise_t(end))
}

// SetBgAngles sets the start/end angles (in degrees) of the arc's
// background track.
func (a *Arc) SetBgAngles(start, end float32) {
	C.lv_arc_set_bg_angles(a.c, C.lv_value_precise_t(start), C.lv_value_precise_t(end))
}

// SetRotation rotates the whole arc by the given degrees.
func (a *Arc) SetRotation(rotation int32) { C.lv_arc_set_rotation(a.c, C.int32_t(rotation)) }

// SetMode sets how the value maps to the arc's drawn angle range.
func (a *Arc) SetMode(mode ArcMode) { C.lv_arc_set_mode(a.c, C.lv_arc_mode_t(mode)) }

// SetChangeRate sets how many value units per second the arc can change
// by when dragged (limits perceived "jumpiness").
func (a *Arc) SetChangeRate(rate uint32) { C.lv_arc_set_change_rate(a.c, C.uint32_t(rate)) }

// SetKnobOffset offsets the knob's angular position from the value angle.
func (a *Arc) SetKnobOffset(offset int32) { C.lv_arc_set_knob_offset(a.c, C.int32_t(offset)) }

// Value returns the arc's current value.
func (a *Arc) Value() int32 { return int32(C.lv_arc_get_value(a.c)) }

// MinValue returns the arc's minimum value.
func (a *Arc) MinValue() int32 { return int32(C.lv_arc_get_min_value(a.c)) }

// MaxValue returns the arc's maximum value.
func (a *Arc) MaxValue() int32 { return int32(C.lv_arc_get_max_value(a.c)) }

// AngleStart returns the arc's start angle, in degrees.
func (a *Arc) AngleStart() float32 { return float32(C.lv_arc_get_angle_start(a.c)) }

// AngleEnd returns the arc's end angle, in degrees.
func (a *Arc) AngleEnd() float32 { return float32(C.lv_arc_get_angle_end(a.c)) }

// BgAngleStart returns the background track's start angle, in degrees.
func (a *Arc) BgAngleStart() float32 { return float32(C.lv_arc_get_bg_angle_start(a.c)) }

// BgAngleEnd returns the background track's end angle, in degrees.
func (a *Arc) BgAngleEnd() float32 { return float32(C.lv_arc_get_bg_angle_end(a.c)) }

// Mode returns the arc's current mode.
func (a *Arc) Mode() ArcMode { return ArcMode(C.lv_arc_get_mode(a.c)) }

// Rotation returns the arc's rotation, in degrees.
func (a *Arc) Rotation() int32 { return int32(C.lv_arc_get_rotation(a.c)) }

// KnobOffset returns the knob's angular offset.
func (a *Arc) KnobOffset() int32 { return int32(C.lv_arc_get_knob_offset(a.c)) }

// ChangeRate returns the arc's drag change-rate limit.
func (a *Arc) ChangeRate() uint32 { return uint32(C.lv_arc_get_change_rate(a.c)) }

// BindValue links the arc's value to an int Subject.
func (a *Arc) BindValue(subject *Subject) *Observer {
	return &Observer{c: C.lv_arc_bind_value(a.c, subject.c)}
}
