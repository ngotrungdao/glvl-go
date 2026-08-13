package lvgl

/*
#include <stdlib.h>
#include <lvgl.h>
*/
import "C"
import "unsafe"

// Part mirrors lv_part_t.
type Part uint32

var (
	PartMain      = Part(C.LV_PART_MAIN)
	PartScrollbar = Part(C.LV_PART_SCROLLBAR)
	PartIndicator = Part(C.LV_PART_INDICATOR)
	PartKnob      = Part(C.LV_PART_KNOB)
	PartSelected  = Part(C.LV_PART_SELECTED)
	PartItems     = Part(C.LV_PART_ITEMS)
	PartCursor    = Part(C.LV_PART_CURSOR)
	PartAny       = Part(C.LV_PART_ANY)
)

// Selector mirrors lv_style_selector_t: a Part OR'd with zero or more
// State flags, e.g. Selector(PartKnob) | Selector(StatePressed).
type Selector uint32

// Sel builds a Selector from a Part and zero or more States.
func Sel(part Part, states ...State) Selector {
	sel := Selector(part)
	for _, st := range states {
		sel |= Selector(st)
	}
	return sel
}

// AddStyle attaches a style to the object for the given part/state
// selector. The style must remain valid (not Delete()d) for as long as it
// stays attached.
func (o *Obj) AddStyle(s *Style, sel Selector) {
	C.lv_obj_add_style(o.c, s.c, C.lv_style_selector_t(sel))
}

// RemoveStyle detaches a previously attached style from the given
// part/state selector.
func (o *Obj) RemoveStyle(s *Style, sel Selector) {
	C.lv_obj_remove_style(o.c, s.c, C.lv_style_selector_t(sel))
}

// RemoveStyleAll detaches every style from the object.
func (o *Obj) RemoveStyleAll() {
	C.lv_obj_remove_style_all(o.c)
}

// The setters below mirror core_style.go's coverage (124 of 129 property
// setters — see its doc comment for the 5 not covered), but apply
// directly to the object for a given part/state selector
// (lv_obj_set_style_*, core/lv_obj_style_gen.h) instead of going through a
// reusable Style. Same extension pattern: mechanical 1:1 translation.
//
// One exception: the grid track descriptor array setters
// (lv_obj_set_style_grid_column/row_dsc_array) aren't wrapped here, since
// Obj.SetGridDscArray (layouts_grid.go) already covers the same need via
// the more common direct (non-style) API and adding a second, parallel
// pinning path for the rarely-used style-based variant wasn't worth the
// complexity.

func (o *Obj) SetStyleWidth(v int32, sel Selector) {
	C.lv_obj_set_style_width(o.c, C.int32_t(v), C.lv_style_selector_t(sel))
}
func (o *Obj) SetStyleHeight(v int32, sel Selector) {
	C.lv_obj_set_style_height(o.c, C.int32_t(v), C.lv_style_selector_t(sel))
}
func (o *Obj) SetStylePadTop(v int32, sel Selector) {
	C.lv_obj_set_style_pad_top(o.c, C.int32_t(v), C.lv_style_selector_t(sel))
}
func (o *Obj) SetStylePadBottom(v int32, sel Selector) {
	C.lv_obj_set_style_pad_bottom(o.c, C.int32_t(v), C.lv_style_selector_t(sel))
}
func (o *Obj) SetStylePadLeft(v int32, sel Selector) {
	C.lv_obj_set_style_pad_left(o.c, C.int32_t(v), C.lv_style_selector_t(sel))
}
func (o *Obj) SetStylePadRight(v int32, sel Selector) {
	C.lv_obj_set_style_pad_right(o.c, C.int32_t(v), C.lv_style_selector_t(sel))
}
func (o *Obj) SetStylePadRow(v int32, sel Selector) {
	C.lv_obj_set_style_pad_row(o.c, C.int32_t(v), C.lv_style_selector_t(sel))
}
func (o *Obj) SetStylePadColumn(v int32, sel Selector) {
	C.lv_obj_set_style_pad_column(o.c, C.int32_t(v), C.lv_style_selector_t(sel))
}

// SetStylePadAll sets all four padding sides to the same value.
func (o *Obj) SetStylePadAll(v int32, sel Selector) {
	o.SetStylePadTop(v, sel)
	o.SetStylePadBottom(v, sel)
	o.SetStylePadLeft(v, sel)
	o.SetStylePadRight(v, sel)
}

