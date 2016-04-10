package x11ui

import (
	"image"
	"image/color"
	"log"
	"math"
	"math/rand"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/llgcode/draw2d"

	"github.com/BurntSushi/freetype-go/freetype/truetype"

	"github.com/BurntSushi/xgb/xproto"
	"github.com/BurntSushi/xgbutil"
	"github.com/BurntSushi/xgbutil/ewmh"
	"github.com/BurntSushi/xgbutil/mousebind"
	"github.com/BurntSushi/xgbutil/xevent"
	"github.com/BurntSushi/xgbutil/xgraphics"
	"github.com/BurntSushi/xgbutil/xwindow"

	ttf "github.com/golang/freetype/truetype"
	"github.com/llgcode/draw2d/draw2dimg"
	"github.com/lucasb-eyer/go-colorful"
)

var (
	// systemBG = colorful.Hsv(0, 0, 48) //LinearRgb(.125, .125, .125)

	systemBG = colorful.LinearRgb(.5, .3, .3)
	// systemBG = colorful.Color{48, 48, 48}
	systemFG = colorful.LinearRgb(.5, .9, .5)
	// The color to use for the background.
	systemFont  *truetype.Font
	systemFData = draw2d.FontData{Name: "luxi", Family: draw2d.FontFamilyMono, Style: draw2d.FontStyleBold | draw2d.FontStyleItalic}
	dsFont      *ttf.Font
)

func init() {
	fontReader, err := os.Open("FreeMonoBold.ttf")
	if err != nil {
		log.Printf("INIT : Load X11 Font ", err)
	}
	draw2d.SetFontFolder(".")
	// Now parse the font.
	systemFont, err = xgraphics.ParseFont(fontReader)
	dsFont = draw2d.GetFont(systemFData)

	// gc.SetFontData(draw2d.FontData{Name: "luxi", Family: draw2d.FontFamilyMono, Style: draw2d.FontStyleBold | draw2d.FontStyleItalic})
	if err != nil {
		log.Printf("INIT : Load Draw2D Font ", err)
	}
}

func SetResourcePath(path string) {

	draw2d.SetFontFolder(path)
	fontReader, err := os.Open(path + "/./FreeMonoBold.ttf")
	if err != nil {
		log.Printf("Load Resource ", err)
	}
	// Now parse the font.
	systemFont, err = xgraphics.ParseFont(fontReader)
	dsFont = draw2d.GetFont(systemFData)

	systemFData = draw2d.FontData{Name: "luxi", Family: draw2d.FontFamilyMono, Style: draw2d.FontStyleBold | draw2d.FontStyleItalic}

	// gc.SetFontData(draw2d.FontData{Name: "luxi", Family: draw2d.FontFamilyMono, Style: draw2d.FontStyleBold | draw2d.FontStyleItalic})
	if err != nil {
		log.Printf("Load Font ", err)
	}
}

type OnClickFn func(w *Window, x, y int)

type Window struct {
	//parent *xwindow.Window
	*xwindow.Window
	clkAdv OnClickFn
	clk    func()
	Rect
	title      string
	background colorful.Color
	bgcolor    color.Color
	view       *xwindow.Window
	isButton   bool
	isCheckBox bool
	checkState bool
	wg         sync.Mutex
}

func (w *Window) Title() string {
	return w.title
}

func (w *Window) Click() {
	if w.clkAdv != nil {
		w.clkAdv(w, w.CenterX(), w.CenterY())
	}

	if w.clk != nil {
		w.clk()
	}
}

func (w *Window) OnClick(fn func()) {
	w.clk = fn
	log.Println("Registering ", w.Title(), "Click to ", fn)
}

func (w *Window) OnClickAdv(fn OnClickFn) {
	w.clkAdv = fn
	log.Println("Registering Adv Click ", w.Title(), "Click to ", fn)
}

func (w *Window) onHoverEvent(X *xgbutil.XUtil, e xevent.EnterNotifyEvent) {

	w.rePaint(StateHovered)
	log.Println("On Hover ", w.Title())
}
func (w *Window) onLeaveEvent(X *xgbutil.XUtil, e xevent.LeaveNotifyEvent) {

	if w.isCheckBox {
		if w.checkState {
			w.rePaint(StateSpecial)
		} else {
			w.rePaint(StateNormal)
		}
		return
	}
	w.rePaint(StateNormal)

}

