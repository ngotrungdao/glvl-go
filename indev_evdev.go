//go:build evdev

// This file needs LV_USE_EVDEV=1 in the lvgl-c build it links against —
// not the case for this repo's current prebuilt liblvgl.a (checked via
// `nm liblvgl.a | grep lv_evdev_create`: absent). It's gated behind the
// "evdev" build tag so a normal `go build` keeps working against that
// build; pass `-tags evdev` once lvgl-c has been rebuilt with evdev
// enabled (e.g. for a headless target pairing this with DRMDisplay).
package lvgl

/*
#include <stdlib.h>
#include <lvgl.h>

extern void goEvdevDiscoveryTrampoline(lv_indev_t *indev, lv_evdev_type_t type, void *user_data);
*/
import "C"
import (
	"errors"
	"runtime/cgo"
	"unsafe"
)

// EvdevType mirrors lv_evdev_type_t, describing what kind of device an
// evdev discovery callback found.
type EvdevType uint32

var (
	EvdevTypeRel = EvdevType(C.LV_EVDEV_TYPE_REL) // mice
	EvdevTypeAbs = EvdevType(C.LV_EVDEV_TYPE_ABS) // touchscreens, touchpads
	EvdevTypeKey = EvdevType(C.LV_EVDEV_TYPE_KEY) // keyboards, keypads, buttons
)

var errEvdev = errors.New("lvgl: evdev device open failed")

// EvdevCreate opens a Linux evdev input device (e.g. "/dev/input/event0")
// directly — no X11/Wayland/libinput involved — the natural input source
// to pair with DRMDisplay on a headless embedded target. indevType should
// be IndevTypePointer or IndevTypeKeypad.
func EvdevCreate(indevType IndevType, devPath string) (*Indev, error) {
	cPath := C.CString(devPath)
	defer C.free(unsafe.Pointer(cPath))
	c := C.lv_evdev_create(C.lv_indev_type_t(indevType), cPath)
	if c == nil {
		return nil, errEvdev
	}
	return wrapIndev(c), nil
}

// EvdevCreateFD is like EvdevCreate, but takes ownership of an
// already-open file descriptor instead of a path.
func EvdevCreateFD(indevType IndevType, fd int) (*Indev, error) {
	c := C.lv_evdev_create_fd(C.lv_indev_type_t(indevType), C.int(fd))
	if c == nil {
		return nil, errEvdev
	}
	return wrapIndev(c), nil
}

// Delete closes and frees an evdev input device. Only valid for an Indev
// created via EvdevCreate/EvdevCreateFD — SDL/Wayland-backed Indev values
// (e.g. from Display.Pointer/Keyboard) don't need or support this.
func (i *Indev) Delete() { C.lv_evdev_delete(i.c) }

// SetSwapAxes sets whether an evdev pointer device's X/Y coordinates
// should be swapped (defaults to false).
func (i *Indev) SetSwapAxes(swap bool) { C.lv_evdev_set_swap_axes(i.c, C.bool(swap)) }

// SetCalibration configures a coordinate transformation for an evdev
// pointer device, applied after axis swap (if any): the raw range
// [minX,maxX]x[minY,maxY] is mapped onto the display's own coordinates.
func (i *Indev) SetCalibration(minX, minY, maxX, maxY int) {
	C.lv_evdev_set_calibration(i.c, C.int(minX), C.int(minY), C.int(maxX), C.int(maxY))
}

// EvdevIsRawKey reports whether e's key event carries a raw keycode LVGL
// couldn't map to one of its own LV_KEY_* constants.
func EvdevIsRawKey(e *Event) bool { return bool(C.lv_evdev_is_raw_key(e.c)) }

// EvdevRawKey returns the raw keycode for an event where EvdevIsRawKey is
// true.
func EvdevRawKey(e *Event) uint16 { return uint16(C.lv_evdev_get_raw_key(e.c)) }

type evdevDiscoveryCallback struct {
	fn func(indev *Indev, kind EvdevType)
}

var evdevDiscoveryHandle cgo.Handle

// EvdevDiscoveryStart begins automatically creating an Indev for every
// evdev device already present in /dev/input/, and for every new one that
// appears afterwards. fn (optional, may be nil) is called once per
// discovered device. Only one discovery process can run at a time.
func EvdevDiscoveryStart(fn func(indev *Indev, kind EvdevType)) error {
	h := cgo.NewHandle(&evdevDiscoveryCallback{fn: fn})
	res := C.lv_evdev_discovery_start(C.lv_evdev_discovery_cb_t(C.goEvdevDiscoveryTrampoline), unsafeHandlePointer(h))
	if res != C.LV_RESULT_OK {
		h.Delete()
		return errors.New("lvgl: evdev discovery already running or failed to start")
	}
	evdevDiscoveryHandle = h
	return nil
}

// EvdevDiscoveryStop stops automatic evdev discovery started by
// EvdevDiscoveryStart. Safe to call from within the discovery callback.
func EvdevDiscoveryStop() error {
	if C.lv_evdev_discovery_stop() != C.LV_RESULT_OK {
		return errors.New("lvgl: evdev discovery not running")
	}
	if evdevDiscoveryHandle != 0 {
		evdevDiscoveryHandle.Delete()
		evdevDiscoveryHandle = 0
	}
	return nil
}

//export goEvdevDiscoveryTrampoline
func goEvdevDiscoveryTrampoline(indev *C.lv_indev_t, kind C.lv_evdev_type_t, userData unsafe.Pointer) {
	h := cgo.Handle(uintptr(userData))
	cb, ok := h.Value().(*evdevDiscoveryCallback)
	if !ok || cb.fn == nil {
		return
	}
	cb.fn(wrapIndev(indev), EvdevType(kind))
}
