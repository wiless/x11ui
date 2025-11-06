package x11ui

import "log"

type Button struct {
	*Window
}

func NewButton(title string, p *Window, dims ...int) *Window {
	// var btn Window
	if p == nil {
		log.Fatal("Cannot Create Widget without Application")
	}
	btn := newWindow(p.Window.X, p, title, dims...)
	// btn.SetTitle(title)
	btn.isButton = true
	btn.rePaint(StateNormal)
	return btn
}

func NewToggleButton(title string, p *Window, dims ...int) *Window {
	// var btn Window
	btn := newWindow(p.Window.X, p, title, dims...)
	btn.isButton = true
	btn.isCheckBox = true
	btn.rePaint(StateNormal)
	// btn.SetTitle(title)
	// btn.Rect = newRect(dims...)

	// sshot, gerr := xgraphics.NewDrawable(X, xproto.Drawable(btn.Window.Id))
	// if gerr != nil {
	// 	log.Println("Error Loading Drawable Image ", gerr)
	// }
	// sshot.XShowExtra("nothing", )
	// log.Println("Trying to save ", btn.Title()+".png")
	// sshot.SavePng(btn.Title() + ".png")

	return btn
}