func (w *Window) IsChecked() bool {

	if w.isCheckBox {
		return w.checkState
	}
	return false

}

func (w *Window) Toggle() {

	w.wg.Lock()
	if w.isCheckBox {
		w.checkState = !w.checkState
	}
	w.wg.Unlock()

}
func (w *Window) mouseReleaseHandler(X *xgbutil.XUtil, e xevent.ButtonReleaseEvent) {

	switch e.Detail {
	case 1: // left click
		if !w.isCheckBox {
			w.rePaint(StateReleased)
		}

		if w.clkAdv != nil {
			go w.clkAdv(w, int(e.EventX), int(e.EventY))
		}

		if w.clk != nil {
			w.clk()
		}

	default:
		// log.Println(w.Title(), "Some Button Clicked() ", e.Detail)
	}

}

func (w *Window) mouseHandler(X *xgbutil.XUtil, e xevent.ButtonPressEvent) {

	switch e.Detail {
	case 1: // left click
		w.rePaint(StatePressed)
		if w.clk == nil {
			log.Println(w.Title(), " Clicked() ", e.String())
		} else {
			// log.Println("Window CallBack ", w.clk)
			go w.clk()

		}
		if w.clkAdv != nil {
			w.clkAdv(w, int(e.EventX), int(e.EventY))
		}
	default:
		// log.Println(w.Title(), "Some Button Clicked() ", e.Detail)
	}

}

func (w *Window) SetTitle(t string) {
	w.title = t
	ewmh.WmNameSet(w.Window.X, w.Id, w.title)
}

func (w *Window) SetBackGround(c colorful.Color) {

	hcolor, err := strconv.ParseUint(c.Hex(), 16, 32)
	if err == nil {
		w.background = c
		w.Change(xproto.CwBackPixel, uint32(hcolor))
		w.ClearAll()
	}

}
func (w *Window) SetBGcolor(c color.Color) {
	w.bgcolor = c
	g := w.drawView(StateNormal)
	w.finishPaint(g)
}
func (w *Window) X() *xgbutil.XUtil {
	if w.Window == nil {
		return nil
	}
	return w.Window.X
}

func (w *Window) drawView(s WidgetState) *xgraphics.Image {
	r := w.ImageRect()
	dest := image.NewRGBA(r)
	gc := draw2dimg.NewGraphicContext(dest)

	// bg := colorful.LinearRgb(.025, .025, .025)
	switch s {
	case StateNormal, StateReleased:
		gc.SetFillColor(w.bgcolor)
		gc.SetStrokeColor(systemFG)
	case StateHovered:
		gc.SetFillColor(color.RGBA{0x35, 0x20, 0x20, 0x20})
		gc.SetStrokeColor(systemFG)
	case StatePressed:
		gc.SetFillColor(color.RGBA{0x20, 0x30, 0x20, 0x20})
		gc.SetStrokeColor(systemFG)
	case StateSpecial:
		gc.SetFillColor(color.RGBA{0x20, 0x80, 0x20, 0x80})
		gc.SetStrokeColor(systemFG)
	}

	gc.SetLineWidth(1)

	// gc.SetLineJoin(draw2d.RoundJoin)
	// gc.Rotate(math.Pi / 4.0)

	ww, hh := float64(w.Width), float64(w.Height)
	margin := 3.0
	ww, hh = ww-margin, hh-margin
	// cx, cy := ww/2, hh/2
	// Draw a closed shape

	// if xpressed {
	// 	gc.QuadCurveTo(cx, cy, ww, margin)
	// 	gc.QuadCurveTo(cx, cy, ww, hh)
	// 	gc.QuadCurveTo(cx, cy, margin, hh)
	// 	gc.QuadCurveTo(cx, cy, margin, margin)
	// 	// gc.QuadCurveTo(ww-5*margin, hh-5*margin, ww, hh)
	// } else {

	gc.BeginPath()
	gc.MoveTo(margin, margin)
	gc.LineTo(ww, margin)
	gc.LineTo(ww, hh)
	gc.LineTo(margin, hh)
	gc.LineTo(margin, margin)
	gc.FillStroke()
	gc.Close()

	g := xgraphics.NewConvert(w.X(), dest)
	w.drawLabel(g, w.title, margin, margin)
	return g
}

