package x11ui

import (
	"fmt"
	"image/color"
	"log"
	"os"
	"strings"
	"time"

	"github.com/lucasb-eyer/go-colorful"

	"github.com/BurntSushi/xgb/xproto"
	"github.com/BurntSushi/xgbutil"
	"github.com/BurntSushi/xgbutil/ewmh"
	"github.com/BurntSushi/xgbutil/keybind"
	"github.com/BurntSushi/xgbutil/mousebind"
	"github.com/BurntSushi/xgbutil/xevent"
	"github.com/BurntSushi/xgbutil/xwindow"
)

// LayoutDirection specifies the direction for automatic widget layout.
type LayoutDirection int

const (
	// LayoutNone indicates no automatic layout direction.
	LayoutNone LayoutDirection = iota
	// LayoutHor indicates horizontal layout.
	LayoutHor
	// LayoutVer indicates vertical layout.
	LayoutVer
)

// Handler is a function type for event handlers.
type Handler func()

// Application represents an X11 application, managing windows and events.
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
	hpadding, vpadding int
	moveTimer          *time.Timer
	resizeTimer        *time.Timer
}

// X returns the underlying *xgbutil.XUtil connection for the application.
func (a *Application) X() *xgbutil.XUtil {
	return a.xu
}

// Width returns the width of the application's screen in pixels.
func (a *Application) Width() int {
	return int(a.xu.Screen().WidthInPixels)

}

// Height returns the height of the application's screen in pixels.
func (a *Application) Height() int {
	return int(a.xu.Screen().HeightInPixels)
}

// AutoLayout sets the automatic layout direction for new child windows.
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

// calculateChildrenBoundingBox calculates the bounding box of all child windows.
// func (a *Application) calculateChildrenBoundingBox() Rect {
// 	if len(a.childrenWindows) == 0 {
// 		return Rect{}
// 	}

// 	minX, minY := a.childrenWindows[0].Rect.X, a.childrenWindows[0].Rect.Y
// 	maxX, maxY := a.childrenWindows[0].Rect.X+a.childrenWindows[0].Rect.Width,
// 		a.childrenWindows[0].Rect.Y+a.childrenWindows[0].Rect.Height

// 	for _, child := range a.childrenWindows {
// 		childRect := child.Rect
// 		if childRect.X < minX {
// 			minX = childRect.X
// 		}
// 		if childRect.Y < minY {
// 			minY = childRect.Y
// 		}
// 		if childRect.X+childRect.Width > maxX {
// 			maxX = childRect.X + childRect.Width
// 		}
// 		if childRect.Y+childRect.Height > maxY {
// 			maxY = childRect.Y + childRect.Height
// 		}
// 	}
// 	return Rect{minX, minY, maxX - minX, maxY - minY}
// }

// NewApp creates a new application with a default title and specified dimensions.
func NewApp(fullscreen bool, width, height int) *Application {

	return NewApplication("X11 Application", width, height, true, fullscreen)
}

// NewApplication creates a new X11 application with the given title, dimensions, and resize/fullscreen options.
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

	// var mkey xgbutil.KeyKey

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

// AppWin returns the main application window.
func (a *Application) AppWin() *Window {
	return a.appWin
}

// mouseHandler processes mouse button press events for the application.
func (s *Application) mouseHandler(X *xgbutil.XUtil, e xevent.ButtonPressEvent) {
	// log.Println("Application ", s.appWin.Title(), " clicked at ", e.EventX, e.EventY)

	if s.appWin.clkAdv != nil {
		s.appWin.clkAdv(s.appWin, int(e.EventX), int(e.EventY))
	}
	if s.appWin.clk != nil {
		s.appWin.clk()
	}

}

func (s *Application) mouseReleaseHandler(X *xgbutil.XUtil, e xevent.ButtonReleaseEvent) {
	// No-op, as isDragging is removed
}

// keybHandler processes keyboard press events for the application.
func (s *Application) keybHandler(X *xgbutil.XUtil, e xevent.KeyPressEvent) {
	if s.xu != X {

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

		}
		fn()
	} else {
		log.Printf("Caught Key : %s", finalstr)
	}

}

// RegisterGlobalKey registers a key to be handled globally by the application.
// This will grab the keyboard, meaning no other application will receive keyboard input.
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

// RegisterKey registers a key with a function to be handled by the application.
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

