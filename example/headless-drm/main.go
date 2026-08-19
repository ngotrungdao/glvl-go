// Command headless-drm runs lvgl-go through LVGL's native Linux DRM/KMS
// driver instead of SDL2 — for a genuinely headless target with no
// X11/Wayland/desktop at all (e.g. a Raspberry Pi driving its own display
// directly). Unlike SDLWindowCreate, DRMDisplay takes over the whole
// physical display: there's no window, no title bar, no decorations, and
// no window manager to draw them even if there were one — DRM/KMS renders
// straight into the display's own framebuffer.
//
// Needs lvgl-c built with LV_USE_LINUX_DRM=1 (this repo's build has it).
// For real input, also build lvgl-c with LV_USE_EVDEV=1 and run
// `go run -tags evdev ./example/headless-drm` — without that tag this
// still builds and runs, just with no input wired up (see input_stub.go).
//
// Not yet run against real DRM hardware from this pass: this dev
// machine's display is already owned by a running desktop compositor,
// and DRM needs exclusive access to the device (SET_MASTER), so running
// this here would either fail outright or fight the live session for the
// screen. Verified to compile and link against this repo's current
// (x86-64) liblvgl.a — see the README's "Headless embedded targets"
// section for the full verification story.
package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"lvgl"
)

func main() {
	// lv_linux_drm_find_device_path just returns *a* DRM device file, with
	// no guarantee it has a connected output (its own doc comment doesn't
	// claim otherwise) -- on a machine with more than one GPU (e.g. an
	// Intel iGPU plus a discrete card), it can easily pick the one with
	// nothing plugged in. Confirmed on real hardware: check
	// /sys/class/drm/card*-*/status for "connected" to find the right one
	// if auto-detection picks wrong, and pass it explicitly here.
	devicePath := flag.String("device", "", `DRM device file, e.g. "/dev/dri/card1" (default: auto-detect via lv_linux_drm_find_device_path)`)
	flag.Parse()

	lvgl.Init()

	disp := lvgl.DRMDisplayCreate()
	if disp == nil {
		fmt.Println("FAIL: could not create a DRM display")
		os.Exit(1)
	}

	path := *devicePath
	if path == "" {
		var err error
		path, err = lvgl.FindDRMDevicePath()
		if err != nil {
			fmt.Println("FindDRMDevicePath failed, falling back to /dev/dri/card0:", err)
			path = "/dev/dri/card0"
		}
	}
	if err := disp.SetFile(path, -1); err != nil { // -1: auto-select the connector
		fmt.Println("FAIL: SetFile:", path, err)
		fmt.Println(`try -device explicitly, e.g. -device=/dev/dri/card1 -- check`,
			`/sys/class/drm/card*-*/status for which one says "connected"`)
		os.Exit(1)
	}
	fmt.Println("DRM display attached:", path)

	screen := disp.ScreenActive()

	label := lvgl.NewLabel(screen)
	label.SetText("lvgl-go on DRM/KMS -- no window, no desktop, no compositor")
	label.SetWidth(lvgl.Pct(80))
	label.SetLongMode(lvgl.LabelLongModeWrap)
	label.Center()

	btn := lvgl.NewButton(screen)
	btn.SetSize(160, 50)
	btn.AlignObj(lvgl.AlignBottomMid, 0, -20)
	btnLabel := lvgl.NewLabel(btn.Obj)
	btnLabel.SetText("Press me")
	btnLabel.Center()

	clicks := 0
	btn.AddEventCB(lvgl.EventClicked, func(e *lvgl.Event) {
		clicks++
		btnLabel.SetText(fmt.Sprintf("Pressed %d", clicks))
	})

	setupInput(btn.Obj)

	// DRM has no "window close" concept -- no compositor, no [x] button --
	// so this runs until killed. Listen for Ctrl+C / SIGTERM so a normal
	// stop doesn't just get silently swallowed.
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	last := time.Now()
	for {
		select {
		case <-stop:
			fmt.Println("shutting down")
			return
		default:
		}

		now := time.Now()
		// A no-op once the DRM driver's own tick source is live (it calls
		// lv_tick_set_cb internally, same as the SDL driver) -- kept for
		// backend-agnosticism, same reasoning as Run in app.go.
		lvgl.TickInc(uint32(now.Sub(last).Milliseconds()))
		last = now

		sleepMs := lvgl.TimerHandler()
		if sleepMs == 0 {
			sleepMs = 1
		}
		time.Sleep(time.Duration(sleepMs) * time.Millisecond)
	}
}
