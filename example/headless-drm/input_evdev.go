//go:build evdev

package main

import (
	"fmt"

	"lvgl"
)

// setupInput auto-attaches an Indev for every evdev device already
// present in /dev/input/, and any that appear later, instead of
// hardcoding device paths -- so this runs on whatever happens to be
// plugged in (USB keyboard, USB mouse, touchscreen, ...). A deployment
// with fixed, known hardware may prefer lvgl.EvdevCreate("/dev/input/eventN")
// directly instead.
//
// focusable is the widget keypad/encoder-type devices should be able to
// reach via a Group, since there's no pointer-driven click to fall back
// on if the only input is a keyboard.
func setupInput(focusable *lvgl.Obj) {
	group := lvgl.NewGroup()
	group.AddObj(focusable)

	err := lvgl.EvdevDiscoveryStart(func(indev *lvgl.Indev, kind lvgl.EvdevType) {
		fmt.Println("evdev device discovered, kind:", kind)
		if kind == lvgl.EvdevTypeKey {
			indev.SetGroup(group)
		}
	})
	if err != nil {
		fmt.Println("evdev discovery failed to start:", err)
	}
}
