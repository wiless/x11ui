package x11ui

import (
	"image/color"
	"log"

	"github.com/BurntSushi/xgbutil/xgraphics"
)

type AlignMode int

const (
	AlignTopLeft AlignMode = iota
	AlignHVCenter
	AlignHCenter
	AlignVCenter
)

type Label struct {
	*Widget
	// title      string
	autoresize bool
	align      AlignMode
	fsize      float64
}

func NewLabel(title string, p *Window, dims ...int) *Label {
	if p == nil {
		log.Fatal("Cannot Create Widget without Application")
	}
	lbl := new(Label)
	lbl.Widget = WidgetFactory(p, dims...)
	lbl.init()
	lbl.SetAlignMode(AlignHVCenter)
	lbl.SetLabel(title)

	// tbox.Create(p, dims...)
	// tbox.loadTheme()
	// pbar.SetValue(0.5)
	return lbl
}

func (l *Label) AutoResize(auto bool) *Label {
	l.autoresize = auto
	l.SetLabel(l.title)
	return l
}

func (l *Label) SetAlignMode(align AlignMode) *Label {
	l.align = align
	l.SetLabel(l.title)
	return l
}

func (l *Label) SetFontSize(size float64) *Label {
	l.fsize = size
	l.SetLabel(l.title)
	return l
}

func (l *Label) init() {
	// xgraphics.TextMaxExtents(systemFont, 12, l.title)
	// log.Println("extends ", cw, ch)
	l.SetFontSize(12)
	// t.AddRulers()
	l.updateCanvas()
	// go t.ShowIBeam()
	// l.registerHandlers()

}

func (l *Label) SetBackground(c color.Color) {
	l.Widget.SetBackground(c)
	l.SetLabel(l.title)

}

func (l *Label) SetLabel(lbl string) {
	l.title = lbl
	if l.autoresize {
		w, h := xgraphics.Extents(systemFont, l.fsize, l.title)
		l.Widget.xwin.Resize(w+3, h+3)
	}
	l.updateLabel(StateNormal)
}

func (t *Label) updateLabel(state WidgetState) {
	// log.Println("updateing text")
	// W, H := float64(t.Width()), float64(t.Height())

	// // gc := t.Context()
	// // gc.SetFillColor(color.RGBA{0, 255, 200, 255})
	// // gc.SetStrokeColor(color.RGBA{255, 200, 0, 255})
	// // draw2dkit.Rectangle(gc, 10, 10, W, H)
	// // gc.StrokeStringAt(t.title, 0, 0)
	// // gc.FillStroke()
	// // gc.Close()
	t.drawBackground()
	var xpos, ypos int
	if t.align == AlignHVCenter {
		tw, th := xgraphics.Extents(systemFont, t.fsize, t.title)
		xpos, ypos = t.Width()-tw, t.Height()-th
		xpos, ypos = xpos/2, ypos/2
	}
	t.canvas.Text(xpos, ypos, t.txtColor, t.fsize, systemFont, t.title)

	// t.canvas.XSurfaceSet(t.xwin.Id)
	t.updateCanvas()

}
