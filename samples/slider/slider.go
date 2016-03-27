package main

import (
	"image/color"
	"log"

	"github.com/BurntSushi/xgb/xproto"
	"github.com/BurntSushi/xgbutil"

	"github.com/BurntSushi/xgbutil/mousebind"

	"github.com/wiless/x11ui"
)

type Slider struct {
	*x11ui.ProgressBar
	scaler    float64
	maxValue  float64
	viewWidth float64
}

func NewSlider(title string, p *x11ui.Window, dims ...int) *Slider {
	if p == nil {
		log.Fatal("Cannot Create Widget without Application")
	}
	slider := new(Slider)
	slider.maxValue = 100

	pbar := x11ui.NewProgressBar(title, p, dims...)
	pbar.SetBarColor(color.White)
	pbar.SetTextColor(color.Black)
	// pbar.SetValue(0.5)

	pbar.Widget().OnClickAdv(slider.drawBar)

	slider.ProgressBar = pbar
	slider.SetValue(0)
	slider.viewWidth = float64(slider.Widget().Rect.Width)
	slider.SetMaxValue(100)
	/// add event listeners
	slider.AddListeners()

	return slider
}

func (s *Slider) dragFunction(X *xgbutil.XUtil, rootX, rootY int, eventX, eventY int) {
	// X.Ungrab()
	log.Println(rootX, rootY, eventX, eventY)
}

func (s *Slider) AddListeners() {
	w := s.ProgressBar.Widget().Window
	mousebind.Drag(w.X, w.Id, w.Id, "1", false,
		func(X *xgbutil.XUtil, rx, ry, ex, ey int) (bool, xproto.Cursor) {
			s.drawBar(s.Widget(), ex, ey)
			return true, 0
		},
		func(X *xgbutil.XUtil, rx, ry, ex, ey int) {
			s.drawBar(s.Widget(), ex, ey)
		},
		func(X *xgbutil.XUtil, rx, ry, ex, ey int) {})

}

func (s *Slider) SetMaxValue(v float64) {
	s.maxValue = v
	s.scaler = 1.0 / s.viewWidth
	s.ProgressBar.SetFmtString("%5.2f")
	s.ProgressBar.SetDisplayScale(v)
}

func (s *Slider) MouseDrag() {

}

func (s *Slider) drawBar(w *x11ui.Window, x, y int) {
	v := float64(x) * s.scaler
	//mousebind.GrabPointer(xu, win, confine, cursor)
	s.SetValue(v)

}
