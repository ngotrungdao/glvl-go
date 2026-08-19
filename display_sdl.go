package lvgl

/*
#include <stdlib.h>
#include <lvgl.h>

extern void goEventTrampoline(lv_event_t *e);
*/
import "C"
import (
	"runtime/cgo"
	"unsafe"
)

// Display wraps an lv_display_t created by the SDL2 driver.
//
// Closing the window (its own [x] button) has LVGL delete the underlying
// lv_display_t, which this type observes via an LV_EVENT_DELETE hook to
// flip IsOpen() to false — there is no lv_sdl_window_is_open equivalent
// in the SDL driver's API the way there was for Wayland. Once every SDL
// window is gone, SDL posts SDL_QUIT, and this build has
// LV_SDL_DIRECT_EXIT=1 (see lv_conf.h), so LVGL's own SDL event pump
// calls exit(0) directly from inside that C callback — Run's loop (and
// any other deferred Go cleanup) never gets a chance to observe IsOpen()
// turning false for that path. IsOpen()/Close() remain meaningful for a
// programmatic close or a multi-window app; see "Known limitations" in
// the README.
type Display struct {
	c    *C.lv_display_t
	open bool
}

// SDLWindowCreate opens an SDL2 window of the given resolution and title,
// and returns the LVGL display associated with it. Init must have been
// called first.
func SDLWindowCreate(width, height int32, title string) *Display {
	c := C.lv_sdl_window_create(C.int32_t(width), C.int32_t(height))
	if c == nil {
		return nil
	}

	cTitle := C.CString(title)
	defer C.free(unsafe.Pointer(cTitle))
	C.lv_sdl_window_set_title(c, cTitle)

	d := &Display{c: c, open: true}

	var h cgo.Handle
	h = cgo.NewHandle(&eventCallback{fn: func(*Event) {
		d.open = false
		h.Delete()
	}})
	C.lv_display_add_event_cb(c, C.lv_event_cb_t(C.goEventTrampoline),
		C.lv_event_code_t(EventDelete), unsafeHandlePointer(h))

	return d
}

// ScreenActive returns the display's currently active screen object.
func (d *Display) ScreenActive() *Obj {
	return wrapObj(C.lv_display_get_screen_active(d.c))
}

// Close closes the window programmatically, deleting the underlying
// lv_display_t. Safe to call more than once.
func (d *Display) Close() {
	if !d.open {
		return
	}
	C.lv_display_delete(d.c)
}

// IsOpen reports whether the window is still open. See the Display doc
// comment for the one case (the user closing the last window via [x])
// where this never gets observed turning false.
func (d *Display) IsOpen() bool {
	return d.open
}

// SetResizeable enables/disables resizing the window via the window
// manager.
func (d *Display) SetResizeable(resizeable bool) {
	C.lv_sdl_window_set_resizeable(d.c, C.bool(resizeable))
}

// SetSize resizes the window.
func (d *Display) SetSize(width, height int32) {
	C.lv_sdl_window_set_size(d.c, C.int32_t(width), C.int32_t(height))
}

// SetZoom sets the window's zoom factor (1.0 = no zoom).
func (d *Display) SetZoom(zoom float32) {
	C.lv_sdl_window_set_zoom(d.c, C.float(zoom))
}

// Zoom returns the window's current zoom factor.
func (d *Display) Zoom() float32 {
	return float32(C.lv_sdl_window_get_zoom(d.c))
}

// SetTitle changes the window's title.
func (d *Display) SetTitle(title string) {
	cTitle := C.CString(title)
	defer C.free(unsafe.Pointer(cTitle))
	C.lv_sdl_window_set_title(d.c, cTitle)
}

// Pointer creates and returns the SDL mouse input device. Unlike the
// Wayland driver's per-display getters, SDL's indev sources are created
// process-wide the first time they're requested — call this after the
// first Display exists.
func (d *Display) Pointer() *Indev { return wrapIndev(C.lv_sdl_mouse_create()) }

// Keyboard creates and returns the SDL keyboard input device. See the
// Pointer doc comment about SDL's indev sources being process-wide.
func (d *Display) Keyboard() *Indev { return wrapIndev(C.lv_sdl_keyboard_create()) }

// RotatePointCCW rotates p in place opposite to the display's configured
// software rotation and returns the result.
func (d *Display) RotatePointCCW(p Point) Point {
	cp := C.lv_point_t{x: C.int32_t(p.X), y: C.int32_t(p.Y)}
	C.lv_display_rotate_point_ccw(d.c, &cp)
	return Point{X: int32(cp.x), Y: int32(cp.y)}
}