// Show starts the application's event loop and displays the main window.
func (s *Application) Show() {

	xevent.Main(s.xu)

}

// DefaultKeys enables or disables default keybindings for the application (e.g., 'q' for close, 'f' for fullscreen).
func (s *Application) DefaultKeys(enable bool) {
	if enable {
		s.RegisterKey("q", s.Close)
		s.RegisterKey("Q", s.Close)
		s.RegisterKey("f", s.FullScreen)
		s.RegisterKey("F", s.FullScreen)
	}
}

// Close terminates the application.
func (s *Application) Close() {
	xevent.Quit(s.xu)
}

// XWin returns the underlying *xwindow.Window of the application's main window.
func (s *Application) XWin() *xwindow.Window {
	return s.appWin.Window
}

// Empty returns an empty Window associated with the application's root window.
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

	s.appWin.Create(s.xu.RootWin(), 0, 0, w, h, xproto.CwBackPixel, 0x000000)
	s.appWin.prevRect = s.appWin.Rect
	ewmh.WmAllowedActionsSet(s.xu, s.appWin.Id, []string{
		"_NET_WM_ACTION_MOVE",
		"_NET_WM_ACTION_RESIZE",
		"_NET_WM_ACTION_SHADE",
	})

	// Listen for Key{Press,Release} events.
	s.appWin.Listen(xproto.EventMaskKeyPress, xproto.EventMaskKeyRelease, xproto.EventMaskButtonPress, xproto.EventMaskButtonPress, xproto.EventMaskSubstructureNotify, xproto.EventMaskStructureNotify, xproto.EventMaskFocusChange)

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
			currentRect := Rect{int(e.X), int(e.Y), int(e.Width), int(e.Height)}

			// Check for size change
			sizeChanged := currentRect.Width != s.appWin.prevRect.Width || currentRect.Height != s.appWin.prevRect.Height
			if sizeChanged {
				if s.appWin.onSizeChange != nil {
					s.appWin.onSizeChange(currentRect.Width, currentRect.Height)
				}
				// Reset or start the resize timer
				if s.resizeTimer != nil {
					s.resizeTimer.Stop()
				}
				s.resizeTimer = time.AfterFunc(100*time.Millisecond, func() {
					// log.Println("Window resize complete")
					if s.appWin.onResizeComplete != nil {
						s.appWin.onResizeComplete(currentRect.Width, currentRect.Height)
					}
				})
			}

			// Check for position change
			positionChanged := currentRect.X != s.appWin.prevRect.X || currentRect.Y != s.appWin.prevRect.Y
			if positionChanged {
				// log.Println("Window is moving")
				if s.appWin.onMove != nil {
					s.appWin.onMove(currentRect.X, currentRect.Y)
				}

				// Reset or start the move timer
				if s.moveTimer != nil {
					s.moveTimer.Stop()
				}
				s.moveTimer = time.AfterFunc(100*time.Millisecond, func() {
					// log.Println("Window move complete")
					if s.appWin.onMoveComplete != nil {
						s.appWin.onMoveComplete(currentRect.X, currentRect.Y)
					}
				})
			}

			s.appWin.prevRect = currentRect
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

func (s *Application) Restore() {
	// _NET_WM_ACTION_SHADE

	// err := ewmh.WmStateReq(s.xu, s.appWin.Id, ewmh.StateToggle, "_NET_WM_ACTION_SHADE")
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// Restore the window by setting an empty state.
	err := ewmh.WmStateSet(s.xu, s.appWin.Id, []string{})
	if err != nil {
		panic(err)
	}

}

// https://specifications.freedesktop.org/wm-spec/1.3/ar01s05.html
func (s *Application) Hide() {
	// _NET_WM_STATE_HIDDEN


	// Get the _NET_WM_STATE_HIDDEN atom.
	hiddenAtom, err := ewmh.WmStateGet(s.xu, s.appWin.Id)
	if err != nil {
		panic(err)
	}


	ewmh.WmNameSet(s.xu, s.appWin.Id, "Minimize Example")

	// Set the WM_STATE property to hidden.
	err = ewmh.WmStateSet(s.xu, s.appWin.Id, []string{"_NET_WM_STATE_FULLSCREEN"})

	// err = ewmh.WmStateReq(s.xu, s.appWin.Id, ewmh.StateToggle, "_NET_WM_STATE_FOCUSED")

	if err != nil {
		panic(err)
	}
	// err := ewmh.WmStateReq(s.xu, s.appWin.Id, ewmh.StateToggle, "_NET_WM_ACTION_MINIMIZE")
	// if err != nil {
	// 	log.Fatal(err)
	// }

}

