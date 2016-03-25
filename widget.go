package x11ui

import (
	"log"

	"github.com/BurntSushi/xgb/xproto"
	"github.com/BurntSushi/xgbutil"
	"github.com/BurntSushi/xgbutil/mousebind"
	"github.com/BurntSushi/xgbutil/xevent"
	"github.com/BurntSushi/xgbutil/xwindow"
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
	win.Create(parent, r.X, r.Y, r.Width, r.Height, xproto.CwBackPixel, 0xfffff)

	// win.Create(parent, r.X, r.Y, r.Width, r.Height, mask, values...)
	if err != nil {
		log.Fatal(err)
	}
	// win.MoveResize(r.X, r.Y, r.Width, r.Height)
	if p == nil {
		win.Change(xproto.CwBackPixel, 0x684426)

	} else {
		// win.Change(xproto.CwBackPixel, 0xFFAA00)
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
