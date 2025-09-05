package x11ui

import (
	"embed"
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

	"fmt"
	"path/filepath"

	ttf "github.com/golang/freetype/truetype"

	"github.com/BurntSushi/xgb/xproto"
	"github.com/BurntSushi/xgbutil"
	"github.com/BurntSushi/xgbutil/ewmh"
	"github.com/BurntSushi/xgbutil/mousebind"
	"github.com/BurntSushi/xgbutil/xevent"
	"github.com/BurntSushi/xgbutil/xgraphics"
	"github.com/BurntSushi/xgbutil/xwindow"

	"github.com/llgcode/draw2d/draw2dimg"
	"github.com/lucasb-eyer/go-colorful"
)

//go:embed fonts/luxisr.ttf fonts/FreeMonoBold.ttf
var content embed.FS

var (
	// systemBG is the default background color for UI elements.
	systemBG = colorful.LinearRgb(.5, .3, .3)
	// systemFG is the default foreground color for UI elements.
	systemFG = colorful.LinearRgb(.5, .9, .5)
	// systemFont is the default TrueType font used for rendering text.
	systemFont *truetype.Font
	// systemFData holds the font data for the default system font.

	systemFData = draw2d.FontData{Name: "luxi", Family: draw2d.FontFamilyMono, Style: draw2d.FontStyleBold | draw2d.FontStyleItalic}
	// dsFont is the draw2d font instance derived from systemFData.
	dsFont *truetype.Font
)

// init initializes the UI package, loading default resources.
func init() {

	err := SetResourcePath("./fonts", "")
	// systemFData = draw2d.FontData{Name: "luxi", Family: draw2d.FontFamilyMono, Style: draw2d.FontStyleBold | draw2d.FontStyleItalic}
	if err != nil {
		log.Printf("INIT : Load X11 Font %v", err)
	}

}

// getExecutableDir returns the directory of the currently running executable.
func getExecutableDir() (string, error) {
	exePath, err := os.Executable()
	if err != nil {
		return "", err
	}
	return filepath.Dir(exePath), nil
}

// Register Fonts
func RegisterFont(path, name string) {
	fontbytes, err := os.ReadFile(path)
	if err != nil {
		log.Printf("Failed to locate font file: %v", err)
	}
	font, err := ttf.Parse(fontbytes)
	if err != nil {
		log.Printf("Failed to parse font: %v", err)
	} else {
		draw2d.RegisterFont(draw2d.FontData{Name: name}, font)
	}
}

// SetResourcePath sets the path for UI resources like fonts
// The path parameter is the directory containing the font files.
func SetResourcePath(path, fontName string) error {
	if fontName == "" {
		fontName = "FreeMonoBold.ttf"
		log.Printf("Setting Default Font %v", fontName)
	}
	// Opening from embedded FS
	fr, e := content.Open(filepath.Join("fonts", fontName))
	if e != nil {
		return fmt.Errorf("could not open embedded font %s: %w", fontName, e)
	}
	defer fr.Close()
	systemFont, e = xgraphics.ParseFont(fr)
	if e != nil {
		return fmt.Errorf("could not parse embedded font %s: %w", fontName, e)
	}
	draw2d.SetFontFolder(path) // No need for a folder, fonts are embedded

	// draw2d.SetFontCache(fc)
	// dsFont = draw2d.GetFont(systemFData)
	log.Printf("Default draw2D FontFolder %v", draw2d.GetFontFolder())

	return nil
}

// OnClickFn defines the signature for advanced click handlers that receive window and coordinate information.
type OnClickFn func(w *Window, x, y int)

// Window represents an X11 window with UI capabilities.
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
	margin     float64
	rawimage   *image.RGBA
	ximg       *xgraphics.Image
}

// Title returns the title of the window.
func (w *Window) Title() string {
	return w.title
}

// Click programmatically triggers the click handlers associated with the window.
func (w *Window) Click() {
	if w.clkAdv != nil {
		w.clkAdv(w, w.CenterX(), w.CenterY())
	}

	if w.clk != nil {
		w.clk()
	}
}

// OnClick sets a simple click handler for the window.
func (w *Window) OnClick(fn func()) {
	w.clk = fn
	log.Println("Registering ", w.Title(), "Click to ", fn)
}

// OnClickAdv sets an advanced click handler that receives the window and click coordinates.
func (w *Window) OnClickAdv(fn OnClickFn) {
	w.clkAdv = fn
	log.Println("Registering Adv Click ", w.Title(), "Click to ", fn)
}

// onHoverEvent is the event handler for when the mouse enters the window's area.
func (w *Window) onHoverEvent(X *xgbutil.XUtil, e xevent.EnterNotifyEvent) {

	w.rePaint(StateHovered)
	log.Println("On Hover ", w.Title())
}

// onLeaveEvent is the event handler for when the mouse leaves the window's area.
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

// IsChecked returns the checked state of the window if it's a checkbox.
func (w *Window) IsChecked() bool {

	if w.isCheckBox {
		return w.checkState
	}
	return false

}

