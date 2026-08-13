package lvgl

/*
#include <lvgl.h>

extern void goTimerTrampoline(lv_timer_t *timer);
*/
import "C"
import "runtime/cgo"

type timerCallback struct {
	fn func(*Timer)
}

// Timer wraps an lv_timer_t: a recurring callback driven by LVGL's own
// timer loop (the same one Run/TimerHandler drives), as opposed to a
// plain Go time.Ticker which wouldn't be synchronized with LVGL's
// rendering cycle.
//
// Unlike Anim, lv_timer_t has no "deleted" callback hook, so if
// SetAutoDelete is left at its default (on) and SetRepeatCount is
// finite, LVGL will delete the underlying C timer itself once it's spent
// — with no notification back to Go to free the associated cgo.Handle.
// That's a small, permanent handle leak per such timer. For a
// long-running or frequently created timer, prefer SetAutoDelete(false)
// and call Delete() explicitly once you're done with it.
//
// Never do both: calling Delete() yourself (e.g. from within the timer's
// own callback) while SetAutoDelete is still on and SetRepeatCount is
// finite races LVGL's own auto-delete-on-exhaustion logic against your
// manual delete — a double-free that was observed to hang the process
// rather than crash cleanly, not just a theoretical risk.
type Timer struct {
	c *C.lv_timer_t
	h cgo.Handle
}

// NewTimer creates and starts a timer that calls fn every period
// milliseconds. Call Delete to stop and free it.
func NewTimer(period uint32, fn func(*Timer)) *Timer {
	h := cgo.NewHandle(&timerCallback{fn: fn})
	c := C.lv_timer_create(C.lv_timer_cb_t(C.goTimerTrampoline), C.uint32_t(period), unsafeHandlePointer(h))
	return &Timer{c: c, h: h}
}

// Delete stops and frees the timer. Safe to call from within the timer's
// own callback.
func (t *Timer) Delete() {
	if t.c == nil {
		return
	}
	C.lv_timer_delete(t.c)
	t.h.Delete()
	t.c = nil
}

// Pause stops the timer from firing without deleting it.
func (t *Timer) Pause() { C.lv_timer_pause(t.c) }

// Resume resumes a paused timer.
func (t *Timer) Resume() { C.lv_timer_resume(t.c) }

// SetPeriod changes how often the timer fires, in milliseconds.
func (t *Timer) SetPeriod(period uint32) { C.lv_timer_set_period(t.c, C.uint32_t(period)) }

// SetRepeatCount sets how many times the timer fires before
// auto-deleting, or a negative count to repeat forever.
func (t *Timer) SetRepeatCount(count int32) {
	C.lv_timer_set_repeat_count(t.c, C.int32_t(count))
}

// SetAutoDelete enables/disables automatically deleting the timer once
// its repeat count is exhausted (enabled by default).
func (t *Timer) SetAutoDelete(enable bool) { C.lv_timer_set_auto_delete(t.c, C.bool(enable)) }

// Ready makes the timer fire on the next timer-handler cycle regardless
// of its period/elapsed time.
func (t *Timer) Ready() { C.lv_timer_ready(t.c) }

// Reset restarts the timer's elapsed-time counter.
func (t *Timer) Reset() { C.lv_timer_reset(t.c) }

// Paused reports whether the timer is currently paused.
func (t *Timer) Paused() bool { return bool(C.lv_timer_get_paused(t.c)) }

//export goTimerTrampoline
func goTimerTrampoline(c *C.lv_timer_t) {
	h := cgo.Handle(uintptr(C.lv_timer_get_user_data(c)))
	cb := h.Value().(*timerCallback)
	cb.fn(&Timer{c: c, h: h})
}
