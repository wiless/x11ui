package main

import (
	"github.com/wiless/x11ui"
)

func main() {
	app := x11ui.NewApplication("Auto Layout Example", 800, 600, true, false)
	if app == nil {
		panic("could not create application")
	}
	defer app.Close()

	win := app.AppWin()

	// Create a vertical container
	vcontainer := app.NewContainer(x11ui.LayoutVer, 10, 10, 380, 580)
	vcontainer.SetSpacing(5, 5)

	// Add buttons to the vertical container
	for i := 0; i < 5; i++ {
		btn := vcontainer.AddButton("Button")
		btn.SetLabel("Vertical " + string(i+'0'))
	}

	// Create a horizontal container
	hcontainer := app.NewContainer(x11ui.LayoutHor, 410, 10, 380, 580)
	hcontainer.SetSpacing(5, 5)

	// Add buttons to the horizontal container
	for i := 0; i < 3; i++ {
		btn := hcontainer.AddButton("Button")
		btn.SetLabel("Horizontal " + string(i+'0'))
	}

	win.Show()
	app.Show()
}
