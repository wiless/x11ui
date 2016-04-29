package x11ui

import (
	"image"

	"github.com/BurntSushi/xgbutil/xgraphics"
	"github.com/llgcode/draw2d/draw2dimg"

	"log"
)

type ImgButton struct {
	*Widget

	// Custom properties
	text         string
	lastx, lasty float64
	cursor       float64
	line         float64
	fname        string
}

func NewImgButton(title string, p *Window, dims ...int) *ImgButton {
	if p == nil {
		log.Fatal("Cannot Create Widget without Application")
	}
	tbox := new(ImgButton)
	tbox.Widget = WidgetFactory(p, dims...)
	tbox.fname = "hg.png"
	tbox.init()

	// tbox.Create(p, dims...)
	// tbox.loadTheme()
	// pbar.SetValue(0.5)
	return tbox
}

func (t *ImgButton) registerHandlers() {
	// xevent.KeyPressFun(t.keybHandler).Connect(t.xu, t.xwin.Id)
}

func (i *ImgButton) SetPicture(fname string) {
	i.fname = fname
	i.addPicture()

	i.updateCanvas()
}

func min(x, y int) int {
	if x < y {
		return x
	} else {
		return y
	}
}

func (t *ImgButton) addPicture() {

	img, err := draw2dimg.LoadFromPngFile(t.fname)
	if err != nil {
		deBug("Background Image", err)

		return
	}

	irect := image.Rectangle{image.Point{0, 0}, image.Point{t.Width(), t.Height()}}

	inset := irect.Inset(2)
	log.Println(irect, inset)

	mx := min(inset.Dx(), inset.Dy())
	simg := xgraphics.Scale(img, mx, mx)
	log.Print(inset, irect)

	// si := t.canvas.SubImage(inset).(*xgraphics.Image)
	// xg := xgraphics.NewConvert(t.xu, )
	// xg.XDraw()
	// xg.XPaintRects(t.xwin.Id, inset)

	si := t.canvas.SubImage(inset).(*xgraphics.Image)
	xgraphics.Blend(si, simg, image.Point{0, 0})
	// si.CreatePixmap()
	// si.XDraw()
	// si.XPaint(t.xwin.Id)
	t.canvas.XSurfaceSet(t.xwin.Id)
	t.updateCanvas()
}

func (t *ImgButton) init() {

	// t.drawBorder(StateNormal)
	// t.updateCanvas()
	t.registerHandlers()

}