// Toggle changes the checked state of the window if it's a checkbox.
func (w *Window) Toggle() {

	w.wg.Lock()
	if w.isCheckBox {
		w.checkState = !w.checkState
	}
	w.wg.Unlock()

}

// mouseReleaseHandler handles the mouse button release events.
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

// mouseHandler handles the mouse button press events.
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

// SetTitle sets the title of the window.
func (w *Window) SetTitle(t string) {
	w.title = t
	ewmh.WmNameSet(w.Window.X, w.Id, w.title)
}

// SetBackGround sets the background color of the window using a colorful.Color.
func (w *Window) SetBackGround(c colorful.Color) {

	hcolor, err := strconv.ParseUint(c.Hex(), 16, 32)
	if err == nil {
		w.background = c
		w.Change(xproto.CwBackPixel, uint32(hcolor))
		w.ClearAll()
	}

}

// SetBGcolor sets the background color of the window using a color.Color and redraws it.
func (w *Window) SetBGcolor(c color.Color) {
	w.bgcolor = c
	g := w.drawView(StateNormal)
	w.finishPaint(g)
}

// X returns the underlying *xgbutil.XUtil connection.
func (w *Window) X() *xgbutil.XUtil {
	if w.Window == nil {
		return nil
	}
	return w.Window.X
}

// SetMargin sets the margin for the window content.
func (w *Window) SetMargin(m float64) {
	w.margin = m
}

// drawView draws the main view of the window based on its state.
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

	gc.SetLineWidth(0)

	// gc.SetLineJoin(draw2d.RoundJoin)
	// gc.Rotate(math.Pi / 4.0)

	ww, hh := float64(w.Width), float64(w.Height)
	margin := w.margin
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

