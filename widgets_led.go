package lvgl

/*
#include <lvgl.h>
*/
import "C"

// Led wraps an lv_led widget: a colored circle with adjustable brightness.
type Led struct{ *Obj }

// NewLed creates a LED as a child of parent.
func NewLed(parent *Obj) *Led {
	return &Led{wrapObj(C.lv_led_create(parent.c))}
}

// SetColor sets the LED's color.
func (l *Led) SetColor(c Color) { C.lv_led_set_color(l.c, c.toC()) }

// SetBrightness sets the LED's brightness, 0 (off) to 255 (max).
func (l *Led) SetBrightness(b uint8) { C.lv_led_set_brightness(l.c, C.uint8_t(b)) }

// On sets the LED to full brightness.
func (l *Led) On() { C.lv_led_on(l.c) }

// Off turns the LED off.
func (l *Led) Off() { C.lv_led_off(l.c) }

// Toggle switches the LED between on and off.
func (l *Led) Toggle() { C.lv_led_toggle(l.c) }

// Color returns the LED's current color.
func (l *Led) Color() Color { return colorFromC(C.lv_led_get_color(l.c)) }

// Brightness returns the LED's current brightness.
func (l *Led) Brightness() uint8 { return uint8(C.lv_led_get_brightness(l.c)) }
