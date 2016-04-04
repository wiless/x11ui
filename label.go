package x11ui

import (
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
	title      string
	autoresize bool
	align      AlignMode
}

func NewLabel(title string, p *Window, dims ...int) *Label {
	if p == nil {
		log.Fatal("Cannot Create Widget without Application")
	}
	lbl := new(Label)
	lbl.Widget = WidgetFactory(p, dims...)
	lbl.init()
	lbl.SetLabel(title)
	// tbox.Create(p, dims...)
	// tbox.loadTheme()
	// pbar.SetValue(0.5)
	return lbl
}

func (l *Label) AutoResize(auto bool) {
	l.autoresize = auto
}

func (l *Label) SetAlignMode(align AlignMode) {
	l.align = align
}

func (l *Label) init() {
	cw, ch := xgraphics.TextMaxExtents(systemFont, 12, l.title)
	log.Println("extends ", cw, ch)

	// t.AddRulers()
	l.updateCanvas()
	// go t.ShowIBeam()
	// l.registerHandlers()

}

func (l *Label) SetLabel(lbl string) {
	l.title = lbl
	if l.autoresize {
		w, h := xgraphics.Extents(systemFont, 12, l.title)
		l.Widget.xwin.Resize(w, h)
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
		tw, th := xgraphics.Extents(systemFont, 12, t.title)
		xpos, ypos = t.Width()-tw, t.Height()-th
		xpos, ypos = xpos/2, ypos/2
	}
	t.canvas.Text(xpos, ypos, t.txtColor, 12, systemFont, t.title)

	// t.canvas.XSurfaceSet(t.xwin.Id)
	t.updateCanvas()

}
