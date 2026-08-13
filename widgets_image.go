package lvgl

/*
#include <stdlib.h>
#include <lvgl.h>
*/
import "C"
import "unsafe"

// BlendMode mirrors lv_blend_mode_t.
type BlendMode uint32

var (
	BlendModeNormal      = BlendMode(C.LV_BLEND_MODE_NORMAL)
	BlendModeAdditive    = BlendMode(C.LV_BLEND_MODE_ADDITIVE)
	BlendModeSubtractive = BlendMode(C.LV_BLEND_MODE_SUBTRACTIVE)
	BlendModeMultiply    = BlendMode(C.LV_BLEND_MODE_MULTIPLY)
	BlendModeDifference  = BlendMode(C.LV_BLEND_MODE_DIFFERENCE)
)

// ImageAlign mirrors lv_image_align_t: how the image is positioned/scaled
// within the widget's area.
type ImageAlign uint32

var (
	ImageAlignDefault          = ImageAlign(C.LV_IMAGE_ALIGN_DEFAULT)
	ImageAlignTopLeft          = ImageAlign(C.LV_IMAGE_ALIGN_TOP_LEFT)
	ImageAlignTopMid           = ImageAlign(C.LV_IMAGE_ALIGN_TOP_MID)
	ImageAlignTopRight         = ImageAlign(C.LV_IMAGE_ALIGN_TOP_RIGHT)
	ImageAlignBottomLeft       = ImageAlign(C.LV_IMAGE_ALIGN_BOTTOM_LEFT)
	ImageAlignBottomMid        = ImageAlign(C.LV_IMAGE_ALIGN_BOTTOM_MID)
	ImageAlignBottomRight      = ImageAlign(C.LV_IMAGE_ALIGN_BOTTOM_RIGHT)
	ImageAlignLeftMid          = ImageAlign(C.LV_IMAGE_ALIGN_LEFT_MID)
	ImageAlignRightMid         = ImageAlign(C.LV_IMAGE_ALIGN_RIGHT_MID)
	ImageAlignCenter           = ImageAlign(C.LV_IMAGE_ALIGN_CENTER)
	ImageAlignStretch          = ImageAlign(C.LV_IMAGE_ALIGN_STRETCH)
	ImageAlignTile             = ImageAlign(C.LV_IMAGE_ALIGN_TILE)
	ImageAlignContain          = ImageAlign(C.LV_IMAGE_ALIGN_CONTAIN)
	ImageAlignContainDownscale = ImageAlign(C.LV_IMAGE_ALIGN_CONTAIN_DOWNSCALE)
)

// Image wraps an lv_image widget.
//
// Unlike widget text setters (e.g. Label.SetText), lv_image_set_src does
// not copy a path/symbol string it's given: it stores the pointer as-is.
// SetSrcPath therefore keeps the C string alive on the C heap for as long
// as it's the current source, freeing the previous one (if any) when
// replaced.
type Image struct {
	*Obj
	srcPath *C.char
}

// NewImage creates an image as a child of parent.
func NewImage(parent *Obj) *Image {
	return &Image{Obj: wrapObj(C.lv_image_create(parent.c))}
}

// SetSrcPath sets the image source to a file path or a built-in symbol
// string (e.g. an LV_SYMBOL_* constant). Image sources backed by a
// compiled-in lv_image_dsc_t aren't supported by this wrapper yet.
func (im *Image) SetSrcPath(path string) {
	cPath := C.CString(path)
	C.lv_image_set_src(im.c, unsafe.Pointer(cPath))
	if im.srcPath != nil {
		C.free(unsafe.Pointer(im.srcPath))
	}
	im.srcPath = cPath
}

// SrcPath returns the current source string (path or symbol).
func (im *Image) SrcPath() string {
	return C.GoString((*C.char)(C.lv_image_get_src(im.c)))
}

// SetOffset sets the pixel offset applied to the image content.
func (im *Image) SetOffset(x, y int32) {
	C.lv_image_set_offset_x(im.c, C.int32_t(x))
	C.lv_image_set_offset_y(im.c, C.int32_t(y))
}