func (w *Window) drawBackground(s WidgetState) *xgraphics.Image {

	r := w.ImageRect()
	dest := image.NewRGBA(r)

	gc := draw2dimg.NewGraphicContext(dest)

	// bg := colorful.LinearRgb(.025, .025, .025)

	switch s {
	case StateNormal, StateReleased:
		gc.SetFillColor(color.RGBA{0x20, 0x20, 0x20, 20})
		gc.SetStrokeColor(systemFG)
	case StateHovered:
		gc.SetFillColor(color.RGBA{0x35, 0x20, 0x20, 20})
		gc.SetStrokeColor(systemFG)
	case StatePressed:
		gc.SetFillColor(color.RGBA{0x20, 0x30, 0x20, 20})
		gc.SetStrokeColor(systemFG)
	case StateSpecial:
		gc.SetFillColor(color.RGBA{0x20, 0x80, 0x20, 0x80})
		gc.SetStrokeColor(systemFG)
	}

	gc.SetLineWidth(1)

	// gc.SetLineJoin(draw2d.RoundJoin)
	// gc.Rotate(math.Pi / 4.0)

	ww, hh := float64(w.Width), float64(w.Height)
	margin := 1.0
	ww, hh = ww-margin, hh-margin
	// cx, cy := ww/2, hh/2
	// Draw a closed shape

	// if xpressed {
	// 	gc.QuadCurveTo(cx, cy, ww, margin)
	// 	gc.QuadCurveTo(cx, cy, ww, hh)
	// 	gc.QuadCurveTo(cx, cy, margin, hh)
	// 	gc.QuadCurveTo(cx, cy, margin, margin)

	// } else {

	gc.BeginPath()
	gc.MoveTo(margin, margin)
	gc.LineTo(ww, margin)
	gc.LineTo(ww, hh)
	gc.LineTo(margin, hh)
	gc.LineTo(margin, margin)
	gc.FillStroke()
	gc.Close()
	// }

	g := xgraphics.NewConvert(w.X(), dest)
	w.drawLabel(g, w.title)
	return g
}

func (w *Window) drawLabel(g *xgraphics.Image, str string, pos ...float64) {
	// r := w.Rect.ImageRect()
	tw, th := xgraphics.Extents(systemFont, 13, w.title)
	x, y := (w.Width-tw)/2, (w.Height-th)/2
	if len(pos) == 2 {
		x, y = int(pos[0]), int(pos[1])
	}
	g.Text(x, y, systemFG, 13, systemFont, w.title)
}

type WidgetState int

const (
	// BevelJoin represents cut segments joint
	StateNormal WidgetState = iota
	// RoundJoin represents rounded segments joint
	StatePressed
	// MiterJoin represents peaker segments joint
	StateReleased
	StateHovered
	StateSpecial
)

func (w *Window) rePaint(s WidgetState) {
	if w.isButton == true {
		g := w.drawBackground(s)
		w.finishPaint(g)
	}

}

func (w *Window) update(g *xgraphics.Image, margin ...int) {
	// w.drawLabel(g, w.title)
	g.XSurfaceSet(w.Id)
	g.XDraw()

}
func (w *Window) finalize(g *xgraphics.Image, margin ...int) {
	r := w.Rect
	if len(margin) == 2 {
		r.Grow(-2, -2)
		r.X = 2
		r.Y = 2
	}
	g.XPaintRects(w.Id, r.ImageRect())
}

func (w *Window) finishPaint(g *xgraphics.Image, margin ...int) {
	w.update(g, margin...)
	w.finalize(g, margin...)
}

type Property struct {
	X, Y int
	W, H int
}

func (p *Property) Step(dp Property) {

}

// return the distance from the src property
func (p Property) Delta(src Property) (dp Property) {
	dp.X = p.X - src.X
	dp.Y = p.Y - src.Y
	dp.W = p.W - src.W
	dp.H = p.H - src.H

	return dp
}

func (p *Property) Scale(steps int) {

}

func (w *Window) AnimateProperty(d time.Duration, start, stop Property) {

}

func (w *Window) Animate(t int) {
	tt := time.NewTicker(10 * time.Millisecond)
	ww := 10
	hh := 10
	for range tt.C {
		if ww > w.Rect.Width {
			tt.Stop()
			break
		} else {
			// w.WMResize(ww, hh)
			w.Resize(ww, hh)
			ww += 10
			hh += 10

		}

	}

}

