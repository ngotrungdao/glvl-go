package lvgl

import "time"

// Run drives LVGL's tick and timer loop until the display is closed (e.g.
// via the window's close button, or a call to Display.Close). It measures
// real elapsed time between iterations to feed TickInc accurately, and
// sleeps for the duration TimerHandler reports before the next timer is
// due, rather than a fixed interval.
func Run(d *Display) {
	last := time.Now()
	for d.IsOpen() {
		now := time.Now()
		TickInc(uint32(now.Sub(last).Milliseconds()))
		last = now

		sleepMs := TimerHandler()
		if sleepMs == 0 {
			sleepMs = 1
		}
		time.Sleep(time.Duration(sleepMs) * time.Millisecond)
	}
}
