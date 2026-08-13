package lvgl

/*
#include <stdlib.h>
#include <lvgl.h>
*/
import "C"
import (
	"errors"
	"runtime"
	"unsafe"
)

// Font wraps an lv_font_t, used with Style.SetTextFont /
// Obj.SetStyleTextFont to control both font family and size (LVGL has no
// separate "set font size" call: size is a property of which Font you
// use). Built-in fonts (FontMontserrat*) are static LVGL data and need no
// cleanup; fonts created with FreeTypeFontCreate or TinyTTFCreateFile/Data
// own C-heap resources and must be freed with Delete once no longer used
// by any style/object, matching Style's memory-ownership rules.
type Font struct {
	c    *C.lv_font_t
	kind fontKind
	data unsafe.Pointer // pinned raw font bytes for fontKindTinyTTF fonts created via TinyTTFCreateData; nil otherwise
}

type fontKind uint8

const (
	fontKindStatic fontKind = iota
	fontKindFreeType
	fontKindTinyTTF
)

func staticFont(c *C.lv_font_t) *Font { return &Font{c: c, kind: fontKindStatic} }

// Built-in Montserrat bitmap fonts bundled with LVGL, at every size
// enabled in this build's lv_conf (8-48px). These need no Delete call.
var (
	FontMontserrat8  = staticFont(&C.lv_font_montserrat_8)
	FontMontserrat10 = staticFont(&C.lv_font_montserrat_10)
	FontMontserrat12 = staticFont(&C.lv_font_montserrat_12)
	FontMontserrat14 = staticFont(&C.lv_font_montserrat_14)
	FontMontserrat16 = staticFont(&C.lv_font_montserrat_16)
	FontMontserrat18 = staticFont(&C.lv_font_montserrat_18)
	FontMontserrat20 = staticFont(&C.lv_font_montserrat_20)
	FontMontserrat22 = staticFont(&C.lv_font_montserrat_22)
	FontMontserrat24 = staticFont(&C.lv_font_montserrat_24)
	FontMontserrat26 = staticFont(&C.lv_font_montserrat_26)
	FontMontserrat28 = staticFont(&C.lv_font_montserrat_28)
	FontMontserrat30 = staticFont(&C.lv_font_montserrat_30)
	FontMontserrat32 = staticFont(&C.lv_font_montserrat_32)
	FontMontserrat34 = staticFont(&C.lv_font_montserrat_34)
	FontMontserrat36 = staticFont(&C.lv_font_montserrat_36)
	FontMontserrat38 = staticFont(&C.lv_font_montserrat_38)
	FontMontserrat40 = staticFont(&C.lv_font_montserrat_40)
	FontMontserrat42 = staticFont(&C.lv_font_montserrat_42)
	FontMontserrat44 = staticFont(&C.lv_font_montserrat_44)
	FontMontserrat46 = staticFont(&C.lv_font_montserrat_46)
	FontMontserrat48 = staticFont(&C.lv_font_montserrat_48)
)

// Delete frees a dynamically loaded font's C resources (FreeType/TinyTTF
// fonts only; a no-op for built-in static fonts like FontMontserrat16).
// Only call this once the font is no longer set on any style/object.
func (f *Font) Delete() {
	if f.c == nil {
		return
	}
	switch f.kind {
	case fontKindFreeType:
		C.lv_freetype_font_delete(f.c)
	case fontKindTinyTTF:
		C.lv_tiny_ttf_destroy(f.c)
	default:
		return
	}
	f.c = nil
	if f.data != nil {
		C.free(f.data)
		f.data = nil
	}
	runtime.SetFinalizer(f, nil)
}

func (f *Font) leakWarn() {
	println("lvgl: Font finalized without an explicit Delete() call; its C memory was leaked")
}

// FreeTypeRenderMode mirrors lv_freetype_font_render_mode_t.
type FreeTypeRenderMode uint32

var (
	FreeTypeRenderModeBitmap  = FreeTypeRenderMode(C.LV_FREETYPE_FONT_RENDER_MODE_BITMAP)
	FreeTypeRenderModeOutline = FreeTypeRenderMode(C.LV_FREETYPE_FONT_RENDER_MODE_OUTLINE)
)

