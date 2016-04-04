package x11ui

import (
	"fmt"
	"image"
	"image/color"
	"log"

	"github.com/BurntSushi/xgb/xproto"
	"github.com/BurntSushi/xgbutil/xgraphics"

	"github.com/llgcode/draw2d/draw2dkit"

	"github.com/llgcode/draw2d/draw2dimg"
)

type CheckBox struct {
	*Widget

	// Custom properties
	cb, txttb boxRegion
	label     string
	state     bool
}

func NewCheckBox(title string, p *Window, dims ...int) *CheckBox {
	if p == nil {
		log.Fatal("Cannot Create Widget without Application")
	}
	cbox := new(CheckBox)
	cbox.Widget = WidgetFactory(p, dims...)
	cbox.init()

	// tbox.Create(p, dims...)
	// tbox.loadTheme()
	// pbar.SetValue(0.5)
	return cbox
}

func (c *CheckBox) init() {
	c.state = false
	c.title = "Check Box"
	// c.drawBackground()
	// c.drawTextBox(StateNormal)
	// c.AddRulers()
	// c.updateCanvas()
	// c.registerHandlers()

	c.cb.Region = NewRegion(Rect{0, 0, 50, 50})
	c.cb.BG = color.RGBA{155, 0, 0, 255}
	c.cb.TC = color.RGBA{0, 255, 0, 255}

	c.txttb.Region = NewRegion(Rect{0, 0, 150, 50})
	c.txttb.SetText("Transmit Pilot ?")
	c.txttb.BG = color.RGBA{155, 0, 0, 255}
	c.txttb.TC = color.RGBA{0, 255, 0, 255}

	c.cb.filename = "hg.png"
	c.txttb.filename = "Well_003.png"
	// c.Layout.ox = midpoint.X
	// c.Layout.oy = midpoint.Y
	midpoint := image.Point{c.Width()/2 - 75, c.Height()/2 - 25}
	// log.Println("midpoint", midpoint, r.CenterX(), r.CenterY())
	_, _ = midpoint.X, midpoint.Y
	c.Layout = CreateLayout(midpoint.X, midpoint.Y, 150, 50)
	// c.AddRegionAt(tr, 0, 0+midpoint.Y).AddRegionAt(big, tr.Width, 0+midpoint.Y)
	c.AddRegionAt(c.cb, 0, 0).AddRegionAt(c.txttb, c.cb.Width, 0)
	// c.xwin.Detach()
	c.HoverFn = c.onHover
	c.HandlerFunctions.LeaveFn = c.onLeave
	c.HandlerFunctions.ClkFn = c.onClick
	// c.Layout=new LayoutNewLayout(w, x0, y0)

}

func (c *CheckBox) onClick() {

	ximg, err := xgraphics.NewDrawable(c.xu, xproto.Drawable(c.xwin.Id))
	ximg.SavePng("something.png")

	log.Println("File saved", err)
}

type boxRegion struct {
	Region
	Current    float64
	Coulomb    float64
	PvsCurrent float64
	Voltage    float64
	state      bool
	filename   string
	Caption    string
}

func (b boxRegion) GetRegion() *Region {
	return &b.Region
}

func (b *boxRegion) SetText(s string) {
	b.Caption = s
}

func (c *CheckBox) CopyPaste(r image.Rectangle) {
	c.canvas.For(func(x, y int) xgraphics.BGRA {

		pt := image.Point{x, y}
		inside := pt.In(r)

		if inside {
			// log.Println("Found true ", c.rawimg.At(x, y))
			// bg := xgraphics.BGRA{255, 20, 0, 255}
			bg := toBGRA(c.rawimg.At(x, y))
			return bg

		} else {
			return c.canvas.At(x, y).(xgraphics.BGRA)
		}

	})
	c.canvas.XDraw()
	c.canvas.XPaint(c.xwin.Id)

}

func (c *CheckBox) onLeave() {
	fmt.Printf("Leaving ")
	// r := GetIRect(50, 30)
	// img := image.NewRGBA(r)
	// gc := draw2dimg.NewGraphicContext(img)
	// irect := GetIRect(50, 50)
	gc := c.gc
	gc.SetFillColor(color.RGBA{0, 255, 10, 0})
	gc.SetStrokeColor(color.RGBA{20, 20, 255, 100})
	draw2dkit.Rectangle(gc, 0, 0, 50, 30)
	gc.FillStroke()
	gc.Close()
	// c.CopyPaste(r)
	// xgraphics.Blend(c.canvas, img, image.Point{0, 0})

	// xg := xgraphics.NewConvert(c.xu, img)
	// xg.
	// xg.XDraw()
	// xg.XPaintRects(c.xwin.Id, irect)

	// for x := 0; x < 30; x++ {
	// 	c.canvas.Set(x, x, color.RGBA{255, 0, 0, 255})
	// }
	// c.canvas.XDraw()
	// c.canvas.XPaint(c.xwin.Id)
}
func (c *CheckBox) onHover() {
	fmt.Printf("Entering")
	r := GetIRect(50, 30)
	// img := image.NewRGBA(r)
	// gc := draw2dimg.NewGraphicContext(img)
	gc := c.gc
	gc.SetFillColor(color.RGBA{0, 0, 0, 0})
	gc.SetStrokeColor(color.RGBA{245, 0, 0, 0})
	// draw2dkit.Rectangle(gc, 0, 0, 50, 30)
	// gc.FillStroke()
	draw2dkit.Circle(gc, 20, 20, 20)
	gc.FillStroke()
	gc.Close()

	// c.canvas.XDraw()
	// c.canvas.XPaint(c.xwin.Id)

	c.CopyPaste(r)

	// for x := 0; x < 30; x++ {
	// 	c.canvas.Set(x, x, color.RGBA{255, 0, 0, 255})
	// }
	// c.canvas.XDraw()
	// c.canvas.XPaint(c.xwin.Id)
}

