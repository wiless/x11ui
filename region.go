package x11ui

import (
	"image"
	"image/color"

	"github.com/BurntSushi/xgbutil/xgraphics"
)

type RegionPainterFn func() *image.RGBA

type Region struct {
	Rect
	BG      color.Color
	FG      color.Color
	TC      color.Color
	Margin  float64
	PaintFn RegionPainterFn
}

func NewRegion(r Rect) Region {
	return Region{Rect: r, BG: systemBG, FG: systemFG, TC: systemFG}
}

type RegionPainter interface {
	GetRegion() *Region
	// SetRegion(*Region)
	PaintRegion() *image.RGBA
}

type Layout struct {
	regions []RegionPainter
	ox, oy  int
	w, h    int
}

func NewLayout(w *Window, x0, y0 int) *Layout {
	l := new(Layout)
	l.ox, l.oy = x0, y0
	l.w, l.h = w.Rect.Width, w.Rect.Height

	return l
}

func (l *Layout) AddRegion(r RegionPainter) {
	l.regions = append(l.regions, r)
}

func (l *Layout) DrawOnWindow(w *Window) {
	r := w.Rect

	pixmap := l.regions[0].PaintRegion()
	g := xgraphics.NewConvert(w.X(), pixmap)
	g.XSurfaceSet(w.Id)
	g.XDraw()
	g.XPaintRects(w.Id, r.ImageRect())
}
