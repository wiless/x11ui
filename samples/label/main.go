package main

import (
	"time"

	"github.com/wiless/x11ui"
)

var l *x11ui.Label

func main() {
	x11ui.SetResourcePath("../../fonts/")
	ap := x11ui.NewApplication("Hello World", 500, 500, false, false)
	l = x11ui.NewLabel("Welcome to 1st Pin ", ap.AppWin(), 0, 0, 150, 30)
	// l.ClkFn = hello
	go hello()
	l.SetAlignMode(x11ui.AlignHVCenter)
	// l.AutoResize(true)
	// l.HoverFn = hello
	ap.Show()
}

func hello() {
	t := time.Tick(1 * time.Second)
	for range t {
		l.SetLabel(time.Now().Format("3:04:05pm"))
	}

}
