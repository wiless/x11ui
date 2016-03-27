package x11ui

import (
	"fmt"
	"image"
	"image/color"
	"log"

	"github.com/llgcode/draw2d/draw2dkit"

	"github.com/BurntSushi/xgbutil"
	"github.com/BurntSushi/xgbutil/xgraphics"
	"github.com/llgcode/draw2d/draw2dimg"
	"github.com/lucasb-eyer/go-colorful"
)

type ProgressBar struct {
	p  *Window //ParentWindow
	me *Window

	// Custom properties
	dispScaler float64
	val        float64
	bgColor    color.Color
	barColor   color.Color
	txtColor   color.Color
	fmtString  string
	margin     float64
	border     float64
}

func NewProgressBar(title string, p *Window, dims ...int) *ProgressBar {
	if p == nil {
		log.Fatal("Cannot Create Widget without Application")
	}
	pbar := new(ProgressBar)
	pbar.me = NewWidget(p.Window.X, p, title, dims...)
	pbar.SetDisplayScale(100.0)
	pbar.ResetFmtString()
	pbar.loadTheme()
	// pbar.SetValue(0.5)
	return pbar
}

func (p *ProgressBar) X() *xgbutil.XUtil {
	return p.me.X()
}

func (p *ProgressBar) loadTheme() {
	p.bgColor = systemBG
	p.barColor = colorful.LinearRgb(.4, .6, .1)
	p.txtColor = color.RGBA{20, 200, 30, 200}
}

func (p *ProgressBar) SetBackGroundColor(bc color.Color) {
	p.bgColor = bc
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

func (p *ProgressBar) SetDisplayScale(s float64) {
	p.dispScaler = s
}

// Sets the value in the fraction 0 to 1
func (p *ProgressBar) SetValue(v float64) {
	p.val = v
	p.reDrawBar()

}

func (p *ProgressBar) Parent() *Window {
	return p.p
}

func (p *ProgressBar) Widget() *Window {
	return p.me
}
func (p *ProgressBar) reDrawBar() {
	/// ignoring widgetState
	p.drawBackground(StateNormal)

}

func (p *ProgressBar) SetBorderWidth(bw float64) {
	p.border = bw
}

func (p *ProgressBar) Margin() float64 {
	return p.margin
}
func (p *ProgressBar) SetMargin(m float64) {
	p.margin = m
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
	gc.SetLineWidth(p.border)
	gc.SetFillColor(p.bgColor)
	draw2dkit.Rectangle(gc, 0, 0, WW, HH)
	gc.FillStroke()

	if p.val > 0 {
		// Draw the BAR
		NSegments := 1.0
		ww := float64(r.Width) * p.val
		margin := p.margin
		ww = ww - margin
		hh := HH - margin

		gc.SetLineWidth(0)
		sw := (ww - (NSegments-1)*margin) / NSegments
		gc.SetFillColor(p.barColor)
		// for i := 0.0; i < NSegments; i++ {
		draw2dkit.Rectangle(gc, margin, margin, sw, hh)
		gc.FillStroke()

		/// Draw Grid Lines
		NGrids := 10.0
		GridWidth := (WW - 2*margin) / NGrids
		lc, ok := p.barColor.(color.RGBA)
		if ok {
			lc.A = 10
			gc.SetStrokeColor(lc)
		}
		xpos := margin
		gc.SetLineWidth(1)

		for i := 0.0; i < NGrids; i++ {
			if xpos < ww {
				gc.MoveTo(xpos, 0)
				gc.LineTo(xpos, hh)
				xpos += GridWidth
				gc.Close()
				gc.Stroke()
			}

		}

		// }
	}

	// gc.SetFillColor(color.RGBA{130, 120, 30, 110})
	// gc.SetStrokeColor(color.RGBA{30, 120, 130, 110})
	// draw2dkit.Circle(gc, 250, 10, 30)
	// gc.FillStroke()

	// c := color.RGBA{200, 200, 200, 255}
	// gc.SetStrokeColor(image.White)
	// c := color.RGBA{200, 200, 200, 255}

	gc.SetFillColor(p.txtColor)
	gc.SetLineWidth(0)
	cx, cy := r.Center()
	str := fmt.Sprintf(p.fmtString, (p.val * p.dispScaler))
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
