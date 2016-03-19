package main

import (
	"log"
	"time"

	"github.com/BurntSushi/xgbutil/ewmh"

	"github.com/wiless/x11ui"
)

var child, child2, dlg *x11ui.Window
var done chan bool
var app *x11ui.Application

func init() {
	// glfw.Init()
}

func SayHello() {
	log.Print("I am welcome")
	// c := glfw.GetX11Display()

	// ewmh.WmStateReq(app.AppWin().X(), dlg.Id, ewmh.StateToggle, "_NET_WM_STATE_FULLSCREEN")

	// ewmh.ActiveWindowSet(app.AppWin().X(), dlg.Id) //_NET_ACTIVE_WINDOW
	// ewmh.WmStateReq(app.AppWin().X(), dlg.Id, ewmh.StateToggle, "_NET_WM_STATE_MODAL")
	// ewmh.WmStateReq(app.AppWin().X(), dlg.Id, ewmh.StateToggle, "_NET_WM_STATE_ABOVE")
	// ewmh.WmStateReq(app.AppWin().X(), dlg.Id, ewmh.StateToggle, "_NET_WM_STATE_DEMANDS_ATTENTION")
	// ewmh.WmStateReq(app.AppWin().X(), dlg.Id, ewmh.StateToggle, "_NET_WM_ACTION_MOVE")

	// s, _ := ewmh.SupportedGet(app.AppWin().X())
	// log.Println("Supported set ", s)
	// ewmh.WmStateReq(app.AppWin().X(), dlg.Id, ewmh.t, "_NET_WM_WINDOW_TYPE_DIALOG")
	// ewmh.WmWindowOpacitySet(app.AppWin().X(), dlg.Id, .5)

}

func main() {

	app = x11ui.NewApplication("Hello World", 1000, 250, false, false)
	app.InvertBackGround()
	// X := s.AppWin().X()
	// win := s.AppWin().Id
	// // mg, _ := motif.WmHintsGet(X, s.AppWin().Id)
	// wh, _ := ewmh.WmStrutGet(X, win)
	// iwh, e := icccm.WmHintsGet(X, win)

	// log.Println("Window Manager supports ", wh)
	// log.Println("Window Manager supports ", iwh, e)

	// s.RegisterGlobalKey("control-mod2-f", DoThis)
	app.RegisterKey("q", app.Close)
	app.RegisterKey("i", app.InvertBackGround)
	app.RegisterKey("f", app.FullScreen)
	app.RegisterKey("shift-?", SayHello)
	app.AppWin().OnClick(mainfn)

	// dlg := s.NewFloatingWindow("File Open")
	// dlg := newWindow(s.xu, nil, 0, 0, 400, 800)
	// dlg.SetTitle("Dialog")
	// dlg.OnClick(floatwin)
	// dlg = app.NewFloatingWindow("New Dialog", 0, 300)
	done = make(chan bool)

	// child.OnClick(childfn)
	app.SetLayoutSpacing(5, 5)

	app.AutoLayout(x11ui.LayoutHor, 10, 10)

	app.AddButton("Hello", 0, 0).OnClick(mainfn)
	app.AddButton("Sensor 1", 0, 0).OnClick(extrafn1)
	app.AddButton("Long Process", 100, 0).OnClick(extrafn2)
	app.AddToggleBtn("Pause", 200, 0).OnClickAdv(toggle)

	app.AutoLayout(x11ui.LayoutVer, 10, 200)

	app.AddButton("Pink 4", 300, 0).OnClick(extrafn1)
	app.AddButton("graph", 400, 0).OnClick(plotgraph)
	app.AddButton("Pink 4", 300, 0)
	app.AddButton("graph", 400, 0)
	app.AddButton("Pink 4", 300, 0)
	app.AddButton("graph", 400, 0)
	app.AutoLayout(x11ui.LayoutHor, 120, 200, 602, 400)
	child = app.NewChildWindow("View 1")
	child2 = app.NewChildWindow("View 2", 100, 100, 400, 500)

	// obtn := NewButton("Orange ", dlg, 0, 400)
	// obtn.SetBackGround(0xFF00FF)
	// // obtn.StackSibling(dlg.Id, xproto.StackModeAbove)
	// obtn.OnClick(okbutton)
	// btn.OnClick(childfn)
	// btn.SetBackGround(0xFFAA00)
	go startPlot()
	app.Show()

}

func startPlot() {
	var pause bool
	for {
		select {
		case pause = <-done:

		default:
			if !pause {
				child.Plot()
				child2.Plot()
				time.Sleep(100 * time.Millisecond)
			}
		}

	}

}
func toggle(w *x11ui.Window, _, _ int) {

	w.Toggle()
	// ewmh.WmStateReq(child2.X(), child2.Id, ewmh.StateToggle, "_NET_WM_STATE_HIDDEN")
	// state, _ := icccm.WmStateGet(child2.X(), child2.Id)
	// state.State = icccm.StateIconic

	ewmh.WmStateSet(child2.X(), child2.Id, []string{"WM_DELETE"})

	done <- w.IsChecked()

}

func mainfn() {

}

func plotgraph() {

	child.Plot()

}
func floatwin() {

}
func childfn() {

}
func extrafn(w *x11ui.Window, x, y int) {

}
func extrafn1() {

}
func extrafn2() {

	for i := 0; i < 5; i++ {
		// log.Println("\r ===== Counter ==== ", w.Title(), "@", x, y, "===========", i)
		time.Sleep(1 * time.Second)
	}

}
func DoThis() {
	log.Printf("Hello Sendil")
}

func okbutton() bool {
	// log.Println("===== PUSH BUTTON  ==== ", w.Title(), "@", x, y, "===========")
	return true
}