// drawBackground draws the background of the window based on its state.
func (w *Window) drawBackground(s WidgetState) *xgraphics.Image {

	r := w.ImageRect()
	dest := image.NewRGBA(r)

	gc := draw2dimg.NewGraphicContext(dest)

	// fontFamily := "CustomFont" // A name you give to your font

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

// drawLabel draws the window's title on the provided xgraphics.Image.
func (w *Window) drawLabel(g *xgraphics.Image, str string, pos ...float64) {
	// r := w.Rect.ImageRect()
	tw, th := xgraphics.Extents(systemFont, 13, w.title)
	x, y := (w.Width-tw)/2, (w.Height-th)/2
	if len(pos) == 2 {
		x, y = int(pos[0]), int(pos[1])
	}
	g.Text(x, y, systemFG, 13, systemFont, w.title)
}

// WidgetState represents the current state of a UI widget
type WidgetState int

const (
	// StateNormal indicates the normal, default state of a widget.
	StateNormal WidgetState = iota
	// StatePressed indicates the state when a widget is being pressed (e.g., a button).
	StatePressed
	// StateReleased indicates the state when a widget has just been released after being pressed.
	StateReleased
	// StateHovered indicates the state when the mouse cursor is hovering over a widget.
	StateHovered
	// StateSpecial indicates a special or custom state for a widget.
	StateSpecial
)

// rePaint redraws the window with a given state, if it's a button.
func (w *Window) rePaint(s WidgetState) {
	if w.isButton == true {
		g := w.drawBackground(s)
		w.finishPaint(g)
	}

}

// update applies the drawing from the xgraphics.Image to the window's surface.
func (w *Window) update(g *xgraphics.Image, margin ...int) {
	// w.drawLabel(g, w.title)
	g.XSurfaceSet(w.Id)
	g.XDraw()

}

// finalize paints the updated surface onto the window.
func (w *Window) finalize(g *xgraphics.Image, margin ...int) {
	r := w.Rect
	if len(margin) == 2 {
		r.Grow(-2, -2)
		r.X = 2
		r.Y = 2
	}
	g.XPaintRects(w.Id, r.ImageRect())
}

// finishPaint is a convenience function to update and finalize the window painting.
func (w *Window) finishPaint(g *xgraphics.Image, margin ...int) {
	w.update(g, margin...)
	w.finalize(g, margin...)
}

// Property represents geometric properties of a window for animation.
type Property struct {
	X, Y int
	W, H int
}

// Step is a placeholder for a step function in animation. (Currently a stub)
func (p *Property) Step(dp Property) {

}

// Delta returns the difference between two Property structs.
func (p Property) Delta(src Property) (dp Property) {
	dp.X = p.X - src.X
	dp.Y = p.Y - src.Y
	dp.W = p.W - src.W
	dp.H = p.H - src.H

	return dp
}

// Scale scales the property values. (Currently a stub)
func (p *Property) Scale(steps int) {

}

// AnimateProperty animates a window's properties over a duration. (Currently a stub)
func (w *Window) AnimateProperty(d time.Duration, start, stop Property) {

}

// Animate performs a simple resize animation.
func (w *Window) Animate(t int) {
	log.Println("Animate Window ", w.Title())
	tt := time.NewTicker(1000 * time.Millisecond)
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

// Draw handles expose events to redraw parts of the window.
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

// Filler is a function used to fill the window background with a pattern.
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

// newWindow creates a new basic window as a child of p.
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
	win.Listen(xproto.EventMaskKeyPress, xproto.EventMaskKeyRelease, xproto.EventMaskButtonPress, xproto.EventMaskButtonRelease, xproto.EventMaskExposure, xproto.EventMaskEnterWindow, xproto.EventMaskLeaveWindow, xproto.EventMaskKeyPress)
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

// Plot draws a random graph on the window. Used for demonstration.
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

// RandomGraph generates an image with a random sine wave graph.
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

// PaintOnce draws the window's background and label a single time.
func (w *Window) PaintOnce() {
	g := w.drawBackground(StateNormal)
	w.finishPaint(g)
}

// XWin returns the underlying *xwindow.Window.
func (w *Window) XWin() *xwindow.Window {
	return w.Window
}

// XProtoWin returns the X protocol window ID.
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

// CreateRawImage creates a new RGBA image associated with the window for direct drawing.
func (w *Window) CreateRawImage(x, y, ww, hh int) *image.RGBA {
	w.rawimage = image.NewRGBA(image.Rect(x, y, ww, hh))

	// w.rawimage.SetRGBA(x int, y int, c color.RGBA)
	w.ximg = xgraphics.NewConvert(w.X(), w.rawimage)
	err := w.ximg.XSurfaceSet(w.Id)
	log.Print(err)
	return w.rawimage
}

// RawImage returns the raw image buffer of the window.
func (w Window) RawImage() *image.RGBA {
	return w.rawimage
}

// UpdatePlot updates the window content from an image.
func (w *Window) UpdatePlot(img image.Image) {

	// log.Printf("Bounds ", img.Bounds(), "window rects", w.Rect)
	// ox, oy := w.Rect.X, w.Rect.Y
	rr := w.ImageRect()

	xgraphics.Blend(w.ximg, w.rawimage, image.Point{0, 0})
	// ximg = xgraphics.NewConvert(w.X(), xgraphics.Scale(img, rr.Dx(), rr.Dy()))

	w.ximg.XSurfaceSet(w.Id)
	log.Print("Redrawing ", rr)
	w.ximg.XPaintRects(w.Id, image.Rect(0, 0, rr.Dx(), rr.Dy()))
	w.ximg.XDraw()

}

// ReDrawImage redraws the window from its raw image buffer.
func (w *Window) ReDrawImage() {
	if w.rawimage == nil {
		log.Print("Create Raw Image .. first")
		return
	}
	// log.Printf("Bounds ", img.Bounds(), "window rects", w.Rect)
	// ox, oy := w.Rect.X, w.Rect.Y
	// rr := w.rawimage.Bounds()
	rr, _ := w.Geometry()
	// log.Print("PLOT WIN RESIZED", rr)
	// // w.ximg := xgraphics.NewConvert(w.X(),xgraphics.Scale(w.rawimage, rr.Dx(), rr.Dy()))
	// // w.ximg. = xgraphics.NewConvert(w.X(), w.rawimage)
	//
	// // xgraphics.Scale(w.rawimage, rr.Width(), rr.Height())
	// w.ximg = w.ximg.Scale(rr.Width(), rr.Height())
	// err := w.ximg.XSurfaceSet(w.Id)
	// log.Print(err)
	// log.Print("ximg ", w.ximg.Bounds())
	// log.Print("rawimg ", w.rawimage.Rect)
	// di := xgraphics.Scale(w.rawimage, rr.Width(), rr.Height())
	// log.Print("scaled raw ", di.Bounds())

	rc := color.RGBA{255, 255, 0, 255}
	red := image.NewRGBA(w.rawimage.Bounds())
	for x := 0; x < red.Bounds().Dx(); x++ {
		for y := 0; y < red.Bounds().Dy(); y++ {
			red.SetRGBA(x, y, rc)
		}

	}
	xgraphics.Blend(w.ximg, red, image.Point{10, 10})
	xgraphics.Blend(w.ximg, w.rawimage, image.Point{0, 0})
	// xgraphics.Blend(w.ximg, red, image.Point{10, 10})

	// xgraphics.BlendBgColor(w.ximg, rc) //(w.ximg, red, image.Point{0, 0})

	// log.Print("Xpaint raw ", rr.Width(), rr.Height())
	// w.ximg.For(func(x int, y int) xgraphics.BGRA {
	// 	var c xgraphics.BGRA
	// 	ic := w.rawimage.At(x, y)
	// 	r, g, b, _ := ic.RGBA()
	// 	c.R = uint8(r)
	// 	c.G = uint8(g)
	// 	c.B = uint8(b)
	// 	c.A = 255
	// 	// log.Print("Replacing @ ", x, y, ic, c)
	// 	return c
	// })
	// ximg = xgraphics.NewConvert(w.X(), xgraphics.Scale(img, rr.Dx(), rr.Dy()))

	//w.ximg.XSurfaceSet(w.Id)
	w.ximg.XPaintRects(w.Id, image.Rect(0, 0, rr.Width(), rr.Height()))
	w.ximg.XDraw()

}
