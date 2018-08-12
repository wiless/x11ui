package x11ui

import (
	"fmt"
	"image/color"
	"log"
	"os"
	"strings"

	"github.com/lucasb-eyer/go-colorful"

	"github.com/BurntSushi/xgb/xproto"
	"github.com/BurntSushi/xgbutil"
	"github.com/BurntSushi/xgbutil/ewmh"
	"github.com/BurntSushi/xgbutil/keybind"
	"github.com/BurntSushi/xgbutil/mousebind"
	"github.com/BurntSushi/xgbutil/xevent"
	"github.com/BurntSushi/xgbutil/xwindow"
)

type LayoutDirection int

const (
	LayoutNone LayoutDirection = iota
	LayoutHor
	LayoutVer
)

type Handler func()
type Application struct {
	xu                 *xgbutil.XUtil
	appWin             *Window
	NWindows           int
	Childs             []*xwindow.Window
	keycallbacks       keybind.KeyPressFun
	KeyMaps            map[string]Handler
	ClickMaps          map[xproto.Window]OnClickFn
	w, h               int
	title              string
	Debug              bool
	Dark               bool
	bg, fg             color.Color
	l                  LayoutDirection
	pvsChildRect       Rect
	defW, defH         int
	hspacing, vspacing int
}

func (a *Application) X() *xgbutil.XUtil {
	return a.xu
}

func (a *Application) Width() int {
	return int(a.xu.Screen().WidthInPixels)

}
func (a *Application) Height() int {
	return int(a.xu.Screen().HeightInPixels)
}

func (a *Application) AutoLayout(l LayoutDirection, newpos ...int) {
	a.l = l
	if len(newpos) > 0 {
		a.pvsChildRect = newRect(newpos...)
		if l == LayoutHor {
			a.pvsChildRect.ShiftRight(-(a.pvsChildRect.Width + a.hspacing))
		}
		if l == LayoutVer {
			a.pvsChildRect.ShiftDown(-(a.pvsChildRect.Height + a.vspacing))
		}
	}
}

func NewApp(fullscreen bool, width, height int) *Application {

	return NewApplication("X11 Application", width, height, true, fullscreen)
}

func NewApplication(title string, width, height int, resizeable, fullApplication bool) *Application {
	s := Application{w: width, h: height, title: title, Dark: false}
	s.KeyMaps = make(map[string]Handler)
	s.title = title
	s.Debug = true
	s.defW, s.defH = width, height
	var err error
	s.xu, err = xgbutil.NewConn()
	if err != nil {
		log.Fatal("Application: ", err)
	}

	keybind.Initialize(s.xu)
	mousebind.Initialize(s.xu)

	s.defaultWindow()
	if fullApplication {
		s.FullScreen()
	}

	s.title = title
	s.appWin.SetTitle(title)

	s.keycallbacks = s.keybHandler
	xevent.KeyPressFun(s.keybHandler).Connect(s.xu, s.AppWin().Id)

	mousebind.ButtonPressFun(s.mouseHandler).Connect(s.xu, s.appWin.Id, "1", false, false)
	// cb1 := keybind.KeyPressFun(
	// 	func(X *xgbutil.XUtil, e xevent.KeyPressEvent) {
	// 		log.Println("Key press!")
	// 	})

	// root, err := xwindow.Generate(X)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	return &s
}

func (a *Application) AppWin() *Window {
	return a.appWin
}

func (s *Application) mouseHandler(X *xgbutil.XUtil, e xevent.ButtonPressEvent) {
	// log.Println("Application ", s.appWin.Title(), " clicked at ", e.EventX, e.EventY)

	if s.appWin.clkAdv != nil {
		s.appWin.clkAdv(s.appWin, int(e.EventX), int(e.EventY))
	}
	if s.appWin.clk != nil {
		s.appWin.clk()
	}

}

