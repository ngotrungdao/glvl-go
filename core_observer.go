package lvgl

/*
#include <stdlib.h>
#include <lvgl.h>

extern void goObserverTrampoline(lv_observer_t *observer, lv_subject_t *subject);
*/
import "C"
import (
	"runtime"
	"runtime/cgo"
	"unsafe"
)

type subjectKind uint8

const (
	subjectKindInt subjectKind = iota
	subjectKindFloat
	subjectKindColor
	subjectKindString
	subjectKindPointer
)

// Subject wraps an lv_subject_t: an observable value (int, float, color,
// string, or pointer) that widgets or your own callbacks can react to
// automatically when it changes, via AddObserver or a widget's BindValue/
// BindText/BindChecked method (LVGL 9's reactive data-binding system).
//
// LVGL stores a raw, non-refcounted pointer to a Subject's C memory for
// as long as anything observes or binds to it, so — like Style — it's
// allocated on the C heap rather than as ordinary Go memory. Call Delete
// once nothing observes/binds to it anymore.
type Subject struct {
	c            *C.lv_subject_t
	kind         subjectKind
	buf, prevBuf *C.char // string subjects only
}

func newSubject(kind subjectKind) *Subject {
	c := (*C.lv_subject_t)(C.malloc(C.sizeof_lv_subject_t))
	s := &Subject{c: c, kind: kind}
	runtime.SetFinalizer(s, (*Subject).leakWarn)
	return s
}

func (s *Subject) leakWarn() {
	println("lvgl: Subject finalized without an explicit Delete() call; its C memory was leaked")
}

// NewSubjectInt creates an observable int32 value.
func NewSubjectInt(initial int32) *Subject {
	s := newSubject(subjectKindInt)
	C.lv_subject_init_int(s.c, C.int32_t(initial))
	return s
}

// NewSubjectFloat creates an observable float value (needs LV_USE_FLOAT,
// enabled in this build).
func NewSubjectFloat(initial float32) *Subject {
	s := newSubject(subjectKindFloat)
	C.lv_subject_init_float(s.c, C.float(initial))
	return s
}

// NewSubjectColor creates an observable Color value.
func NewSubjectColor(initial Color) *Subject {
	s := newSubject(subjectKindColor)
	C.lv_subject_init_color(s.c, initial.toC())
	return s
}

// NewSubjectString creates an observable string value. bufSize must be
// large enough for the longest string ever stored, including the NUL
// terminator; LVGL keeps both the current and previous value in two
// buffers of this size.
func NewSubjectString(initial string, bufSize int) *Subject {
	s := newSubject(subjectKindString)
	s.buf = (*C.char)(C.malloc(C.size_t(bufSize)))
	s.prevBuf = (*C.char)(C.malloc(C.size_t(bufSize)))
	cInitial := C.CString(initial)
	defer C.free(unsafe.Pointer(cInitial))
	C.lv_subject_init_string(s.c, s.buf, s.prevBuf, C.size_t(bufSize), cInitial)
	return s
}

// Delete deinitializes and frees the subject's C memory (and its string
// buffers, if any). Only call this once nothing observes or is bound to
// it anymore.
func (s *Subject) Delete() {
	if s.c == nil {
		return
	}
	C.lv_subject_deinit(s.c)
	C.free(unsafe.Pointer(s.c))
	s.c = nil
	if s.buf != nil {
		C.free(unsafe.Pointer(s.buf))
		C.free(unsafe.Pointer(s.prevBuf))
		s.buf, s.prevBuf = nil, nil
	}
	runtime.SetFinalizer(s, nil)
}

// SetInt sets the subject's value and notifies observers if it changed.
func (s *Subject) SetInt(v int32) { C.lv_subject_set_int(s.c, C.int32_t(v)) }

// Int returns the subject's current int value.
func (s *Subject) Int() int32 { return int32(C.lv_subject_get_int(s.c)) }

// PreviousInt returns the subject's value before the last change.
func (s *Subject) PreviousInt() int32 { return int32(C.lv_subject_get_previous_int(s.c)) }

// SetFloat sets the subject's value and notifies observers if it changed.
func (s *Subject) SetFloat(v float32) { C.lv_subject_set_float(s.c, C.float(v)) }

// SetColor sets the subject's value and notifies observers if it changed.
func (s *Subject) SetColor(v Color) { C.lv_subject_set_color(s.c, v.toC()) }

// Color returns the subject's current color value.
func (s *Subject) Color() Color { return colorFromC(C.lv_subject_get_color(s.c)) }

// PreviousColor returns the subject's color before the last change.
func (s *Subject) PreviousColor() Color { return colorFromC(C.lv_subject_get_previous_color(s.c)) }

// SetString sets the subject's value (copied into its internal buffer)
// and notifies observers if it changed.
func (s *Subject) SetString(v string) {
	cV := C.CString(v)
	defer C.free(unsafe.Pointer(cV))
	C.lv_subject_copy_string(s.c, cV)
}

// String returns the subject's current string value.
func (s *Subject) String() string { return C.GoString(C.lv_subject_get_string(s.c)) }

// PreviousString returns the subject's string before the last change.
func (s *Subject) PreviousString() string {
	return C.GoString(C.lv_subject_get_previous_string(s.c))
}

// Notify forces all observers to be called, even if the value hasn't
// changed.
func (s *Subject) Notify() { C.lv_subject_notify(s.c) }

type observerCallback struct {
	fn func(*Subject)
}

// Observer wraps an lv_observer_t, returned from AddObserver or a
// widget's Bind* method. Call Remove to stop observing.
type Observer struct {
	c *C.lv_observer_t
	h cgo.Handle
}

// Remove stops the observer from being called and frees its resources.
func (o *Observer) Remove() {
	if o.c == nil {
		return
	}
	C.lv_observer_remove(o.c)
	if o.h != 0 {
		o.h.Delete()
	}
	o.c = nil
}

// AddObserver registers fn to be called whenever the subject's value
// changes (and once immediately upon registration). Call Remove on the
// returned Observer to stop.
func (s *Subject) AddObserver(fn func(subject *Subject)) *Observer {
	h := cgo.NewHandle(&observerCallback{fn: fn})
	c := C.lv_subject_add_observer(s.c, C.lv_observer_cb_t(C.goObserverTrampoline), unsafeHandlePointer(h))
	return &Observer{c: c, h: h}
}

//export goObserverTrampoline
func goObserverTrampoline(observer *C.lv_observer_t, subject *C.lv_subject_t) {
	ud := C.lv_observer_get_user_data(observer)
	h := cgo.Handle(uintptr(ud))
	cb, ok := h.Value().(*observerCallback)
	if !ok {
		return
	}
	// The Subject wrapper reconstructed here is a thin view onto the same
	// C memory the caller's own *Subject already points at; it does not
	// take ownership (no finalizer), so it's safe to let it be GC'd
	// without ever calling Delete.
	cb.fn(&Subject{c: subject})
}

// BindChecked links the object's StateChecked (e.g. a Checkbox or
// Switch) to an int Subject: nonzero means checked.
func (o *Obj) BindChecked(subject *Subject) *Observer {
	return &Observer{c: C.lv_obj_bind_checked(o.c, subject.c)}
}
