// Command gallery opens a Wayland window showing every widget the lvgl
// package wraps, for visual regression checking after changes.
package main

import "lvgl"

func main() {
	lvgl.Init()

	disp := lvgl.WaylandWindowCreate(700, 480, "lvgl-go gallery")
	if disp == nil {
		panic("failed to create Wayland window")
	}
	screen := disp.ScreenActive()
	screen.SetFlexFlow(lvgl.FlexFlowRowWrap)
	screen.SetStylePadAll(10, lvgl.Selector(lvgl.PartMain))

	arc := lvgl.NewArc(screen)
	arc.SetRange(0, 100)
	arc.SetValue(60)

	spinnerAnimTarget := lvgl.NewSpinner(screen)
	spinnerAnimTarget.SetSize(50, 50)

	bar := lvgl.NewBar(screen)
	bar.SetSize(150, 20)
	bar.SetRange(0, 100)

	// Animate the bar's value from 0 to 100 to prove the animation
	// trampoline drives a real widget.
	a := lvgl.NewAnim()
	a.SetValues(0, 100)
	a.SetDuration(2000)
	a.SetRepeatCount(lvgl.AnimRepeatInfinite)
	a.SetExecCB(func(v int32) {
		bar.SetValue(v, false)
	})
	a.Start()

	dd := lvgl.NewDropdown(screen)
	dd.SetOptions("Option 1\nOption 2\nOption 3")

	roller := lvgl.NewRoller(screen)
	roller.SetOptions("Red\nGreen\nBlue", lvgl.RollerModeNormal)
	roller.SetVisibleRowCount(3)

	ta := lvgl.NewTextArea(screen)
	ta.SetSize(150, 40)
	ta.SetOneLine(true)
	ta.SetPlaceholderText("Type here")

	table := lvgl.NewTable(screen)
	table.SetRowCount(2)
	table.SetColumnCount(2)
	table.SetCellValue(0, 0, "A1")
	table.SetCellValue(0, 1, "B1")
	table.SetCellValue(1, 0, "A2")
	table.SetCellValue(1, 1, "B2")

	tabs := lvgl.NewTabView(screen)
	tabs.SetSize(200, 150)
	tab1 := tabs.AddTab("Tab 1")
	lvgl.NewLabel(tab1).SetText("Content 1")
	tab2 := tabs.AddTab("Tab 2")
	lvgl.NewLabel(tab2).SetText("Content 2")

	sb := lvgl.NewSpinBox(screen)
	sb.SetRange(0, 999)
	sb.SetValue(42)

	cal := lvgl.NewCalendar(screen)
	cal.SetSize(220, 200)
	cal.SetTodayDate(2026, 8, 11)
	cal.SetMonthShown(2026, 8)

	chart := lvgl.NewChart(screen)
	chart.SetSize(200, 120)
	chart.SetPointCount(10)
	series := chart.AddSeries(lvgl.ColorHex(0x3388ff), lvgl.ChartAxisPrimaryY)
	for _, v := range []int32{10, 40, 25, 60, 35, 70, 50, 80, 45, 90} {
		chart.SetNextValue(series, v)
	}

	list := lvgl.NewList(screen)
	list.SetSize(160, 150)
	list.AddText("Section")
	list.AddButton("", "Item 1")
	list.AddButton("", "Item 2")

	lvgl.Run(disp)
}
