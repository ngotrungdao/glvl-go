package lvgl

/*
#include <stdlib.h>
#include <lvgl.h>
*/
import "C"
import "unsafe"

// TranslationInit initializes LVGL's translation subsystem. Call once
// before SetLanguage/Translate.
//
// Registering translation packs themselves (lv_translation_add_language,
// building the underlying lv_translation_pack_t/tag tables) isn't
// wrapped — those are normally built as static data in C and would need
// their own dedicated supporting API this pass didn't build. This file
// covers the read/use side: switching language and looking up tags in a
// pack registered elsewhere (e.g. from C, or a future addition).
func TranslationInit() { C.lv_translation_init() }

// TranslationDeinit frees the translation subsystem's resources.
func TranslationDeinit() { C.lv_translation_deinit() }

// SetLanguage sets the active language (an IETF-style tag like "en" or
// "de", matching however the registered pack's languages were named).
func SetLanguage(lang string) {
	cLang := C.CString(lang)
	defer C.free(unsafe.Pointer(cLang))
	C.lv_translation_set_language(cLang)
}

// Language returns the current active language tag.
func Language() string { return C.GoString(C.lv_translation_get_language()) }

// Translate looks up tag in the active language, returning tag itself
// unchanged if no translation is found.
func Translate(tag string) string {
	cTag := C.CString(tag)
	defer C.free(unsafe.Pointer(cTag))
	return C.GoString(C.lv_translation_get(cTag))
}
