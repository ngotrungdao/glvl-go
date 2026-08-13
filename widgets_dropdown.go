package lvgl

/*
#include <stdlib.h>
#include <lvgl.h>
*/
import "C"
import "unsafe"

// DropdownPosLast passed to AddOption inserts the new option at the end
// (LV_DROPDOWN_POS_LAST).
const DropdownPosLast uint32 = 0xFFFF

// Dropdown wraps an lv_dropdown widget.
//
// Like Image.SetSrcPath, the symbol passed to SetSymbol is stored by LVGL
// as a pointer, not copied, so it's pinned on the C heap for the
// dropdown's lifetime, freeing any previous one on replacement.
type Dropdown struct {
	*Obj
	symbol *C.char
}

// NewDropdown creates a dropdown as a child of parent.
func NewDropdown(parent *Obj) *Dropdown {
	return &Dropdown{Obj: wrapObj(C.lv_dropdown_create(parent.c))}
}

// SetText sets the text shown on the closed dropdown (independent of the
// selected option; leave unset to show the selected option's text).
func (d *Dropdown) SetText(text string) {
	cText := C.CString(text)
	defer C.free(unsafe.Pointer(cText))
	C.lv_dropdown_set_text(d.c, cText)
}

// SetOptions sets the list of selectable options, newline-separated.
func (d *Dropdown) SetOptions(options string) {
	cOpt := C.CString(options)
	defer C.free(unsafe.Pointer(cOpt))
	C.lv_dropdown_set_options(d.c, cOpt)
}

// AddOption inserts a single option at pos (or DropdownPosLast to append).
func (d *Dropdown) AddOption(option string, pos uint32) {
	cOpt := C.CString(option)
	defer C.free(unsafe.Pointer(cOpt))
	C.lv_dropdown_add_option(d.c, cOpt, C.uint32_t(pos))
}

// ClearOptions removes all options.
func (d *Dropdown) ClearOptions() { C.lv_dropdown_clear_options(d.c) }

// SetSelected selects the option at the given index.
func (d *Dropdown) SetSelected(idx uint32) {
	C.lv_dropdown_set_selected(d.c, C.uint32_t(idx))
}

// SetDir sets which direction the option list opens in.
func (d *Dropdown) SetDir(dir Dir) { C.lv_dropdown_set_dir(d.c, C.lv_dir_t(dir)) }

// SetSymbol sets the icon shown on the closed dropdown (e.g. a Symbol*
// constant), or "" to remove it.
func (d *Dropdown) SetSymbol(symbol string) {
	var cSymbol *C.char
	if symbol != "" {
		cSymbol = C.CString(symbol)
	}
	C.lv_dropdown_set_symbol(d.c, unsafe.Pointer(cSymbol))

	if d.symbol != nil {
		C.free(unsafe.Pointer(d.symbol))
	}
	d.symbol = cSymbol
}

// SetSelectedHighlight enables/disables highlighting the selected option
// in the open list.
func (d *Dropdown) SetSelectedHighlight(enable bool) {
	C.lv_dropdown_set_selected_highlight(d.c, C.bool(enable))
}

// List returns the popup list object shown while the dropdown is open.
func (d *Dropdown) List() *Obj { return wrapObj(C.lv_dropdown_get_list(d.c)) }

// Text returns the dropdown's fixed text (set via SetText), if any.
func (d *Dropdown) Text() string { return C.GoString(C.lv_dropdown_get_text(d.c)) }

// Options returns the newline-separated option list.
func (d *Dropdown) Options() string {
	return C.GoString(C.lv_dropdown_get_options(d.c))
}

// Selected returns the index of the currently selected option.
func (d *Dropdown) Selected() uint32 {
	return uint32(C.lv_dropdown_get_selected(d.c))
}

// OptionCount returns how many options there are.
func (d *Dropdown) OptionCount() uint32 { return uint32(C.lv_dropdown_get_option_count(d.c)) }

// SelectedStr returns the text of the currently selected option.
func (d *Dropdown) SelectedStr() string {
	const bufSize = 128
	buf := (*C.char)(C.malloc(bufSize))
	defer C.free(unsafe.Pointer(buf))
	C.lv_dropdown_get_selected_str(d.c, buf, bufSize)
	return C.GoString(buf)
}

// OptionIndex returns the index of the given option text, or -1 if not found.
func (d *Dropdown) OptionIndex(option string) int32 {
	cOpt := C.CString(option)
	defer C.free(unsafe.Pointer(cOpt))
	return int32(C.lv_dropdown_get_option_index(d.c, cOpt))
}

// Symbol returns the icon set via SetSymbol, if any.
func (d *Dropdown) Symbol() string { return C.GoString(C.lv_dropdown_get_symbol(d.c)) }

// SelectedHighlight reports whether the selected option is highlighted in
// the open list.
func (d *Dropdown) SelectedHighlight() bool {
	return bool(C.lv_dropdown_get_selected_highlight(d.c))
}

// Open opens the option list.
func (d *Dropdown) Open() { C.lv_dropdown_open(d.c) }

// Close closes the option list.
func (d *Dropdown) Close() { C.lv_dropdown_close(d.c) }

// IsOpen reports whether the option list is currently open.
func (d *Dropdown) IsOpen() bool { return bool(C.lv_dropdown_is_open(d.c)) }

// BindValue links the dropdown's selected index to an int Subject.
func (d *Dropdown) BindValue(subject *Subject) *Observer {
	return &Observer{c: C.lv_dropdown_bind_value(d.c, subject.c)}
}
