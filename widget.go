package x11ui

import (
	"image"
	"image/color"
	"log"

	"github.com/BurntSushi/xgb/xproto"
	"github.com/BurntSushi/xgbutil"
	"github.com/BurntSushi/xgbutil/mousebind"
	"github.com/BurntSushi/xgbutil/xevent"
	"github.com/BurntSushi/xgbutil/xgraphics"
	"github.com/BurntSushi/xgbutil/xwindow"
	"github.com/llgcode/draw2d/draw2dimg"
	"github.com/llgcode/draw2d/draw2dkit"
)

func NewWidget(X *xgbutil.XUtil, p *Window, t string, dims ...int) *Window {
	w := new(Window)
	w.title = t
	w.background = systemBG
	var parent xproto.Window
	if p == nil {
		parent = X.RootWin()
	} else {
		parent = p.Id
	}
	mousebind.Initialize(X)
	r := newRect(dims...)
	win, err := xwindow.Generate(X)
	// s := X.Screen()

	// mask := xproto.GcForeground | xproto.GcGraphicsExposures
	// values := []uint32{s.BlackPixel, 0}
	win.Create(parent, r.X, r.Y, r.Width, r.Height, xproto.CwBackingPixel, 0xffffff)

	// win.Create(parent, r.X, r.Y, r.Width, r.Height, mask, values...)
	if err != nil {
		log.Fatal(err)
	}
	// win.MoveResize(r.X, r.Y, r.Width, r.Height)
	if p == nil {
		win.Change(xproto.CwBackPixel, 0xFF00FF, 0xFF0000)

	} else {

		// win.Change(xproto.CwBackPixel, 0x101010) // dark shade
		// win.Change(xproto.CwBackPixel, 0x0000FF, 0xFF0000)
		// win.Change(xproto.CwBorderPixel, 0xFF0000)

		// win.Change(xproto.CwBackPixel, 0x00000)
		log.Println("I am not nill")

	}

	//if p == nil {
	win.Listen(xproto.EventMaskKeyPress, xproto.EventMaskKeyRelease, xproto.EventMaskButtonPress, xproto.EventMaskButtonRelease, xproto.EventMaskExposure, xproto.EventMaskEnterWindow, xproto.EventMaskLeaveWindow)
	//}

	w.Rect = r

	w.Window = win
	if p == nil {
		// xevent.ButtonPressFun(w.mouseHandler).Connect(X, win.Id, "1", false, true)
		xevent.ButtonPressFun(w.mouseHandler).Connect(X, win.Id)
	} else {
		mousebind.ButtonPressFun(w.mouseHandler).Connect(X, win.Id, "1", false, true)
		mousebind.ButtonReleaseFun(w.mouseReleaseHandler).Connect(X, win.Id, "1", false, true)
		xevent.EnterNotifyFun(w.onHoverEvent).Connect(X, win.Id)
		xevent.LeaveNotifyFun(w.onLeaveEvent).Connect(X, win.Id)
		mousebind.ButtonPressFun(w.mouseHandler).Connect(X, win.Id, "2", false, true)
	}
	// xevent.ExposeFun(w.Draw).Connect(X, win.Id)

	win.WMGracefulClose(
		func(w *xwindow.Window) {
			// Detach all event handlers.
			// This should always be done when a window can no longer
			// receive events.
			log.Printf("Window destroyed %d ", w.Id)
			xevent.Detach(w.X, w.Id)
			mousebind.Detach(w.X, w.Id)
			w.Destroy()
			// Exit if there are no more windows left.

		})

	// It's important that the map comes after setting WMGracefulClose, since
	// the WM isn't obliged to watch updates to the WM_PROTOCOLS property.
	win.Map()

	// xevent.ButtonPressFun(w.mouseHandler).Connect(X, win.Id)
	return w
}

func DrawDummy(w *Window, s WidgetState) {
	r := w.Rect
	r.MoveTo(0, 0)
	r.ImageRect()
	dest := image.NewRGBA(r.ImageRect())

	gc := draw2dimg.NewGraphicContext(dest)

	// bg := colorful.LinearRgb(.025, .025, .025)
	switch s {
	case StateNormal, StateReleased:
		gc.SetFillColor(color.RGBA{0x20, 0x20, 0x20, 20})
		gc.SetStrokeColor(systemFG)
	case StateHovered:
		gc.SetFillColor(color.RGBA{0x35, 0x20, 0x20, 20})
		gc.SetStrokeColor(systemFG)
	case StatePressed:
		gc.SetFillColor(color.RGBA{0x20, 0x30, 0x20, 20})
		gc.SetStrokeColor(systemFG)
	case StateSpecial:
		gc.SetFillColor(color.RGBA{0x20, 0x80, 0x20, 0x80})
		gc.SetStrokeColor(systemFG)
	}

	// // gc.SetLineJoin(draw2d.RoundJoin)
	// // gc.Rotate(math.Pi / 4.0)
	WW := float64(r.Width)
	HH := float64(r.Height)

	// Draw Background
	gc.SetLineWidth(0)
	gc.SetFillColor(color.RGBA{130, 120, 30, 10})
	gc.SetStrokeColor(color.RGBA{30, 120, 130, 10})
	draw2dkit.Rectangle(gc, 0, 0, WW, HH)
	gc.FillStroke()

	gc.FillStroke()
	gc.SetFillColor(color.RGBA{130, 120, 30, 10})
	gc.SetStrokeColor(color.RGBA{30, 120, 130, 10})
	draw2dkit.Circle(gc, 250, 10, 30)
	gc.FillStroke()

	gc.Close()
	g := xgraphics.NewConvert(w.X(), dest)

	// w.drawLabel(g, w.title)
	g.XSurfaceSet(w.Id)
	g.XDraw()
	g.XPaintRects(w.Id, r.ImageRect())
	// return g
	// return g
}
