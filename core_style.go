package lvgl

/*
#include <stdlib.h>
#include <lvgl.h>
*/
import "C"
import (
	"runtime"
	"unsafe"
)

// Opa mirrors lv_opa_t: 0 (transparent) to 255 (fully covering).
type Opa uint8

const (
	OpaTransp Opa = 0
	OpaCover  Opa = 255
)

// BorderSide mirrors lv_border_side_t; values can be OR'd together.
type BorderSide uint8

const (
	BorderSideNone   BorderSide = 0x00
	BorderSideBottom BorderSide = 0x01
	BorderSideTop    BorderSide = 0x02
	BorderSideLeft   BorderSide = 0x04
	BorderSideRight  BorderSide = 0x08
	BorderSideFull   BorderSide = 0x0F
)

// GradDir mirrors lv_grad_dir_t.
type GradDir uint32

var (
	GradDirNone    = GradDir(C.LV_GRAD_DIR_NONE)
	GradDirVer     = GradDir(C.LV_GRAD_DIR_VER)
	GradDirHor     = GradDir(C.LV_GRAD_DIR_HOR)
	GradDirLinear  = GradDir(C.LV_GRAD_DIR_LINEAR)
	GradDirRadial  = GradDir(C.LV_GRAD_DIR_RADIAL)
	GradDirConical = GradDir(C.LV_GRAD_DIR_CONICAL)
)

// BlurQuality mirrors lv_blur_quality_t.
type BlurQuality uint32

var (
	BlurQualityAuto      = BlurQuality(C.LV_BLUR_QUALITY_AUTO)
	BlurQualitySpeed     = BlurQuality(C.LV_BLUR_QUALITY_SPEED)
	BlurQualityPrecision = BlurQuality(C.LV_BLUR_QUALITY_PRECISION)
)

// TextDecor mirrors lv_text_decor_t; values can be OR'd together.
type TextDecor uint8

const (
	TextDecorNone          TextDecor = 0x00
	TextDecorUnderline     TextDecor = 0x01
	TextDecorStrikethrough TextDecor = 0x02
)

// TextLeadingTrim mirrors lv_text_leading_trim_t.
type TextLeadingTrim uint32

var (
	TextLeadingTrimNone            = TextLeadingTrim(C.LV_TEXT_LEADING_TRIM_NONE)
	TextLeadingTrimCapitalBaseline = TextLeadingTrim(C.LV_TEXT_LEADING_TRIM_CAPITAL_BASELINE)
	TextLeadingTrimLowerBaseline   = TextLeadingTrim(C.LV_TEXT_LEADING_TRIM_LOWER_BASELINE)
	TextLeadingTrimCapital         = TextLeadingTrim(C.LV_TEXT_LEADING_TRIM_CAPITAL)
	TextLeadingTrimLower           = TextLeadingTrim(C.LV_TEXT_LEADING_TRIM_LOWER)
)

// BaseDir mirrors lv_base_dir_t.
type BaseDir uint32

var (
	BaseDirLTR     = BaseDir(C.LV_BASE_DIR_LTR)
	BaseDirRTL     = BaseDir(C.LV_BASE_DIR_RTL)
	BaseDirAuto    = BaseDir(C.LV_BASE_DIR_AUTO)
	BaseDirNeutral = BaseDir(C.LV_BASE_DIR_NEUTRAL)
	BaseDirWeak    = BaseDir(C.LV_BASE_DIR_WEAK)
)

