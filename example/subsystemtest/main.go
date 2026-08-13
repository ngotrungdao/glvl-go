// Command subsystemtest visually checks Observer/Subject data binding
// (a slider driving a label via a shared Subject, with no direct
// event-callback wiring between them) and a scrollable container.
package main

import "lvgl"

func main() {
	lvgl.Init()
	disp := lvgl.WaylandWindowCreate(400, 300, "lvgl-go subsystem test")
	if disp == nil {
		panic("failed to create Wayland window")
	}
	screen := disp.ScreenActive()

	label := lvgl.NewLabel(screen)
	label.SetPos(20, 20)

	subject := lvgl.NewSubjectInt(30)
	label.BindText(subject, "Value: %d")

	slider := lvgl.NewSlider(screen)
	slider.SetPos(20, 60)
	slider.SetWidth(200)
	slider.SetRange(0, 100)
	slider.BindValue(subject)
	slider.SetValue(30, false)

	scrollBox := lvgl.NewObj(screen)
	scrollBox.SetPos(20, 110)
	scrollBox.SetSize(150, 100)
	scrollBox.SetScrollDir(lvgl.DirVer)
	for i := 1; i <= 6; i++ {
		row := lvgl.NewButton(scrollBox)
		row.SetSize(lvgl.Pct(90), 30)
		l := lvgl.NewLabel(row.Obj)
		l.SetText("Row")
	}
	scrollBox.ScrollTo(0, 60, false)

	lvgl.Run(disp)
}
