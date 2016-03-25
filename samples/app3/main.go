// sample code to test progressbar
package main

import (
	"fmt"
	"image/color"
	"time"

	"github.com/wiless/x11ui"
)

func main() {
	fmt.Printf("Hello")
	app := x11ui.NewApplication("Hello World", 800, 600, false, false)
	p := x11ui.NewProgressBar("Simulation", app.AppWin(), 10, 10, 620, 30)
	r := color.RGBA{125, 0, 0, 0}

	go func() {
		time.Sleep(2 * time.Second)
		t := time.NewTicker(1 * time.Second)
		val := 0.0
		for range t.C {
			p.SetValue(val)

			if val >= .5 {
				// var x colorful
				// r.G += 20
				p.SetBarColor(r)
				// color.RGBA{100, 200, 100, 100}
			}
			val += .1

			if val > 1 {
				t.Stop()
			}
		}
	}()
	app.Show()
}
