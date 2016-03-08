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
	btn.isButton = true

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

func NewToggleButton(title string, p *Window, dims ...int) *Window {
	// var btn Window
	btn := newWindow(p.Window.X, p, title, dims...)
	btn.isButton = true
	btn.isCheckBox = true

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

// // drawGopher draws the gopher image to the canvas.
// func drawGopher(canvas *xgraphics.Image, gopher image.Image,
// 	win *xwindow.Window) {

// 	// Find the rectangle of the canvas where we're going to draw the gopher.

// 	gopherRect := image.Rect(50, 50, 200, 200) // midRect(x, y, gopherWidth, gopherHeight, width, height)

// 	// If the rectangle contains no pixels, don't draw anything.
// 	if gopherRect.Empty() {
// 		return
// 	}

// 	// Output a little message.
// 	// log.Printf("Drawing gopher at (%d, %d)", x, y)

// 	// Get a subimage of the gopher that's in sync with gopherRect.
// 	gopherPt := image.Pt(0, 0)
// 	// gopherPt := image.Pt(gopher.Bounds().Min.X, gopher.Bounds().Min.Y)

// 	// if gopherRect.Min.X == 0 {
// 	// 	gopherPt.X = gopherWidth - gopherRect.Dx()
// 	// }
// 	// if gopherRect.Min.Y == 0 {
// 	// 	gopherPt.Y = gopherHeight - gopherRect.Dy()
// 	// }

// 	// Create the canvas subimage.
// 	subCanvas := canvas.SubImage(gopherRect).(*xgraphics.Image)

// 	// Blend the gopher image into the sub-canvas.
// 	// This does alpha blending.
// 	xgraphics.Blend(subCanvas, gopher, gopherPt)

// 	// Now draw the changes to the pixmap.
// 	subCanvas.XDraw()

// 	// And paint them to the window.
// 	subCanvas.XPaint(win.Id)
// }