// Style is a reusable set of appearance properties that can be attached to
// any number of objects via Obj.AddStyle.
//
// LVGL stores a raw, non-refcounted pointer to a Style's C memory for as
// long as it's attached to any object, so it is allocated on the C heap
// rather than as ordinary Go memory. Call Delete once a Style is detached
// from every object that used it.
//
// This file hand-wraps 123 of LVGL's 129 lv_style_set_* property setters
// (core/lv_style_gen.h). Not covered: 5 that take pointer types needing
// their own dedicated supporting API this pass didn't build — SetAnim
// (lv_anim_t*), SetTransition (lv_style_transition_dsc_t*), SetBgGrad
// (lv_grad_dsc_t*), SetColorFilterDsc (lv_color_filter_dsc_t*),
// SetImageColorkey (lv_image_colorkey_t*) — plus SetLayout (a low-level
// raw layout-ID selector normally set automatically by SetFlexFlow/
// SetGridDscArray, not something user code sets directly). To add more,
// follow the same pattern: each existing setter is a direct, mechanical
// translation of
// the matching lv_style_set_* signature.
//
// Like Image.SetSrcPath, image-source and bitmap-mask paths are stored by
// LVGL as pointers, not copied, so they're pinned on the C heap for the
// style's lifetime (freed by Delete). Grid track descriptors keep a Go
// reference to the *GridDsc for the same reason, but ownership/Delete of
// the GridDsc itself stays with the caller, same as Obj.SetGridDscArray.
type Style struct {
	c *C.lv_style_t

	bgImageSrc, arcImageSrc, bitmapMaskSrc *C.char
	gridCols, gridRows                     *GridDsc
}

// NewStyle allocates and initializes a new, empty style.
func NewStyle() *Style {
	c := (*C.lv_style_t)(C.malloc(C.sizeof_lv_style_t))
	C.lv_style_init(c)
	s := &Style{c: c}
	runtime.SetFinalizer(s, (*Style).leakWarn)
	return s
}

func (s *Style) leakWarn() {
	println("lvgl: Style finalized without an explicit Delete() call; its C memory was leaked")
}

// Delete resets and frees the style's C memory, including any pinned
// image-source/bitmap-mask strings. Only call this once the style is no
// longer attached to any object (see Obj.RemoveStyle).
func (s *Style) Delete() {
	if s.c == nil {
		return
	}
	C.lv_style_reset(s.c)
	C.free(unsafe.Pointer(s.c))
	s.c = nil
	for _, p := range []*C.char{s.bgImageSrc, s.arcImageSrc, s.bitmapMaskSrc} {
		if p != nil {
			C.free(unsafe.Pointer(p))
		}
	}
	s.bgImageSrc, s.arcImageSrc, s.bitmapMaskSrc = nil, nil, nil
	runtime.SetFinalizer(s, nil)
}

func (s *Style) SetWidth(v int32)     { C.lv_style_set_width(s.c, C.int32_t(v)) }
func (s *Style) SetMinWidth(v int32)  { C.lv_style_set_min_width(s.c, C.int32_t(v)) }
func (s *Style) SetMaxWidth(v int32)  { C.lv_style_set_max_width(s.c, C.int32_t(v)) }
func (s *Style) SetHeight(v int32)    { C.lv_style_set_height(s.c, C.int32_t(v)) }
func (s *Style) SetMinHeight(v int32) { C.lv_style_set_min_height(s.c, C.int32_t(v)) }
func (s *Style) SetMaxHeight(v int32) { C.lv_style_set_max_height(s.c, C.int32_t(v)) }
func (s *Style) SetLength(v int32)    { C.lv_style_set_length(s.c, C.int32_t(v)) }
func (s *Style) SetX(v int32)         { C.lv_style_set_x(s.c, C.int32_t(v)) }
func (s *Style) SetY(v int32)         { C.lv_style_set_y(s.c, C.int32_t(v)) }
func (s *Style) SetAlign(v Align)     { C.lv_style_set_align(s.c, C.lv_align_t(v)) }

func (s *Style) SetTransformWidth(v int32)    { C.lv_style_set_transform_width(s.c, C.int32_t(v)) }
func (s *Style) SetTransformHeight(v int32)   { C.lv_style_set_transform_height(s.c, C.int32_t(v)) }
func (s *Style) SetTranslateX(v int32)        { C.lv_style_set_translate_x(s.c, C.int32_t(v)) }
func (s *Style) SetTranslateY(v int32)        { C.lv_style_set_translate_y(s.c, C.int32_t(v)) }
func (s *Style) SetTranslateRadial(v int32)   { C.lv_style_set_translate_radial(s.c, C.int32_t(v)) }
func (s *Style) SetTransformScaleX(v int32)   { C.lv_style_set_transform_scale_x(s.c, C.int32_t(v)) }
func (s *Style) SetTransformScaleY(v int32)   { C.lv_style_set_transform_scale_y(s.c, C.int32_t(v)) }
func (s *Style) SetTransformRotation(v int32) { C.lv_style_set_transform_rotation(s.c, C.int32_t(v)) }
func (s *Style) SetTransformPivotX(v int32)   { C.lv_style_set_transform_pivot_x(s.c, C.int32_t(v)) }
func (s *Style) SetTransformPivotY(v int32)   { C.lv_style_set_transform_pivot_y(s.c, C.int32_t(v)) }
func (s *Style) SetTransformSkewX(v int32)    { C.lv_style_set_transform_skew_x(s.c, C.int32_t(v)) }
func (s *Style) SetTransformSkewY(v int32)    { C.lv_style_set_transform_skew_y(s.c, C.int32_t(v)) }

