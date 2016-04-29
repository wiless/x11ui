package main

import "github.com/wiless/x11ui"

func main() {

	x11ui.SetResourcePath("../../fonts/")
	app := x11ui.NewApp(false, 500, 500)
	w := app.NewChildWindow("Something ", 50, 50, 350, 350)

	sview := x11ui.NewScrollView("View My", w, 0, 0, 350, 350)
	// sview.View().SetBGcolor(color.RGBA{255, 200, 100, 30})
	lbl := x11ui.NewLabel("VBar 1", sview.View(), 30, 30, 100, 30)
	sview.AddWidget(lbl.Widget)
	lbl = x11ui.NewLabel("VBar 2", sview.View(), 30, 60, 100, 30)
	sview.AddWidget(lbl.Widget)
	lbl = x11ui.NewLabel("VBar 3", sview.View(), 30, 90, 100, 30)
	sview.AddWidget(lbl.Widget)

	app.Show()
}
