//go:build !evdev

package main

import (
	"fmt"

	"lvgl"
)

// setupInput is a no-op here: this build wasn't compiled with -tags evdev
// (or lvgl-c wasn't built with LV_USE_EVDEV=1), so there's no input
// driver to wire up. The display still renders -- useful as a display-only
// smoke test -- but nothing will react to a keyboard/touchscreen/mouse.
// See input_evdev.go for the real implementation.
func setupInput(focusable *lvgl.Obj) {
	fmt.Println("input: built without -tags evdev, no input device wired up (display-only)")
}
