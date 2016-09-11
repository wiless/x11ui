// sample code to test progressbar
package main

import (
	"time"

	"github.com/wiless/x11ui"
)

func main() {

	x11ui.SetResourcePath("../../fonts/")
	app := x11ui.NewApplication("Hello World", 800, 600, false, false)
	// var p x11ui.TextBox
	x11ui.DEBUG_LEVEL = 2
	x11ui.WidgetFactory(app.AppWin(), 0, 0, 100, 100)
	x11ui.NewTextBox("Hello ", app.AppWin(), 0, 400, 300, 50)
	//

	a := x11ui.NewImgButton("OK", app.AppWin(), 0, 200, 200, 200)
	t := x11ui.NewImgButton("OK", app.AppWin(), 200, 500, 200, 200)
	go func() {
		a.SetPicture("hg.png")
		t.SetPicture("hg.png")
		time.Sleep(1 * time.Second)
		t.SetPicture("Well_003.png")

	}()

	/// show a checkboxs
	cb := x11ui.NewCheckBox("Disable", app.AppWin(), 40, 40, 200, 200)
	cb.RePaint()
	app.Show()

}