func (s *Application) keybHandler(X *xgbutil.XUtil, e xevent.KeyPressEvent) {
	if s.xu != X {
		log.Println("\n I am not the write handler !!")
		return
	}
	modStr := keybind.ModifierString(e.State)
	keyStr := keybind.LookupString(X, e.State, e.Detail)

	modStr = strings.Replace(modStr, "lock-", "", -1)
	modStr = strings.Replace(modStr, "mod2", "", -1)
	finalstr := keyStr
	if modStr != "" {
		finalstr = fmt.Sprint(modStr, keyStr)
	}
	// log.Println("Event code is ", e.Detail)
	// log.Printf("%s MAPS  to Keycode %v ", finalstr, keybind.StrToKeycodes(s.xu, finalstr))

	if fn, ok := s.KeyMaps[finalstr]; ok {
		if s.Debug {
			log.Printf("Caught Key : %s", finalstr)
		}
		fn()
	} else {
		log.Printf("Caught Key : %s", finalstr)
	}

}

func (s *Application) RegisterGlobalKey(keyname string, fn Handler) bool {

	if s.RegisterKey(keyname, fn) {
		if err := keybind.GrabKeyboard(s.xu, s.xu.RootWin()); err != nil {
			log.Fatalf("Could not grab keyboard: %s", err)
			return false
		}
		log.Println("WARNING: We are taking *complete* control of the root " +
			"window. The only way out is to press 'Control + Escape' or to " +
			"close the window with the mouse.")
	} else {
		return false
	}
	return true

}

//RegisterKey registers a key with a function
func (s *Application) RegisterKey(keyname string, fn Handler) bool {
	// log.Println("Current Root ", s.xu)

	s.KeyMaps[keyname] = fn
	if s.Debug {
		log.Println("Registering ", keyname, " with this ", fn)
	}
	// if len(s.KeyMaps) == 0 {
	// 	xevent.KeyPressFun(s.keybHandler).Connect(s.xu, s.bgWin.Id)
	// }

	return true
}

func (s *Application) Show() {

	xevent.Main(s.xu)

}

func (s *Application) DefaultKeys(enable bool) {
	if enable {
		s.RegisterKey("q", s.Close)
		s.RegisterKey("f", s.FullScreen)
	}
}

func (s *Application) Close() {
	xevent.Quit(s.xu)
}
func (s *Application) XWin() *xwindow.Window {
	return s.appWin.Window
}

func (s *Application) Empty() *Window {
	var w Window
	w.Id = s.xu.RootWin()
	return &w
}

