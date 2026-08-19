// Command gallery2 visually checks the second batch of widgets: Win,
// Keyboard+TextArea, ButtonMatrix, Led, Line, TileView, Scale, SpanGroup,
// QRCode, Barcode.
package main

import "lvgl"

func main() {
	lvgl.Init()

	disp := lvgl.SDLWindowCreate(760, 520, "lvgl-go gallery 2")
	if disp == nil {
		panic("failed to create SDL window")
	}
	screen := disp.ScreenActive()
	screen.SetFlexFlow(lvgl.FlexFlowRowWrap)
	screen.SetStylePadAll(10, lvgl.Selector(lvgl.PartMain))

	win := lvgl.NewWin(screen)
	win.SetSize(200, 150)
	win.AddTitle("My Window")
	win.AddButton(lvgl.SymbolClose, 40)
	winLabel := lvgl.NewLabel(win.Content())
	winLabel.SetText("Window content")

	led := lvgl.NewLed(screen)
	led.SetSize(30, 30)
	led.SetColor(lvgl.ColorHex(0x00ff00))
	led.SetBrightness(255)

	line := lvgl.NewLine(screen)
	line.SetSize(100, 60)
	line.SetPoints([]lvgl.Point{{X: 0, Y: 0}, {X: 50, Y: 60}, {X: 100, Y: 0}})
	line.SetStyleLineWidth(3, lvgl.Selector(lvgl.PartMain))
	line.SetStyleLineColor(lvgl.ColorHex(0xff8800), lvgl.Selector(lvgl.PartMain))

	scale := lvgl.NewScale(screen)
	scale.SetSize(150, 60)
	scale.SetRange(0, 100)
	scale.SetTotalTickCount(11)
	scale.SetMajorTickEvery(5)

	bm := lvgl.NewButtonMatrix(screen)
	bm.SetSize(180, 100)
	bm.SetMap([][]string{{"1", "2", "3"}, {"4", "5", "6"}, {"0"}})

	sg := lvgl.NewSpanGroup(screen)
	sg.SetSize(180, 60)
	span1 := sg.AddSpan()
	span1.SetText("Rich ")
	span2 := sg.AddSpan()
	span2.SetText("text!")
	sg.Refresh()

	qr := lvgl.NewQRCode(screen)
	qr.SetSize(100)
	qr.SetDarkColor(lvgl.ColorBlack())
	qr.SetLightColor(lvgl.ColorWhite())
	qr.SetData("https://lvgl.io")

	bc := lvgl.NewBarcode(screen)
	bc.SetSize(180, 60)
	bc.SetScale(2)
	bc.Update("HELLO")

	ta := lvgl.NewTextArea(screen)
	ta.SetSize(200, 40)
	ta.SetOneLine(true)
	kb := lvgl.NewKeyboard(screen)
	kb.SetSize(300, 150)
	kb.SetTextArea(ta)

	lvgl.Run(disp)
}
