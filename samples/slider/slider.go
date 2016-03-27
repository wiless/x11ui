package main

import (
	"log"
	"math"

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
	stepSize  float64
}

func NewSlider(title string, p *x11ui.Window, dims ...int) *Slider {
	if p == nil {
		log.Fatal("Cannot Create Widget without Application")
	}
	slider := new(Slider)
	slider.maxValue = 100

	pbar := x11ui.NewProgressBar(title, p, dims...)
	// pbar.SetBarColor(color.White)
	// pbar.SetTextColor(color.Black)
	// pbar.SetValue(0.5)
	pbar.SetMargin(5)
	pbar.SetBorderWidth(2)
	pbar.Widget().OnClickAdv(slider.drawBar)

	slider.ProgressBar = pbar
	slider.viewWidth = float64(slider.Widget().Rect.Width) - 2*slider.ProgressBar.Margin()
	// slider.SetFmtString("%5.2f")

	slider.SetMaxValue(100)
	slider.SetStepSize(10)
	slider.SetValue(0)

	/// add event listeners
	slider.AddListeners()

	return slider
}

func TrimValues(v float64) float64 {
	switch {
	case v < 0:
		return 0
	case v > 1:
		return 1
	default:
		return v
	}
}

func (s *Slider) Value() float64 {

	return s.ProgressBar.Value() * s.maxValue

}

func (s *Slider) SetValue(v float64) {
	vv := TrimValues(v / s.maxValue)

	s.ProgressBar.SetValue(vv)
}

func (s *Slider) SetStepSize(v float64) {
	s.stepSize = v
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

		}, func(X *xgbutil.XUtil, rx, ry, ex, ey int) {})

}

func (s *Slider) SetMaxValue(v float64) {
	s.maxValue = v
	s.scaler = 1.0 / s.viewWidth
	s.ProgressBar.SetDisplayScale(v)
}

func (s *Slider) drawBar(w *x11ui.Window, x, y int) {
	v := float64(x) * s.scaler
	if s.stepSize > 0 {
		pixelunitStep := s.viewWidth / (s.maxValue / s.stepSize)
		snappedX := math.Ceil(float64(x) / pixelunitStep)
		v = snappedX * pixelunitStep / s.viewWidth
		v = TrimValues(v)
	}
	//mousebind.GrabPointer(xu, win, confine, cursor)
	s.ProgressBar.SetValue(v)

}
