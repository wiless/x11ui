// sample code to test progressbar
package main

import (
	"fmt"

	"github.com/wiless/x11ui"
)

func main() {
	fmt.Printf("Hello")
	app := x11ui.NewApplication("Hello World", 800, 600, false, false)
	s := NewSlider("Simulation", app.AppWin(), 10, 10, 620, 50)
	s.SetMaxValue(60)

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
