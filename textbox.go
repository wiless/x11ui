package x11ui

import (
	"bytes"
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
	readOnly  bool
}

//NewTextBox creates a child TextBox widget in the roots of Window p
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
	for range c {
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

func (t *TextBox) SetReadOnly(readonly bool) {

	if !t.readOnly && readonly {
		xevent.Detach(t.xu, t.xwin.Id)
		t.readOnly = readonly
	}
	if t.readOnly && !readonly {
		t.readOnly = readonly
		xevent.KeyPressFun(t.keybHandler).Connect(t.xu, t.xwin.Id)
	}
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

func (t *TextBox) SetText(txt string) {
	ar := strings.Split(txt, "\n")
	t.line = 0
	t.cursor = 0
	t.bgColor = color.RGBA{0, 0, 255, 0}

	t.drawBackground()
	t.drawTextBox(StateNormal)
	for _, s := range ar {
		t.AppendLine(s)
	}

}

// lineWidth returns the number of characters that can be accomodated in a line
func (t *TextBox) lineWidth() int {
	return (t.Width() - int(t.margin))
}

func (t *TextBox) AppendLine(str string) {

	if strings.Contains(str, "\n") {
		ar := strings.Split(str, "\n")

		for _, s := range ar {
			t.AppendLine(s)
		}
		return
	}

	t.line += t.linespace
	t.cursor = 0

	charsPerLine := 29
	if len(str) > charsPerLine {

		splits := SplitSubN(str, charsPerLine)
		for _, ss := range splits {
			t.AppendLine(ss)
		}
		return
	}

	nx, ny, _ := t.canvas.Text(int(t.cursor), int(t.line), t.txtColor, 12, systemFont, str)
	_ = ny
	// log.Println("nx,ny", nx, ny)
	t.cursor = nx

	if t.cursor > (t.Width() - int(t.margin)) {
		t.line += t.linespace // line spacing
		t.cursor = 0
	}

	t.updateCanvas()
}

func SplitSubN(s string, n int) []string {
	sub := ""
	subs := []string{}

	runes := bytes.Runes([]byte(s))
	l := len(runes)
	for i, r := range runes {
		sub = sub + string(r)
		if (i+1)%n == 0 {
			subs = append(subs, sub)
			sub = ""
		} else if (i + 1) == l {
			subs = append(subs, sub)
		}
	}

	return subs
}
func (t *TextBox) AddRulers() {
	t.gc.SetFillColor(color.RGBA{0, 128, 0, 0})
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
	if str == "Return" || str == "KP_Enter" || str == "\n" {
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
		log.Println("I am returning more than 1 Char")
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
	t.readOnly = false

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

}