func (w *Window) Draw(X *xgbutil.XUtil, e xevent.ExposeEvent) {
	if w.title == "graph" {
		log.Println(w.Title(), "======== need to draw ========", e.String())
		//	w.rePaint()
	}
	// draw.Draw(X.Screen().	, r, src, sp, op)
	// c.PolyRectangle(win, fg, rectangles)
	// c.ImageText8(win, bg, 20, 20, strBytes)

	// if w.title == "graph" {
	// 	var graph *image.RGBA

	// 	r := w.Rect.ImageRect()
	// 	graph = RandomGraph(r)
	// 	g := xgraphics.NewConvert(w.X(), graph)

	// 	g.XPaint(w.view.Id)
	// 	// g.XPaintRects(w.view.Id, 0, 0)
	// 	// g.XSurfaceSet(w.view.Id)

	// }

}

func (w *Window) Filler(x, y int) xgraphics.BGRA {
	margin := 1
	// borderPixel := 1
	var r xgraphics.BGRA

	if (x > margin && x < w.Width-margin) && (y > margin && y < w.Height-margin) {
		// result := systemBG.BlendRgb(systemFG, float64(w.Width)/float64(x)).Clamped()
		r.A, r.B, r.G = systemBG.RGB255()
		return r
		// return bg
	} else {
		// result := systemBG.BlendRgb(systemFG, .9).Clamped()
		r.A, r.B, r.G = systemFG.RGB255()
		return r
		// return fg
	}

	// if x > y {
	// 	return bg
	// } else {
	// 	return fg
	// }
}

func newWindow(X *xgbutil.XUtil, p *Window, t string, dims ...int) *Window {
	w := new(Window)
	w.title = t
	w.background = systemBG
	var parent xproto.Window
	if p == nil {
		parent = X.RootWin()
	} else {
		parent = p.Id
	}
	mousebind.Initialize(X)
	r := newRect(dims...)
	win, err := xwindow.Generate(X)
	// s := X.Screen()
	w.bgcolor = color.RGBA{0x20, 0x20, 0x20, 0xFF}

	// mask := xproto.GcForeground | xproto.GcGraphicsExposures
	// values := []uint32{s.BlackPixel, 0}
	win.Create(parent, r.X, r.Y, r.Width, r.Height, xproto.CwBackPixel, 0xfffff)

	// win.Create(parent, r.X, r.Y, r.Width, r.Height, mask, values...)
	if err != nil {
		log.Fatal(err)
	}
	// win.MoveResize(r.X, r.Y, r.Width, r.Height)
	if p == nil {
		win.Change(xproto.CwBackPixel, 0x684426)

	} else {
		// win.Change(xproto.CwBackPixel, 0xFFAA00)
	}

	//if p == nil {
	win.Listen(xproto.EventMaskKeyPress, xproto.EventMaskKeyRelease, xproto.EventMaskButtonPress, xproto.EventMaskButtonRelease, xproto.EventMaskExposure, xproto.EventMaskEnterWindow, xproto.EventMaskLeaveWindow)
	//}

	w.Rect = r

	w.Window = win
	if p == nil {
		// xevent.ButtonPressFun(w.mouseHandler).Connect(X, win.Id, "1", false, true)
		xevent.ButtonPressFun(w.mouseHandler).Connect(X, win.Id)
	} else {
		mousebind.ButtonPressFun(w.mouseHandler).Connect(X, win.Id, "1", false, true)
		mousebind.ButtonReleaseFun(w.mouseReleaseHandler).Connect(X, win.Id, "1", false, true)
		xevent.EnterNotifyFun(w.onHoverEvent).Connect(X, win.Id)
		xevent.LeaveNotifyFun(w.onLeaveEvent).Connect(X, win.Id)
		mousebind.ButtonPressFun(w.mouseHandler).Connect(X, win.Id, "2", false, true)
	}
	// xevent.ExposeFun(w.Draw).Connect(X, win.Id)

	win.WMGracefulClose(
		func(w *xwindow.Window) {
			// Detach all event handlers.
			// This should always be done when a window can no longer
			// receive events.
			log.Printf("Window destroyed %d ", w.Id)
			xevent.Detach(w.X, w.Id)
			mousebind.Detach(w.X, w.Id)
			w.Destroy()
			// Exit if there are no more windows left.

		})

	// It's important that the map comes after setting WMGracefulClose, since
	// the WM isn't obliged to watch updates to the WM_PROTOCOLS property.
	win.Map()
	w.PaintOnce()

	// xevent.ButtonPressFun(w.mouseHandler).Connect(X, win.Id)
	return w
}

