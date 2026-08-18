package lvgl

/*
#include <lvgl.h>
*/
import "C"

// TickInc advances LVGL's internal tick counter by tickPeriodMs
// milliseconds, for backends that need the host to manually drive the
// tick source (LVGL doesn't auto-manage one universally).
//
// The Wayland driver used by this package registers a real tick callback
// of its own (verified: TickGet() advances with real wall-clock time
// even if TickInc is never called), so once a Display exists via
// WaylandWindowCreate, TickInc calls are effectively no-ops — LVGL is
// already reading a real clock. Run (app.go) still calls it every
// iteration for backend-agnosticism, but its correctness in practice
// comes from pacing with time.Sleep between TimerHandler calls, not from
// TickInc. Don't rely on TickInc to control time once a Wayland window
// is open — poll/sleep in real time instead (see Run for the pattern).
func TickInc(tickPeriodMs uint32) {
	C.lv_tick_inc(C.uint32_t(tickPeriodMs))
}

// TickGet returns LVGL's current tick count in milliseconds. With the
// Wayland backend (see TickInc), this reflects real wall-clock time.
func TickGet() uint32 {
	return uint32(C.lv_tick_get())
}

// TimerHandler runs LVGL's timers (animations, redraws, input processing,
// ...) once and returns the number of milliseconds until it should be
// called again for the highest-priority pending timer.
func TimerHandler() uint32 {
	return uint32(C.lv_timer_handler())
}

// TimerTimeToNext returns the number of milliseconds until the
// soonest-due LVGL timer needs to run, without actually running any
// timers (unlike TimerHandler's return value, which is a side effect of
// running them).
func TimerTimeToNext() uint32 {
	return uint32(C.lv_timer_get_time_to_next())
}
