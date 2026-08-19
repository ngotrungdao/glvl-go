// Command fonttest visually checks font family/size/color control.
package main

import "lvgl"

func main() {
	lvgl.Init()
	disp := lvgl.SDLWindowCreate(400, 300, "lvgl-go font test")
	if disp == nil {
		panic("failed to create SDL window")
	}
	screen := disp.ScreenActive()
	screen.SetFlexFlow(lvgl.FlexFlowColumn)
	screen.SetStylePadAll(10, lvgl.Selector(lvgl.PartMain))

	sizes := []*lvgl.Font{lvgl.FontMontserrat12, lvgl.FontMontserrat20, lvgl.FontMontserrat32}
	colors := []lvgl.Color{lvgl.ColorHex(0xff4444), lvgl.ColorHex(0x44ff44), lvgl.ColorHex(0x4488ff)}

	for i, f := range sizes {
		l := lvgl.NewLabel(screen)
		l.SetText("Size and color")
		l.SetStyleTextFont(f, lvgl.Selector(lvgl.PartMain))
		l.SetStyleTextColor(colors[i], lvgl.Selector(lvgl.PartMain))
	}

	lvgl.Run(disp)
}
