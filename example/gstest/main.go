// Command gstest tries to play a real video file through the GStreamer
// widget in a visible window, to check whether Play() works for a real
// file source (as opposed to videotestsrc, which was found to hang).
package main

import (
	"fmt"
	"os"
	"time"

	"lvgl"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: gstest <video-file-path>")
		os.Exit(1)
	}
	path := os.Args[1]

	lvgl.Init()
	disp := lvgl.WaylandWindowCreate(480, 360, "lvgl-go gstreamer test")
	if disp == nil {
		panic("failed to create Wayland window")
	}
	screen := disp.ScreenActive()

	label := lvgl.NewLabel(screen)
	label.SetText("loading...")
	label.SetPos(10, 10)

	gs := lvgl.NewGStreamer(screen)
	gs.SetSize(400, 300)
	gs.SetPos(10, 40)
	gs.SetInnerAlign(lvgl.ImageAlignContain)

	fmt.Println("calling SetSrc...")
	if err := gs.SetSrc(lvgl.GStreamerFactoryFile, lvgl.GStreamerPropertyFile, path); err != nil {
		fmt.Println("SetSrc failed:", err)
		os.Exit(1)
	}
	fmt.Println("SetSrc ok, calling Play...")

	done := make(chan struct{})
	go func() {
		gs.Play()
		close(done)
	}()

	select {
	case <-done:
		fmt.Println("Play() returned normally")
	case <-time.After(5 * time.Second):
		fmt.Println("Play() did NOT return within 5s (hung)")
		os.Exit(2)
	}

	label.SetText("playing")

	deadline := time.Now().Add(20 * time.Second)
	last := time.Now()
	lastReport := time.Now()
	for time.Now().Before(deadline) {
		now := time.Now()
		lvgl.TickInc(uint32(now.Sub(last).Milliseconds()))
		last = now
		lvgl.TimerHandler()
		if time.Since(lastReport) >= time.Second {
			fmt.Printf("state=%d position=%d duration=%d\n", gs.State(), gs.Position(), gs.Duration())
			lastReport = time.Now()
		}
		time.Sleep(16 * time.Millisecond)
	}

	fmt.Println("done")
}
