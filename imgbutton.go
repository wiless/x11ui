package x11ui

import (
	"bufio"
	"bytes"
	"image"
	"image/draw"
	"image/jpeg"
	"os"
	"strings"

	"github.com/BurntSushi/xgbutil/xevent"
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
	imageAlign   Alignment // New field for image alignment
}

func NewImgButton(title string, p *Window, dims ...int) *ImgButton {
	if p == nil {
		log.Fatal("Cannot Create Widget without Application")
	}
	tbox := new(ImgButton)
	tbox.Widget = WidgetFactory(p, dims...)
	// tbox.fname = "hg.png"
	tbox.init()
	xevent.KeyPressFun(tbox.keybHandler).Connect(tbox.xu, tbox.xwin.Id)

	// xevent.KeyPressFun(s.keybHandler).Connect(s.xu, s.AppWin().Id)

	// tbox.Create(p, dims...)
	// tbox.loadTheme()
	// pbar.SetValue(0.5)
	return tbox
}

func (t *ImgButton) ResizeWidget(w, h int) {
	t.Win().ReSize(w, h)
	t.Widget.Resize(0, 0, w, h)
	t.setupCanvas()

}

func (t *ImgButton) registerHandlers() {
	// xevent.KeyPressFun(t.keybHandler).Connect(t.xu, t.xwin.Id)
}

func (t *ImgButton) SetImageAlign(alignment Alignment) {
	t.imageAlign = alignment
	t.updateCanvas()
}

func (i *ImgButton) SetPicture(fname string) {
	i.fname = fname
	i.addPicture()

	i.updateCanvas()
}

func (t *ImgButton) DrawImage(img image.Image) {
	t.DrawImageEx(img, false)
}

func (t *ImgButton) DrawImageEx(img image.Image, eraseBg bool) {

	if eraseBg {
		t.drawBackground()
	}
	rgb := t.canvas

	// Calculate image position based on alignment
	imgWidth := img.Bounds().Dx()
	imgHeight := img.Bounds().Dy()
	btnWidth := t.Width()
	btnHeight := t.Height()

	var offsetX, offsetY int

	switch t.imageAlign {
	case AlignTopLeft:
		offsetX = 0
		offsetY = 0
	case AlignTopCenter:
		offsetX = (btnWidth - imgWidth) / 2
		offsetY = 0
	case AlignTopRight:
		offsetX = btnWidth - imgWidth
		offsetY = 0
	case AlignMiddleLeft:
		offsetX = 0
		offsetY = (btnHeight - imgHeight) / 2
	case AlignCenter:
		offsetX = (btnWidth - imgWidth) / 2
		offsetY = (btnHeight - imgHeight) / 2
	case AlignMiddleRight:
		offsetX = btnWidth - imgWidth
		offsetY = (btnHeight - imgHeight) / 2
	case AlignBottomLeft:
		offsetX = 0
		offsetY = btnHeight - imgHeight
	case AlignBottomCenter:
		offsetX = (btnWidth - imgWidth) / 2
		offsetY = btnHeight - imgHeight
	case AlignBottomRight:
		offsetX = btnWidth - imgWidth
		offsetY = btnHeight - imgHeight
	default: // Default to top-left if no alignment or unknown alignment
		offsetX = 0
		offsetY = 0
	}

	draw.Draw(rgb, rgb.Bounds().Bounds(), img, image.Point{-offsetX, -offsetY}, draw.Src)

	t.canvas.XDraw()
	t.canvas.XPaint(t.xwin.Id)
}

func (t *ImgButton) DrawJpg(jpgdata []byte) {
	r := bytes.NewReader(jpgdata)
	img, _ := jpeg.Decode(r)

	irect := image.Rectangle{image.Point{0, 0}, image.Point{t.Width(), t.Height()}}

	// inset := irect.Inset(2)
	// log.Println(irect, inset)

	mx := min(irect.Dx(), irect.Dy())
	simg := xgraphics.Scale(img, mx, mx)

	// log.Print(inset, irect)

	// si := t.canvas.SubImage(inset).(*xgraphics.Image)
	// xg := xgraphics.NewConvert(t.xu, )
	// xg.XDraw()
	// xg.XPaintRects(t.xwin.Id, inset)

	si := t.canvas.SubImage(irect).(*xgraphics.Image)
	xgraphics.Blend(si, simg, image.Point{0, 0})

	// si.CreatePixmap()
	// si.XDraw()
	// si.XPaint(t.xwin.Id)
	t.canvas.XSurfaceSet(t.xwin.Id)
	t.updateCanvas()

}
func min(x, y int) int {
	if x < y {
		return x
	} else {
		return y
	}
}

func (t *ImgButton) addPicture() {
	if t.fname == "" || t.fname == "none" || t.fname == "null" {
		log.Println("No image file specified")
		return
	}
	var img image.Image
	var err error
	if strings.HasSuffix(t.fname, ".png") {
		img, err = draw2dimg.LoadFromPngFile(t.fname)
		if err != nil {
			deBug("Background Image", err)

			return
		}
	}

	if strings.HasSuffix(t.fname, ".jpg") || strings.HasSuffix(t.fname, ".jpeg") {

		f, err := os.OpenFile(t.fname, 0, 0)
		if err != nil {
			deBug("Background Image", err)
			return
		}
		defer f.Close()
		b := bufio.NewReader(f)
		img, err = jpeg.Decode(b)
		if err != nil {
			deBug("Background Image", err)
			return
		}
	}

	irect := image.Rectangle{image.Point{0, 0}, image.Point{t.Width(), t.Height()}}

	inset := irect.Inset(2)

	mx := min(inset.Dx(), inset.Dy())
	simg := xgraphics.Scale(img, mx, mx)

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
