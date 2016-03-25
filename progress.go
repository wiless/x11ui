package x11ui

import (
	"fmt"
	"image"
	"image/color"
	"log"

	"github.com/llgcode/draw2d/draw2dkit"

	"github.com/BurntSushi/xgbutil/xgraphics"
	"github.com/llgcode/draw2d/draw2dimg"
	"github.com/lucasb-eyer/go-colorful"
)

type ProgressBar struct {
	p  *Window //ParentWindow
	me *Window

	// Custom properties
	scale     float64
	val       float64
	barColor  color.Color
	txtColor  color.Color
	fmtString string
}

func NewProgressBar(title string, p *Window, dims ...int) *ProgressBar {
	if p == nil {
		log.Fatal("Cannot Create Widget without Application")
	}
	pbar := new(ProgressBar)
	pbar.me = NewWidget(p.Window.X, p, title, dims...)

	pbar.ResetFmtString()
	pbar.loadTheme()
	// pbar.SetValue(0.5)
	return pbar
}

func (p *ProgressBar) loadTheme() {
	p.barColor = colorful.LinearRgb(.4, .6, .1)
	p.txtColor = color.RGBA{20, 200, 30, 200}
}

func (p *ProgressBar) SetBarColor(bc color.Color) {
	p.barColor = bc
}
func (p *ProgressBar) SetTextColor(tc color.Color) {
	p.txtColor = tc
}

func (p *ProgressBar) SetFmtString(s string) {
	p.fmtString = s
}
func (p *ProgressBar) FmtString() string {
	return p.fmtString
}
func (p *ProgressBar) ResetFmtString() {
	p.fmtString = "%3.0f%%"
}

// return the value in the scale 0 to 1
func (p ProgressBar) Value() float64 {
	return p.val
}

// Sets the value in the fraction 0 to 1
func (p *ProgressBar) SetValue(v float64) {
	p.val = v
	p.reDrawBar()

}

func (p *ProgressBar) reDrawBar() {
	/// ignoring widgetState
	p.drawBackground(StateNormal)

}

func (p *ProgressBar) drawBackground(s WidgetState) {
	r := p.me.Rect
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
	draw2dkit.Rectangle(gc, 0, 0, WW, HH)
	gc.FillStroke()

	// Draw the BAR
	ww := float64(r.Width) * p.val
	margin := 2.0
	ww = ww - margin
	gc.SetLineWidth(0)
	gc.SetFillColor(p.barColor)
	draw2dkit.Rectangle(gc, margin, margin, ww, HH-1*margin)
	gc.FillStroke()

	// c := color.RGBA{200, 200, 200, 255}
	// gc.SetStrokeColor(image.White)
	// c := color.RGBA{200, 200, 200, 255}
	gc.SetFillColor(p.txtColor)
	gc.SetLineWidth(1)
	cx, cy := r.Center()
	str := fmt.Sprintf(p.fmtString, (p.val * 100.0))
	x0, y0, w0, h0 := gc.GetStringBounds(str)
	// log.Println("Required dimension for string ", x0, y0, w0, h0, cx, cy)
	tx, ty := float64(cx)-w0/2.0-x0, float64(cy)-h0/2.0-y0/2.0
	// gc.StrokeStringAt(str, tx, ty)
	gc.FillStringAt(str, tx, ty)
	// gc.FillStroke()
	gc.Close()
	g := xgraphics.NewConvert(p.me.X(), dest)

	// w.drawLabel(g, w.title)
	g.XSurfaceSet(p.me.Id)
	g.XDraw()
	g.XPaintRects(p.me.Id, r.ImageRect())
	// return g
	// return g
}
