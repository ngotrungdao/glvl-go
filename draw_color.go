package lvgl

/*
#include <lvgl.h>
*/
import "C"

// Color mirrors lv_color_t: a plain RGB888 triple, no alpha channel.
type Color struct {
	Blue, Green, Red uint8
}

func (c Color) toC() C.lv_color_t {
	return C.lv_color_t{blue: C.uint8_t(c.Blue), green: C.uint8_t(c.Green), red: C.uint8_t(c.Red)}
}

func colorFromC(c C.lv_color_t) Color {
	return Color{Blue: uint8(c.blue), Green: uint8(c.green), Red: uint8(c.red)}
}

// ColorHex builds a Color from a 0xRRGGBB value, e.g. ColorHex(0xff8800).
func ColorHex(rgb uint32) Color {
	return colorFromC(C.lv_color_hex(C.uint32_t(rgb)))
}

// ColorMake builds a Color from separate red, green, blue components.
func ColorMake(r, g, b uint8) Color {
	return colorFromC(C.lv_color_make(C.uint8_t(r), C.uint8_t(g), C.uint8_t(b)))
}

// ColorWhite returns pure white.
func ColorWhite() Color { return colorFromC(C.lv_color_white()) }

// ColorBlack returns pure black.
func ColorBlack() Color { return colorFromC(C.lv_color_black()) }

// Point mirrors lv_point_t.
type Point struct {
	X, Y int32
}

// Area mirrors lv_area_t: a rectangle given by its two opposite corners.
type Area struct {
	X1, Y1, X2, Y2 int32
}
