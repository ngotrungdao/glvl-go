// Command styletest visually checks a sample of the new style properties:
// gradient, transform/rotation, margin, text decoration.
package main

import "lvgl"

func main() {
	lvgl.Init()
	disp := lvgl.WaylandWindowCreate(400, 300, "lvgl-go style test")
	if disp == nil {
		panic("failed to create Wayland window")
	}
	screen := disp.ScreenActive()

	card := lvgl.NewObj(screen)
	card.SetSize(200, 100)
	card.SetPos(20, 20)
	s := lvgl.NewStyle()
	s.SetBgColor(lvgl.ColorHex(0x3366ff))
	s.SetBgGradColor(lvgl.ColorHex(0xff3366))
	s.SetBgGradDir(lvgl.GradDirHor)
	s.SetRadius(12)
	s.SetMarginAll(5)
	card.AddStyle(s, lvgl.Selector(lvgl.PartMain))

	label := lvgl.NewLabel(screen)
	label.SetText("Underlined text")
	label.SetPos(20, 140)
	label.SetStyleTextDecor(lvgl.TextDecorUnderline, lvgl.Selector(lvgl.PartMain))
	label.SetStyleTextColor(lvgl.ColorWhite(), lvgl.Selector(lvgl.PartMain))

	rotated := lvgl.NewObj(screen)
	rotated.SetSize(80, 40)
	rotated.SetPos(150, 180)
	rs := lvgl.NewStyle()
	rs.SetBgColor(lvgl.ColorHex(0x00cc88))
	rs.SetTransformRotation(300) // 30 degrees, tenths of a degree
	rs.SetTransformPivotX(40)
	rs.SetTransformPivotY(20)
	rotated.AddStyle(rs, lvgl.Selector(lvgl.PartMain))

	lvgl.Run(disp)
}
