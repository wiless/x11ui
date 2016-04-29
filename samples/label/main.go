package main

import (
	"image/color"
	"log"
	"time"

	"github.com/wiless/x11ui"
)

var l *x11ui.Label

func main() {
	x11ui.SetResourcePath("../../fonts/")
	ap := x11ui.NewApplication("Hello World", 500, 500, false, false)
	child := ap.NewChildWindow("Hello world", 0, 0, 300, 100)
	child.SetBGcolor(color.RGBA{100, 0, 0, 30})
	l = x11ui.NewLabel("Welcome to 1st Pin ", child, 10, 10, 150, 30)
	l.SetAlignMode(x11ui.AlignHVCenter)
	l.AutoResize(true)
	// l.Ho = hello
	go hello()

	l.ClkFn = SayHello
	ap.Show()
}

func SayHello() {
	log.Print("Hi Clicked me")
}

func hello() {
	t := time.Tick(1 * time.Second)
	for range t {
		l.SetLabel(time.Now().Format("3:04:05pm"))
	}

}
