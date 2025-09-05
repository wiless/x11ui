// sample code to test progressbar
package main

import (
	"fmt"
	"image/color"

	"github.com/wiless/x11ui"
)

func main() {
	fmt.Printf("Hello")
	x11ui.RegisterFont("./fonts/Pacifico-Regular.ttf", "Pacifico-Regular")
	x11ui.RegisterFont("./fonts/fawesome.ttf", "Font-Awesome")
	// x11ui.RegisterFont("./fonts/fontawesome-solid.ttf", "Font-Awesome")

	// x11ui.SetResourcePath("./fonts/", "Pacifico-Regular.ttf")
	app := x11ui.NewApplication("\uf095 Hello World", 800, 600, false, false)
	s := NewSlider("Simulation", app.AppWin(), 10, 10, 620, 50)

	s.SetBackGroundColor(color.RGBA{0x10, 0x30, 0x10, 0xff})
	s.SetBarColor(color.RGBA{0x10, 0x40, 0x10, 0x10})
	s.SetTextColor(color.White)
	s.SetMaxValue(100)
	s.SetStepSize(5)
	s.SetValue(0)
	s.SetFmtString("%2.0f mins")
	app.SetDefaultKeys()

	// x11ui.DrawDummy(w, x11ui.StateNormal)
	// r := color.RGBA{125, 0, 0, 250}

	// go func() {

	// 	t := time.NewTicker(1 * time.Second)
	// 	val := 0.0
	// 	for range t.C {
	// 		p.SetValue(val)
	// 		p.SetBarColor(r)
	// 		if val == .5 {
	// 			// var x colorful
	// 			// r.G += 20
	// 			p.SetBarColor(r)

	// 			// wd.MoveResize(30, 30, 200, 300)
	// 			// color.RGBA{100, 200, 100, 100}
	// 		}
	// 		val += .1

	// 		if val > 1 {
	// 			t.Stop()
	// 		}
	// 	}
	// }()
	app.Show()
}
