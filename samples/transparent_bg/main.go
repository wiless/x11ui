package main

import (
	"image/color"

	"github.com/wiless/x11ui"
)

func main() {
	app := x11ui.NewApplication("Transparent Background Example", 800, 600, true, false)
	if app == nil {
		panic("could not create application")
	}
	defer app.Close()

	win := app.AppWin()

	// Create a child window with a semi-transparent red background
	childWin := app.NewChildWindow("Child Window", 50, 50, 300, 300)
	childWin.SetBGcolor(color.RGBA{255, 0, 0, 128}) // Semi-transparent red

	// Create a container with a semi-transparent blue background
	container := app.NewContainer(x11ui.LayoutVer, 400, 50, 300, 300)
	container.SetBGcolor(color.RGBA{0, 0, 255, 128}) // Semi-transparent blue

	// Add a button to the container
	btn := container.AddButton("My Button")
	btn.SetLabel("Click Me")
	btn.SetOnClick(func() {
		println("Button clicked!")
	})

	win.Show()
	app.Show()
}