func (o *Obj) SetStyleBgColor(v Color, sel Selector) {
	C.lv_obj_set_style_bg_color(o.c, v.toC(), C.lv_style_selector_t(sel))
}
func (o *Obj) SetStyleBgOpa(v Opa, sel Selector) {
	C.lv_obj_set_style_bg_opa(o.c, C.lv_opa_t(v), C.lv_style_selector_t(sel))
}
func (o *Obj) SetStyleBgGradColor(v Color, sel Selector) {
	C.lv_obj_set_style_bg_grad_color(o.c, v.toC(), C.lv_style_selector_t(sel))
}

func (o *Obj) SetStyleBorderColor(v Color, sel Selector) {
	C.lv_obj_set_style_border_color(o.c, v.toC(), C.lv_style_selector_t(sel))
}
func (o *Obj) SetStyleBorderOpa(v Opa, sel Selector) {
	C.lv_obj_set_style_border_opa(o.c, C.lv_opa_t(v), C.lv_style_selector_t(sel))
}
func (o *Obj) SetStyleBorderWidth(v int32, sel Selector) {
	C.lv_obj_set_style_border_width(o.c, C.int32_t(v), C.lv_style_selector_t(sel))
}
func (o *Obj) SetStyleBorderSide(v BorderSide, sel Selector) {
	C.lv_obj_set_style_border_side(o.c, C.lv_border_side_t(v), C.lv_style_selector_t(sel))
}

func (o *Obj) SetStyleOutlineWidth(v int32, sel Selector) {
	C.lv_obj_set_style_outline_width(o.c, C.int32_t(v), C.lv_style_selector_t(sel))
}
func (o *Obj) SetStyleOutlineColor(v Color, sel Selector) {
	C.lv_obj_set_style_outline_color(o.c, v.toC(), C.lv_style_selector_t(sel))
}
func (o *Obj) SetStyleOutlineOpa(v Opa, sel Selector) {
	C.lv_obj_set_style_outline_opa(o.c, C.lv_opa_t(v), C.lv_style_selector_t(sel))
}
func (o *Obj) SetStyleOutlinePad(v int32, sel Selector) {
	C.lv_obj_set_style_outline_pad(o.c, C.int32_t(v), C.lv_style_selector_t(sel))
}

func (o *Obj) SetStyleShadowWidth(v int32, sel Selector) {
	C.lv_obj_set_style_shadow_width(o.c, C.int32_t(v), C.lv_style_selector_t(sel))
}
func (o *Obj) SetStyleShadowSpread(v int32, sel Selector) {
	C.lv_obj_set_style_shadow_spread(o.c, C.int32_t(v), C.lv_style_selector_t(sel))
}
func (o *Obj) SetStyleShadowOfsX(v int32, sel Selector) {
	C.lv_obj_set_style_shadow_offset_x(o.c, C.int32_t(v), C.lv_style_selector_t(sel))
}
func (o *Obj) SetStyleShadowOfsY(v int32, sel Selector) {
	C.lv_obj_set_style_shadow_offset_y(o.c, C.int32_t(v), C.lv_style_selector_t(sel))
}
func (o *Obj) SetStyleShadowColor(v Color, sel Selector) {
	C.lv_obj_set_style_shadow_color(o.c, v.toC(), C.lv_style_selector_t(sel))
}
func (o *Obj) SetStyleShadowOpa(v Opa, sel Selector) {
	C.lv_obj_set_style_shadow_opa(o.c, C.lv_opa_t(v), C.lv_style_selector_t(sel))
}

func (o *Obj) SetStyleTextColor(v Color, sel Selector) {
	C.lv_obj_set_style_text_color(o.c, v.toC(), C.lv_style_selector_t(sel))
}
func (o *Obj) SetStyleTextOpa(v Opa, sel Selector) {
	C.lv_obj_set_style_text_opa(o.c, C.lv_opa_t(v), C.lv_style_selector_t(sel))
}
func (o *Obj) SetStyleTextLetterSpace(v int32, sel Selector) {
	C.lv_obj_set_style_text_letter_space(o.c, C.int32_t(v), C.lv_style_selector_t(sel))
}
func (o *Obj) SetStyleTextLineSpace(v int32, sel Selector) {
	C.lv_obj_set_style_text_line_space(o.c, C.int32_t(v), C.lv_style_selector_t(sel))
}
func (o *Obj) SetStyleTextAlign(v TextAlign, sel Selector) {
	C.lv_obj_set_style_text_align(o.c, C.lv_text_align_t(v), C.lv_style_selector_t(sel))
}

