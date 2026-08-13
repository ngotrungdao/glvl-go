package lvgl

/*
#include <stdlib.h>
#include <lvgl.h>
*/
import "C"
import (
	"errors"
	"unsafe"
)

// GStreamer source factory/property pairs, reimplemented from LVGL's
// LV_GSTREAMER_FACTORY_*/LV_GSTREAMER_PROPERTY_* string macros (cgo can't
// reach #define string constants). Pass the Factory/Property pair
// together to SetSrc, e.g. SetSrc(GStreamerFactoryFile,
// GStreamerPropertyFile, "/path/to/video.mp4").
const (
	GStreamerFactoryURIDecode  = "uridecodebin"
	GStreamerPropertyURIDecode = "uri"

	GStreamerFactoryFile  = "filesrc"
	GStreamerPropertyFile = "location"

	GStreamerFactoryHTTP  = "souphttpsrc"
	GStreamerPropertyHTTP = "location"

	GStreamerFactoryHTTPS  = "souphttpsrc"
	GStreamerPropertyHTTPS = "location"

	GStreamerFactoryV4L2Camera  = "v4l2src"
	GStreamerPropertyV4L2Camera = "device"

	GStreamerFactoryAlsaAudio  = "alsasrc"
	GStreamerPropertyAlsaAudio = "device"

	GStreamerFactoryPulseAudio  = "pulsesrc"
	GStreamerPropertyPulseAudio = "device"

	GStreamerFactoryTestAudio = "audiotestsrc"
	GStreamerFactoryTestVideo = "videotestsrc"
	GStreamerFactoryApp       = "appsrc"
)

// GStreamerState mirrors lv_gstreamer_state_t.
type GStreamerState uint32

var (
	GStreamerStateNull    = GStreamerState(C.LV_GSTREAMER_STATE_NULL)
	GStreamerStateReady   = GStreamerState(C.LV_GSTREAMER_STATE_READY)
	GStreamerStatePaused  = GStreamerState(C.LV_GSTREAMER_STATE_PAUSED)
	GStreamerStatePlaying = GStreamerState(C.LV_GSTREAMER_STATE_PLAYING)
)

var errGStreamerSetSrc = errors.New("lvgl: gstreamer source setup failed")

// GStreamer wraps an lv_gstreamer widget: plays audio/video through a
// GStreamer pipeline, rendering video frames into the widget like an
// Image.
type GStreamer struct{ *Obj }

// NewGStreamer creates a GStreamer player as a child of parent.
func NewGStreamer(parent *Obj) *GStreamer {
	return &GStreamer{wrapObj(C.lv_gstreamer_create(parent.c))}
}

// SetSrc configures the pipeline's source element. factoryName and
// property are typically one of the GStreamerFactory*/GStreamerProperty*
// pairs above (property may be "" to create the source without setting a
// property). source is the property's value, e.g. a file path or URI.
func (g *GStreamer) SetSrc(factoryName, property, source string) error {
	cFactory := C.CString(factoryName)
	defer C.free(unsafe.Pointer(cFactory))

	var cProp, cSrc *C.char
	if property != "" {
		cProp = C.CString(property)
		defer C.free(unsafe.Pointer(cProp))
	}
	if source != "" {
		cSrc = C.CString(source)
		defer C.free(unsafe.Pointer(cSrc))
	}

	if C.lv_gstreamer_set_src(g.c, cFactory, cProp, cSrc) != C.LV_RESULT_OK {
		return errGStreamerSetSrc
	}
	return nil
}

// Play starts/resumes playback.
func (g *GStreamer) Play() { C.lv_gstreamer_play(g.c) }

// Pause pauses playback.
func (g *GStreamer) Pause() { C.lv_gstreamer_pause(g.c) }

// Stop stops playback.
func (g *GStreamer) Stop() { C.lv_gstreamer_stop(g.c) }

// SetPosition seeks to the given position, in milliseconds.
func (g *GStreamer) SetPosition(ms uint32) { C.lv_gstreamer_set_position(g.c, C.uint32_t(ms)) }

// Duration returns the stream's total duration, in milliseconds.
func (g *GStreamer) Duration() uint32 { return uint32(C.lv_gstreamer_get_duration(g.c)) }

// Position returns the current playback position, in milliseconds.
func (g *GStreamer) Position() uint32 { return uint32(C.lv_gstreamer_get_position(g.c)) }

// State returns the pipeline's current state.
func (g *GStreamer) State() GStreamerState {
	return GStreamerState(C.lv_gstreamer_get_state(g.c))
}

// SetVolume sets the playback volume, 0-100 (values above 100 are clamped).
func (g *GStreamer) SetVolume(volume uint8) { C.lv_gstreamer_set_volume(g.c, C.uint8_t(volume)) }

// Volume returns the current playback volume.
func (g *GStreamer) Volume() uint8 { return uint8(C.lv_gstreamer_get_volume(g.c)) }

// SetRate sets the playback speed: 256 = 1x, 128 = 0.5x, 512 = 2x, etc.
func (g *GStreamer) SetRate(rate uint32) { C.lv_gstreamer_set_rate(g.c, C.uint32_t(rate)) }

// SetInnerAlign sets how the decoded video frame is positioned/scaled
// within the widget's own object box (set via SetSize) — e.g.
// ImageAlignContain to scale the video to fit the widget's layout size
// while preserving its aspect ratio (letterboxed). GStreamer is an
// lv_image subclass at the C level, so this reuses Image's
// lv_image_align_t API directly.
func (g *GStreamer) SetInnerAlign(align ImageAlign) {
	C.lv_image_set_inner_align(g.c, C.lv_image_align_t(align))
}

// InnerAlign returns the current inner alignment mode.
func (g *GStreamer) InnerAlign() ImageAlign { return ImageAlign(C.lv_image_get_inner_align(g.c)) }
