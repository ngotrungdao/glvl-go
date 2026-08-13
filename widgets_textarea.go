package lvgl

/*
#include <stdlib.h>
#include <lvgl.h>
*/
import "C"
import "unsafe"

// TextAlign mirrors lv_text_align_t.
type TextAlign uint32

var (
	TextAlignAuto   = TextAlign(C.LV_TEXT_ALIGN_AUTO)
	TextAlignLeft   = TextAlign(C.LV_TEXT_ALIGN_LEFT)
	TextAlignCenter = TextAlign(C.LV_TEXT_ALIGN_CENTER)
	TextAlignRight  = TextAlign(C.LV_TEXT_ALIGN_RIGHT)
)

// TextArea wraps an lv_textarea widget.
type TextArea struct{ *Obj }

// NewTextArea creates a text area as a child of parent.
func NewTextArea(parent *Obj) *TextArea {
	return &TextArea{wrapObj(C.lv_textarea_create(parent.c))}
}

// SetText replaces the text area's content.
func (t *TextArea) SetText(text string) {
	cText := C.CString(text)
	defer C.free(unsafe.Pointer(cText))
	C.lv_textarea_set_text(t.c, cText)
}

// AddChar inserts a single character (a Unicode code point) at the cursor.
func (t *TextArea) AddChar(c rune) { C.lv_textarea_add_char(t.c, C.uint32_t(c)) }

// AddText inserts text at the cursor.
func (t *TextArea) AddText(text string) {
	cText := C.CString(text)
	defer C.free(unsafe.Pointer(cText))
	C.lv_textarea_add_text(t.c, cText)
}

// DeleteChar deletes the character before the cursor (backspace).
func (t *TextArea) DeleteChar() { C.lv_textarea_delete_char(t.c) }

// DeleteCharForward deletes the character after the cursor (delete key).
func (t *TextArea) DeleteCharForward() { C.lv_textarea_delete_char_forward(t.c) }

// SetPlaceholderText sets the text shown when the text area is empty.
func (t *TextArea) SetPlaceholderText(text string) {
	cText := C.CString(text)
	defer C.free(unsafe.Pointer(cText))
	C.lv_textarea_set_placeholder_text(t.c, cText)
}

// SetCursorPos moves the cursor to the given character index.
func (t *TextArea) SetCursorPos(pos int32) { C.lv_textarea_set_cursor_pos(t.c, C.int32_t(pos)) }

// SetCursorClickPos enables/disables moving the cursor by clicking/tapping.
func (t *TextArea) SetCursorClickPos(enable bool) {
	C.lv_textarea_set_cursor_click_pos(t.c, C.bool(enable))
}

// CursorRight moves the cursor one character right.
func (t *TextArea) CursorRight() { C.lv_textarea_cursor_right(t.c) }

// CursorLeft moves the cursor one character left.
func (t *TextArea) CursorLeft() { C.lv_textarea_cursor_left(t.c) }

// CursorDown moves the cursor one line down.
func (t *TextArea) CursorDown() { C.lv_textarea_cursor_down(t.c) }

// CursorUp moves the cursor one line up.
func (t *TextArea) CursorUp() { C.lv_textarea_cursor_up(t.c) }

// SetPasswordMode enables/disables masking entered characters.
func (t *TextArea) SetPasswordMode(enable bool) {
	C.lv_textarea_set_password_mode(t.c, C.bool(enable))
}

// SetPasswordBullet sets the string used to mask characters in password mode.
func (t *TextArea) SetPasswordBullet(bullet string) {
	cBullet := C.CString(bullet)
	defer C.free(unsafe.Pointer(cBullet))
	C.lv_textarea_set_password_bullet(t.c, cBullet)
}

// SetPasswordShowTime sets how long, in milliseconds, a newly typed
// character stays visible before masking in password mode.
func (t *TextArea) SetPasswordShowTime(ms uint32) {
	C.lv_textarea_set_password_show_time(t.c, C.uint32_t(ms))
}

// SetOneLine restricts the text area to a single line.
func (t *TextArea) SetOneLine(enable bool) {
	C.lv_textarea_set_one_line(t.c, C.bool(enable))
}

// SetAcceptedChars restricts input to only the characters in list.
func (t *TextArea) SetAcceptedChars(list string) {
	cList := C.CString(list)
	defer C.free(unsafe.Pointer(cList))
	C.lv_textarea_set_accepted_chars(t.c, cList)
}

// SetMaxLength limits how many characters can be entered.
func (t *TextArea) SetMaxLength(n uint32) {
	C.lv_textarea_set_max_length(t.c, C.uint32_t(n))
}

// SetInsertReplace sets a string to substitute for the next AddText/typed
// insertion (consumed after one use).
func (t *TextArea) SetInsertReplace(text string) {
	cText := C.CString(text)
	defer C.free(unsafe.Pointer(cText))
	C.lv_textarea_set_insert_replace(t.c, cText)
}

// SetTextSelection enables/disables text selection (e.g. via shift+arrow
// or click-drag).
func (t *TextArea) SetTextSelection(enable bool) {
	C.lv_textarea_set_text_selection(t.c, C.bool(enable))
}

// SetAlign sets the text's horizontal alignment.
func (t *TextArea) SetAlign(align TextAlign) {
	C.lv_textarea_set_align(t.c, C.lv_text_align_t(align))
}

// Text returns the text area's current content.
func (t *TextArea) Text() string {
	return C.GoString(C.lv_textarea_get_text(t.c))
}

// PlaceholderText returns the placeholder text.
func (t *TextArea) PlaceholderText() string {
	return C.GoString(C.lv_textarea_get_placeholder_text(t.c))
}

// Label returns the internal label object holding the text area's text.
func (t *TextArea) Label() *Obj { return wrapObj(C.lv_textarea_get_label(t.c)) }

// CursorPos returns the cursor's current character index.
func (t *TextArea) CursorPos() uint32 { return uint32(C.lv_textarea_get_cursor_pos(t.c)) }

// CursorClickPos reports whether clicking/tapping moves the cursor.
func (t *TextArea) CursorClickPos() bool { return bool(C.lv_textarea_get_cursor_click_pos(t.c)) }

// PasswordMode reports whether password masking is enabled.
func (t *TextArea) PasswordMode() bool { return bool(C.lv_textarea_get_password_mode(t.c)) }

// PasswordBullet returns the string used to mask characters.
func (t *TextArea) PasswordBullet() string {
	return C.GoString(C.lv_textarea_get_password_bullet(t.c))
}

// OneLine reports whether the text area is restricted to one line.
func (t *TextArea) OneLine() bool { return bool(C.lv_textarea_get_one_line(t.c)) }

// AcceptedChars returns the accepted character set, if any was set.
func (t *TextArea) AcceptedChars() string {
	return C.GoString(C.lv_textarea_get_accepted_chars(t.c))
}

// MaxLength returns the maximum character count, if one was set.
func (t *TextArea) MaxLength() uint32 { return uint32(C.lv_textarea_get_max_length(t.c)) }

// TextSelection reports whether text selection is enabled.
func (t *TextArea) TextSelection() bool { return bool(C.lv_textarea_get_text_selection(t.c)) }

// PasswordShowTime returns how long a typed character stays unmasked, in
// milliseconds.
func (t *TextArea) PasswordShowTime() uint32 {
	return uint32(C.lv_textarea_get_password_show_time(t.c))
}

// CurrentChar returns the Unicode code point immediately before the cursor.
func (t *TextArea) CurrentChar() rune { return rune(C.lv_textarea_get_current_char(t.c)) }
