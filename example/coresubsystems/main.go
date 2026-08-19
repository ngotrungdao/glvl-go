// Command coresubsystems is a manual regression check for Timer, eased
// Anim, Group/Indev, and Observer/Subject data binding — split out from
// example/headless because those checks were found to hang unreliably
// once run against example/headless's heavily accumulated (80+ widget)
// scene, but are completely reliable against the small, fresh scene this
// program builds (see README's Known limitations for the underlying
// redraw-stall issue this sidesteps rather than fixes).
package main

import (
	"fmt"
	"os"
	"runtime"
	"time"

	"lvgl"
)

// waitUntil polls cond (calling TimerHandler each time) until it's true
// or the timeout elapses, sleeping briefly between polls instead of
// busy-spinning. Deliberately synchronous (not goroutine-isolated): LVGL
// isn't thread-safe, so a background goroutine still calling
// TimerHandler() after an abandoned timeout would race the main
// goroutine's later calls — worse than a hang. See
// example/headless/main.go's waitUntil for the same reasoning in more
// detail.
func waitUntil(timeout time.Duration, cond func() bool) bool {
	deadline := time.Now().Add(timeout)
	for !cond() {
		if time.Now().After(deadline) {
			return false
		}
		lvgl.TimerHandler()
		time.Sleep(2 * time.Millisecond)
	}
	return true
}

// pumpFor calls TimerHandler, paced with real sleeps, for the given
// duration (unlike waitUntil, there's no early-exit condition — used to
// prove nothing happens over a window of real time).
func pumpFor(d time.Duration) {
	deadline := time.Now().Add(d)
	for time.Now().Before(deadline) {
		lvgl.TimerHandler()
		time.Sleep(2 * time.Millisecond)
	}
}

func check(name string, cond bool) {
	if !cond {
		fmt.Printf("FAIL: %s\n", name)
		os.Exit(1)
	}
	fmt.Printf("ok:   %s\n", name)
}

func main() {
	lvgl.Init()

	disp := lvgl.SDLWindowCreate(1, 1, "lvgl-go coresubsystems check")
	if disp == nil {
		fmt.Println("FAIL: could not create an SDL window (is a display available?)")
		os.Exit(1)
	}
	screen := disp.ScreenActive()

	// --- custom timer ---
	// Auto-delete path: leave SetAutoDelete on its default (true) and let
	// LVGL delete the timer itself once the repeat count is exhausted;
	// don't also call Delete() manually (that would double-free it).
	autoFires := 0
	autoTimer := lvgl.NewTimer(5, func(t *lvgl.Timer) {
		autoFires++
	})
	autoTimer.SetRepeatCount(3)
	waitUntil(300*time.Millisecond, func() bool { return autoFires >= 3 })
	check("auto-delete timer fires its full repeat count", autoFires == 3)

	// Manual-delete path: disable auto-delete, call Delete() ourselves
	// once satisfied, matching the doc comment's recommended pattern.
	manualFires := 0
	manualTimer := lvgl.NewTimer(5, func(t *lvgl.Timer) {
		manualFires++
	})
	manualTimer.SetAutoDelete(false)
	waitUntil(100*time.Millisecond, func() bool { return manualFires > 0 })
	manualTimer.Delete()
	firesAtDelete := manualFires
	pumpFor(100 * time.Millisecond)
	check("manually deleted timer stops firing", manualFires == firesAtDelete && firesAtDelete > 0)

	// --- eased animation ---
	easedAnim := lvgl.NewAnim()
	easedAnim.SetValues(0, 100)
	easedAnim.SetDuration(10)
	easedAnim.SetPath(lvgl.AnimPathEaseInOut)
	easedCompleted := false
	easedAnim.SetCompletedCB(func() { easedCompleted = true })
	easedAnim.Start()
	waitUntil(300*time.Millisecond, func() bool { return easedCompleted })
	check("eased animation completes", easedCompleted)

	// --- group + indev ---
	g := lvgl.NewGroup()
	gObj1 := lvgl.NewObj(screen)
	gObj2 := lvgl.NewObj(screen)
	g.AddObj(gObj1)
	g.AddObj(gObj2)
	check("group obj count round-trips", g.ObjCount() == 2)

	lvgl.FocusObj(gObj1)
	check("group focus round-trips", g.Focused().Same(gObj1))

	g.FocusNext()
	check("group focus next round-trips", g.Focused().Same(gObj2))

	kbIndev := disp.Keyboard()
	if kbIndev != nil {
		kbIndev.SetGroup(g)
		check("indev group round-trips", kbIndev.Group() != nil)
		check("indev type is keypad", kbIndev.Type() == lvgl.IndevTypeKeypad)
	} else {
		check("keyboard indev accessible", false)
	}

	ptrIndev := disp.Pointer()
	check("pointer indev accessible and is pointer type", ptrIndev != nil && ptrIndev.Type() == lvgl.IndevTypePointer)

	g.Delete()

	// --- observer/subject data binding ---
	intSubject := lvgl.NewSubjectInt(10)
	check("subject int round-trips", intSubject.Int() == 10)

	observerCalls := 0
	var lastObserved int32
	obs := intSubject.AddObserver(func(s *lvgl.Subject) {
		observerCalls++
		lastObserved = s.Int()
	})
	intSubject.SetInt(20)
	check("observer fires on subject change", observerCalls >= 1 && lastObserved == 20)
	obs.Remove()

	intSubject.SetInt(30)
	callsAfterRemove := observerCalls
	check("removed observer stops firing", callsAfterRemove == observerCalls)

	boundSlider := lvgl.NewSlider(screen)
	boundSlider.SetRange(0, 100)
	sliderSubject := lvgl.NewSubjectInt(0)
	sliderObs := boundSlider.BindValue(sliderSubject)
	sliderSubject.SetInt(55)
	runtime.GC()
	check("slider BindValue round-trips via subject", boundSlider.Value() == 55)
	sliderObs.Remove()

	strSubject := lvgl.NewSubjectString("hello", 32)
	boundLabel := lvgl.NewLabel(screen)
	labelObs := boundLabel.BindText(strSubject, "")
	strSubject.SetString("world")
	check("label BindText round-trips via subject", boundLabel.Text() == "world")
	labelObs.Remove()

	checkBoxObj := lvgl.NewCheckbox(screen)
	checkedSubject := lvgl.NewSubjectInt(0)
	checkedObs := checkBoxObj.BindChecked(checkedSubject)
	checkedSubject.SetInt(1)
	check("BindChecked round-trips via subject", checkBoxObj.HasState(lvgl.StateChecked))
	checkedObs.Remove()

	fmt.Println("all checks passed")
}
