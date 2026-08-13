package lvgl

/*
#include <lvgl.h>
*/
import "C"

// KeyboardMode mirrors lv_keyboard_mode_t.
type KeyboardMode uint32

var (
	KeyboardModeTextLower = KeyboardMode(C.LV_KEYBOARD_MODE_TEXT_LOWER)
	KeyboardModeTextUpper = KeyboardMode(C.LV_KEYBOARD_MODE_TEXT_UPPER)
	KeyboardModeSpecial   = KeyboardMode(C.LV_KEYBOARD_MODE_SPECIAL)
	KeyboardModeNumber    = KeyboardMode(C.LV_KEYBOARD_MODE_NUMBER)
)

// Keyboard wraps an lv_keyboard widget: an on-screen keyboard that types
// into a linked TextArea.
type Keyboard struct{ *Obj }

// NewKeyboard creates a keyboard as a child of parent.
func NewKeyboard(parent *Obj) *Keyboard {
	return &Keyboard{wrapObj(C.lv_keyboard_create(parent.c))}
}

// SetTextArea links the keyboard to the text area it types into.
func (k *Keyboard) SetTextArea(ta *TextArea) {
	C.lv_keyboard_set_textarea(k.c, ta.c)
}

// SetMode switches between the lowercase/uppercase/special/number layouts.
func (k *Keyboard) SetMode(mode KeyboardMode) {
	C.lv_keyboard_set_mode(k.c, C.lv_keyboard_mode_t(mode))
}

// SetPopovers enables/disables showing a popover with the pressed key's
// character while typing.
func (k *Keyboard) SetPopovers(enable bool) { C.lv_keyboard_set_popovers(k.c, C.bool(enable)) }

// TextArea returns the linked text area, if any.
func (k *Keyboard) TextArea() *Obj { return wrapObj(C.lv_keyboard_get_textarea(k.c)) }

// Popovers reports whether key-press popovers are enabled.
func (k *Keyboard) Popovers() bool { return bool(C.lv_keyboard_get_popovers(k.c)) }

// SelectedButton returns the index of the currently selected key.
func (k *Keyboard) SelectedButton() uint32 { return uint32(C.lv_keyboard_get_selected_button(k.c)) }

// ButtonText returns the label text of the key at the given index.
func (k *Keyboard) ButtonText(btnID uint32) string {
	return C.GoString(C.lv_keyboard_get_button_text(k.c, C.uint32_t(btnID)))
}
