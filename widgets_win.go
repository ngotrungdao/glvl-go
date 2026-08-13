package lvgl

/*
#include <stdlib.h>
#include <lvgl.h>
*/
import "C"
import "unsafe"

// Win wraps an lv_win widget: a window chrome with a header bar (title +
// buttons) and a content area. LVGL itself deprecates lv_win in favor of
// building a window from a flex column (see layouts_flex.go); this wrapper
// still exposes it since it works and is the simplest option, but a
// SetFlexFlow(FlexFlowColumn) container may be preferable for new code.
type Win struct{ *Obj }

// NewWin creates a window as a child of parent.
func NewWin(parent *Obj) *Win {
	return &Win{wrapObj(C.lv_win_create(parent.c))}
}

// AddTitle adds a title label to the window's header.
func (w *Win) AddTitle(text string) *Obj {
	cText := C.CString(text)
	defer C.free(unsafe.Pointer(cText))
	return wrapObj(C.lv_win_add_title(w.c, cText))
}

// AddButton adds a button (e.g. a symbol icon) to the window's header,
// btnWidth pixels wide.
//
// Like Image.SetSrcPath, icon is passed to LVGL as an image source, which
// stores the pointer rather than copying it; since AddButton has no
// natural point to free it later, the icon string is intentionally never
// freed (a small, bounded leak, traded for correctness).
func (w *Win) AddButton(icon string, btnWidth int32) *Obj {
	return wrapObj(C.lv_win_add_button(w.c, unsafe.Pointer(C.CString(icon)), C.int32_t(btnWidth)))
}

// Content returns the window's content container, to which children
// should be added.
func (w *Win) Content() *Obj { return wrapObj(C.lv_win_get_content(w.c)) }

// Header returns the window's header container.
func (w *Win) Header() *Obj { return wrapObj(C.lv_win_get_header(w.c)) }
