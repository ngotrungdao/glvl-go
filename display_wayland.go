package lvgl

/*
#include <stdlib.h>
#include <lvgl.h>
*/
import "C"
import "unsafe"

// Display wraps an lv_display_t created by the Wayland driver.
type Display struct {
	c *C.lv_display_t
}

// WaylandWindowCreate opens a Wayland window of the given resolution and
// title, and returns the LVGL display associated with it. Init must have
// been called first.
func WaylandWindowCreate(width, height uint32, title string) *Display {
	cTitle := C.CString(title)
	defer C.free(unsafe.Pointer(cTitle))
	c := C.lv_wayland_window_create(C.uint32_t(width), C.uint32_t(height), cTitle, nil)
	if c == nil {
		return nil
	}
	return &Display{c: c}
}

// ScreenActive returns the display's currently active screen object.
func (d *Display) ScreenActive() *Obj {
	return wrapObj(C.lv_display_get_screen_active(d.c))
}

// Close closes the window programmatically.
func (d *Display) Close() {
	C.lv_wayland_window_close(d.c)
}

// IsOpen reports whether the window is still open.
func (d *Display) IsOpen() bool {
	return bool(C.lv_wayland_window_is_open(d.c))
}

// SetFullscreen toggles fullscreen mode for the window.
func (d *Display) SetFullscreen(fullscreen bool) {
	C.lv_wayland_window_set_fullscreen(d.c, C.bool(fullscreen))
}

// SetMaximized toggles the maximized state of the window.
func (d *Display) SetMaximized(maximized bool) {
	C.lv_wayland_window_set_maximized(d.c, C.bool(maximized))
}

// SetMinimized minimizes the window.
func (d *Display) SetMinimized() {
	C.lv_wayland_window_set_minimized(d.c)
}

// Pointer returns the display's mouse/touchpad input device.
func (d *Display) Pointer() *Indev { return wrapIndev(C.lv_wayland_get_pointer(d.c)) }

// Keyboard returns the display's keyboard input device.
func (d *Display) Keyboard() *Indev { return wrapIndev(C.lv_wayland_get_keyboard(d.c)) }

// Touchscreen returns the display's touchscreen input device.
func (d *Display) Touchscreen() *Indev { return wrapIndev(C.lv_wayland_get_touchscreen(d.c)) }