// SetStyleTextFont sets both the font family and size directly on the
// object (see Font's doc comment).
func (o *Obj) SetStyleTextFont(f *Font, sel Selector) {
	C.lv_obj_set_style_text_font(o.c, f.c, C.lv_style_selector_t(sel))
}

func (o *Obj) SetStyleLineColor(v Color, sel Selector) {
	C.lv_obj_set_style_line_color(o.c, v.toC(), C.lv_style_selector_t(sel))
}
func (o *Obj) SetStyleLineWidth(v int32, sel Selector) {
	C.lv_obj_set_style_line_width(o.c, C.int32_t(v), C.lv_style_selector_t(sel))
}
func (o *Obj) SetStyleLineOpa(v Opa, sel Selector) {
	C.lv_obj_set_style_line_opa(o.c, C.lv_opa_t(v), C.lv_style_selector_t(sel))
}

func (o *Obj) SetStyleImageRecolor(v Color, sel Selector) {
	C.lv_obj_set_style_image_recolor(o.c, v.toC(), C.lv_style_selector_t(sel))
}
func (o *Obj) SetStyleImageRecolorOpa(v Opa, sel Selector) {
	C.lv_obj_set_style_image_recolor_opa(o.c, C.lv_opa_t(v), C.lv_style_selector_t(sel))
}
func (o *Obj) SetStyleImageOpa(v Opa, sel Selector) {
	C.lv_obj_set_style_image_opa(o.c, C.lv_opa_t(v), C.lv_style_selector_t(sel))
}

func (o *Obj) SetStyleRadius(v int32, sel Selector) {
	C.lv_obj_set_style_radius(o.c, C.int32_t(v), C.lv_style_selector_t(sel))
}
func (o *Obj) SetStyleOpa(v Opa, sel Selector) {
	C.lv_obj_set_style_opa(o.c, C.lv_opa_t(v), C.lv_style_selector_t(sel))
}

func (o *Obj) SetStyleMinWidth(v int32, sel Selector) {
	C.lv_obj_set_style_min_width(o.c, C.int32_t(v), C.lv_style_selector_t(sel))
}
func (o *Obj) SetStyleMaxWidth(v int32, sel Selector) {
	C.lv_obj_set_style_max_width(o.c, C.int32_t(v), C.lv_style_selector_t(sel))
}
func (o *Obj) SetStyleMinHeight(v int32, sel Selector) {
	C.lv_obj_set_style_min_height(o.c, C.int32_t(v), C.lv_style_selector_t(sel))
}
func (o *Obj) SetStyleMaxHeight(v int32, sel Selector) {
	C.lv_obj_set_style_max_height(o.c, C.int32_t(v), C.lv_style_selector_t(sel))
}
func (o *Obj) SetStyleLength(v int32, sel Selector) {
	C.lv_obj_set_style_length(o.c, C.int32_t(v), C.lv_style_selector_t(sel))
}
func (o *Obj) SetStyleX(v int32, sel Selector) {
	C.lv_obj_set_style_x(o.c, C.int32_t(v), C.lv_style_selector_t(sel))
}
func (o *Obj) SetStyleY(v int32, sel Selector) {
	C.lv_obj_set_style_y(o.c, C.int32_t(v), C.lv_style_selector_t(sel))
}
func (o *Obj) SetStyleAlign(v Align, sel Selector) {
	C.lv_obj_set_style_align(o.c, C.lv_align_t(v), C.lv_style_selector_t(sel))
}