func (s *Style) SetPadTop(v int32)    { C.lv_style_set_pad_top(s.c, C.int32_t(v)) }
func (s *Style) SetPadBottom(v int32) { C.lv_style_set_pad_bottom(s.c, C.int32_t(v)) }
func (s *Style) SetPadLeft(v int32)   { C.lv_style_set_pad_left(s.c, C.int32_t(v)) }
func (s *Style) SetPadRight(v int32)  { C.lv_style_set_pad_right(s.c, C.int32_t(v)) }
func (s *Style) SetPadRow(v int32)    { C.lv_style_set_pad_row(s.c, C.int32_t(v)) }
func (s *Style) SetPadColumn(v int32) { C.lv_style_set_pad_column(s.c, C.int32_t(v)) }
func (s *Style) SetPadRadial(v int32) { C.lv_style_set_pad_radial(s.c, C.int32_t(v)) }

func (s *Style) SetMarginTop(v int32)    { C.lv_style_set_margin_top(s.c, C.int32_t(v)) }
func (s *Style) SetMarginBottom(v int32) { C.lv_style_set_margin_bottom(s.c, C.int32_t(v)) }
func (s *Style) SetMarginLeft(v int32)   { C.lv_style_set_margin_left(s.c, C.int32_t(v)) }
func (s *Style) SetMarginRight(v int32)  { C.lv_style_set_margin_right(s.c, C.int32_t(v)) }

// SetPadAll sets all four padding sides to the same value.
func (s *Style) SetPadAll(v int32) {
	s.SetPadTop(v)
	s.SetPadBottom(v)
	s.SetPadLeft(v)
	s.SetPadRight(v)
}

// SetMarginAll sets all four margin sides to the same value.
func (s *Style) SetMarginAll(v int32) {
	s.SetMarginTop(v)
	s.SetMarginBottom(v)
	s.SetMarginLeft(v)
	s.SetMarginRight(v)
}

func (s *Style) SetBgColor(v Color)        { C.lv_style_set_bg_color(s.c, v.toC()) }
func (s *Style) SetBgOpa(v Opa)            { C.lv_style_set_bg_opa(s.c, C.lv_opa_t(v)) }
func (s *Style) SetBgGradColor(v Color)    { C.lv_style_set_bg_grad_color(s.c, v.toC()) }
func (s *Style) SetBgGradDir(v GradDir)    { C.lv_style_set_bg_grad_dir(s.c, C.lv_grad_dir_t(v)) }
func (s *Style) SetBgMainStop(v int32)     { C.lv_style_set_bg_main_stop(s.c, C.int32_t(v)) }
func (s *Style) SetBgGradStop(v int32)     { C.lv_style_set_bg_grad_stop(s.c, C.int32_t(v)) }
func (s *Style) SetBgMainOpa(v Opa)        { C.lv_style_set_bg_main_opa(s.c, C.lv_opa_t(v)) }
func (s *Style) SetBgGradOpa(v Opa)        { C.lv_style_set_bg_grad_opa(s.c, C.lv_opa_t(v)) }
func (s *Style) SetBgImageOpa(v Opa)       { C.lv_style_set_bg_image_opa(s.c, C.lv_opa_t(v)) }
func (s *Style) SetBgImageRecolor(v Color) { C.lv_style_set_bg_image_recolor(s.c, v.toC()) }
func (s *Style) SetBgImageRecolorOpa(v Opa) {
	C.lv_style_set_bg_image_recolor_opa(s.c, C.lv_opa_t(v))
}
func (s *Style) SetBgImageTiled(v bool) { C.lv_style_set_bg_image_tiled(s.c, C.bool(v)) }

