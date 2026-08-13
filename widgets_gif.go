package lvgl

/*
#include <stdlib.h>
#include <lvgl.h>
*/
import "C"
import "unsafe"

// ColorFormat mirrors lv_color_format_t (a subset of the common ones).
type ColorFormat uint32

var (
	ColorFormatA8       = ColorFormat(C.LV_COLOR_FORMAT_A8)
	ColorFormatRGB565   = ColorFormat(C.LV_COLOR_FORMAT_RGB565)
	ColorFormatRGB888   = ColorFormat(C.LV_COLOR_FORMAT_RGB888)
	ColorFormatARGB8888 = ColorFormat(C.LV_COLOR_FORMAT_ARGB8888)
	ColorFormatNative   = ColorFormat(C.LV_COLOR_FORMAT_NATIVE)
)

// Gif wraps an lv_gif widget: plays an animated GIF.
//
// Like Image.SetSrcPath, the source path is stored by LVGL as a pointer,
// not copied, so SetSrcPath keeps it pinned on the C heap for the
// widget's lifetime, freeing any previous one on replacement.
type Gif struct {
	*Obj
	src *C.char
}

// NewGif creates a GIF player as a child of parent.
func NewGif(parent *Obj) *Gif {
	return &Gif{Obj: wrapObj(C.lv_gif_create(parent.c))}
}

// SetSrcPath sets the GIF source to a file path.
func (g *Gif) SetSrcPath(path string) {
	cPath := C.CString(path)
	C.lv_gif_set_src(g.c, unsafe.Pointer(cPath))

	if g.src != nil {
		C.free(unsafe.Pointer(g.src))
	}
	g.src = cPath
}

// Restart restarts the animation from the first frame.
func (g *Gif) Restart() { C.lv_gif_restart(g.c) }

// Pause pauses the animation.
func (g *Gif) Pause() { C.lv_gif_pause(g.c) }

// Resume resumes a paused animation.
func (g *Gif) Resume() { C.lv_gif_resume(g.c) }

// SetLoopCount sets how many times the GIF repeats (a negative count
// loops forever).
func (g *Gif) SetLoopCount(count int32) { C.lv_gif_set_loop_count(g.c, C.int32_t(count)) }

// SetColorFormat sets the pixel format used to decode frames into.
func (g *Gif) SetColorFormat(format ColorFormat) {
	C.lv_gif_set_color_format(g.c, C.lv_color_format_t(format))
}

// SetAutoPauseInvisible enables/disables automatically pausing playback
// while the widget is not visible on screen.
func (g *Gif) SetAutoPauseInvisible(enable bool) {
	C.lv_gif_set_auto_pause_invisible(g.c, C.bool(enable))
}
