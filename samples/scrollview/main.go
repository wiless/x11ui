package main

import (
	"image/color"

	"github.com/wiless/x11ui"
)

var sview *x11ui.ScrollView

func main() {

	x11ui.SetResourcePath("../../fonts/")
	app := x11ui.NewApp(false, 500, 500)
	w := app.NewChildWindow("Something ", 50, 50, 350, 350)

	sview = x11ui.NewScrollView("View My", w, 0, 0, 350, 350)
	// sview.View().SetBGcolor(color.RGBA{255, 200, 100, 30})
	sview.SetBackground(color.RGBA{255, 0, 0, 255})
	lbl := x11ui.NewLabel("VBar 1", sview.View(), 30, 30, 100, 30)
	sview.AddWidget(lbl.Widget)
	lbl.SetBackground(color.RGBA{255, 0, 255, 255})

	lbl = x11ui.NewLabel("VBar 2", sview.View(), 30, 60, 100, 30)
	sview.AddWidget(lbl.Widget)
	lbl.SetBackground(color.RGBA{0, 0, 255, 255})

	lbl = x11ui.NewLabel("VBar 3", sview.View(), 30, 90, 100, 30)
	sview.AddWidget(lbl.Widget)

	btn := x11ui.NewLabel("Ok", app.AppWin(), 0, 0, 100, 30) // app.NewChildWindow("Ok", 0, 0, 100, 30)
	btn.ClkFn = HideView

	lbl.SetBackground(color.RGBA{255, 120, 255, 255})

	btn1 := x11ui.NewLabel("Up", app.AppWin(), 0, 30, 100, 30) // app.NewChildWindow("Ok", 0, 0, 100, 30)
	btn1.ClkFn = sview.ScrollUp
	btn2 := x11ui.NewLabel("Down", app.AppWin(), 0, 60, 100, 30) // app.NewChildWindow("Ok", 0, 0, 100, 30)
	btn2.ClkFn = sview.ScrollDown

	sview.SetBackground(color.RGBA{255, 0, 0, 55})
	// sview.ShowScrollBars(true)

	app.Show()
}

func HideView() {
	// log.Println("Toggle Scroll view")
	sview.ShowScrollBars(!sview.IsVisible())
}