// SetBgImageSrc sets a background image (file path or symbol string). Like
// Image.SetSrcPath, the string is pinned on the C heap for the style's
// lifetime rather than copied by LVGL.
func (s *Style) SetBgImageSrc(path string) {
	c := C.CString(path)
	C.lv_style_set_bg_image_src(s.c, unsafe.Pointer(c))
	if s.bgImageSrc != nil {
		C.free(unsafe.Pointer(s.bgImageSrc))
	}
	s.bgImageSrc = c
}

func (s *Style) SetBorderColor(v Color) { C.lv_style_set_border_color(s.c, v.toC()) }
func (s *Style) SetBorderOpa(v Opa)     { C.lv_style_set_border_opa(s.c, C.lv_opa_t(v)) }
func (s *Style) SetBorderWidth(v int32) { C.lv_style_set_border_width(s.c, C.int32_t(v)) }
func (s *Style) SetBorderSide(v BorderSide) {
	C.lv_style_set_border_side(s.c, C.lv_border_side_t(v))
}
func (s *Style) SetBorderPost(v bool) { C.lv_style_set_border_post(s.c, C.bool(v)) }

func (s *Style) SetOutlineWidth(v int32) { C.lv_style_set_outline_width(s.c, C.int32_t(v)) }
func (s *Style) SetOutlineColor(v Color) { C.lv_style_set_outline_color(s.c, v.toC()) }
func (s *Style) SetOutlineOpa(v Opa)     { C.lv_style_set_outline_opa(s.c, C.lv_opa_t(v)) }
func (s *Style) SetOutlinePad(v int32)   { C.lv_style_set_outline_pad(s.c, C.int32_t(v)) }

func (s *Style) SetShadowWidth(v int32)  { C.lv_style_set_shadow_width(s.c, C.int32_t(v)) }
func (s *Style) SetShadowSpread(v int32) { C.lv_style_set_shadow_spread(s.c, C.int32_t(v)) }
func (s *Style) SetShadowOfsX(v int32)   { C.lv_style_set_shadow_offset_x(s.c, C.int32_t(v)) }
func (s *Style) SetShadowOfsY(v int32)   { C.lv_style_set_shadow_offset_y(s.c, C.int32_t(v)) }
func (s *Style) SetShadowColor(v Color)  { C.lv_style_set_shadow_color(s.c, v.toC()) }
func (s *Style) SetShadowOpa(v Opa)      { C.lv_style_set_shadow_opa(s.c, C.lv_opa_t(v)) }

func (s *Style) SetDropShadowRadius(v int32) { C.lv_style_set_drop_shadow_radius(s.c, C.int32_t(v)) }
func (s *Style) SetDropShadowOfsX(v int32)   { C.lv_style_set_drop_shadow_offset_x(s.c, C.int32_t(v)) }
func (s *Style) SetDropShadowOfsY(v int32)   { C.lv_style_set_drop_shadow_offset_y(s.c, C.int32_t(v)) }
func (s *Style) SetDropShadowColor(v Color)  { C.lv_style_set_drop_shadow_color(s.c, v.toC()) }
func (s *Style) SetDropShadowOpa(v Opa)      { C.lv_style_set_drop_shadow_opa(s.c, C.lv_opa_t(v)) }
func (s *Style) SetDropShadowQuality(v BlurQuality) {
	C.lv_style_set_drop_shadow_quality(s.c, C.lv_blur_quality_t(v))
}

func (s *Style) SetBlurRadius(v int32)  { C.lv_style_set_blur_radius(s.c, C.int32_t(v)) }
func (s *Style) SetBlurBackdrop(v bool) { C.lv_style_set_blur_backdrop(s.c, C.bool(v)) }
func (s *Style) SetBlurQuality(v BlurQuality) {
	C.lv_style_set_blur_quality(s.c, C.lv_blur_quality_t(v))
}

