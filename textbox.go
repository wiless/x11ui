package x11ui

import (
	"image"
	"image/color"
	"log"
	"math"
	"strings"

	"github.com/BurntSushi/xgbutil"
	"github.com/BurntSushi/xgbutil/keybind"
	"github.com/BurntSushi/xgbutil/xevent"
	"github.com/BurntSushi/xgbutil/xgraphics"
	"github.com/llgcode/draw2d/draw2dkit"
)

type TextBox struct {
	*Widget

	// Custom properties
	text         string
	lastx, lasty float64
	cursor       float64
	line         float64
	charSpace    float64
	linespace    float64
}

func NewTextBox(title string, p *Window, dims ...int) *TextBox {
	if p == nil {
		log.Fatal("Cannot Create Widget without Application")
	}
	tbox := new(TextBox)
	tbox.Widget = WidgetFactory(p, dims...)
	tbox.init()

	// tbox.Create(p, dims...)
	// tbox.loadTheme()
	// pbar.SetValue(0.5)
	return tbox
}

func (t *TextBox) registerHandlers() {
	xevent.KeyPressFun(t.keybHandler).Connect(t.xu, t.xwin.Id)
}
func (t *TextBox) keybHandler(X *xgbutil.XUtil, e xevent.KeyPressEvent) {

	modStr := keybind.ModifierString(e.State)
	keyStr := keybind.LookupString(X, e.State, e.Detail)

	modStr = strings.Replace(modStr, "lock-", "", -1)
	modStr = strings.Replace(modStr, "mod2", "", -1)

	if modStr != "" {
		// finalstr := keyStr
		// finalstr = fmt.Sprint(modStr, keyStr)
		log.Printf("\n 1 MODSTR %v : KEYSTR '%v' ", modStr, keyStr)
		t.handleKeyboard(keyStr)
	} else {
		// finalstr = fmt.Sprint(modStr, keyStr)
		log.Printf("\n 2 MODSTR %v : KEYSTR '%v' ", modStr, keyStr)
		t.handleKeyboard(keyStr)
	}
	// log.Println("Event code is ", e.Detail)
	// log.Printf("%s MAPS  to Keycode %v ", finalstr, keybind.StrToKeycodes(s.xu, finalstr))

	// if fn, ok := s.KeyMaps[finalstr]; ok {
	// 	if s.Debug {
	// 		log.Printf("Caught Key : %s", finalstr)
	// 	}
	// 	fn()
	// }

}

func (t *TextBox) AddRulers() {

	t.gc.SetStrokeColor(color.RGBA{30, 30, 100, 255})
	width := float64(t.Width())
	t.gc.SetLineDash([]float64{10, 5, 10, 5}, 0)
	Nlines := float64(t.Height()) / t.linespace
	for i := 0.0; i < Nlines; i++ {
		t.gc.MoveTo(0, i*t.linespace)
		t.gc.LineTo(width, i*t.linespace)
		t.gc.Stroke()
	}
	t.gc.SetLineDash([]float64{}, 0)
}
func (t *TextBox) handleKeyboard(str string) {
	if str == "Return" {
		t.line += t.linespace
		t.cursor = 0
		return

	}

	if len(str) != 1 {
		log.Println("I am returning")
		return
	}
	if str == " " {
		t.cursor += 10 // should be set to char space
		return
	}
	log.Println("Width", t.cursor)
	_, _, w, _ := t.gc.GetStringBounds(str)
	log.Println("Width", t.cursor, w)
	t.gc.SetStrokeColor(t.txtColor)
	t.gc.FillStringAt(str, t.cursor, t.line+t.linespace)
	t.cursor += w
	if t.cursor > (float64(t.Width()) - t.margin) {
		t.line += t.linespace // line spacing
		t.cursor = 0
	}
	t.gc.Close()

	t.updateCanvas()

}

func backgroundFn(x, y int) xgraphics.BGRA {

	r, g, b, a := systemBG.RGBA()
	return xgraphics.BGRA{uint8(b), uint8(g), uint8(r), uint8(a)}
}

func (t *TextBox) drawBackground() {

	irect := image.Rectangle{image.Point{0, 0}, image.Point{t.Width(), t.Height()}}
	subimg := t.canvas.SubImage(irect).(*xgraphics.Image)
	subimg.For(backgroundFn)
	// g := xgraphics.NewConvert(t.xu, img)
	// g.XDraw()
	// g.XPaint(t.xwin.Id)

	// xgraphics.Blend(si, simg, image.Point{0, 0})
	subimg.XDraw()
	subimg.XPaint(t.xwin.Id)
}

func (t *TextBox) init() {
	t.drawBackground()
	var x0, y0 float64
	x0, y0, t.charSpace, t.linespace = t.gc.GetStringBounds("W")

	log.Println("spacing ", t.charSpace, t.linespace, x0, y0)
	t.linespace = math.Abs(y0) + t.linespace
	t.linespace *= 2

	t.drawTextBox(StateNormal)

	t.AddRulers()

	t.updateCanvas()
	t.registerHandlers()

}

func (t *TextBox) drawTextBox(s WidgetState) {
	t.gc.SetFontSize(15)
	// r.MoveTo(0, 0)
	// r.ImageRect()
	W, H := float64(t.Width()), float64(t.Height())
	gc := t.Context()
	gc.SetFillColor(t.fgColor)
	gc.SetStrokeColor(t.lineColor)
	draw2dkit.Rectangle(gc, 0, 0, W, H)
	gc.FillStroke()
	gc.Close()

	// t.updateCanvas()
	// // bg := colorful.LinearRgb(.025, .025, .025)
	// switch s {
	// case StateNormal, StateReleased:
	// 	gc.SetFillColor(color.RGBA{0x20, 0x20, 0x20, 20})
	// 	gc.SetStrokeColor(systemFG)
	// case StateHovered:
	// 	gc.SetFillColor(color.RGBA{0x35, 0x20, 0x20, 20})
	// 	gc.SetStrokeColor(systemFG)
	// case StatePressed:
	// 	gc.SetFillColor(color.RGBA{0x20, 0x30, 0x20, 20})
	// 	gc.SetStrokeColor(systemFG)
	// case StateSpecial:
	// 	gc.SetFillColor(color.RGBA{0x20, 0x80, 0x20, 0x80})
	// 	gc.SetStrokeColor(systemFG)
	// }

}
