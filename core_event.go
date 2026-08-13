package lvgl

/*
#include <lvgl.h>

extern void goEventTrampoline(lv_event_t *e);
*/
import "C"
import (
	"runtime/cgo"
	"unsafe"
)

// EventCode mirrors lv_event_code_t.
type EventCode uint32

var (
	EventAll            = EventCode(C.LV_EVENT_ALL)
	EventPressed        = EventCode(C.LV_EVENT_PRESSED)
	EventPressing       = EventCode(C.LV_EVENT_PRESSING)
	EventPressLost      = EventCode(C.LV_EVENT_PRESS_LOST)
	EventShortClicked   = EventCode(C.LV_EVENT_SHORT_CLICKED)
	EventSingleClicked  = EventCode(C.LV_EVENT_SINGLE_CLICKED)
	EventDoubleClicked  = EventCode(C.LV_EVENT_DOUBLE_CLICKED)
	EventLongPressed    = EventCode(C.LV_EVENT_LONG_PRESSED)
	EventLongPressedRep = EventCode(C.LV_EVENT_LONG_PRESSED_REPEAT)
	EventClicked        = EventCode(C.LV_EVENT_CLICKED)
	EventReleased       = EventCode(C.LV_EVENT_RELEASED)
	EventScrollBegin    = EventCode(C.LV_EVENT_SCROLL_BEGIN)
	EventScrollEnd      = EventCode(C.LV_EVENT_SCROLL_END)
	EventScroll         = EventCode(C.LV_EVENT_SCROLL)
	EventGesture        = EventCode(C.LV_EVENT_GESTURE)
	EventKey            = EventCode(C.LV_EVENT_KEY)
	EventFocused        = EventCode(C.LV_EVENT_FOCUSED)
	EventDefocused      = EventCode(C.LV_EVENT_DEFOCUSED)
	EventValueChanged   = EventCode(C.LV_EVENT_VALUE_CHANGED)
	EventInsert         = EventCode(C.LV_EVENT_INSERT)
	EventRefresh        = EventCode(C.LV_EVENT_REFRESH)
	EventReady          = EventCode(C.LV_EVENT_READY)
	EventCancel         = EventCode(C.LV_EVENT_CANCEL)
	EventStateChanged   = EventCode(C.LV_EVENT_STATE_CHANGED)
	EventCreate         = EventCode(C.LV_EVENT_CREATE)
	EventDelete         = EventCode(C.LV_EVENT_DELETE)
	EventChildChanged   = EventCode(C.LV_EVENT_CHILD_CHANGED)
	EventScreenLoaded   = EventCode(C.LV_EVENT_SCREEN_LOADED)
	EventScreenUnloaded = EventCode(C.LV_EVENT_SCREEN_UNLOADED)
	EventSizeChanged    = EventCode(C.LV_EVENT_SIZE_CHANGED)
	EventStyleChanged   = EventCode(C.LV_EVENT_STYLE_CHANGED)
	EventLayoutChanged  = EventCode(C.LV_EVENT_LAYOUT_CHANGED)
)

// Event wraps an lv_event_t passed into an event callback. It is only
// valid for the duration of the callback that received it.
type Event struct {
	c *C.lv_event_t
}

// Code returns the event's code.
func (e *Event) Code() EventCode {
	return EventCode(C.lv_event_get_code(e.c))
}

// Target returns the object the event was originally sent to.
func (e *Event) Target() *Obj {
	return wrapObj((*C.lv_obj_t)(C.lv_event_get_target(e.c)))
}

// CurrentTarget returns the object whose event handler is currently
// running (relevant when events bubble up through parents).
func (e *Event) CurrentTarget() *Obj {
	return wrapObj((*C.lv_obj_t)(C.lv_event_get_current_target(e.c)))
}

type eventCallback struct {
	fn func(*Event)
}

// AddEventCB registers fn to be called whenever an event matching code (or
// EventAll) occurs on the object. The callback, and any others registered
// on this object, are automatically released when the object is deleted.
func (o *Obj) AddEventCB(code EventCode, fn func(*Event)) {
	o.ensureAutoCleanup()
	h := cgo.NewHandle(&eventCallback{fn: fn})
	C.lv_obj_add_event_cb(o.c, C.lv_event_cb_t(C.goEventTrampoline),
		C.lv_event_code_t(code), unsafeHandlePointer(h))
	o.handles = append(o.handles, h)
}

// ensureAutoCleanup registers a one-time LV_EVENT_DELETE hook (bypassing
// AddEventCB itself, to avoid recursing) that frees every cgo.Handle this
// object has accumulated once LVGL deletes it.
func (o *Obj) ensureAutoCleanup() {
	if o.autoCleanup {
		return
	}
	o.autoCleanup = true

	var deleteHookHandle cgo.Handle
	deleteHookHandle = cgo.NewHandle(&eventCallback{
		fn: func(e *Event) {
			for _, h := range o.handles {
				h.Delete()
			}
			o.handles = nil
			deleteHookHandle.Delete()
		},
	})
	C.lv_obj_add_event_cb(o.c, C.lv_event_cb_t(C.goEventTrampoline),
		C.lv_event_code_t(EventDelete), unsafeHandlePointer(deleteHookHandle))
}

// SendEvent synthesizes and dispatches an event of the given code to the
// object, as if LVGL itself had triggered it (lv_obj_send_event).
func (o *Obj) SendEvent(code EventCode, param unsafe.Pointer) {
	C.lv_obj_send_event(o.c, C.lv_event_code_t(code), param)
}

//export goEventTrampoline
func goEventTrampoline(e *C.lv_event_t) {
	ud := C.lv_event_get_user_data(e)
	h := cgo.Handle(uintptr(ud))
	cb := h.Value().(*eventCallback)
	cb.fn(&Event{c: e})
}
