package lvgl

/*
#include <stdlib.h>
#include <lvgl.h>
*/
import "C"
import "unsafe"

// ImageButtonState mirrors lv_imagebutton_state_t.
type ImageButtonState uint32

var (
	ImageButtonStateReleased        = ImageButtonState(C.LV_IMAGEBUTTON_STATE_RELEASED)
	ImageButtonStatePressed         = ImageButtonState(C.LV_IMAGEBUTTON_STATE_PRESSED)
	ImageButtonStateDisabled        = ImageButtonState(C.LV_IMAGEBUTTON_STATE_DISABLED)
	ImageButtonStateCheckedReleased = ImageButtonState(C.LV_IMAGEBUTTON_STATE_CHECKED_RELEASED)
	ImageButtonStateCheckedPressed  = ImageButtonState(C.LV_IMAGEBUTTON_STATE_CHECKED_PRESSED)
	ImageButtonStateCheckedDisabled = ImageButtonState(C.LV_IMAGEBUTTON_STATE_CHECKED_DISABLED)
)

// ImageButton wraps an lv_imagebutton widget: a button rendered from up to
// three source images (left/mid/right) per state, instead of a styled
// rectangle.
//
// Like Image.SetSrcPath, source paths are stored by LVGL as pointers, not
// copied, so SetSrc keeps its C strings pinned on the C heap for the
// button's lifetime rather than freeing them after the call (a small,
// bounded leak if SetSrc is called repeatedly for the same state, traded
// for correctness).
type ImageButton struct {
	*Obj
	pins []*C.char
}

// NewImageButton creates an image button as a child of parent.
func NewImageButton(parent *Obj) *ImageButton {
	return &ImageButton{Obj: wrapObj(C.lv_imagebutton_create(parent.c))}
}

// SetSrc sets the left/mid/right source images (file paths or symbol
// strings) used while the button is in the given state. Pass "" for any
// of left/mid/right that isn't used in that state.
func (b *ImageButton) SetSrc(state ImageButtonState, left, mid, right string) {
	lp, mp, rp := b.pin(left), b.pin(mid), b.pin(right)
	C.lv_imagebutton_set_src(b.c, C.lv_imagebutton_state_t(state), lp, mp, rp)
}

func (b *ImageButton) pin(s string) unsafe.Pointer {
	if s == "" {
		return nil
	}
	c := C.CString(s)
	b.pins = append(b.pins, c)
	return unsafe.Pointer(c)
}

// SetSrcLeft sets just the left source image for the given state.
func (b *ImageButton) SetSrcLeft(state ImageButtonState, src string) {
	C.lv_imagebutton_set_src_left(b.c, C.lv_imagebutton_state_t(state), b.pin(src))
}

// SetSrcMid sets just the middle source image for the given state.
func (b *ImageButton) SetSrcMid(state ImageButtonState, src string) {
	C.lv_imagebutton_set_src_mid(b.c, C.lv_imagebutton_state_t(state), b.pin(src))
}

// SetSrcRight sets just the right source image for the given state.
func (b *ImageButton) SetSrcRight(state ImageButtonState, src string) {
	C.lv_imagebutton_set_src_right(b.c, C.lv_imagebutton_state_t(state), b.pin(src))
}

// SetState forces the button into a specific visual state (by default it
// follows the widget's own pressed/checked/disabled state automatically).
func (b *ImageButton) SetState(state ImageButtonState) {
	C.lv_imagebutton_set_state(b.c, C.lv_imagebutton_state_t(state))
}

// SrcLeft returns the left source image set for the given state.
func (b *ImageButton) SrcLeft(state ImageButtonState) string {
	return C.GoString((*C.char)(C.lv_imagebutton_get_src_left(b.c, C.lv_imagebutton_state_t(state))))
}

// SrcMid returns the middle source image set for the given state.
func (b *ImageButton) SrcMid(state ImageButtonState) string {
	return C.GoString((*C.char)(C.lv_imagebutton_get_src_middle(b.c, C.lv_imagebutton_state_t(state))))
}

// SrcRight returns the right source image set for the given state.
func (b *ImageButton) SrcRight(state ImageButtonState) string {
	return C.GoString((*C.char)(C.lv_imagebutton_get_src_right(b.c, C.lv_imagebutton_state_t(state))))
}
