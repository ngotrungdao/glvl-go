package lvgl

/*
#include <lvgl.h>
*/
import "C"

// Dir mirrors lv_dir_t; values can be OR'd together (e.g. DirHor).
type Dir uint8

var (
	DirNone   = Dir(C.LV_DIR_NONE)
	DirLeft   = Dir(C.LV_DIR_LEFT)
	DirRight  = Dir(C.LV_DIR_RIGHT)
	DirTop    = Dir(C.LV_DIR_TOP)
	DirBottom = Dir(C.LV_DIR_BOTTOM)
	DirHor    = Dir(C.LV_DIR_HOR)
	DirVer    = Dir(C.LV_DIR_VER)
	DirAll    = Dir(C.LV_DIR_ALL)
)
