// Command ffmpegtest tries to play a real video file through the
// FFmpeg widget in a visible window, as a manual regression check for
// the FFmpeg wrapper (widgets_ffmpeg.go).
package main

import (
	"fmt"
	"os"
	"time"

	"lvgl"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: ffmpegtest <video-file-path>")
		os.Exit(1)
	}
	path := os.Args[1]

	lvgl.Init()
	lvgl.FFmpegInit()

	n, err := lvgl.FFmpegGetFrameNum(path)
	if err != nil {
		fmt.Println("FFmpegGetFrameNum failed:", err)
		os.Exit(1)
	}
	fmt.Println("frame count:", n)

	disp := lvgl.SDLWindowCreate(480, 360, "lvgl-go ffmpeg test")
	if disp == nil {
		panic("failed to create SDL window")
	}
	screen := disp.ScreenActive()

	label := lvgl.NewLabel(screen)
	label.SetText("loading...")
	label.SetPos(10, 10)

	fp := lvgl.NewFFmpeg(screen)
	fp.SetSize(400, 300)
	fp.SetPos(10, 40)
	fp.SetInnerAlign(lvgl.ImageAlignContain) // scale video to fit the widget's box, keeping aspect ratio

	fmt.Println("calling SetSrc...")
	if err := fp.SetSrc(path); err != nil {
		fmt.Println("SetSrc failed:", err)
		os.Exit(1)
	}
	fmt.Println("SetSrc ok, calling SetCmd(Start)...")

	// Play once through to the end rather than looping forever, so the
	// run actually completes.
	fp.SetAutoRestart(false)
	fp.SetCmd(lvgl.FFmpegCmdStart)
	label.SetText("playing")

	// The wrapper doesn't expose the video's frame rate (lv_ffmpeg.h has
	// no public getter for it), so assume a conservative lower-bound fps
	// and pad generously — this only needs to run at least as long as
	// the real video, not exactly as long.
	const assumedMinFPS = 15.0
	const safetyMargin = 1.5
	runFor := time.Duration(float64(n)/assumedMinFPS*safetyMargin*float64(time.Second)) + 2*time.Second
	fmt.Printf("playing for up to %s (frame count %d)...\n", runFor, n)

	// Report actual elapsed time to LVGL's tick, not a hardcoded 16ms.
	// time.Sleep is only a minimum, not a guarantee — under system load
	// the OS can wake this goroutine tens of milliseconds late, and if
	// we still told LVGL "16ms passed" every time, its internal clock
	// (and therefore the video's playback timer) would drift further
	// and further behind real time. Measuring and reporting the true
	// gap keeps LVGL's clock self-correcting instead of compounding
	// scheduler jitter into visible drift.
	deadline := time.Now().Add(runFor)
	last := time.Now()
	for time.Now().Before(deadline) {
		now := time.Now()
		lvgl.TickInc(uint32(now.Sub(last).Milliseconds()))
		last = now
		lvgl.TimerHandler()
		time.Sleep(16 * time.Millisecond)
	}

	fmt.Println("done")
}