func (o *Obj) SetStyleTransformWidth(v int32, sel Selector) {
	C.lv_obj_set_style_transform_width(o.c, C.int32_t(v), C.lv_style_selector_t(sel))
}
func (o *Obj) SetStyleTransformHeight(v int32, sel Selector) {
	C.lv_obj_set_style_transform_height(o.c, C.int32_t(v), C.lv_style_selector_t(sel))
}
func (o *Obj) SetStyleTranslateX(v int32, sel Selector) {
	C.lv_obj_set_style_translate_x(o.c, C.int32_t(v), C.lv_style_selector_t(sel))
}
func (o *Obj) SetStyleTranslateY(v int32, sel Selector) {
	C.lv_obj_set_style_translate_y(o.c, C.int32_t(v), C.lv_style_selector_t(sel))
}
func (o *Obj) SetStyleTranslateRadial(v int32, sel Selector) {
	C.lv_obj_set_style_translate_radial(o.c, C.int32_t(v), C.lv_style_selector_t(sel))
}
func (o *Obj) SetStyleTransformScaleX(v int32, sel Selector) {
	C.lv_obj_set_style_transform_scale_x(o.c, C.int32_t(v), C.lv_style_selector_t(sel))
}
func (o *Obj) SetStyleTransformScaleY(v int32, sel Selector) {
	C.lv_obj_set_style_transform_scale_y(o.c, C.int32_t(v), C.lv_style_selector_t(sel))
}
func (o *Obj) SetStyleTransformRotation(v int32, sel Selector) {
	C.lv_obj_set_style_transform_rotation(o.c, C.int32_t(v), C.lv_style_selector_t(sel))
}
func (o *Obj) SetStyleTransformPivotX(v int32, sel Selector) {
	C.lv_obj_set_style_transform_pivot_x(o.c, C.int32_t(v), C.lv_style_selector_t(sel))
}
func (o *Obj) SetStyleTransformPivotY(v int32, sel Selector) {
	C.lv_obj_set_style_transform_pivot_y(o.c, C.int32_t(v), C.lv_style_selector_t(sel))
}
func (o *Obj) SetStyleTransformSkewX(v int32, sel Selector) {
	C.lv_obj_set_style_transform_skew_x(o.c, C.int32_t(v), C.lv_style_selector_t(sel))
}
func (o *Obj) SetStyleTransformSkewY(v int32, sel Selector) {
	C.lv_obj_set_style_transform_skew_y(o.c, C.int32_t(v), C.lv_style_selector_t(sel))
}

func (o *Obj) SetStylePadRadial(v int32, sel Selector) {
	C.lv_obj_set_style_pad_radial(o.c, C.int32_t(v), C.lv_style_selector_t(sel))
}
func (o *Obj) SetStyleMarginTop(v int32, sel Selector) {
	C.lv_obj_set_style_margin_top(o.c, C.int32_t(v), C.lv_style_selector_t(sel))
}
func (o *Obj) SetStyleMarginBottom(v int32, sel Selector) {
	C.lv_obj_set_style_margin_bottom(o.c, C.int32_t(v), C.lv_style_selector_t(sel))
}
func (o *Obj) SetStyleMarginLeft(v int32, sel Selector) {
	C.lv_obj_set_style_margin_left(o.c, C.int32_t(v), C.lv_style_selector_t(sel))
}
func (o *Obj) SetStyleMarginRight(v int32, sel Selector) {
	C.lv_obj_set_style_margin_right(o.c, C.int32_t(v), C.lv_style_selector_t(sel))
}

// SetStyleMarginAll sets all four margin sides to the same value.
func (o *Obj) SetStyleMarginAll(v int32, sel Selector) {
	o.SetStyleMarginTop(v, sel)
	o.SetStyleMarginBottom(v, sel)
	o.SetStyleMarginLeft(v, sel)
	o.SetStyleMarginRight(v, sel)
}

func (o *Obj) SetStyleBgGradDir(v GradDir, sel Selector) {
	C.lv_obj_set_style_bg_grad_dir(o.c, C.lv_grad_dir_t(v), C.lv_style_selector_t(sel))
}
func (o *Obj) SetStyleBgMainStop(v int32, sel Selector) {
	C.lv_obj_set_style_bg_main_stop(o.c, C.int32_t(v), C.lv_style_selector_t(sel))
}
func (o *Obj) SetStyleBgGradStop(v int32, sel Selector) {
	C.lv_obj_set_style_bg_grad_stop(o.c, C.int32_t(v), C.lv_style_selector_t(sel))
}
func (o *Obj) SetStyleBgMainOpa(v Opa, sel Selector) {
	C.lv_obj_set_style_bg_main_opa(o.c, C.lv_opa_t(v), C.lv_style_selector_t(sel))
}
func (o *Obj) SetStyleBgGradOpa(v Opa, sel Selector) {
	C.lv_obj_set_style_bg_grad_opa(o.c, C.lv_opa_t(v), C.lv_style_selector_t(sel))
}
func (o *Obj) SetStyleBgImageOpa(v Opa, sel Selector) {
	C.lv_obj_set_style_bg_image_opa(o.c, C.lv_opa_t(v), C.lv_style_selector_t(sel))
}
func (o *Obj) SetStyleBgImageRecolor(v Color, sel Selector) {
	C.lv_obj_set_style_bg_image_recolor(o.c, v.toC(), C.lv_style_selector_t(sel))
}
func (o *Obj) SetStyleBgImageRecolorOpa(v Opa, sel Selector) {
	C.lv_obj_set_style_bg_image_recolor_opa(o.c, C.lv_opa_t(v), C.lv_style_selector_t(sel))
}
func (o *Obj) SetStyleBgImageTiled(v bool, sel Selector) {
	C.lv_obj_set_style_bg_image_tiled(o.c, C.bool(v), C.lv_style_selector_t(sel))
}

