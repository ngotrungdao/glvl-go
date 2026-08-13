package lvgl

/*
#include <stdlib.h>
#include <lvgl.h>
*/
import "C"
import "unsafe"

// LabelLongMode mirrors lv_label_long_mode_t.
type LabelLongMode uint32

var (
	LabelLongModeWrap           = LabelLongMode(C.LV_LABEL_LONG_MODE_WRAP)
	LabelLongModeDots           = LabelLongMode(C.LV_LABEL_LONG_MODE_DOTS)
	LabelLongModeScroll         = LabelLongMode(C.LV_LABEL_LONG_MODE_SCROLL)
	LabelLongModeScrollCircular = LabelLongMode(C.LV_LABEL_LONG_MODE_SCROLL_CIRCULAR)
	LabelLongModeClip           = LabelLongMode(C.LV_LABEL_LONG_MODE_CLIP)
)

// Label wraps an lv_label widget.
type Label struct{ *Obj }

// NewLabel creates a label as a child of parent.
func NewLabel(parent *Obj) *Label {
	return &Label{wrapObj(C.lv_label_create(parent.c))}
}

// SetText sets the label's text.
func (l *Label) SetText(text string) {
	cText := C.CString(text)
	defer C.free(unsafe.Pointer(cText))
	C.lv_label_set_text(l.c, cText)
}

// Text returns the label's current text.
func (l *Label) Text() string {
	return C.GoString(C.lv_label_get_text(l.c))
}

// SetRecolor enables/disables inline `#color# text#` recoloring.
func (l *Label) SetRecolor(enable bool) {
	C.lv_label_set_recolor(l.c, C.bool(enable))
}

// Recolor reports whether inline `#color# text#` recoloring is enabled.
func (l *Label) Recolor() bool { return bool(C.lv_label_get_recolor(l.c)) }

// SetLongMode sets how the label handles text that doesn't fit its size.
func (l *Label) SetLongMode(mode LabelLongMode) {
	C.lv_label_set_long_mode(l.c, C.lv_label_long_mode_t(mode))
}

// SetMaxLines caps how many lines of text are shown (0 means unlimited).
func (l *Label) SetMaxLines(lines int32) { C.lv_label_set_max_lines(l.c, C.int32_t(lines)) }

// MaxLines returns the current maximum line count (0 means unlimited).
func (l *Label) MaxLines() int32 { return int32(C.lv_label_get_max_lines(l.c)) }

// SetTextSelectionStart sets the start character index of a text
// selection highlight (needs LV_LABEL_TEXT_SELECTION enabled in lv_conf).
func (l *Label) SetTextSelectionStart(index uint32) {
	C.lv_label_set_text_selection_start(l.c, C.uint32_t(index))
}

// SetTextSelectionEnd sets the end character index of a text selection highlight.
func (l *Label) SetTextSelectionEnd(index uint32) {
	C.lv_label_set_text_selection_end(l.c, C.uint32_t(index))
}

// TextSelectionStart returns the selection highlight's start character index.
func (l *Label) TextSelectionStart() uint32 {
	return uint32(C.lv_label_get_text_selection_start(l.c))
}

// TextSelectionEnd returns the selection highlight's end character index.
func (l *Label) TextSelectionEnd() uint32 {
	return uint32(C.lv_label_get_text_selection_end(l.c))
}

// BindText links the label's text to a Subject: whenever it changes, the
// label's text is set to fmt with the value substituted in (a printf-style
// format string with exactly one matching verb, e.g. "%d" for an int
// Subject or "%s" for a string one), or fmt's value directly (no
// substitution) if fmt is "".
func (l *Label) BindText(subject *Subject, fmt string) *Observer {
	var cFmt *C.char
	if fmt != "" {
		cFmt = C.CString(fmt)
		defer C.free(unsafe.Pointer(cFmt))
	}
	return &Observer{c: C.lv_label_bind_text(l.c, subject.c, cFmt)}
}
