package x11ui

import (
	"image"
	"image/color"
	"log"
	"strings"
	"time"

	"github.com/BurntSushi/xgbutil"
	"github.com/BurntSushi/xgbutil/keybind"
	"github.com/BurntSushi/xgbutil/xevent"
	"github.com/BurntSushi/xgbutil/xgraphics"
	"github.com/llgcode/draw2d/draw2dkit"
)

type TextBox struct {
	*Widget

	// Custom properties
	text      string
	cursor    int
	line      int
	charSpace int
	linespace int
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

func (t *TextBox) ShowIBeam() {
	// for {
	// for x := t.cursor; x < t.cursor+10; x++ {
	var toggle bool
	c := time.Tick(500 * time.Millisecond)
	for _ = range c {
		for y := t.line; y < t.line+t.linespace; y++ {
			if toggle {
				t.canvas.SetBGRA(t.cursor, y, toBGRA(t.bgColor))
			} else {
				t.canvas.SetBGRA(t.cursor, y, toBGRA(color.RGBA{250, 0, 0, 255}))
			}
		}
		err := t.canvas.XDrawChecked()
		if err == nil {
			t.canvas.XDraw()
			t.canvas.XPaint(t.xwin.Id)
		}

		toggle = !toggle
	}

	// }
	// }

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
	t.gc.SetFillColor(color.RGBA{0, 0, 0, 0})
	t.gc.SetStrokeColor(color.RGBA{30, 30, 100, 0})
	width := float64(t.Width())
	t.gc.SetLineDash([]float64{10, 5, 10, 5}, 0)
	Nlines := t.Height() / t.linespace
	for i := 0; i < Nlines; i++ {
		t.gc.MoveTo(0, float64(i*t.linespace))
		t.gc.LineTo(width, float64(i*t.linespace))
		t.gc.Stroke()
	}
	t.gc.SetLineDash([]float64{}, 0)
}
func (t *TextBox) handleKeyboard(str string) {
	if str == "Return" || str == "KP_Enter" {
		t.line += t.linespace
		t.cursor = 0
		return

	}

	if str == "BackSpace" {
		t.cursor -= 14
		nx, _, _ := t.canvas.Text(int(t.cursor), int(t.line), t.txtColor, 12, systemFont, " ")

		t.cursor = nx
		t.updateCanvas()
		return
	}
	if len(str) != 1 {
		log.Println("I am returning")
		return
	}

	// if str == " " {
	// 	// t.cursor += t.charSpace // should be set to char space
	// 	return
	// }
	// log.Println("Width", t.cursor)
	// _, _, w, _ := t.gc.GetStringBounds(str)
	// log.Println("Width", t.cursor, w)
	// t.gc.SetStrokeColor(t.txtColor)
	// t.gc.SetFillColor(color.Black)
	// t.gc.FillStringAt(str, t.cursor, t.line+t.linespace)

	nx, ny, _ := t.canvas.Text(int(t.cursor), int(t.line), t.txtColor, 12, systemFont, str)
	log.Println("nx,ny", nx, ny)
	t.cursor = nx

	if t.cursor > (t.Width() - int(t.margin)) {
		t.line += t.linespace // line spacing
		t.cursor = 0
	}

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
	// t.drawBackground()
	// cw, ch := xgraphics.Extents(systemFont, 12, "W")
	// log.Println("extends ", cw, ch)
	cw, ch := xgraphics.TextMaxExtents(systemFont, 12, "W")
	log.Println("extends ", cw, ch)
	t.linespace = ch
	t.charSpace = cw

	t.drawTextBox(StateNormal)
	// t.AddRulers()
	t.updateCanvas()
	// go t.ShowIBeam()
	t.registerHandlers()

}

func (t *TextBox) drawTextBox(s WidgetState) {
	t.gc.SetFontSize(15)
	// r.MoveTo(0, 0)
	// r.ImageRect()
	W, H := float64(t.Width()), float64(t.Height())
	gc := t.Context()
	gc.SetFillColor(t.bgColor)
	gc.SetStrokeColor(t.lineColor)
	draw2dkit.Rectangle(gc, 0, 0, W, H)
	gc.FillStroke()
	gc.Close()

	t.canvas.XSurfaceSet(t.xwin.Id)

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
