package lvgl

/*
#include <lvgl.h>

extern void goAnimExecTrampoline(lv_anim_t *a, int32_t v);
extern void goAnimCompletedTrampoline(lv_anim_t *a);
extern void goAnimDeletedTrampoline(lv_anim_t *a);

// cgo can't take the address of a named C function as a Go value, so
// built-in easing curves are selected through this small shim instead of
// exposing lv_anim_path_* function pointers directly.
static void lvgl_go_anim_set_path(lv_anim_t *a, int selector) {
    switch (selector) {
        case 0: lv_anim_set_path_cb(a, lv_anim_path_linear); break;
        case 1: lv_anim_set_path_cb(a, lv_anim_path_ease_in); break;
        case 2: lv_anim_set_path_cb(a, lv_anim_path_ease_out); break;
        case 3: lv_anim_set_path_cb(a, lv_anim_path_ease_in_out); break;
        case 4: lv_anim_set_path_cb(a, lv_anim_path_overshoot); break;
        case 5: lv_anim_set_path_cb(a, lv_anim_path_bounce); break;
        case 6: lv_anim_set_path_cb(a, lv_anim_path_step); break;
        default: break;
    }
}
*/
import "C"
import "runtime/cgo"

// AnimRepeatInfinite passed to SetRepeatCount makes an animation repeat
// forever (LV_ANIM_REPEAT_INFINITE).
const AnimRepeatInfinite uint32 = 0xFFFFFFFF

// AnimPath selects one of LVGL's built-in animation easing curves, for
// use with Anim.SetPath.
type AnimPath int32

const (
	AnimPathLinear AnimPath = iota
	AnimPathEaseIn
	AnimPathEaseOut
	AnimPathEaseInOut
	AnimPathOvershoot
	AnimPathBounce
	AnimPathStep
)

type animCallbacks struct {
	exec      func(value int32)
	completed func()
}

// Anim is a builder for an LVGL animation: configure it with the Set*
// methods, then call Start. Once started, LVGL owns and drives its own
// internal copy of the animation state; this Anim value can be discarded.
type Anim struct {
	c C.lv_anim_t
	h cgo.Handle
}

// NewAnim creates a new animation builder.
func NewAnim() *Anim {
	a := &Anim{}
	C.lv_anim_init(&a.c)
	return a
}

// SetValues sets the start and end values passed to the exec callback.
func (a *Anim) SetValues(start, end int32) {
	C.lv_anim_set_values(&a.c, C.int32_t(start), C.int32_t(end))
}

// SetDuration sets how long the animation takes, in milliseconds.
func (a *Anim) SetDuration(ms uint32) {
	C.lv_anim_set_duration(&a.c, C.uint32_t(ms))
}

// SetDelay sets a delay before the animation starts, in milliseconds.
func (a *Anim) SetDelay(ms uint32) {
	C.lv_anim_set_delay(&a.c, C.uint32_t(ms))
}

// SetPath sets the easing curve used to interpolate between the start
// and end values (linear by default).
func (a *Anim) SetPath(path AnimPath) {
	C.lvgl_go_anim_set_path(&a.c, C.int(path))
}

// SetRepeatCount sets how many times to repeat, or AnimRepeatInfinite.
func (a *Anim) SetRepeatCount(count uint32) {
	C.lv_anim_set_repeat_count(&a.c, C.uint32_t(count))
}

// SetExecCB sets the function called on every animation step with the
// current interpolated value.
func (a *Anim) SetExecCB(fn func(value int32)) {
	a.ensureHandle().exec = fn
	C.lv_anim_set_custom_exec_cb(&a.c, C.lv_anim_custom_exec_cb_t(C.goAnimExecTrampoline))
}

// SetCompletedCB sets the function called once the animation fully
// completes (not called if it's deleted early or repeats infinitely).
func (a *Anim) SetCompletedCB(fn func()) {
	a.ensureHandle().completed = fn
	C.lv_anim_set_completed_cb(&a.c, C.lv_anim_completed_cb_t(C.goAnimCompletedTrampoline))
}

func (a *Anim) ensureHandle() *animCallbacks {
	if a.h == 0 {
		cb := &animCallbacks{}
		a.h = cgo.NewHandle(cb)
		C.lv_anim_set_user_data(&a.c, unsafeHandlePointer(a.h))
		C.lv_anim_set_deleted_cb(&a.c, C.lv_anim_deleted_cb_t(C.goAnimDeletedTrampoline))
		return cb
	}
	return a.h.Value().(*animCallbacks)
}

// Start begins running the animation.
func (a *Anim) Start() {
	C.lv_anim_start(&a.c)
}

//export goAnimExecTrampoline
func goAnimExecTrampoline(a *C.lv_anim_t, v C.int32_t) {
	h := cgo.Handle(uintptr(C.lv_anim_get_user_data(a)))
	if cb, ok := h.Value().(*animCallbacks); ok && cb.exec != nil {
		cb.exec(int32(v))
	}
}

//export goAnimCompletedTrampoline
func goAnimCompletedTrampoline(a *C.lv_anim_t) {
	h := cgo.Handle(uintptr(C.lv_anim_get_user_data(a)))
	if cb, ok := h.Value().(*animCallbacks); ok && cb.completed != nil {
		cb.completed()
	}
}

//export goAnimDeletedTrampoline
func goAnimDeletedTrampoline(a *C.lv_anim_t) {
	h := cgo.Handle(uintptr(C.lv_anim_get_user_data(a)))
	h.Delete()
}
