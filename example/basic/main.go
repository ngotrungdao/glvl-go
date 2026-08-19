// Command basic opens an SDL window exercising the widgets and event
// system built so far: a label, a button that updates it on click, and a
// slider that reports its value.
package main

import (
	"strconv"

	"lvgl"
)

func main() {
	lvgl.Init()

	disp := lvgl.SDLWindowCreate(480, 320, "lvgl-go")
	if disp == nil {
		panic("failed to create SDL window")
	}
	screen := disp.ScreenActive()

	card := lvgl.NewObj(screen)
	card.SetSize(240, 300)
	card.SetPos(220, 10)
	card.SetFlexFlow(lvgl.FlexFlowColumn)
	card.SetFlexAlign(lvgl.FlexAlignStart, lvgl.FlexAlignCenter, lvgl.FlexAlignStart)
	cardStyle := lvgl.NewStyle()
	cardStyle.SetBgColor(lvgl.ColorHex(0x2b2f3a))
	cardStyle.SetRadius(12)
	cardStyle.SetPadAll(12)
	cardStyle.SetPadRow(10)
	cardStyle.SetBorderWidth(0)
	card.AddStyle(cardStyle, lvgl.Selector(lvgl.PartMain))

	cardTitle := lvgl.NewLabel(card)
	cardTitle.SetText("Styled flex card")
	cardTitle.SetStyleTextColor(lvgl.ColorWhite(), lvgl.Selector(lvgl.PartMain))

	for i := 1; i <= 3; i++ {
		row := lvgl.NewObj(card)
		row.SetSize(lvgl.Pct(100), 36)
		rowStyle := lvgl.NewStyle()
		rowStyle.SetBgColor(lvgl.ColorHex(0x3a3f4d))
		rowStyle.SetRadius(6)
		rowStyle.SetBorderWidth(0)
		row.AddStyle(rowStyle, lvgl.Selector(lvgl.PartMain))
		rowLabel := lvgl.NewLabel(row)
		rowLabel.SetStyleTextColor(lvgl.ColorHex(0xdddddd), lvgl.Selector(lvgl.PartMain))
		rowLabel.SetText("Row " + strconv.Itoa(i))
		rowLabel.Center()
	}

	label := lvgl.NewLabel(screen)
	label.SetText("Hello from Go + LVGL!")
	label.SetPos(20, 20)

	clicks := 0
	btn := lvgl.NewButton(screen)
	btn.SetPos(20, 60)
	btnLabel := lvgl.NewLabel(btn.Obj)
	btnLabel.SetText("Click me")
	btn.AddEventCB(lvgl.EventClicked, func(e *lvgl.Event) {
		clicks++
		label.SetText("Clicked " + strconv.Itoa(clicks) + " time(s)")
	})

	checkbox := lvgl.NewCheckbox(screen)
	checkbox.SetText("Enable something")
	checkbox.SetPos(20, 110)

	sw := lvgl.NewSwitch(screen)
	sw.SetPos(20, 150)

	slider := lvgl.NewSlider(screen)
	slider.SetPos(20, 190)
	slider.SetWidth(200)
	slider.SetRange(0, 100)
	slider.SetValue(30, false)
	slider.AddEventCB(lvgl.EventValueChanged, func(e *lvgl.Event) {
		label.SetText("Slider: " + strconv.Itoa(int(slider.Value())))
	})

	lvgl.Run(disp)
}