func (s *Application) Maximize() {

	err := ewmh.WmStateReq(s.xu, s.appWin.Id, ewmh.StateToggle, "_NET_WM_STATE_MAXIMIZED_HORZ")

	// ewmh.WmStateGet(s.xu, s.appWin.Id)

	// s.appWin.ClearAll()
	if err != nil {
		log.Fatal(err)
	}
}

func (s *Application) FullScreen() {

	err := ewmh.WmStateReq(s.xu, s.appWin.Id, ewmh.StateToggle, "_NET_WM_STATE_FULLSCREEN")
	// err := ewmh.WmStateReq(s.xu, s.appWin.Id, ewmh.StateAdd, "_NET_WM_STATE_HIDDEN")
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

func (a *Application) NewContainer(layout LayoutDirection, dims ...int) *Container {
	tempWindow := a.NewChildWindow("", dims...)

	container := &Container{
		Window:          *tempWindow, // Dereference the pointer to assign to the embedded value
		layoutDirection: layout,
		children:        make([]WindowProvider, 0),
		hspacing:        a.hspacing, // Default spacing
		vspacing:        a.vspacing,

		pvsChildRect: tempWindow.Rect, // Use the Rect from the temporary window
	}

	container.SetBGcolor(color.RGBA{125, 125, 0, 0}) // Fully transparent black

	container.OnResizeComplete(func(w, h int) {
		container.RelayoutChildren()
	})

	return container
}

func (a *Application) AddButton(caption string, geo ...int) *Window {
	var r Rect
	if len(geo) > 0 {
		r = newRect(geo...)
	} else {
		switch a.l {
		case LayoutVer:
			a.pvsChildRect.ShiftDown(a.pvsChildRect.Height + a.vspacing)
			r = a.pvsChildRect
		case LayoutHor:
			a.pvsChildRect.ShiftRight(a.pvsChildRect.Width + a.hspacing)
			r = a.pvsChildRect
		}
	}

	paddedRect := r.WithPadding(a.hpadding, a.vpadding)
	obj := NewButton(caption, a.AppWin(), paddedRect.array()...)
	a.pvsChildRect = r
	return obj
}

func (a *Application) AddToggleBtn(caption string, geo ...int) *Window {
	var r Rect
	if len(geo) > 0 {
		r = newRect(geo...)
	} else {
		switch a.l {
		case LayoutVer:
			a.pvsChildRect.ShiftDown(a.pvsChildRect.Height + a.vspacing)
			r = a.pvsChildRect
		case LayoutHor:
			a.pvsChildRect.ShiftRight(a.pvsChildRect.Width + a.hspacing)
			r = a.pvsChildRect
		}
	}

	paddedRect := r.WithPadding(a.hpadding, a.vpadding)
	obj := NewToggleButton(caption, a.AppWin(), paddedRect.array()...)
	a.pvsChildRect = r
	return obj
}
func (a *Application) SetLayoutSpacing(dx, dy int) {
	a.hspacing, a.vspacing = dx, dy
}

func (a *Application) SetLayoutPadding(dx, dy int) {
	a.hpadding, a.vpadding = dx, dy
}
func (a *Application) NewChildWindow(title string, dims ...int) *Window {
	var r Rect
	if len(dims) > 0 {
		r = newRect(dims...)
	} else {
		switch a.l {
		case LayoutVer:
			a.pvsChildRect.ShiftDown(a.pvsChildRect.Height + a.vspacing)
			r = a.pvsChildRect
		case LayoutHor:
			a.pvsChildRect.ShiftRight(a.pvsChildRect.Width + a.hspacing)
			r = a.pvsChildRect
		}
	}

	paddedRect := r.WithPadding(a.hpadding, a.vpadding)
	w := a.newWindow(a.appWin.Id, paddedRect)
	a.pvsChildRect = r

	// w.SetBackGround(colorful.LinearRgb(0, 0, 0))
	w.bgcolor = color.RGBA{0, 0, 0, 255}
	w.drawView(StateNormal)
	w.finishPaint(w.ximg)
	// w.SetTitle(title)
	w.Detach()
	return w
}
