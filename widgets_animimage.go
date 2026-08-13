package lvgl

/*
#include <stdlib.h>
#include <lvgl.h>
*/
import "C"
import "unsafe"

// AnimImage wraps an lv_animimg widget: a flipbook-style animation shown
// by cycling through a list of image sources.
//
// Like Image.SetSrcPath, the source array is stored by LVGL as a pointer,
// not copied, so SetSrc keeps it pinned on the C heap for the widget's
// lifetime, freeing any previous array on replacement.
type AnimImage struct {
	*Obj
	srcPtrs []*C.char
	srcArr  unsafe.Pointer
}

// NewAnimImage creates an animated image as a child of parent.
func NewAnimImage(parent *Obj) *AnimImage {
	return &AnimImage{Obj: wrapObj(C.lv_animimg_create(parent.c))}
}

// SetSrc sets the ordered list of image sources (file paths) to cycle
// through.
func (a *AnimImage) SetSrc(paths []string) {
	cStrs := make([]*C.char, len(paths))
	for i, p := range paths {
		cStrs[i] = C.CString(p)
	}

	arr := (**C.char)(C.malloc(C.size_t(len(cStrs)) * C.size_t(unsafe.Sizeof((*C.char)(nil)))))
	copy(unsafe.Slice(arr, len(cStrs)), cStrs)

	C.lv_animimg_set_src(a.c, (*unsafe.Pointer)(unsafe.Pointer(arr)), C.size_t(len(cStrs)))

	a.freeSrc()
	a.srcPtrs = cStrs
	a.srcArr = unsafe.Pointer(arr)
}

func (a *AnimImage) freeSrc() {
	for _, p := range a.srcPtrs {
		C.free(unsafe.Pointer(p))
	}
	if a.srcArr != nil {
		C.free(a.srcArr)
	}
	a.srcPtrs = nil
	a.srcArr = nil
}

// SetDuration sets how long, in milliseconds, one full cycle through the
// image list takes.
func (a *AnimImage) SetDuration(ms uint32) { C.lv_animimg_set_duration(a.c, C.uint32_t(ms)) }

// SetRepeatCount sets how many times the animation repeats, or
// AnimRepeatInfinite.
func (a *AnimImage) SetRepeatCount(count uint32) {
	C.lv_animimg_set_repeat_count(a.c, C.uint32_t(count))
}

// SetReverseDuration sets how long, in milliseconds, the reverse playback
// leg takes (0 disables reverse playback).
func (a *AnimImage) SetReverseDuration(ms uint32) {
	C.lv_animimg_set_reverse_duration(a.c, C.uint32_t(ms))
}

// SetReverseDelay sets the delay, in milliseconds, before reverse playback
// starts.
func (a *AnimImage) SetReverseDelay(ms uint32) {
	C.lv_animimg_set_reverse_delay(a.c, C.uint32_t(ms))
}

// Duration returns the current cycle duration, in milliseconds.
func (a *AnimImage) Duration() uint32 { return uint32(C.lv_animimg_get_duration(a.c)) }

// RepeatCount returns the current repeat count.
func (a *AnimImage) RepeatCount() uint32 { return uint32(C.lv_animimg_get_repeat_count(a.c)) }

// Start begins playing the animation (lv_animimg widgets auto-start on
// source assignment in some LVGL versions, but calling this is always
// safe).
func (a *AnimImage) Start() { C.lv_animimg_start(a.c) }
