package lvgl

/*
#include <stdlib.h>
#include <lvgl.h>
*/
import "C"
import (
	"errors"
	"unsafe"
)

var errDRM = errors.New("lvgl: DRM display operation failed")

// DRMDisplay wraps an lv_display_t created by LVGL's native Linux DRM/KMS
// driver (LV_USE_LINUX_DRM). Unlike Display (the SDL2 driver), it talks
// directly to /dev/dri/cardN via DRM/KMS + GBM — no X11, Wayland, or SDL2
// involved — making it the right choice for a headless embedded target
// (e.g. a Raspberry Pi with no desktop environment) rather than
// SDLWindowCreate.
//
// Custom display-mode selection (lv_linux_drm_set_mode_cb) isn't wrapped:
// its callback signature has no user_data slot to smuggle a Go closure
// through, and the default behavior — select the connector's native/
// preferred mode — is exactly what a fixed embedded panel wants anyway.
type DRMDisplay struct {
	c *C.lv_display_t
}

// DRMDisplayCreate creates a new DRM display object. Call SetFile
// afterwards to attach it to a physical DRM device/connector before using
// it — mirroring lv_linux_drm_create()'s own two-step C API.
func DRMDisplayCreate() *DRMDisplay {
	c := C.lv_linux_drm_create()
	if c == nil {
		return nil
	}
	return &DRMDisplay{c: c}
}

// SetFile attaches the display to a DRM device file (e.g. "/dev/dri/card0")
// and connector. Pass connectorID -1 to auto-select the first available
// connector.
func (d *DRMDisplay) SetFile(file string, connectorID int64) error {
	cFile := C.CString(file)
	defer C.free(unsafe.Pointer(cFile))
	if C.lv_linux_drm_set_file(d.c, cFile, C.int64_t(connectorID)) != C.LV_RESULT_OK {
		return errDRM
	}
	return nil
}

// ScreenActive returns the display's currently active screen object.
func (d *DRMDisplay) ScreenActive() *Obj {
	return wrapObj(C.lv_display_get_screen_active(d.c))
}

// FindDRMDevicePath scans the system for a suitable DRM device and
// returns its path (e.g. "/dev/dri/card0"), for use with SetFile. Returns
// an error if no suitable device is found.
func FindDRMDevicePath() (string, error) {
	cPath := C.lv_linux_drm_find_device_path()
	if cPath == nil {
		return "", errDRM
	}
	defer C.lv_free(unsafe.Pointer(cPath))
	return C.GoString(cPath), nil
}