// SetStyleBgImageSrc sets a background image directly on the object. Like
// Image.SetSrcPath, the string is pinned on the C heap for the object's
// lifetime (freed by Obj.Delete) rather than copied by LVGL.
func (o *Obj) SetStyleBgImageSrc(path string, sel Selector) {
	c := C.CString(path)
	C.lv_obj_set_style_bg_image_src(o.c, unsafe.Pointer(c), C.lv_style_selector_t(sel))
	if o.bgImageSrc != nil {
		C.free(unsafe.Pointer(o.bgImageSrc))
	}
	o.bgImageSrc = c
}

func (o *Obj) SetStyleBorderPost(v bool, sel Selector) {
	C.lv_obj_set_style_border_post(o.c, C.bool(v), C.lv_style_selector_t(sel))
}

func (o *Obj) SetStyleDropShadowRadius(v int32, sel Selector) {
	C.lv_obj_set_style_drop_shadow_radius(o.c, C.int32_t(v), C.lv_style_selector_t(sel))
}
func (o *Obj) SetStyleDropShadowOfsX(v int32, sel Selector) {
	C.lv_obj_set_style_drop_shadow_offset_x(o.c, C.int32_t(v), C.lv_style_selector_t(sel))
}
func (o *Obj) SetStyleDropShadowOfsY(v int32, sel Selector) {
	C.lv_obj_set_style_drop_shadow_offset_y(o.c, C.int32_t(v), C.lv_style_selector_t(sel))
}
func (o *Obj) SetStyleDropShadowColor(v Color, sel Selector) {
	C.lv_obj_set_style_drop_shadow_color(o.c, v.toC(), C.lv_style_selector_t(sel))
}
func (o *Obj) SetStyleDropShadowOpa(v Opa, sel Selector) {
	C.lv_obj_set_style_drop_shadow_opa(o.c, C.lv_opa_t(v), C.lv_style_selector_t(sel))
}
func (o *Obj) SetStyleDropShadowQuality(v BlurQuality, sel Selector) {
	C.lv_obj_set_style_drop_shadow_quality(o.c, C.lv_blur_quality_t(v), C.lv_style_selector_t(sel))
}

func (o *Obj) SetStyleBlurRadius(v int32, sel Selector) {
	C.lv_obj_set_style_blur_radius(o.c, C.int32_t(v), C.lv_style_selector_t(sel))
}
func (o *Obj) SetStyleBlurBackdrop(v bool, sel Selector) {
	C.lv_obj_set_style_blur_backdrop(o.c, C.bool(v), C.lv_style_selector_t(sel))
}
func (o *Obj) SetStyleBlurQuality(v BlurQuality, sel Selector) {
	C.lv_obj_set_style_blur_quality(o.c, C.lv_blur_quality_t(v), C.lv_style_selector_t(sel))
}

