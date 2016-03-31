package x11ui

import (
	"image"
	"image/color"

	"github.com/BurntSushi/xgbutil/xgraphics"
	"github.com/BurntSushi/xgbutil/xrect"
	"github.com/llgcode/draw2d/draw2dimg"
	"github.com/llgcode/draw2d/draw2dkit"
)

type Rect struct {
	X, Y          int
	Width, Height int
}

func (r *Rect) MoveTo(x, y int) *Rect {
	r.X = x
	r.Y = y
	return r
}

func (r *Rect) ShiftBy(dx, dy int) *Rect {
	r.ShiftRight(dx)
	r.ShiftDown(dy)
	return r
}

func (r *Rect) ShiftRight(dx int) *Rect {
	r.X += dx
	return r

}
func (r *Rect) ShiftDown(dy int) *Rect {
	r.Y += dy
	return r
}
func (r *Rect) ReSize(w, h int) *Rect {
	r.Width, r.Height = w, h
	return r
}
func (r *Rect) Grow(dw, dh int) *Rect {
	r.Width += dw
	r.Height += dh

	return r
}

func (r *Rect) array() []int {
	result := []int{r.X, r.Y, r.Width, r.Height}
	return result
}

func newRect(dims ...int) Rect {
	r := Rect{0, 0, 100, 50}
	for i, v := range dims {
		switch i {
		case 0:
			r.X = v
		case 1:
			r.Y = v
		case 2:
			r.Width = v
		case 3:
			r.Height = v
		}
	}
	return r
}

func (r *Rect) CenterX() int {
	return r.X + r.Width/2
}

func (r *Rect) CenterY() int {
	return r.Y + r.Height/2
}

func (r *Rect) Center() (x, y int) {
	return r.CenterX(), r.CenterY()
}

func (r *Rect) ImageRect() image.Rectangle {
	return image.Rect(0, 0, r.Width, r.Height)

}

func XRectToImageRect(r xrect.Rect) image.Rectangle {
	return image.Rect(r.X(), r.Y(), r.Width(), r.Height())
}

func DrawDummy(w *Window, s WidgetState) {
	r := w.Rect
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
	gc.SetFillColor(color.RGBA{130, 120, 30, 10})
	gc.SetStrokeColor(color.RGBA{30, 120, 130, 10})
	draw2dkit.Rectangle(gc, 0, 0, WW, HH)
	gc.FillStroke()

	gc.FillStroke()
	gc.SetFillColor(color.RGBA{130, 120, 30, 10})
	gc.SetStrokeColor(color.RGBA{30, 120, 130, 10})
	draw2dkit.Circle(gc, 250, 10, 30)
	gc.FillStroke()

	gc.Close()
	g := xgraphics.NewConvert(w.X(), dest)

	// w.drawLabel(g, w.title)
	g.XSurfaceSet(w.Id)
	g.XDraw()
	g.XPaintRects(w.Id, r.ImageRect())
	// return g
	// return g
}