func (s *Style) SetTextColor(v Color)       { C.lv_style_set_text_color(s.c, v.toC()) }
func (s *Style) SetTextOpa(v Opa)           { C.lv_style_set_text_opa(s.c, C.lv_opa_t(v)) }
func (s *Style) SetTextLetterSpace(v int32) { C.lv_style_set_text_letter_space(s.c, C.int32_t(v)) }
func (s *Style) SetTextLineSpace(v int32)   { C.lv_style_set_text_line_space(s.c, C.int32_t(v)) }
func (s *Style) SetTextAlign(v TextAlign)   { C.lv_style_set_text_align(s.c, C.lv_text_align_t(v)) }
func (s *Style) SetTextDecor(v TextDecor)   { C.lv_style_set_text_decor(s.c, C.lv_text_decor_t(v)) }
func (s *Style) SetTextLeadingTrim(v TextLeadingTrim) {
	C.lv_style_set_text_leading_trim(s.c, C.lv_text_leading_trim_t(v))
}
func (s *Style) SetTextOutlineStrokeColor(v Color) {
	C.lv_style_set_text_outline_stroke_color(s.c, v.toC())
}
func (s *Style) SetTextOutlineStrokeOpa(v Opa) {
	C.lv_style_set_text_outline_stroke_opa(s.c, C.lv_opa_t(v))
}
func (s *Style) SetTextOutlineStrokeWidth(v int32) {
	C.lv_style_set_text_outline_stroke_width(s.c, C.int32_t(v))
}

// SetTextFont sets both the font family and size (see Font's doc comment
// — LVGL has no separate size property, it's baked into the Font).
func (s *Style) SetTextFont(f *Font) { C.lv_style_set_text_font(s.c, f.c) }

func (s *Style) SetLineColor(v Color)     { C.lv_style_set_line_color(s.c, v.toC()) }
func (s *Style) SetLineWidth(v int32)     { C.lv_style_set_line_width(s.c, C.int32_t(v)) }
func (s *Style) SetLineOpa(v Opa)         { C.lv_style_set_line_opa(s.c, C.lv_opa_t(v)) }
func (s *Style) SetLineDashWidth(v int32) { C.lv_style_set_line_dash_width(s.c, C.int32_t(v)) }
func (s *Style) SetLineDashGap(v int32)   { C.lv_style_set_line_dash_gap(s.c, C.int32_t(v)) }
func (s *Style) SetLineRounded(v bool)    { C.lv_style_set_line_rounded(s.c, C.bool(v)) }

func (s *Style) SetArcColor(v Color)  { C.lv_style_set_arc_color(s.c, v.toC()) }
func (s *Style) SetArcOpa(v Opa)      { C.lv_style_set_arc_opa(s.c, C.lv_opa_t(v)) }
func (s *Style) SetArcWidth(v int32)  { C.lv_style_set_arc_width(s.c, C.int32_t(v)) }
func (s *Style) SetArcRounded(v bool) { C.lv_style_set_arc_rounded(s.c, C.bool(v)) }

// SetArcImageSrc sets an image used to draw the arc instead of a solid
// color. Like Image.SetSrcPath, the string is pinned on the C heap for
// the style's lifetime rather than copied by LVGL.
func (s *Style) SetArcImageSrc(path string) {
	c := C.CString(path)
	C.lv_style_set_arc_image_src(s.c, unsafe.Pointer(c))
	if s.arcImageSrc != nil {
		C.free(unsafe.Pointer(s.arcImageSrc))
	}
	s.arcImageSrc = c
}

func (s *Style) SetImageRecolor(v Color) { C.lv_style_set_image_recolor(s.c, v.toC()) }
func (s *Style) SetImageRecolorOpa(v Opa) {
	C.lv_style_set_image_recolor_opa(s.c, C.lv_opa_t(v))
}
func (s *Style) SetImageOpa(v Opa) { C.lv_style_set_image_opa(s.c, C.lv_opa_t(v)) }