func (o *Obj) SetStyleTextDecor(v TextDecor, sel Selector) {
	C.lv_obj_set_style_text_decor(o.c, C.lv_text_decor_t(v), C.lv_style_selector_t(sel))
}
func (o *Obj) SetStyleTextLeadingTrim(v TextLeadingTrim, sel Selector) {
	C.lv_obj_set_style_text_leading_trim(o.c, C.lv_text_leading_trim_t(v), C.lv_style_selector_t(sel))
}
func (o *Obj) SetStyleTextOutlineStrokeColor(v Color, sel Selector) {
	C.lv_obj_set_style_text_outline_stroke_color(o.c, v.toC(), C.lv_style_selector_t(sel))
}
func (o *Obj) SetStyleTextOutlineStrokeOpa(v Opa, sel Selector) {
	C.lv_obj_set_style_text_outline_stroke_opa(o.c, C.lv_opa_t(v), C.lv_style_selector_t(sel))
}
func (o *Obj) SetStyleTextOutlineStrokeWidth(v int32, sel Selector) {
	C.lv_obj_set_style_text_outline_stroke_width(o.c, C.int32_t(v), C.lv_style_selector_t(sel))
}

func (o *Obj) SetStyleLineDashWidth(v int32, sel Selector) {
	C.lv_obj_set_style_line_dash_width(o.c, C.int32_t(v), C.lv_style_selector_t(sel))
}
func (o *Obj) SetStyleLineDashGap(v int32, sel Selector) {
	C.lv_obj_set_style_line_dash_gap(o.c, C.int32_t(v), C.lv_style_selector_t(sel))
}
func (o *Obj) SetStyleLineRounded(v bool, sel Selector) {
	C.lv_obj_set_style_line_rounded(o.c, C.bool(v), C.lv_style_selector_t(sel))
}

func (o *Obj) SetStyleArcColor(v Color, sel Selector) {
	C.lv_obj_set_style_arc_color(o.c, v.toC(), C.lv_style_selector_t(sel))
}
func (o *Obj) SetStyleArcOpa(v Opa, sel Selector) {
	C.lv_obj_set_style_arc_opa(o.c, C.lv_opa_t(v), C.lv_style_selector_t(sel))
}
func (o *Obj) SetStyleArcWidth(v int32, sel Selector) {
	C.lv_obj_set_style_arc_width(o.c, C.int32_t(v), C.lv_style_selector_t(sel))
}
func (o *Obj) SetStyleArcRounded(v bool, sel Selector) {
	C.lv_obj_set_style_arc_rounded(o.c, C.bool(v), C.lv_style_selector_t(sel))
}

// SetStyleArcImageSrc sets an image used to draw the arc directly on the
// object. Like Image.SetSrcPath, the string is pinned on the C heap for
// the object's lifetime (freed by Obj.Delete) rather than copied by LVGL.
func (o *Obj) SetStyleArcImageSrc(path string, sel Selector) {
	c := C.CString(path)
	C.lv_obj_set_style_arc_image_src(o.c, unsafe.Pointer(c), C.lv_style_selector_t(sel))
	if o.arcImageSrc != nil {
		C.free(unsafe.Pointer(o.arcImageSrc))
	}
	o.arcImageSrc = c
}

func (o *Obj) SetStyleRadialOffset(v int32, sel Selector) {
	C.lv_obj_set_style_radial_offset(o.c, C.int32_t(v), C.lv_style_selector_t(sel))
}
func (o *Obj) SetStyleClipCorner(v bool, sel Selector) {
	C.lv_obj_set_style_clip_corner(o.c, C.bool(v), C.lv_style_selector_t(sel))
}
func (o *Obj) SetStyleOpaLayered(v Opa, sel Selector) {
	C.lv_obj_set_style_opa_layered(o.c, C.lv_opa_t(v), C.lv_style_selector_t(sel))
}
func (o *Obj) SetStyleColorFilterOpa(v Opa, sel Selector) {
	C.lv_obj_set_style_color_filter_opa(o.c, C.lv_opa_t(v), C.lv_style_selector_t(sel))
}
func (o *Obj) SetStyleRecolor(v Color, sel Selector) {
	C.lv_obj_set_style_recolor(o.c, v.toC(), C.lv_style_selector_t(sel))
}
func (o *Obj) SetStyleRecolorOpa(v Opa, sel Selector) {
	C.lv_obj_set_style_recolor_opa(o.c, C.lv_opa_t(v), C.lv_style_selector_t(sel))
}
func (o *Obj) SetStyleAnimDuration(v uint32, sel Selector) {
	C.lv_obj_set_style_anim_duration(o.c, C.uint32_t(v), C.lv_style_selector_t(sel))
}
func (o *Obj) SetStyleBlendMode(v BlendMode, sel Selector) {
	C.lv_obj_set_style_blend_mode(o.c, C.lv_blend_mode_t(v), C.lv_style_selector_t(sel))
}
func (o *Obj) SetStyleBaseDir(v BaseDir, sel Selector) {
	C.lv_obj_set_style_base_dir(o.c, C.lv_base_dir_t(v), C.lv_style_selector_t(sel))
}
func (o *Obj) SetStyleRotarySensitivity(v uint32, sel Selector) {
	C.lv_obj_set_style_rotary_sensitivity(o.c, C.uint32_t(v), C.lv_style_selector_t(sel))
}

