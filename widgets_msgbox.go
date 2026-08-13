package lvgl

/*
#include <stdlib.h>
#include <lvgl.h>
*/
import "C"
import "unsafe"

// MsgBox wraps an lv_msgbox widget. Content is built compositionally:
// create it, then add a title/text/buttons as needed.
//
// Like Image.SetSrcPath, the icon passed to AddHeaderButton is stored by
// LVGL as a pointer, not copied; since a msgbox can have multiple header
// buttons, each icon string is intentionally never freed (a small,
// bounded leak per call, traded for correctness — see List.AddButton).
type MsgBox struct{ *Obj }

// NewMsgBox creates a message box as a child of parent.
func NewMsgBox(parent *Obj) *MsgBox {
	return &MsgBox{wrapObj(C.lv_msgbox_create(parent.c))}
}

// AddTitle adds a title to the message box's header.
func (m *MsgBox) AddTitle(title string) *Obj {
	cTitle := C.CString(title)
	defer C.free(unsafe.Pointer(cTitle))
	return wrapObj(C.lv_msgbox_add_title(m.c, cTitle))
}

// AddText adds a paragraph of body text.
func (m *MsgBox) AddText(text string) *Obj {
	cText := C.CString(text)
	defer C.free(unsafe.Pointer(cText))
	return wrapObj(C.lv_msgbox_add_text(m.c, cText))
}

// AddFooterButton adds a labeled button to the message box's footer.
func (m *MsgBox) AddFooterButton(text string) *Obj {
	cText := C.CString(text)
	defer C.free(unsafe.Pointer(cText))
	return wrapObj(C.lv_msgbox_add_footer_button(m.c, cText))
}

// AddCloseButton adds a close ("X") button to the header.
func (m *MsgBox) AddCloseButton() *Obj {
	return wrapObj(C.lv_msgbox_add_close_button(m.c))
}

// AddHeaderButton adds an icon-only button (e.g. a Symbol* constant) to
// the message box's header.
func (m *MsgBox) AddHeaderButton(icon string) *Obj {
	return wrapObj(C.lv_msgbox_add_header_button(m.c, unsafe.Pointer(C.CString(icon))))
}

// Header returns the message box's header container.
func (m *MsgBox) Header() *Obj { return wrapObj(C.lv_msgbox_get_header(m.c)) }

// Footer returns the message box's footer container.
func (m *MsgBox) Footer() *Obj { return wrapObj(C.lv_msgbox_get_footer(m.c)) }

// Content returns the message box's body content container.
func (m *MsgBox) Content() *Obj { return wrapObj(C.lv_msgbox_get_content(m.c)) }

// Title returns the message box's title label.
func (m *MsgBox) Title() *Obj { return wrapObj(C.lv_msgbox_get_title(m.c)) }

// Close closes and deletes the message box immediately.
func (m *MsgBox) Close() { C.lv_msgbox_close(m.c) }

// CloseAsync schedules the message box to close on the next event-loop
// iteration (safe to call from within one of its own event callbacks).
func (m *MsgBox) CloseAsync() { C.lv_msgbox_close_async(m.c) }