// Offset returns the current pixel offset.
func (im *Image) Offset() (x, y int32) {
	return int32(C.lv_image_get_offset_x(im.c)), int32(C.lv_image_get_offset_y(im.c))
}

// SetRotation rotates the image by the given tenths of a degree.
func (im *Image) SetRotation(angle int32) { C.lv_image_set_rotation(im.c, C.int32_t(angle)) }

// Rotation returns the current rotation, in tenths of a degree.
func (im *Image) Rotation() int32 { return int32(C.lv_image_get_rotation(im.c)) }

// SetPivot sets the point rotation/scaling is applied around, relative
// to the image's top-left corner.
func (im *Image) SetPivot(x, y int32) { C.lv_image_set_pivot(im.c, C.int32_t(x), C.int32_t(y)) }

// Pivot returns the current rotation/scaling pivot point.
func (im *Image) Pivot() Point {
	var p C.lv_point_t
	C.lv_image_get_pivot(im.c, &p)
	return Point{X: int32(p.x), Y: int32(p.y)}
}

// SetScale sets the image's zoom factor; 256 means no scaling.
func (im *Image) SetScale(zoom uint32) { C.lv_image_set_scale(im.c, C.uint32_t(zoom)) }

// Scale returns the current zoom factor.
func (im *Image) Scale() int32 { return int32(C.lv_image_get_scale(im.c)) }

// SetScaleX sets the horizontal-only zoom factor.
func (im *Image) SetScaleX(zoom uint32) { C.lv_image_set_scale_x(im.c, C.uint32_t(zoom)) }

// SetScaleY sets the vertical-only zoom factor.
func (im *Image) SetScaleY(zoom uint32) { C.lv_image_set_scale_y(im.c, C.uint32_t(zoom)) }

// ScaleX returns the current horizontal zoom factor.
func (im *Image) ScaleX() int32 { return int32(C.lv_image_get_scale_x(im.c)) }

// ScaleY returns the current vertical zoom factor.
func (im *Image) ScaleY() int32 { return int32(C.lv_image_get_scale_y(im.c)) }

// SetBlendMode sets how the image's pixels blend with what's beneath them.
func (im *Image) SetBlendMode(mode BlendMode) {
	C.lv_image_set_blend_mode(im.c, C.lv_blend_mode_t(mode))
}

// BlendMode returns the current blend mode.
func (im *Image) BlendMode() BlendMode { return BlendMode(C.lv_image_get_blend_mode(im.c)) }

// SetAntialias enables/disables antialiasing when the image is
// rotated/scaled.
func (im *Image) SetAntialias(enable bool) { C.lv_image_set_antialias(im.c, C.bool(enable)) }

// Antialias reports whether antialiasing is enabled.
func (im *Image) Antialias() bool { return bool(C.lv_image_get_antialias(im.c)) }

// SetInnerAlign sets how the image content is positioned/scaled within
// the widget's own area.
func (im *Image) SetInnerAlign(align ImageAlign) {
	C.lv_image_set_inner_align(im.c, C.lv_image_align_t(align))
}

// InnerAlign returns the current inner alignment mode.
func (im *Image) InnerAlign() ImageAlign { return ImageAlign(C.lv_image_get_inner_align(im.c)) }

// SrcWidth returns the source image's native (unscaled) width.
func (im *Image) SrcWidth() int32 { return int32(C.lv_image_get_src_width(im.c)) }

// SrcHeight returns the source image's native (unscaled) height.
func (im *Image) SrcHeight() int32 { return int32(C.lv_image_get_src_height(im.c)) }

// TransformedWidth returns the image's width after rotation/scaling.
func (im *Image) TransformedWidth() int32 { return int32(C.lv_image_get_transformed_width(im.c)) }

// TransformedHeight returns the image's height after rotation/scaling.
func (im *Image) TransformedHeight() int32 {
	return int32(C.lv_image_get_transformed_height(im.c))
}