func (s *Application) defaultWindow() {
	// s.NWindows++
	// s.appWin = s.NewFloatingWindow("Application ")
	// s.appWin.WMMoveResize(0, 0, 1024, 523)
	var err error

	s.appWin = new(Window)
	s.appWin.Window, err = xwindow.Generate(s.xu)

	if err != nil {
		log.Fatalln(err)
	}
	sinfo := s.xu.Screen()
	w, h := int(sinfo.WidthInPixels)/2, int(sinfo.HeightInPixels)/2
	w, h = s.defW, s.defH

	s.appWin.Create(s.xu.RootWin(), 0, 0, w, h, xproto.CwBackPixel, 0x101010)

	// Listen for Key{Press,Release} events.
	s.appWin.Listen(xproto.EventMaskKeyPress, xproto.EventMaskKeyRelease, xproto.EventMaskButtonPress, xproto.EventMaskButtonPress, xproto.EventMaskSubstructureNotify, xproto.EventMaskStructureNotify)

	// xevent.ResizeRequestFun(
	// 	func(p *xgbutil.XUtil, e xevent.ResizeRequestEvent) {
	// 		log.Printf("2. Received event ", e)
	//
	// 	}).Connect(s.xu, s.appWin.Id)
	//
	xevent.ConfigureRequestFun(
		func(p *xgbutil.XUtil, e xevent.ConfigureRequestEvent) {
			//
			// log.Printf("4. CONFIGURE REEQUEST ", e)

		}).Connect(s.xu, s.appWin.Id)

	xevent.ConfigureNotifyFun(
		func(p *xgbutil.XUtil, e xevent.ConfigureNotifyEvent) {

			if e.SequenceId() != 126 {
				// log.Printf("3. Received CONFIGNOTIFICATION ", e)
			}

		}).Connect(s.xu, s.appWin.Id)

	// Make this window close gracefully.
	s.appWin.WMGracefulClose(
		func(w *xwindow.Window) {
			xevent.Detach(w.X, w.Id)
			keybind.Detach(w.X, w.Id)
			mousebind.Detach(w.X, w.Id)
			w.Destroy()
			// if quit {
			xevent.Quit(w.X)
			// }
		})

	// Map the window.
	s.appWin.Map()

	// // Get a random background color, create the window (ask to receive button
	// // release events while we're at it) and map the window.
	// bgColor := rand.Intn(0xffffff + 1)
	// win.Create(X.RootWin(), 0, 0, 200, 200,
	// 	xproto.CwBackPixel|xproto.CwEventMask,
	// 	uint32(bgColor), xproto.EventMaskButtonRelease)

	// // WMGracefulClose does all of the work for us. It sets the appropriate
	// // values for WM_PROTOCOLS, and listens for ClientMessages that implement
	// // the WM_DELETE_WINDOW protocol. When one is found, the provided callback
	// // is executed.
	// win.WMGracefulClose(
	// 	func(w *xwindow.Window) {
	// 		// Detach all event handlers.
	// 		// This should always be done when a window can no longer
	// 		// receive events.
	// 		xevent.Detach(w.X, w.Id)
	// 		mousebind.Detach(w.X, w.Id)
	// 		w.Destroy()

	// 		// Exit if there are no more windows left.
	// 		counter--
	// 		if counter == 0 {
	// 			os.Exit(0)
	// 		}
	// 	})

	// // It's important that the map comes after setting WMGracefulClose, since
	// // the WM isn't obliged to watch updates to the WM_PROTOCOLS property.
	// win.Map()

	// // A mouse binding so that a left click will spawn a new window.
	// // Note that we don't issue a grab here. Typically, window managers will
	// // grab a button press on the client window (which usually activates the
	// // window), so that we'd end up competing with the window manager if we
	// // tried to grab it.
	// // Instead, we set a ButtonRelease mask when creating the window and attach
	// // a mouse binding *without* a grab.
	// err = mousebind.ButtonReleaseFun(
	// 	func(X *xgbutil.XUtil, ev xevent.ButtonReleaseEvent) {
	// 		newWindow(X)
	// 	}).Connect(X, win.Id, "1", false, false)
	// if err != nil {
	// 	log.Fatal(err)
	// }

}

func (s *Application) NewFloatingWindow(title string, dims ...int) *Window {
	// X := s.xu
	w := s.newWindow(0, newRect(dims...))
	w.SetTitle(title)
	if len(dims) < 3 {
		w.Resize(700, 500)
	}

	return w

}

func (s *Application) FullScreen() {

	err := ewmh.WmStateReq(s.xu, s.appWin.Id, ewmh.StateToggle, "_NET_WM_STATE_FULLSCREEN")

	// ewmh.WmStateGet(s.xu, s.appWin.Id)

	// s.appWin.ClearAll()
	if err != nil {
		log.Fatal(err)
	}
}

func (s *Application) InvertBackGround() {
	if s.appWin != nil {

		if s.Dark {
			s.appWin.Change(xproto.CwBackPixel, 0x265170)

		} else {
			s.appWin.Change(xproto.CwBackPixel, 0x1f1f1f)
		}
		s.Dark = !s.Dark

		s.appWin.ClearAll()

	}

}

func toUint32(colorful.Color) uint32 {
	return 0
}

func (a *Application) SetDefaultKeys() {
	a.RegisterKey("q", a.Close)
	a.RegisterKey("i", a.InvertBackGround)
	a.RegisterKey("f", a.FullScreen)
}