// SetStyleBitmapMaskSrc sets an image used as an alpha mask directly on
// the object. Like Image.SetSrcPath, the string is pinned on the C heap
// for the object's lifetime (freed by Obj.Delete) rather than copied by
// LVGL.
func (o *Obj) SetStyleBitmapMaskSrc(path string, sel Selector) {
	c := C.CString(path)
	C.lv_obj_set_style_bitmap_mask_src(o.c, unsafe.Pointer(c), C.lv_style_selector_t(sel))
	if o.bitmapMaskSrc != nil {
		C.free(unsafe.Pointer(o.bitmapMaskSrc))
	}
	o.bitmapMaskSrc = c
}

// SetStyleFlexFlow sets flex flow as part of a style selector, distinct
// from the direct (non-style) Obj.SetFlexFlow in layouts_flex.go.
func (o *Obj) SetStyleFlexFlow(v FlexFlow, sel Selector) {
	C.lv_obj_set_style_flex_flow(o.c, C.lv_flex_flow_t(v), C.lv_style_selector_t(sel))
}
func (o *Obj) SetStyleFlexMainPlace(v FlexAlign, sel Selector) {
	C.lv_obj_set_style_flex_main_place(o.c, C.lv_flex_align_t(v), C.lv_style_selector_t(sel))
}
func (o *Obj) SetStyleFlexCrossPlace(v FlexAlign, sel Selector) {
	C.lv_obj_set_style_flex_cross_place(o.c, C.lv_flex_align_t(v), C.lv_style_selector_t(sel))
}
func (o *Obj) SetStyleFlexTrackPlace(v FlexAlign, sel Selector) {
	C.lv_obj_set_style_flex_track_place(o.c, C.lv_flex_align_t(v), C.lv_style_selector_t(sel))
}
func (o *Obj) SetStyleFlexGrow(v uint8, sel Selector) {
	C.lv_obj_set_style_flex_grow(o.c, C.uint8_t(v), C.lv_style_selector_t(sel))
}

func (o *Obj) SetStyleGridColumnAlign(v GridAlign, sel Selector) {
	C.lv_obj_set_style_grid_column_align(o.c, C.lv_grid_align_t(v), C.lv_style_selector_t(sel))
}
func (o *Obj) SetStyleGridRowAlign(v GridAlign, sel Selector) {
	C.lv_obj_set_style_grid_row_align(o.c, C.lv_grid_align_t(v), C.lv_style_selector_t(sel))
}
func (o *Obj) SetStyleGridCellColumnPos(v int32, sel Selector) {
	C.lv_obj_set_style_grid_cell_column_pos(o.c, C.int32_t(v), C.lv_style_selector_t(sel))
}
func (o *Obj) SetStyleGridCellColumnSpan(v int32, sel Selector) {
	C.lv_obj_set_style_grid_cell_column_span(o.c, C.int32_t(v), C.lv_style_selector_t(sel))
}
func (o *Obj) SetStyleGridCellRowPos(v int32, sel Selector) {
	C.lv_obj_set_style_grid_cell_row_pos(o.c, C.int32_t(v), C.lv_style_selector_t(sel))
}
func (o *Obj) SetStyleGridCellRowSpan(v int32, sel Selector) {
	C.lv_obj_set_style_grid_cell_row_span(o.c, C.int32_t(v), C.lv_style_selector_t(sel))
}
func (o *Obj) SetStyleGridCellXAlign(v GridAlign, sel Selector) {
	C.lv_obj_set_style_grid_cell_x_align(o.c, C.lv_grid_align_t(v), C.lv_style_selector_t(sel))
}
func (o *Obj) SetStyleGridCellYAlign(v GridAlign, sel Selector) {
	C.lv_obj_set_style_grid_cell_y_align(o.c, C.lv_grid_align_t(v), C.lv_style_selector_t(sel))
}
