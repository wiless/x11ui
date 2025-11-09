package main

import (
	"image"
	"image/color"
	"math/rand"

	"github.com/wiless/x11ui"
)

func main() {
	app := x11ui.NewApplication("Transparent Background Test", 800, 600, true, false)
	if app == nil {
		panic("could not create application")
	}
	defer app.Close()

	win := app.AppWin()

	// Create a random color background for the main window
	img := image.NewRGBA(image.Rect(0, 0, win.Width, win.Height))
	for y := 0; y < win.Height; y++ {
		for x := 0; x < win.Width; x++ {
			img.Set(x, y, color.RGBA{uint8(rand.Intn(256)), uint8(rand.Intn(256)), uint8(rand.Intn(256)), 255})
		}
	}

	win.PaintOnce()
	win.UpdatePlot(img)

	// Create a container with a semi-transparent red background
	container := app.NewContainer(x11ui.LayoutVer, 100, 100, 400, 400)
	container.SetBGcolor(color.RGBA{255, 255, 0, 128}) // 50% transparent yellow
	// x11ui.NewLabel("hello", &container.Window, 0, 0, 100, 100)
	// _ = container
	// Add a button to the container

	win.Show()
	app.Show()
}
