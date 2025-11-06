package x11ui

import (
	"fmt"
	"image"
	"image/color"
	"log"

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
	c.cb.BG = CurrentTheme.CheckboxUncheckedColor
	c.cb.TC = CurrentTheme.TextColor

	c.txttb.Region = NewRegion(Rect{0, 0, 150, 50})
	c.txttb.SetText("Transmit Pilot ?")
	c.txttb.BG = CurrentTheme.BackgroundColor
	c.txttb.TC = CurrentTheme.TextColor

	c.cb.filename = "hg.png"
	c.txttb.filename = "Well_003.png"
	midpoint := image.Point{c.Width()/2 - 75, c.Height()/2 - 25}
	c.Layout = CreateLayout(midpoint.X, midpoint.Y, 150, 50)
	c.AddRegionAt(c.cb, 0, 0).AddRegionAt(c.txttb, c.cb.Width, 0)
	c.HoverFn = c.onHover
	c.HandlerFunctions.LeaveFn = c.onLeave
	c.HandlerFunctions.ClkFn = c.onClick
	c.updateColorsBasedOnState()
	// c.Layout=new LayoutNewLayout(w, x0, y0)

}

func (c *CheckBox) onClick() {
	c.state = !c.state
	c.updateColorsBasedOnState()
	c.RePaint()
}

func (c *CheckBox) updateColorsBasedOnState() {
	if c.state {
		c.cb.BG = CurrentTheme.CheckboxCheckedColor
	} else {
		c.cb.BG = CurrentTheme.CheckboxUncheckedColor
	}
}

func (c *CheckBox) SetChecked(checked bool) {
	c.state = checked
	c.updateColorsBasedOnState()
	c.RePaint()
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
	BG         color.Color // Change to color.Color
	TC         color.Color // Change to color.Color
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

		}  
		return c.canvas.At(x, y).(xgraphics.BGRA)

	})
	c.canvas.XDraw()
	c.canvas.XPaint(c.xwin.Id)

}

func (c *CheckBox) onLeave() {
	fmt.Printf("Leaving ")
	gc := c.gc
	gc.SetFillColor(CurrentTheme.BackgroundColor)
	gc.SetStrokeColor(CurrentTheme.LineColor)
	draw2dkit.Rectangle(gc, 0, 0, 50, 30)
	gc.FillStroke()
	gc.Close()
}
func (c *CheckBox) onHover() {
	fmt.Printf("Entering")
	r := GetIRect(50, 30)
	gc := c.gc
	gc.SetFillColor(CurrentTheme.ForegroundColor)
	gc.SetStrokeColor(CurrentTheme.LineColor)
	draw2dkit.Circle(gc, 20, 20, 20)
	gc.FillStroke()
	gc.Close()

	c.CopyPaste(r)
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

	StrokeBorderImg(img, CurrentTheme.CheckboxBorderColor, 0, 4)

	gc.Close()
	return img
}
