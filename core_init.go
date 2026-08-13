package lvgl

/*
#include <lvgl.h>
*/
import "C"

// Init initializes the LVGL library. It must be called once before any
// other lvgl function.
func Init() {
	C.lv_init()
}

// Deinit deinitializes the LVGL library, freeing all resources it holds.
func Deinit() {
	C.lv_deinit()
}

// IsInitialized reports whether Init has been called.
func IsInitialized() bool {
	return bool(C.lv_is_initialized())
}
