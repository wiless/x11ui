package main

import (
	"fmt"
	"image"
	"image/color"
	"log"
	"math"

	"github.com/wiless/vlib"

	"github.com/llgcode/draw2d"
	"github.com/llgcode/draw2d/draw2dimg"
	"github.com/wiless/x11ui"
)

func main() {
	x11ui.SetResourcePath("../../fonts/")
	app := x11ui.NewApplication("Test Layout", 300, 500, false, false)

	app.AutoLayout(x11ui.LayoutVer, 10, 10, 300, 100)
	t := app.NewChildWindow("Top View")
	app.NewChildWindow("Next View")

	var tr TopRegion
	tr.Region = x11ui.NewRegion(x11ui.Rect{0, 0, t.Width, t.Height})
	tr.Region.BG = color.Black
	tr.Region.TC = color.RGBA{255, 0, 0, 255}
	tr.Region.Margin = 13
	layout := x11ui.NewLayout(t, 0, 0)
	tr.Current = 34.01
	tr.Coulomb = -10.1
	layout.AddRegion(tr)
	layout.DrawOnWindow(t)

	app.Show()
}

type TopRegion struct {
	x11ui.Region
	Current    float64
	Coulomb    float64
	PvsCurrent float64
	Voltage    float64
}

func (t *TopRegion) Update(voltage, current, coulomb float64) {
	t.PvsCurrent = t.Current
	t.Current = current
	t.Voltage = voltage
	t.Coulomb = coulomb
}

func (t TopRegion) GetRegion() *x11ui.Region {
	return &t.Region
}

func (t TopRegion) PaintRegion() *image.RGBA {
	r := t.Rect
	r.MoveTo(0, 0)
	r.ImageRect()
	img := image.NewRGBA(r.ImageRect())

	// thunderimg, err := draw2dimg.LoadFromPngFile("res/charge.png")
	// fd, _ := os.Open("res/charge.png")
	// cimage, err := png.Decode(fd)
	cimage, err := draw2dimg.LoadFromPngFile("res/thunder.png")

	if err != nil {
		log.Print(err)
		return nil
	}
	gc := draw2dimg.NewGraphicContext(img)

	dw, dh := float64(r.Width), float64(r.Height)
	// Size of source image
	sw, sh := float64(cimage.Bounds().Dx()), float64(cimage.Bounds().Dy())
	// Draw image to fit in the frame
	// TODO Seems to have a transform bug here on draw image
	scale := math.Min((dw-0*2)/sw, (dh-0*2)/sh)

	gc.Save()
	// gc.Translate((dw-sw*scale)/2, (dh-sh*scale)/2)
	gc.Scale(scale, scale)
	// gc.Rotate(0.2)

	gc.DrawImage(cimage)
	// gc.Scale(1, 1)
	gc.Restore()

	gc.SetFillColor(t.BG)
	gc.SetStrokeColor(t.FG)

	WW := float64(r.Width)
	HH := float64(r.Height)

	// Draw Background
	// gc.SetLineWidth(0)
	// gc.SetFillColor(t.BG)
	// gc.SetStrokeColor(t.FG)
	// draw2dkit.Rectangle(gc, 0, 0, WW, HH)
	// gc.FillStroke()

	// Show Pvs Current
	gc.SetFontSize(12)
	ft := gc.GetFontData()
	ft.Style = draw2d.FontStyleNormal
	gc.SetFontData(ft)
	gc.SetFillColor(color.RGBA{0x05, 0x8D, 0xBA, 0xff})
	str := fmt.Sprintf("%+ 8.2f", t.PvsCurrent)
	x0, y0, tw, th := gc.GetStringBounds(str)
	log.Println("Current", x0, y0, tw, th)
	oldtw, oldth := tw, th
	x, y := WW-t.Margin-tw, t.Margin-th
	gc.FillStringAt(str, x, y)

	// Show Current at 0,0
	gc.SetFontSize(34)
	ft = gc.GetFontData()
	ft.Style = draw2d.FontStyleBold
	gc.SetFontData(ft)
	gc.SetFillColor(t.TC)
	str = fmt.Sprintf("%+ 8.2f", t.Current)
	x0, y0, tw, th = gc.GetStringBounds(str)
	log.Println("Current", x0, y0, tw, th)
	oldtw, oldth = tw, th
	x, y = WW-t.Margin-tw, HH/2-th/2
	gc.FillStringAt(str, x, y)

	// Show Units
	gc.SetFontSize(10)
	// Light Color
	gc.SetFillColor(color.RGBA{0x05, 0x8D, 0xBA, 0xff})
	str = fmt.Sprintf("mA")
	x0, y0, tw, th = gc.GetStringBounds(str)
	log.Println("Bounds mA", x0, y0, tw, th)
	x += (oldtw - tw)
	y += (oldth + th - y0)
	oldtw, oldth = tw, th
	gc.FillStringAt(str, x, y)

	//Show Columns Unit
	gc.SetFontSize(9)
	ft.Style = draw2d.FontStyleNormal
	gc.SetFontData(ft)
	gc.SetFillColor(color.RGBA{0x05, 0x8D, 0xBA, 0xff})
	str = fmt.Sprintf("  Columns")
	x0, y0, tw, th = gc.GetStringBounds(str)
	x, y = t.Margin-x0, HH-t.Margin-th-y0
	oldtw, oldth = tw, th
	gc.FillStringAt(str, x, y)

	// Show Column Value
	gc.SetFontSize(15)
	gc.SetFillColor(t.TC)
	str = fmt.Sprintf("%+ 8.2f", t.Coulomb)
	x0, y0, tw, th = gc.GetStringBounds(str)
	x, y = t.Margin-x0, HH-t.Margin-th-oldth+y0
	gc.FillStringAt(str, x, y)

	vec := vlib.RandUFVec(40)
	rec := x11ui.Rect{r.Width / 4, 50, 200, 50}
	MiniGraph(gc, vec, rec, 0, 1)

	gc.Close()
	return img
}

// GetRegion() *Region
// SetRegion(*Region)
// PaintRegion() *image.RGBA

func MiniGraph(gc *draw2dimg.GraphicContext, y []float64, r x11ui.Rect, ymin, ymax float64) {

	gc.Save()
	scale := float64(r.Height) / (ymax - ymin)
	xstep := float64(r.Width) / float64(len(y))
	offset := ymin
	log.Println(y, scale, xstep, offset)
	// scale := math.Min((dw-0*2)/(len(y)*10), (dh-0*2)/sh)
	// gc.Translate((dw-sw*scale)/2, (dh-sh*scale)/2)
	// gc.Scale(1, scale)
	// gc.SetFillColor(color.RGBA{0xF0, 0x0, 0, 0x0})
	gc.SetStrokeColor(color.RGBA{0x44, 0xF4, 0x44, 0xff})
	gc.SetLineWidth(2)
	// gc.Rotate(math.Pi / 4.0)
	gc.Translate(float64(r.X), float64(r.Y))
	// Draw a closed shape
	gc.MoveTo(0, 0) // should always be called first for a new path
	x := 0.0
	for _, val := range y {
		// y := (math.Sin(2*math.Pi*5*float64(x)/float64(r.Max.X)) + rand.Float64()*.5) * (float64(r.Max.Y) / 2.0)
		gc.LineTo(x, val*scale)
		x += xstep
	}
	// gc.Close()
	gc.Stroke()

	gc.Restore()
}
