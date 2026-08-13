package lvgl

/*
#include <lvgl.h>
*/
import "C"

// IndevType mirrors lv_indev_type_t.
type IndevType uint32

var (
	IndevTypeNone    = IndevType(C.LV_INDEV_TYPE_NONE)
	IndevTypePointer = IndevType(C.LV_INDEV_TYPE_POINTER)
	IndevTypeKeypad  = IndevType(C.LV_INDEV_TYPE_KEYPAD)
	IndevTypeButton  = IndevType(C.LV_INDEV_TYPE_BUTTON)
	IndevTypeEncoder = IndevType(C.LV_INDEV_TYPE_ENCODER)
)

// Indev wraps an lv_indev_t input device (pointer, keypad/keyboard,
// button, or encoder). This package doesn't wrap creating custom input
// device drivers (LVGL's read-callback registration for that is a
// separate, lower-level API) — Indev values here come from an existing
// driver, e.g. Display.Pointer/Keyboard for the Wayland backend.
type Indev struct {
	c *C.lv_indev_t
}

func wrapIndev(c *C.lv_indev_t) *Indev {
	if c == nil {
		return nil
	}
	return &Indev{c: c}
}

// ActiveIndev returns the input device currently being processed (only
// meaningful from within an event callback triggered by input).
func ActiveIndev() *Indev { return wrapIndev(C.lv_indev_active()) }

// Type returns the device's type (pointer, keypad, button, or encoder).
func (i *Indev) Type() IndevType { return IndevType(C.lv_indev_get_type(i.c)) }

// Enable enables/disables the device.
func (i *Indev) Enable(enable bool) { C.lv_indev_enable(i.c, C.bool(enable)) }

// SetGroup attaches the device (must be a keypad or encoder) to a Group
// for keyboard/encoder-style focus navigation.
func (i *Indev) SetGroup(g *Group) { C.lv_indev_set_group(i.c, g.c) }

// Group returns the device's currently attached Group, if any.
func (i *Indev) Group() *Group {
	c := C.lv_indev_get_group(i.c)
	if c == nil {
		return nil
	}
	return &Group{c: c}
}

// SetCursor sets an object (typically an Image) to render as this
// pointer device's cursor.
func (i *Indev) SetCursor(cursor *Obj) { C.lv_indev_set_cursor(i.c, cursor.c) }

// Point returns the pointer/button device's last known coordinates.
func (i *Indev) Point() Point {
	var p C.lv_point_t
	C.lv_indev_get_point(i.c, &p)
	return Point{X: int32(p.x), Y: int32(p.y)}
}

// Key returns the last key the keypad/keyboard device produced.
func (i *Indev) Key() uint32 { return uint32(C.lv_indev_get_key(i.c)) }

// ScrollDir returns the direction the device is currently scrolling in.
func (i *Indev) ScrollDir() Dir { return Dir(C.lv_indev_get_scroll_dir(i.c)) }

// GestureDir returns the direction of the last detected gesture.
func (i *Indev) GestureDir() Dir { return Dir(C.lv_indev_get_gesture_dir(i.c)) }

// PressMoved reports whether the pointer moved while pressed (i.e. is
// dragging rather than just clicking).
func (i *Indev) PressMoved() bool { return bool(C.lv_indev_get_press_moved(i.c)) }

// Reset resets the device's internal state; if obj is non-nil, only that
// object's related state is reset. Pass nil to reset everything.
func (i *Indev) Reset(obj *Obj) {
	var c *C.lv_obj_t
	if obj != nil {
		c = obj.c
	}
	C.lv_indev_reset(i.c, c)
}

// ResetLongPress resets the device's long-press detection state.
func (i *Indev) ResetLongPress() { C.lv_indev_reset_long_press(i.c) }

// StopProcessing stops the current event chain from propagating further
// for this device.
func (i *Indev) StopProcessing() { C.lv_indev_stop_processing(i.c) }

// WaitRelease makes the device ignore input until the user releases
// whatever is currently pressed.
func (i *Indev) WaitRelease() { C.lv_indev_wait_release(i.c) }

// SetLongPressTime sets how long, in milliseconds, a press must be held
// to register as a long press.
func (i *Indev) SetLongPressTime(ms uint16) {
	C.lv_indev_set_long_press_time(i.c, C.uint16_t(ms))
}

// SetLongPressRepeatTime sets the interval, in milliseconds, between
// repeated long-press events.
func (i *Indev) SetLongPressRepeatTime(ms uint16) {
	C.lv_indev_set_long_press_repeat_time(i.c, C.uint16_t(ms))
}

// SetScrollLimit sets how many pixels of movement are needed before a
// press is treated as a scroll instead of a click.
func (i *Indev) SetScrollLimit(pixels uint8) {
	C.lv_indev_set_scroll_limit(i.c, C.uint8_t(pixels))
}
