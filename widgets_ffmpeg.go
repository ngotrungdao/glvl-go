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

// FFmpegPlayerCmd mirrors lv_ffmpeg_player_cmd_t.
type FFmpegPlayerCmd uint32

const (
	FFmpegCmdStart  = FFmpegPlayerCmd(C.LV_FFMPEG_PLAYER_CMD_START)
	FFmpegCmdStop   = FFmpegPlayerCmd(C.LV_FFMPEG_PLAYER_CMD_STOP)
	FFmpegCmdPause  = FFmpegPlayerCmd(C.LV_FFMPEG_PLAYER_CMD_PAUSE)
	FFmpegCmdResume = FFmpegPlayerCmd(C.LV_FFMPEG_PLAYER_CMD_RESUME)
)

var errFFmpegSetSrc = errors.New("lvgl: ffmpeg player source setup failed (bad path or unreadable file?)")

// FFmpegInit registers LVGL's FFmpeg image decoder and must be called
// once before FFmpegGetFrameNum, NewFFmpeg, or opening an image whose
// path is handled by the FFmpeg decoder. Unlike FreeTypeInit this has no
// non-fatal "already initialized" case to work around: lv_ffmpeg_init
// returns void.
func FFmpegInit() { C.lv_ffmpeg_init() }

// FFmpegDeinit unregisters LVGL's FFmpeg image decoder.
func FFmpegDeinit() { C.lv_ffmpeg_deinit() }

// FFmpegGetFrameNum returns the number of frames in the image or video
// file at path. FFmpegInit must be called first.
func FFmpegGetFrameNum(path string) (int, error) {
	cPath := C.CString(path)
	defer C.free(unsafe.Pointer(cPath))

	n := int(C.lv_ffmpeg_get_frame_num(cPath))
	if n < 0 {
		return 0, errors.New("lvgl: ffmpeg failed to read frame count (bad path or unreadable file?)")
	}
	return n, nil
}

// FFmpeg wraps an lv_ffmpeg_player widget: decodes and plays a video
// file through FFmpeg, rendering frames into the widget like an Image.
// FFmpegInit must be called before NewFFmpeg.
type FFmpeg struct{ *Obj }

// NewFFmpeg creates an FFmpeg player as a child of parent.
func NewFFmpeg(parent *Obj) *FFmpeg {
	return &FFmpeg{wrapObj(C.lv_ffmpeg_player_create(parent.c))}
}

// SetSrc sets the video file to play.
func (f *FFmpeg) SetSrc(path string) error {
	cPath := C.CString(path)
	defer C.free(unsafe.Pointer(cPath))

	if C.lv_ffmpeg_player_set_src(f.c, cPath) != C.LV_RESULT_OK {
		return errFFmpegSetSrc
	}
	return nil
}

// SetCmd controls playback: start, stop, pause, or resume.
func (f *FFmpeg) SetCmd(cmd FFmpegPlayerCmd) {
	C.lv_ffmpeg_player_set_cmd(f.c, C.lv_ffmpeg_player_cmd_t(cmd))
}

// SetAutoRestart sets whether the video automatically replays from the
// start once it finishes.
func (f *FFmpeg) SetAutoRestart(enable bool) {
	C.lv_ffmpeg_player_set_auto_restart(f.c, C.bool(enable))
}

// SetDecoder selects the FFmpeg decoder to use by name (e.g. "h264",
// "hevc"), overriding FFmpeg's automatic decoder selection.
func (f *FFmpeg) SetDecoder(decoderName string) {
	cName := C.CString(decoderName)
	defer C.free(unsafe.Pointer(cName))
	C.lv_ffmpeg_player_set_decoder(f.c, cName)
}

// SetInnerAlign sets how the decoded video frame is positioned/scaled
// within the player's own object box (set via SetSize) — e.g.
// ImageAlignContain to scale the video to fit the widget's layout size
// while preserving its aspect ratio (letterboxed), or ImageAlignStretch
// to fill the box exactly, ignoring aspect ratio. The FFmpeg player is
// an lv_image subclass at the C level, so this reuses Image's
// lv_image_align_t API directly.
func (f *FFmpeg) SetInnerAlign(align ImageAlign) {
	C.lv_image_set_inner_align(f.c, C.lv_image_align_t(align))
}

// InnerAlign returns the current inner alignment mode.
func (f *FFmpeg) InnerAlign() ImageAlign { return ImageAlign(C.lv_image_get_inner_align(f.c)) }
