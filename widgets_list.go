package lvgl

/*
#include <stdlib.h>
#include <lvgl.h>
*/
import "C"
import "unsafe"

// List wraps an lv_list widget. LVGL itself deprecates lv_list in favor of
// building a list from a flex column (see layouts_flex.go); this wrapper
// still exposes it since it works and is the simplest option, but a
// SetFlexFlow(FlexFlowColumn) container may be preferable for new code.
type List struct{ *Obj }

// NewList creates a list as a child of parent.
func NewList(parent *Obj) *List {
	return &List{wrapObj(C.lv_list_create(parent.c))}
}

// AddText adds a non-interactive text entry (e.g. a section header).
func (l *List) AddText(text string) *Obj {
	cText := C.CString(text)
	defer C.free(unsafe.Pointer(cText))
	return wrapObj(C.lv_list_add_text(l.c, cText))
}

// AddButton adds a clickable button entry with optional icon (a symbol
// string, e.g. an LV_SYMBOL_* constant, or "" for none) and text.
//
// Like Image.SetSrcPath, the icon is passed to LVGL as an image source,
// which stores the pointer rather than copying it; since AddButton has no
// natural point to free it later, the icon string (if any) is
// intentionally never freed. This is a small, bounded leak (a handful of
// bytes per call with a non-empty icon), traded for correctness.
func (l *List) AddButton(icon, text string) *Obj {
	cText := C.CString(text)
	defer C.free(unsafe.Pointer(cText))
	var iconPtr unsafe.Pointer
	if icon != "" {
		iconPtr = unsafe.Pointer(C.CString(icon))
	}
	return wrapObj(C.lv_list_add_button(l.c, iconPtr, cText))
}

// ButtonText returns a list button's text.
func (l *List) ButtonText(btn *Obj) string {
	return C.GoString(C.lv_list_get_button_text(l.c, btn.c))
}

// SetButtonText changes a list button's text.
func (l *List) SetButtonText(btn *Obj, text string) {
	cText := C.CString(text)
	defer C.free(unsafe.Pointer(cText))
	C.lv_list_set_button_text(l.c, btn.c, cText)
}
