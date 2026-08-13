// Command gallery3 visually checks the third batch: ArcLabel, Canvas,
// Calendar with a header dropdown addon.
package main

import "lvgl"

func main() {
	lvgl.Init()

	disp := lvgl.WaylandWindowCreate(500, 400, "lvgl-go gallery 3")
	if disp == nil {
		panic("failed to create Wayland window")
	}
	screen := disp.ScreenActive()
	screen.SetFlexFlow(lvgl.FlexFlowRowWrap)
	screen.SetStylePadAll(10, lvgl.Selector(lvgl.PartMain))

	al := lvgl.NewArcLabel(screen)
	al.SetSize(150, 150)
	al.SetText("Arc label text going around")
	al.SetAngleStart(0)
	al.SetAngleSize(270)
	al.SetRadius(60)

	cv := lvgl.NewCanvas(screen)
	cv.SetSize(60, 60)
	cv.SetBuffer(60, 60)
	for y := int32(0); y < 60; y++ {
		for x := int32(0); x < 60; x++ {
			cv.SetPixel(x, y, lvgl.ColorHex(uint32(x*4)<<16|uint32(y*4)<<8), lvgl.OpaCover)
		}
	}

	cal := lvgl.NewCalendar(screen)
	cal.SetSize(220, 200)
	cal.AddHeaderDropdown()
	cal.SetTodayDate(2026, 8, 11)
	cal.SetMonthShown(2026, 8)

	lvgl.Run(disp)
}
