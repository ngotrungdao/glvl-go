package lvgl

/*
#include <stdlib.h>
#include <lvgl.h>
*/
import "C"
import "unsafe"

// SpanOverflow mirrors lv_span_overflow_t.
type SpanOverflow uint32

var (
	SpanOverflowClip     = SpanOverflow(C.LV_SPAN_OVERFLOW_CLIP)
	SpanOverflowEllipsis = SpanOverflow(C.LV_SPAN_OVERFLOW_ELLIPSIS)
)

// SpanMode mirrors lv_span_mode_t.
type SpanMode uint32

var (
	SpanModeFixed  = SpanMode(C.LV_SPAN_MODE_FIXED)
	SpanModeExpand = SpanMode(C.LV_SPAN_MODE_EXPAND)
	SpanModeBreak  = SpanMode(C.LV_SPAN_MODE_BREAK)
)

// SpanGroup wraps an lv_spangroup widget: a container of independently
// styled text runs ("spans") flowed together, e.g. for rich/mixed-style
// text.
type SpanGroup struct{ *Obj }

// NewSpanGroup creates a span group as a child of parent.
func NewSpanGroup(parent *Obj) *SpanGroup {
	return &SpanGroup{wrapObj(C.lv_spangroup_create(parent.c))}
}

// SetAlign sets the horizontal alignment of the flowed text.
func (g *SpanGroup) SetAlign(align TextAlign) {
	C.lv_spangroup_set_align(g.c, C.lv_text_align_t(align))
}

// SetOverflow sets how text longer than the group's area is handled.
func (g *SpanGroup) SetOverflow(overflow SpanOverflow) {
	C.lv_spangroup_set_overflow(g.c, C.lv_span_overflow_t(overflow))
}

// SetIndent sets the first-line indent, in pixels.
func (g *SpanGroup) SetIndent(indent int32) { C.lv_spangroup_set_indent(g.c, C.int32_t(indent)) }

// SetMode sets how the group's own size responds to its text content.
func (g *SpanGroup) SetMode(mode SpanMode) { C.lv_spangroup_set_mode(g.c, C.lv_span_mode_t(mode)) }

// SetMaxLines caps how many lines are shown (only relevant in SpanModeBreak).
func (g *SpanGroup) SetMaxLines(lines int32) { C.lv_spangroup_set_max_lines(g.c, C.int32_t(lines)) }

// DeleteSpan removes a span from the group.
func (g *SpanGroup) DeleteSpan(s *Span) { C.lv_spangroup_delete_span(g.c, s.c) }

// Span wraps an lv_span_t, a single styled text run within a SpanGroup.
// It is owned by the SpanGroup it was added to.
type Span struct {
	c     *C.lv_span_t
	group *SpanGroup
}

// AddSpan appends a new, empty span to the group.
func (g *SpanGroup) AddSpan() *Span {
	return &Span{c: C.lv_spangroup_add_span(g.c), group: g}
}

// SetText sets the span's text.
func (s *Span) SetText(text string) {
	cText := C.CString(text)
	defer C.free(unsafe.Pointer(cText))
	C.lv_span_set_text(s.c, cText)
}

// Text returns the span's current text.
func (s *Span) Text() string { return C.GoString(C.lv_span_get_text(s.c)) }

// SetStyle attaches a style to this span's text run.
func (s *Span) SetStyle(style *Style) {
	C.lv_spangroup_set_span_style(s.group.c, s.c, style.c)
}

// Refresh must be called after modifying spans for the changes to take
// effect (lv_spangroup_refresh).
func (g *SpanGroup) Refresh() {
	C.lv_spangroup_refresh(g.c)
}
