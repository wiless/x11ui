package x11ui

import (
	"image"
	"image/color"
	"log"
	"math"

	"github.com/llgcode/draw2d/draw2dimg"
)

type CheckBox struct {
	*Widget

	// Custom properties
	label string
	state bool
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
	var tr, big boxRegion
	tr.Region = NewRegion(Rect{0, 0, 50, 50})
	tr.BG = color.RGBA{155, 0, 0, 255}
	tr.TC = color.RGBA{0, 255, 0, 255}

	big.Region = NewRegion(Rect{50, 0, 50, 50})

	big.BG = color.RGBA{155, 0, 0, 255}
	big.TC = color.RGBA{0, 255, 0, 255}

	tr.filename = "hg.png"
	big.filename = "Well_003.png"

	c.AddRegionAt(tr, 10, 10).AddRegionAt(big, 100, 100)

	// c.Layout=new LayoutNewLayout(w, x0, y0)

}

type boxRegion struct {
	Region
	Current    float64
	Coulomb    float64
	PvsCurrent float64
	Voltage    float64
	state      bool
	filename   string
}

func (b boxRegion) GetRegion() *Region {
	return &b.Region
}

func (b boxRegion) PaintRegion() *image.RGBA {
	r := b.Rect
	r.MoveTo(0, 0)

	iconsize := GetIRect(b.Width, b.Width)
	img := image.NewRGBA(iconsize)

	// thunderimg, err := draw2dimg.LoadFromPngFile("res/charge.png")
	// fd, _ := os.Open("res/charge.png")
	// cimage, err := png.Decode(fd)

	// not good
	icon, err := draw2dimg.LoadFromPngFile(b.filename)

	if err != nil {
		log.Print(err)
		return nil
	}
	gc := draw2dimg.NewGraphicContext(img)

	dw, dh := float64(r.Width), float64(r.Height)
	// Size of source image
	sw, sh := float64(icon.Bounds().Dx()), float64(icon.Bounds().Dy())
	// Draw image to fit in the frame
	// TODO Seems to have a transform bug here on draw image
	scale := math.Min((dw-0*2)/sw, (dh-0*2)/sh)

	gc.Save()
	// gc.Translate((dw-sw*scale)/2, (dh-sh*scale)/2)
	gc.Scale(scale, scale)
	// gc.Rotate(0.2)

	gc.DrawImage(icon)
	// gc.Scale(1, 1)
	gc.Restore()

	gc.SetFillColor(b.BG)
	gc.SetStrokeColor(b.FG)

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