func (s *Application) newWindow(p xproto.Window, r Rect) *Window {
	w := new(Window)
	var parent xproto.Window
	if p == 0 {
		parent = s.xu.RootWin()
	} else {
		parent = p
	}

	win, err := xwindow.Generate(s.xu)
	win.Create(parent, r.X, r.Y, r.Width, r.Height, xproto.CwBackPixel, 0xfffff)
	if err != nil {
		log.Fatal(err)
	}

	// // win.MoveResize(r.X, r.Y, r.Width, r.Height)
	// if p == 0 {
	// 	win.Change(xproto.CwBackPixel, 0x1f1f1f)
	// } else {
	// 	// win.Change(xproto.CwBackPixel, 0xC0FFC0)
	// 	win.Change(xproto.CwBackPixel)
	// }

	w.Rect = r
	w.Window = win
	cb1 := mousebind.ButtonPressFun(w.mouseHandler)

	if p == 0 {
		log.Println("Registering right click also")
		err = cb1.Connect(s.xu, win.Id, "2", false, true)
	} else {
		err = cb1.Connect(s.xu, win.Id, "1", false, true)

	}
	if err != nil {
		log.Println("Error in binding mouse press ", err)
	}

	win.WMGracefulClose(
		func(w *xwindow.Window) {
			// Detach all event handlers.
			// This should always be done when a window can no longer
			// receive events.
			xevent.Detach(w.X, w.Id)
			mousebind.Detach(w.X, w.Id)
			w.Destroy()

			// Exit if there are no more windows left.
			s.NWindows--
			if s.NWindows == 0 {
				os.Exit(0)
			}
		})

	// It's important that the map comes after setting WMGracefulClose, since
	// the WM isn't obliged to watch updates to the WM_PROTOCOLS property.
	win.Map()

	return w
}

func (a *Application) AddButton(caption string, geo ...int) *Window {
	switch a.l {
	case LayoutVer:
		a.pvsChildRect.ShiftDown(a.pvsChildRect.Height + a.vspacing)
		obj := NewButton(caption, a.AppWin(), a.pvsChildRect.array()...)
		return obj
	case LayoutHor:
		a.pvsChildRect.ShiftRight(a.pvsChildRect.Width + a.hspacing)
		obj := NewButton(caption, a.AppWin(), a.pvsChildRect.array()...)
		return obj
	default:
		obj := NewButton(caption, a.AppWin(), geo...)
		a.pvsChildRect = obj.Rect
		return obj
	}

}

func (a *Application) AddToggleBtn(caption string, geo ...int) *Window {
	switch a.l {
	case LayoutVer:
		a.pvsChildRect.ShiftDown(a.pvsChildRect.Height + a.vspacing)
		obj := NewToggleButton(caption, a.AppWin(), a.pvsChildRect.array()...)
		return obj
	case LayoutHor:
		a.pvsChildRect.ShiftRight(a.pvsChildRect.Width + a.hspacing)
		obj := NewToggleButton(caption, a.AppWin(), a.pvsChildRect.array()...)
		return obj
	default:
		obj := NewToggleButton(caption, a.AppWin(), geo...)
		a.pvsChildRect = obj.Rect
		return obj
	}

}
func (a *Application) SetLayoutSpacing(dx, dy int) {
	a.hspacing, a.vspacing = dx, dy
}
func (a *Application) NewChildWindow(title string, dims ...int) *Window {
	var w *Window
	switch a.l {
	case LayoutVer:
		a.pvsChildRect.ShiftDown(a.pvsChildRect.Height + a.vspacing)
		w = a.newWindow(a.appWin.Id, a.pvsChildRect)

	case LayoutHor:
		a.pvsChildRect.ShiftRight(a.pvsChildRect.Width + a.hspacing)
		w = a.newWindow(a.appWin.Id, a.pvsChildRect)

	default:
		w = a.newWindow(a.appWin.Id, newRect(dims...))

		a.pvsChildRect = w.Rect

	}
	// w.SetBackGround(colorful.LinearRgb(0, 0, 0))
	w.bgcolor = color.RGBA{100, 100, 100, 255}
	g := w.drawView(StateNormal)
	w.finishPaint(g)
	// w.SetTitle(title)
	w.Detach()
	return w
}