func (w *Window) Plot() {
	_, e := w.Parent()
	if e != nil {
		// log.Println("Window is Closed ")
		return
	}

	gimg := RandomGraph(w.ImageRect())
	g := xgraphics.NewConvert(w.X(), gimg)

	// r := w.Rect
	// r.Grow(-100, -100).ShiftDown(10).ShiftRight(10)
	// g := xgraphics.New(w.X(), r.ImageRect())
	// g = g.Scale(r.Width, r.Height)

	// for i := 0; i < 255; i++ {
	// 	// c := systemFG.BlendRgb(systemBG, float64(i)/100)
	// 	c := xgraphics.BGRA{uint8(i), 0, 0, 255}
	// 	for j := 0; j < r.Height; j++ {
	// 		g.Set(i, j, c)
	// 	}

	// }

	// g = w.drawView(StateNormal)
	// w.update(g, 10, 10)
	// w.finalize(g, 2, 2)
	w.finishPaint(g)

}

func RandomGraph(r image.Rectangle) *image.RGBA {
	dest := image.NewRGBA(r)
	gc := draw2dimg.NewGraphicContext(dest)
	// Set some properties
	gc.SetFillColor(color.RGBA{0xF0, 0x0, 0, 0x0})
	gc.SetStrokeColor(color.RGBA{0x44, 0xF4, 0x44, 0xff})
	gc.SetLineWidth(2)
	// gc.Rotate(math.Pi / 4.0)
	gc.Scale(1, .8)
	gc.Translate(0, +float64(r.Max.Y/2.0))
	// Draw a closed shape
	gc.MoveTo(0, 0) // should always be called first for a new path
	for x := 0; x < r.Max.X; x++ {
		y := (math.Sin(2*math.Pi*5*float64(x)/float64(r.Max.X)) + rand.Float64()*.5) * (float64(r.Max.Y) / 2.0)
		gc.LineTo(float64(x), float64(y))
	}
	// gc.Close()
	gc.Stroke()
	return dest
}

func (w *Window) PaintOnce() {
	g := w.drawBackground(StateNormal)
	w.finishPaint(g)
}

func (w *Window) XWin() *xwindow.Window {
	return w.Window
}

func (w *Window) XProtoWin() xproto.Window {
	return w.Window.Id
}

func xNewWidget(X *xgbutil.XUtil, p *Window, t string, dims ...int) *Window {
	w := new(Window)
	w.title = t
	w.background = systemBG
	var parent xproto.Window
	if p == nil {
		parent = X.RootWin()
	} else {
		parent = p.Id
	}
	mousebind.Initialize(X)
	r := newRect(dims...)
	win, err := xwindow.Generate(X)
	if err != nil {
		log.Fatal("NewWidget : Unable to Create ", err)
	}

	///Raw window creation & Manage handlers
	win.Create(parent, r.X, r.Y, r.Width, r.Height, xproto.CwBackPixel, 0x30)
	win.Listen(xproto.EventMaskKeyPress, xproto.EventMaskKeyRelease, xproto.EventMaskButtonPress, xproto.EventMaskButtonRelease, xproto.EventMaskExposure, xproto.EventMaskEnterWindow, xproto.EventMaskLeaveWindow)
	mousebind.ButtonPressFun(w.mouseHandler).Connect(X, win.Id, "1", false, true)
	mousebind.ButtonReleaseFun(w.mouseReleaseHandler).Connect(X, win.Id, "1", false, true)
	xevent.EnterNotifyFun(w.onHoverEvent).Connect(X, win.Id)
	xevent.LeaveNotifyFun(w.onLeaveEvent).Connect(X, win.Id)
	mousebind.ButtonPressFun(w.mouseHandler).Connect(X, win.Id, "2", false, true)

	win.WMGracefulClose(
		func(w *xwindow.Window) {
			// Detach all event handlers.
			// This should always be done when a window can no longer
			// receive events.
			log.Printf("Window destroyed %d ", w.Id)
			xevent.Detach(w.X, w.Id)
			mousebind.Detach(w.X, w.Id)
			w.Destroy()
			// Exit if there are no more windows left.

		})

	// It's important that the map comes after setting WMGracefulClose, since
	// the WM isn't obliged to watch updates to the WM_PROTOCOLS property.
	win.Map()

	w.Rect = r
	// xevent.ButtonPressFun(w.mouseHandler).Connect(X, win.Id)
	return w
}