// FreeTypeStyle mirrors lv_freetype_font_style_t; values can be OR'd
// (e.g. FreeTypeStyleBold|FreeTypeStyleItalic).
type FreeTypeStyle uint32

var (
	FreeTypeStyleNormal = FreeTypeStyle(C.LV_FREETYPE_FONT_STYLE_NORMAL)
	FreeTypeStyleItalic = FreeTypeStyle(C.LV_FREETYPE_FONT_STYLE_ITALIC)
	FreeTypeStyleBold   = FreeTypeStyle(C.LV_FREETYPE_FONT_STYLE_BOLD)
)

var (
	errFreeTypeInit       = errors.New("lvgl: freetype init failed")
	errFreeTypeFontCreate = errors.New("lvgl: freetype font creation failed (bad path or size?)")
	errTinyTTFFontCreate  = errors.New("lvgl: tiny_ttf font creation failed (bad path/data or size?)")
)

// FreeTypeInit initializes LVGL's FreeType integration with room for up
// to maxGlyphCnt cached glyphs. Call once before FreeTypeFontCreate.
//
// LVGL may already have FreeType initialized internally (observed in
// practice, likely for its own image-decoding paths) before this is ever
// called, in which case lv_freetype_init logs "freetype already
// initialized" and returns a non-OK result indistinguishable from a real
// failure — FreeTypeFontCreate still works fine in that case. If this
// returns an error, it's safe to try FreeTypeFontCreate anyway rather
// than treat it as fatal.
func FreeTypeInit(maxGlyphCnt uint32) error {
	if C.lv_freetype_init(C.uint32_t(maxGlyphCnt)) != C.LV_RESULT_OK {
		return errFreeTypeInit
	}
	return nil
}

// FreeTypeFontCreate loads a font (e.g. a .ttf/.otf file) at the given
// pixel size via FreeType. FreeTypeInit must be called first. The
// returned Font must be Delete()d once no longer in use.
func FreeTypeFontCreate(path string, renderMode FreeTypeRenderMode, size uint32, style FreeTypeStyle) (*Font, error) {
	cPath := C.CString(path)
	defer C.free(unsafe.Pointer(cPath))

	c := C.lv_freetype_font_create(cPath, C.lv_freetype_font_render_mode_t(renderMode),
		C.uint32_t(size), C.lv_freetype_font_style_t(style))
	if c == nil {
		return nil, errFreeTypeFontCreate
	}
	f := &Font{c: c, kind: fontKindFreeType}
	runtime.SetFinalizer(f, (*Font).leakWarn)
	return f, nil
}

// TinyTTFCreateFile loads a .ttf file at the given pixel size using
// LVGL's lightweight built-in TrueType renderer (an alternative to
// FreeType that needs no FreeTypeInit call). The returned Font must be
// Delete()d once no longer in use.
func TinyTTFCreateFile(path string, fontSize int32) (*Font, error) {
	cPath := C.CString(path)
	defer C.free(unsafe.Pointer(cPath))

	c := C.lv_tiny_ttf_create_file(cPath, C.int32_t(fontSize))
	if c == nil {
		return nil, errTinyTTFFontCreate
	}
	f := &Font{c: c, kind: fontKindTinyTTF}
	runtime.SetFinalizer(f, (*Font).leakWarn)
	return f, nil
}

// TinyTTFCreateData loads TTF font data already in memory at the given
// pixel size. LVGL's docs don't specify whether the data is copied or
// referenced, so (like Image.SetSrcPath) the Go copy is conservatively
// pinned on the C heap for the font's lifetime and freed by Delete,
// rather than assumed safe to free after this call returns.
func TinyTTFCreateData(data []byte, fontSize int32) (*Font, error) {
	var dataPtr unsafe.Pointer
	if len(data) > 0 {
		dataPtr = C.CBytes(data)
	}

	c := C.lv_tiny_ttf_create_data(dataPtr, C.size_t(len(data)), C.int32_t(fontSize))
	if c == nil {
		if dataPtr != nil {
			C.free(dataPtr)
		}
		return nil, errTinyTTFFontCreate
	}
	f := &Font{c: c, kind: fontKindTinyTTF, data: dataPtr}
	runtime.SetFinalizer(f, (*Font).leakWarn)
	return f, nil
}