func (b boxRegion) PaintRegion() *image.RGBA {
	r := b.Rect
	r.MoveTo(0, 0)

	iconsize := GetIRect(b.Width, b.Height)
	img := image.NewRGBA(iconsize)
	log.Println("Region ", b.Caption, r)
	// thunderimg, err := draw2dimg.LoadFromPngFile("res/charge.png")
	// fd, _ := os.Open("res/charge.png")
	// cimage, err := png.Decode(fd)

	// not good
	gc := draw2dimg.NewGraphicContext(img)

	if b.filename != "" {
		icon, err := draw2dimg.LoadFromPngFile(b.filename)

		if err != nil {
			log.Print(err)
			return nil
		}

		dw, dh := float64(r.Width), float64(r.Height)
		// Size of source image
		sw, sh := float64(icon.Bounds().Dx()), float64(icon.Bounds().Dy())
		// Draw image to fit in the frame
		// TODO Seems to have a transform bug here on draw image
		// scale := math.Min((dw-0*2)/sw, (dh-0*2)/sh)
		scalex, scaley := (dw-0*2)/sw, (dh-0*2)/sh

		gc.Save()
		// gc.Translate((dw-sw*scale)/2, (dh-sh*scale)/2)
		gc.Scale(scalex, scaley)
		// gc.Rotate(0.2)

		gc.DrawImage(icon)
		// gc.Scale(1, 1)
		gc.Restore()

		gc.SetFillColor(b.BG)
		gc.SetStrokeColor(b.FG)
	}

	if b.Caption != "" {
		// w, h := xgraphics.TextMaxExtents(systemFont, 12, b.Caption)
		x0, y0, w, h := gc.GetStringBounds(b.Caption)
		px := float64(r.Width/2) - w/2
		py := float64(r.Height/2) - h/2
		log.Println(x0, y0)
		gc.StrokeStringAt(b.Caption, float64(px), float64(py))
	}

	StrokeBorderImg(img, color.RGBA{200, 0, 0, 255}, 0, 4)
	// WW := float64(r.Width)
	// HH := float64(r.Height)

	// Draw Background
	// gc.SetLineWidth(1)
	// // gc.SetFillColor(b.BG)
	// gc.SetStrokeColor(b.FG)
	// draw2dkit.Rectangle(gc, 0, 0, WW, HH)
	// gc.Stroke()

	// Show Pvs Current
	// gc.SetFontSize(12)
	// ft := gc.GetFontData()
	// ft.Style = draw2d.FontStyleNormal
	// gc.SetFontData(ft)
	// gc.SetFillColor(color.RGBA{0x05, 0x8D, 0xBA, 0xff})
	// str := fmt.Sprintf("%+ 8.2f", b.PvsCurrent)
	// x0, y0, tw, th := gc.GetStringBounds(str)
	// log.Println("Current", x0, y0, tw, th)
	// oldtw, oldth := tw, th
	// x, y := WW-b.Margin-tw, b.Margin-th
	// gc.FillStringAt(str, x, y)

	// // Show Current at 0,0
	// gc.SetFontSize(34)
	// ft = gc.GetFontData()
	// ft.Style = draw2d.FontStyleBold
	// gc.SetFontData(ft)
	// gc.SetFillColor(b.TC)
	// str = fmt.Sprintf("%+ 8.2f", b.Current)
	// x0, y0, tw, th = gc.GetStringBounds(str)
	// log.Println("Current", x0, y0, tw, th)
	// oldtw, oldth = tw, th
	// x, y = WW-b.Margin-tw, HH/2-th/2
	// gc.FillStringAt(str, x, y)

	// // Show Units
	// gc.SetFontSize(10)
	// // Light Color
	// gc.SetFillColor(color.RGBA{0x05, 0x8D, 0xBA, 0xff})
	// str = fmt.Sprintf("mA")
	// x0, y0, tw, th = gc.GetStringBounds(str)
	// log.Println("Bounds mA", x0, y0, tw, th)
	// x += (oldtw - tw)
	// y += (oldth + th - y0)
	// oldtw, oldth = tw, th
	// gc.FillStringAt(str, x, y)

	// //Show Columns Unit
	// gc.SetFontSize(9)
	// ft.Style = draw2d.FontStyleNormal
	// gc.SetFontData(ft)
	// gc.SetFillColor(color.RGBA{0x05, 0x8D, 0xBA, 0xff})
	// str = fmt.Sprintf("  Columns")
	// x0, y0, tw, th = gc.GetStringBounds(str)
	// x, y = b.Margin-x0, HH-b.Margin-th-y0
	// oldtw, oldth = tw, th
	// gc.FillStringAt(str, x, y)

	// // Show Column Value
	// gc.SetFontSize(15)
	// gc.SetFillColor(b.TC)
	// str = fmt.Sprintf("%+ 8.2f", b.Coulomb)
	// x0, y0, tw, th = gc.GetStringBounds(str)
	// x, y = b.Margin-x0, HH-b.Margin-th-oldth+y0
	// gc.FillStringAt(str, x, y)

	// vec := vlib.RandUFVec(40)
	// rec := .Rect{r.Width / 4, 50, 200, 50}
	// MiniGraph(gc, vec, rec, 0, 1)

	gc.Close()
	return img
}