func (s *Style) SetRadius(v int32)        { C.lv_style_set_radius(s.c, C.int32_t(v)) }
func (s *Style) SetRadialOffset(v int32)  { C.lv_style_set_radial_offset(s.c, C.int32_t(v)) }
func (s *Style) SetClipCorner(v bool)     { C.lv_style_set_clip_corner(s.c, C.bool(v)) }
func (s *Style) SetOpa(v Opa)             { C.lv_style_set_opa(s.c, C.lv_opa_t(v)) }
func (s *Style) SetOpaLayered(v Opa)      { C.lv_style_set_opa_layered(s.c, C.lv_opa_t(v)) }
func (s *Style) SetColorFilterOpa(v Opa)  { C.lv_style_set_color_filter_opa(s.c, C.lv_opa_t(v)) }
func (s *Style) SetRecolor(v Color)       { C.lv_style_set_recolor(s.c, v.toC()) }
func (s *Style) SetRecolorOpa(v Opa)      { C.lv_style_set_recolor_opa(s.c, C.lv_opa_t(v)) }
func (s *Style) SetAnimDuration(v uint32) { C.lv_style_set_anim_duration(s.c, C.uint32_t(v)) }
func (s *Style) SetBlendMode(v BlendMode) { C.lv_style_set_blend_mode(s.c, C.lv_blend_mode_t(v)) }
func (s *Style) SetBaseDir(v BaseDir)     { C.lv_style_set_base_dir(s.c, C.lv_base_dir_t(v)) }
func (s *Style) SetRotarySensitivity(v uint32) {
	C.lv_style_set_rotary_sensitivity(s.c, C.uint32_t(v))
}

// SetBitmapMaskSrc sets an image used as an alpha mask for the object.
// Like Image.SetSrcPath, the string is pinned on the C heap for the
// style's lifetime rather than copied by LVGL.
func (s *Style) SetBitmapMaskSrc(path string) {
	c := C.CString(path)
	C.lv_style_set_bitmap_mask_src(s.c, unsafe.Pointer(c))
	if s.bitmapMaskSrc != nil {
		C.free(unsafe.Pointer(s.bitmapMaskSrc))
	}
	s.bitmapMaskSrc = c
}

func (s *Style) SetFlexFlow(v FlexFlow) { C.lv_style_set_flex_flow(s.c, C.lv_flex_flow_t(v)) }
func (s *Style) SetFlexMainPlace(v FlexAlign) {
	C.lv_style_set_flex_main_place(s.c, C.lv_flex_align_t(v))
}
func (s *Style) SetFlexCrossPlace(v FlexAlign) {
	C.lv_style_set_flex_cross_place(s.c, C.lv_flex_align_t(v))
}
func (s *Style) SetFlexTrackPlace(v FlexAlign) {
	C.lv_style_set_flex_track_place(s.c, C.lv_flex_align_t(v))
}
func (s *Style) SetFlexGrow(v uint8) { C.lv_style_set_flex_grow(s.c, C.uint8_t(v)) }

// SetGridDscArray sets the grid's column/row track descriptors, same as
// Obj.SetGridDscArray but scoped to this style. The Style keeps a Go
// reference to cols/rows so they aren't garbage-collected while still
// needed, but (like Obj.SetGridDscArray) does not take ownership: the
// caller is still responsible for calling Delete on them once no longer
// used by any style/object.
func (s *Style) SetGridDscArray(cols, rows *GridDsc) {
	C.lv_style_set_grid_column_dsc_array(s.c, cols.c)
	C.lv_style_set_grid_row_dsc_array(s.c, rows.c)
	s.gridCols, s.gridRows = cols, rows
}

func (s *Style) SetGridColumnAlign(v GridAlign) {
	C.lv_style_set_grid_column_align(s.c, C.lv_grid_align_t(v))
}
func (s *Style) SetGridRowAlign(v GridAlign) {
	C.lv_style_set_grid_row_align(s.c, C.lv_grid_align_t(v))
}
func (s *Style) SetGridCellColumnPos(v int32) { C.lv_style_set_grid_cell_column_pos(s.c, C.int32_t(v)) }
func (s *Style) SetGridCellColumnSpan(v int32) {
	C.lv_style_set_grid_cell_column_span(s.c, C.int32_t(v))
}
func (s *Style) SetGridCellRowPos(v int32)  { C.lv_style_set_grid_cell_row_pos(s.c, C.int32_t(v)) }
func (s *Style) SetGridCellRowSpan(v int32) { C.lv_style_set_grid_cell_row_span(s.c, C.int32_t(v)) }
func (s *Style) SetGridCellXAlign(v GridAlign) {
	C.lv_style_set_grid_cell_x_align(s.c, C.lv_grid_align_t(v))
}
func (s *Style) SetGridCellYAlign(v GridAlign) {
	C.lv_style_set_grid_cell_y_align(s.c, C.lv_grid_align_t(v))
}
