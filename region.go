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
	offsets []image.Point
	ox, oy  int
	w, h    int
}

func NewLayout(w *Window, x0, y0 int) *Layout {
	l := new(Layout)
	l.ox, l.oy = x0, y0
	l.w, l.h = w.Rect.Width, w.Rect.Height

	return l
}
func CreateLayout(x, y, w, h int) *Layout {
	l := new(Layout)

	l.ox, l.oy, l.w, l.h = x, y, w, h
	return l
}

func (l *Layout) Resize(x, y, w, h int) {
	l.ox, l.oy, l.w, l.h = x, y, w, h
}

func (l *Layout) Size() image.Point {
	return image.Point{l.w, l.h}
}
func (l *Layout) AddRegion(r RegionPainter) *Layout {
	l.offsets = append(l.offsets, origin)
	l.regions = append(l.regions, r)
	return l
}

func (l *Layout) AddRegionAt(r RegionPainter, x, y int) *Layout {
	pt := image.Point{x, y}
	l.offsets = append(l.offsets, pt)
	l.regions = append(l.regions, r)
	return l
}
func (l *Layout) ImageRect() image.Rectangle {
	return image.Rectangle{origin, l.Size()}
}
func (l *Layout) Origin() image.Point {
	return image.Point{l.ox, l.oy}
}
func (l *Layout) DrawOnWindow(w *Window) {
	r := w.Rect

	pixmap := l.regions[0].PaintRegion()
	// r0:=
	// log.Printf("======== Regions ", l.regions, w)
	g := xgraphics.NewConvert(w.X(), pixmap)
	g.XSurfaceSet(w.Id)
	g.XDraw()
	g.XPaintRects(w.Id, r.ImageRect())
}

func (l *Layout) DrawOnWidget(w *Widget) {
	// r := w.Rect

	pixmap := l.regions[0].PaintRegion()
	// r0:=
	// log.Printf("======== Regions ", l.regions, w)
	g := xgraphics.NewConvert(w.xu, pixmap)
	g.XSurfaceSet(w.ID())
	g.XDraw()
	g.XPaintRects(w.ID(), w.ImageRect())
}

func (l *Layout) SetRegion(indx int, r RegionPainter) {
	if indx < len(l.regions) {
		l.regions[indx] = r
	}
}

func (l *Layout) SetRegion(indx int, r RegionPainter) {
	if indx < len(l.regions) {
		l.regions[indx] = r
	}
}
