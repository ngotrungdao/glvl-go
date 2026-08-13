package lvgl

/*
#include <lvgl.h>
*/
import "C"

// ScaleMode mirrors lv_scale_mode_t.
type ScaleMode uint32

var (
	ScaleModeHorizontalTop    = ScaleMode(C.LV_SCALE_MODE_HORIZONTAL_TOP)
	ScaleModeHorizontalBottom = ScaleMode(C.LV_SCALE_MODE_HORIZONTAL_BOTTOM)
	ScaleModeVerticalLeft     = ScaleMode(C.LV_SCALE_MODE_VERTICAL_LEFT)
	ScaleModeVerticalRight    = ScaleMode(C.LV_SCALE_MODE_VERTICAL_RIGHT)
	ScaleModeRoundInner       = ScaleMode(C.LV_SCALE_MODE_ROUND_INNER)
	ScaleModeRoundOuter       = ScaleMode(C.LV_SCALE_MODE_ROUND_OUTER)
)

// Scale wraps an lv_scale widget: a ruler/gauge of tick marks and labels,
// often paired with an Arc or Line "needle".
type Scale struct{ *Obj }

// NewScale creates a scale as a child of parent.
func NewScale(parent *Obj) *Scale {
	return &Scale{wrapObj(C.lv_scale_create(parent.c))}
}

// SetMode sets the scale's orientation/shape.
func (s *Scale) SetMode(mode ScaleMode) { C.lv_scale_set_mode(s.c, C.lv_scale_mode_t(mode)) }

// SetRange sets the scale's min/max value range.
func (s *Scale) SetRange(min, max int32) { C.lv_scale_set_range(s.c, C.int32_t(min), C.int32_t(max)) }

// SetTotalTickCount sets how many ticks are drawn in total.
func (s *Scale) SetTotalTickCount(n uint32) { C.lv_scale_set_total_tick_count(s.c, C.uint32_t(n)) }

// SetMajorTickEvery sets how often (every Nth tick) a major tick is drawn.
func (s *Scale) SetMajorTickEvery(n uint32) { C.lv_scale_set_major_tick_every(s.c, C.uint32_t(n)) }

// SetLabelShow enables/disables tick value labels.
func (s *Scale) SetLabelShow(show bool) { C.lv_scale_set_label_show(s.c, C.bool(show)) }

// SetAngleRange sets the angular span, in degrees, for round modes.
func (s *Scale) SetAngleRange(degrees uint32) {
	C.lv_scale_set_angle_range(s.c, C.uint32_t(degrees))
}

// SetRotation rotates a round scale by the given degrees.
func (s *Scale) SetRotation(rotation int32) { C.lv_scale_set_rotation(s.c, C.int32_t(rotation)) }

// SetMinValue sets the scale's minimum value.
func (s *Scale) SetMinValue(min int32) { C.lv_scale_set_min_value(s.c, C.int32_t(min)) }

// SetMaxValue sets the scale's maximum value.
func (s *Scale) SetMaxValue(max int32) { C.lv_scale_set_max_value(s.c, C.int32_t(max)) }

// SetLineNeedleValue points a Line widget's endpoint at the given value,
// as a needle of the given pixel length from the scale's pivot.
func (s *Scale) SetLineNeedleValue(needle *Line, length int32, value int32) {
	C.lv_scale_set_line_needle_value(s.c, needle.c, C.int32_t(length), C.int32_t(value))
}

// SetImageNeedleValue rotates an Image widget around the scale's pivot to
// point at the given value.
func (s *Scale) SetImageNeedleValue(needle *Image, value int32) {
	C.lv_scale_set_image_needle_value(s.c, needle.c, C.int32_t(value))
}

// SetPostDraw enables/disables drawing the scale's ticks/labels after
// (on top of) its children, instead of before.
func (s *Scale) SetPostDraw(enable bool) { C.lv_scale_set_post_draw(s.c, C.bool(enable)) }

// SetDrawTicksOnTop enables/disables drawing ticks above section styling.
func (s *Scale) SetDrawTicksOnTop(enable bool) {
	C.lv_scale_set_draw_ticks_on_top(s.c, C.bool(enable))
}

// TotalTickCount returns the total number of ticks.
func (s *Scale) TotalTickCount() int32 { return int32(C.lv_scale_get_total_tick_count(s.c)) }

// MajorTickEvery returns the major-tick interval.
func (s *Scale) MajorTickEvery() int32 { return int32(C.lv_scale_get_major_tick_every(s.c)) }

// Rotation returns the scale's current rotation, in degrees.
func (s *Scale) Rotation() int32 { return int32(C.lv_scale_get_rotation(s.c)) }

// LabelShow reports whether tick value labels are shown.
func (s *Scale) LabelShow() bool { return bool(C.lv_scale_get_label_show(s.c)) }

// AngleRange returns the angular span, in degrees, for round modes.
func (s *Scale) AngleRange() uint32 { return uint32(C.lv_scale_get_angle_range(s.c)) }

// RangeMinValue returns the scale's minimum value.
func (s *Scale) RangeMinValue() int32 { return int32(C.lv_scale_get_range_min_value(s.c)) }

// RangeMaxValue returns the scale's maximum value.
func (s *Scale) RangeMaxValue() int32 { return int32(C.lv_scale_get_range_max_value(s.c)) }

// ScaleSection wraps an lv_scale_section_t: a value sub-range of the
// scale that can be styled differently (e.g. a red "danger zone"). It is
// owned by the Scale it was added to.
type ScaleSection struct{ c *C.lv_scale_section_t }

// AddSection adds a new, unstyled section.
func (s *Scale) AddSection() *ScaleSection {
	sec := C.lv_scale_add_section(s.c)
	if sec == nil {
		return nil
	}
	return &ScaleSection{c: sec}
}

// SetSectionRange sets a section's value sub-range.
func (s *Scale) SetSectionRange(sec *ScaleSection, min, max int32) {
	C.lv_scale_set_section_range(s.c, sec.c, C.int32_t(min), C.int32_t(max))
}

// SetSectionStyleMain sets the style applied to the scale's main line
// within a section's range.
func (s *Scale) SetSectionStyleMain(sec *ScaleSection, style *Style) {
	C.lv_scale_set_section_style_main(s.c, sec.c, style.c)
}

// SetSectionStyleIndicator sets the style applied to major ticks within a
// section's range.
func (s *Scale) SetSectionStyleIndicator(sec *ScaleSection, style *Style) {
	C.lv_scale_set_section_style_indicator(s.c, sec.c, style.c)
}

// SetSectionStyleItems sets the style applied to minor ticks within a
// section's range.
func (s *Scale) SetSectionStyleItems(sec *ScaleSection, style *Style) {
	C.lv_scale_set_section_style_items(s.c, sec.c, style.c)
}
