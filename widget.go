package x11ui

import (
	"fmt"
	"image"
	"image/color"
	"log"
	"strings"

	"github.com/BurntSushi/xgbutil/keybind"

	"github.com/llgcode/draw2d/draw2dimg"

	"github.com/BurntSushi/xgb/xproto"
	"github.com/BurntSushi/xgbutil"
	"github.com/BurntSushi/xgbutil/ewmh"
	"github.com/BurntSushi/xgbutil/mousebind"
	"github.com/BurntSushi/xgbutil/xevent"
	"github.com/BurntSushi/xgbutil/xgraphics"
	"github.com/BurntSushi/xgbutil/xrect"
	"github.com/BurntSushi/xgbutil/xwindow"
)

type Pen struct {
	color.Color
	Width int
}

var Black = color.RGBA{0, 0, 0, 255}
var DarkGray = color.RGBA{20, 20, 20, 255}
var LightGray = xgraphics.BGRA{200, 200, 200, 255}

var White = color.RGBA{255, 255, 255, 255}
var Green = color.RGBA{0, 255, 0, 255}
var DarkGreen = xgraphics.BGRA{0, 150, 0, 255}
var LightGreen = xgraphics.BGRA{0, 50, 0, 255}

type HandlerFunctions struct {
	ClkAdvFn    func(w *Widget, x, y int)
	ClkFn       func()
	HoverFn     func()
	LeaveFn     func()
	DragStartFn func(w *Widget, global, local image.Point) bool
	DragFn      func(w *Widget, global, local image.Point) bool
	DragEndFn   func(w *Widget, global, local image.Point) bool
	EnableHover bool
}

type FactorWidgets interface {
	Move(x, y int)
}

type Widget struct {
	xwin   *xwindow.Window
	pwinID xproto.Window
	xu     *xgbutil.XUtil
	xrect.Rect
	canvas  *xgraphics.Image
	bgimage *xgraphics.Image
	gc      *draw2dimg.GraphicContext
	rawimg  *image.RGBA
	// Properties
	bgColor   color.Color /// canvas background
	fgColor   color.Color /// fill color
	barColor  color.Color
	txtColor  color.Color // text stroke color
	lineColor color.Color // line color
	margin    float64
	border    float64
	title     string
	*Layout
	fsize  int
	childs []*Widget
	/// Handlers
	HandlerFunctions
	modal bool
}

func (w *Widget) SetModal(modal bool) {
	w.modal = modal
}

func (w *Widget) SetX(X *xgbutil.XUtil) {
	w.xu = X
}

func (w *Widget) Context() *draw2dimg.GraphicContext {
	return w.gc
}

func (w *Widget) SetBackground(c color.Color) {
	w.bgColor = c
}

func (w *Widget) SetFontSize(fsize int) {
	w.fsize = fsize

}

func (ww *Widget) CreateChild(dims ...int) *Widget {
	var w *Widget
	var err error

	w = new(Widget)

	w.xu = ww.xu
	w.title = "Empty Widget"
	r := newRect(dims...)
	w.pwinID = ww.ID()
	mousebind.Initialize(w.xu)

	// Create WINDOW using usual approach
	win, err := xwindow.Generate(w.xu)
	if err != nil {
		log.Fatal(err)
	}
	win.Create(ww.pwinID, r.X, r.Y, r.Width, r.Height, xproto.CwBackPixel, 0)
	win.Listen(xproto.EventMaskKeyPress, xproto.EventMaskKeyRelease, xproto.EventMaskButtonPress, xproto.EventMaskButtonRelease, xproto.EventMaskExposure, xproto.EventMaskEnterWindow, xproto.EventMaskLeaveWindow)
	// Set _NET_WM_NAME so it looks nice.
	// err = ewmh.WmNameSet(w.xu, win.Id, w.title)
	// deBug("Could not set _NET_WM_NAME ", err)

	// Paint our image before mapping.
	w.xwin = win
	w.Rect, err = w.xwin.Geometry()
	// It's important that the map comes after setting WMGracefulClose, since
	// the WM isn't obliged to watch updates to the WM_PROTOCOLS property.
	// w.HandlerFunctions = new
	w.init()
	w.Layout = CreateLayout(0, 0, w.Width(), w.Height())

	win.Map()
	ww.appendChild(w)
	return w
}

func (w *Widget) appendChild(child *Widget) {

	w.childs = append(w.childs, child)
}

func WidgetFactory(p *Window, dims ...int) *Widget {
	var w *Widget
	var err error
	w = new(Widget)
	w.modal = false
	DEBUG_LEVEL = 1
	if p == nil {
		log.Println("Cannot create Widget without Parent Window")
		return nil
	}
	w.xu = p.X()
	w.title = "Empty Widget"
	r := newRect(dims...)
	if len(dims) == 0 || len(dims) == 2 {
		r.Width = p.Width
		r.Height = p.Height
	}
	w.pwinID = p.Id
	mousebind.Initialize(w.xu)

	// w.createRegion()
	// CREATE CANVAS based Window

	// win := w.canvas.XShowExtra("Pointer painting", true)

	// Create WINDOW using usual approach
	win, err := xwindow.Generate(w.xu)
	if err != nil {
		log.Fatal(err)
	}
	win.Create(w.pwinID, r.X, r.Y, r.Width, r.Height, xproto.CwBackPixel, 0)
	win.Listen(xproto.EventMaskKeyPress, xproto.EventMaskKeyRelease, xproto.EventMaskButtonPress, xproto.EventMaskButtonRelease, xproto.EventMaskExposure, xproto.EventMaskEnterWindow, xproto.EventMaskLeaveWindow)

	// Set _NET_WM_NAME so it looks nice.
	err = ewmh.WmNameSet(w.xu, win.Id, w.title)
	deBug("Could not set _NET_WM_NAME ", err)

	// err = ewmh.WmWindowOpacitySet(w.xu, win.Id, .3)
	// deBug("Could not set OPACITY ", err)

	// Paint our image before mapping.

	w.xwin = win
	w.Rect, err = w.xwin.Geometry()

	// It's important that the map comes after setting WMGracefulClose, since
	// the WM isn't obliged to watch updates to the WM_PROTOCOLS property.
	// w.HandlerFunctions = new
	w.init()
	w.Layout = CreateLayout(0, 0, w.Width(), w.Height())

	win.Map()
	return w
}

func (w *Widget) init() {
	w.EnableHover = true
	w.LoadTheme("")
	w.handleClose()

	w.setupCanvas()

	w.AttachHandlers()

}

func (w *Widget) LoadTheme(str string) {
	w.bgColor = color.RGBA{0, 0, 0, 255}
	w.fgColor = color.RGBA{120, 120, 120, 20}
	w.lineColor = color.RGBA{20, 120, 20, 255}
	w.txtColor = color.RGBA{255, 255, 0, 255}

}

func (w *Widget) setupCanvas() {
	r := newRect(0, 0, w.Width(), w.Height())
	w.rawimg = image.NewRGBA(r.ImageRect())
	w.gc = draw2dimg.NewGraphicContext(w.rawimg)

	// w.gc.SetLineWidth(2)
	// w.gc.SetStrokeColor(w.lineColor)
	// draw2dkit.Rectangle(w.gc, 0, 0, float64(w.Width()), float64(w.Width()))
	// w.gc.Stroke()

	w.canvas = xgraphics.NewConvert(w.xu, w.rawimg) // (w.xu, r.ImageRect())
	each := func(x, y int) xgraphics.BGRA {
		// log.Println(x, y, LightGray)
		return xgraphics.BGRA{0, 0, 0, 255}
	}
	w.canvas.For(each)
	w.bgimage = xgraphics.NewConvert(w.xu, w.rawimg)

	w.canvas.XSurfaceSet(w.xwin.Id)
	w.canvas.XDraw()
	w.canvas.XPaint(w.xwin.Id)
}

func (w *Widget) handleClose() {
	w.xwin.WMGracefulClose(
		func(xw *xwindow.Window) {
			// Detach all event handlers.
			// This should always be done when a window can no longer
			// receive events.
			log.Printf(w.title, "  destroyed %d ", xw.Id)
			xevent.Detach(xw.X, xw.Id)
			mousebind.Detach(xw.X, xw.Id)
			xw.Destroy()
			// Exit if there are no more windows left.
		})
}

func (w *Widget) Close() {
	w.xwin.Destroy()
}

func (w *Widget) GetRect() Rect {
	return XRectToRect(w.Rect)
}
func (w *Widget) SetTitle(title string) {
	w.title = title
}

func (w *Widget) keybHandler(X *xgbutil.XUtil, e xevent.KeyPressEvent) {

	modStr := keybind.ModifierString(e.State)
	keyStr := keybind.LookupString(X, e.State, e.Detail)

	modStr = strings.Replace(modStr, "lock-", "", -1)
	modStr = strings.Replace(modStr, "mod2", "", -1)
	finalstr := keyStr
	if modStr != "" {
		finalstr = fmt.Sprint(modStr, keyStr)
	}
	_ = finalstr
	// log.Println("Event code is ", e.Detail)
	// log.Printf("%s MAPS  to Keycode %v ", finalstr, keybind.StrToKeycodes(s.xu, finalstr))

	// if fn, ok := s.KeyMaps[finalstr]; ok {
	// 	if s.Debug {
	// 		log.Printf("Caught Key : %s", finalstr)
	// 	}
	// 	fn()
	// }
	// log.Println("Widgeyt ", finalstr)

}

func (w *Widget) mouseClick(X *xgbutil.XUtil, e xevent.ButtonPressEvent) {
	if w.HandlerFunctions.ClkFn == nil {
		log.Println("Widget :Mouse clicked at ", e.EventX, e.EventY)
	} else {
		if !w.modal {
			go w.ClkFn()
		} else {
			w.ClkFn()
		}
	}

}

func (w *Widget) AttachHandlers() *Widget {
	// Attach Handlers
	// mousebind.ButtonPressFun(w.mouseHandler).Connect(X, win.Id, "1", false, true)
	// mousebind.ButtonReleaseFun(w.mouseReleaseHandler).Connect(X, win.Id, "1", false, true)

	xevent.EnterNotifyFun(w.onHoverEvent).Connect(w.xu, w.xwin.Id)
	xevent.LeaveNotifyFun(w.onLeaveEvent).Connect(w.xu, w.xwin.Id)

	mousebind.ButtonPressFun(w.mouseClick).Connect(w.xu, w.xwin.Id, "1", false, true)
	return w
}

func (w *Widget) UpdateCanvas() {
	if w.canvas == nil {
		return
	}
	w.updateCanvas()
}
func (w *Widget) updateCanvas() {
	// w.canvas.For(func(x, y int) xgraphics.BGRA {
	// 	c := w.rawimg.At(x, y).(color.RGBA)
	// 	return xgraphics.BGRA{c.B, c.G, c.R, c.A}
	// })
	// w.canvas = xgraphics.NewConvert(w.xu, w.rawimg)
	// w.canvas.XSurfaceSet(w.xwin.Id)
	w.canvas.XDraw()
	w.canvas.XPaint(w.xwin.Id)
}

func (w *Widget) ID() xproto.Window {
	return w.xwin.Id
}

func (w *Widget) onHoverEvent(X *xgbutil.XUtil, e xevent.EnterNotifyEvent) {
	if w.EnableHover {
		if w.HoverFn == nil {
			w.drawBorder(StateHovered)
			// w.canvas := xgraphics.NewConvert(X, w.rawimg)
			w.updateCanvas()
		} else {
			// w.drawBorder(StateHovered)
			// // w.canvas := xgraphics.NewConvert(X, w.rawimg)
			// w.updateCanvas()
			// // calling the set hover vunf
			w.HoverFn()
		}
	}

}

func (w *Widget) drawBackground() {
	// irect := w.canvas.Rect
	each := func(x, y int) xgraphics.BGRA {
		return toBGRA(w.bgColor)
	}
	w.canvas.For(each)

}

func (w *Widget) onLeaveEvent(X *xgbutil.XUtil, e xevent.LeaveNotifyEvent) {

	if w.EnableHover {
		if w.LeaveFn == nil {
			w.drawBorder(StateNormal)
			// w.canvas := xgraphics.NewConvert(X, w.rawimg)
			w.updateCanvas()
		} else {
			// w.drawBorder(StateHovered)
			// // w.canvas := xgraphics.NewConvert(X, w.rawimg)
			// w.updateCanvas()
			// // calling the set hover vunf
			w.LeaveFn()
		}
	}

	// log.Println(w.title, " Left Hover")
	w.drawBorder(StateNormal)
	// w.canvas := xgraphics.NewConvert(X, w.rawimg)
	// w.canvas.For(func(x, y int) xgraphics.BGRA {
	// 	c := w.rawimg.At(x, y).(color.RGBA)
	// 	return xgraphics.BGRA{c.B, c.G, c.R, c.A}
	// })
	w.updateCanvas()

	// w.canvas.XSurfaceSet(w.xwin.Id)
	// w.canvas.XDraw()
	// w.canvas.XPaint(w.xwin.Id)

}

func GetIRect(w, h int) image.Rectangle {
	return image.Rectangle{origin, image.Point{w, h}}
}

var origin = image.Point{0, 0}

func (w *Widget) drawBorder(state WidgetState) {
	var clr xgraphics.BGRA
	if state == StateNormal {
		clr = toBGRA(w.bgColor)
	} else {
		clr = toBGRA(w.lineColor)
	}
	xg := w.canvas
	// border image
	outset := w.canvas.Rect
	// outset.Max.Sub(image.Point{5, 5})
	size := outset.Size()
	inset := outset.Inset(2)
	for x := 0; x < size.X; x++ {
		for y := 0; y < size.Y; y++ {
			xcond := (outset.Min.X >= x && inset.Min.X > x) || (inset.Max.X < x)
			ycond := (outset.Min.Y >= y && inset.Min.Y > y) || (inset.Max.Y < y)
			if xcond || ycond {
				xg.SetBGRA(x, y, clr)
			}
		}
	}

	w.canvas.XDraw()
	w.canvas.XPaint(w.xwin.Id)

	// //
	// // r := GetIRect(w.Width(), w.Height())
	// // si := w.canvas.SubImage(r).(*xgraphics.Image)

	// // Fresh OVERWRITE METHOD
	// w.gc.SetLineWidth(2)
	// w.gc.SetFillColor(color.RGBA{0, 0, 0, 0})
	// draw2dkit.Rectangle(w.gc, 0, 0, float64(w.Width()), float64(w.Height()))
	// w.gc.FillStroke()

}
func toBGRA(t color.Color) xgraphics.BGRA {
	c := t.(color.RGBA)
	return xgraphics.BGRA{c.B, c.G, c.R, c.A}
}

/// Layout based Region Painting

func (w *Widget) RePaint() {
	// w.xwin.Detach()

	/// Get the MAIN View
	// r := GetIRect(w.Width(), w.Height())
	// xg := xgraphics.New(w.xu, r)
	ir := GetIRect(w.Layout.w, w.Layout.h)
	ir.Min = image.Point{w.Layout.ox, w.Layout.oy}
	// log.Println("subimage region is ")
	xg := w.canvas

	// for i, reg := range w.Layout.regions {
	// 	_ = i
	// 	pmap := reg.PaintRegion()

	// 	rloc := pmap.Bounds() //.Add(w.Layout.offsets[i])
	// 	log.Println("Plot over this ", rloc, pmap.Bounds())
	// 	//draw.Draw(w.canvas, rloc, pmap, origin, draw.Src)

	// 	// draw2dimg.DrawImage(pmap, w.gc.SubImage(rloc).(*xgraphics.Image), w.gc.Current.Tr, draw.Over, draw2dimg.BilinearFilter)
	// 	// w.gc.DrawImage(pmap)
	// 	// simg := w.canvas.SubImage(rloc).(*xgraphics.Image)
	// 	xgraphics.Blend(w.canvas, pmap, origin)
	// }
	// w.updateCanvas()

	pixmap := w.Layout.regions[0].PaintRegion()
	pixmap2 := w.Layout.regions[1].PaintRegion()

	r0 := pixmap.Bounds()  //.Add(w.Layout.offsets[0])
	r1 := pixmap2.Bounds() //.Add(w.Layout.offsets[1])

	// xg1 := xg.SubImage(r0).(*xgraphics.Image)
	// xgraphics.Blend(xg1, pixmap, origin)

	// xg2 := xg.SubImage(r1).(*xgraphics.Image)
	// xgraphics.Blend(xg2, pixmap2, origin)

	size := r0.Size()
	offset := w.Layout.offsets[0]
	for x := 0; x < size.X; x++ {
		for y := 0; y < size.Y; y++ {
			xg.SetBGRA(x+offset.X+w.Layout.ox, y+offset.Y+w.Layout.oy, toBGRA(pixmap.At(x, y)))
		}
	}

	size = r1.Size()
	offset = w.Layout.offsets[1]
	for x := 0; x < size.X; x++ {
		for y := 0; y < size.Y; y++ {
			xg.SetBGRA(x+offset.X+w.Layout.ox, y+offset.Y+w.Layout.oy, toBGRA(pixmap2.At(x, y)))
		}
	}

	// // border image
	// outset := w.canvas.Rect
	// // outset.Max.Sub(image.Point{5, 5})
	// size = outset.Size()
	// inset := outset.Inset(2)
	// for x := 0; x < size.X; x++ {
	// 	for y := 0; y < size.Y; y++ {
	// 		xcond := (outset.Min.X >= x && inset.Min.X > x) || (inset.Max.X < x)
	// 		ycond := (outset.Min.Y >= y && inset.Min.Y > y) || (inset.Max.Y < y)
	// 		if xcond || ycond {
	// 			xg.SetBGRA(x, y, DarkGreen)
	// 		}
	// 	}
	// }

	// xg.XDraw()
	// xg.XPaint(w.xwin.Id)
	w.canvas.XDraw()
	w.canvas.XPaint(w.xwin.Id)
	// xg.XPaint(w.xwin.Id) //Rects(w.xwin.Id, pixmap.Bounds())

}

func (w *Widget) PaintRegions() {
	// w.xwin.Detach()

	/// Get the MAIN View
	// r := GetIRect(w.Width(), w.Height())
	// xg := xgraphics.New(w.xu, r)
	ir := GetIRect(w.Layout.w, w.Layout.h)
	ir.Min = image.Point{w.Layout.ox, w.Layout.oy}
	// log.Println("subimage region is ")
	xg := w.canvas

	// for i, reg := range w.Layout.regions {
	// 	_ = i
	// 	pmap := reg.PaintRegion()

	// 	rloc := pmap.Bounds() //.Add(w.Layout.offsets[i])
	// 	log.Println("Plot over this ", rloc, pmap.Bounds())
	// 	//draw.Draw(w.canvas, rloc, pmap, origin, draw.Src)

	// 	// draw2dimg.DrawImage(pmap, w.gc.SubImage(rloc).(*xgraphics.Image), w.gc.Current.Tr, draw.Over, draw2dimg.BilinearFilter)
	// 	// w.gc.DrawImage(pmap)
	// 	// simg := w.canvas.SubImage(rloc).(*xgraphics.Image)
	// 	xgraphics.Blend(w.canvas, pmap, origin)
	// }
	// w.updateCanvas()

	pixmap := w.Layout.regions[0].PaintRegion()
	// pixmap2 := w.Layout.regions[1].PaintRegion()

	r0 := pixmap.Bounds() //.Add(w.Layout.offsets[0])
	// r1 := pixmap2.Bounds() //.Add(w.Layout.offsets[1])

	// xg1 := xg.SubImage(r0).(*xgraphics.Image)
	// xgraphics.Blend(xg1, pixmap, origin)

	// xg2 := xg.SubImage(r1).(*xgraphics.Image)
	// xgraphics.Blend(xg2, pixmap2, origin)

	size := r0.Size()
	offset := w.Layout.offsets[0]
	for x := 0; x < size.X; x++ {
		for y := 0; y < size.Y; y++ {
			xg.SetBGRA(x+offset.X+w.Layout.ox, y+offset.Y+w.Layout.oy, toBGRA(pixmap.At(x, y)))
		}
	}

	// size = r1.Size()
	// offset = w.Layout.offsets[1]
	// for x := 0; x < size.X; x++ {
	// 	for y := 0; y < size.Y; y++ {
	// 		xg.SetBGRA(x+offset.X+w.Layout.ox, y+offset.Y+w.Layout.oy, toBGRA(pixmap2.At(x, y)))
	// 	}
	// }

	// // border image
	// outset := w.canvas.Rect
	// // outset.Max.Sub(image.Point{5, 5})
	// size = outset.Size()
	// inset := outset.Inset(2)
	// for x := 0; x < size.X; x++ {
	// 	for y := 0; y < size.Y; y++ {
	// 		xcond := (outset.Min.X >= x && inset.Min.X > x) || (inset.Max.X < x)
	// 		ycond := (outset.Min.Y >= y && inset.Min.Y > y) || (inset.Max.Y < y)
	// 		if xcond || ycond {
	// 			xg.SetBGRA(x, y, DarkGreen)
	// 		}
	// 	}
	// }

	// xg.XDraw()
	// xg.XPaint(w.xwin.Id)
	w.canvas.XDraw()
	w.canvas.XPaint(w.xwin.Id)
	// xg.XPaint(w.xwin.Id) //Rects(w.xwin.Id, pixmap.Bounds())

}

// wrapper changed
func NewWidget(X *xgbutil.XUtil, p *Window, t string, dims ...int) *Widget {
	result := WidgetFactory(p, dims...)
	result.SetTitle(t)
	return result
}

func (w *Widget) Win() *Window {
	win := new(Window)
	win.Window = w.xwin
	win.Rect = XRectToRect(w.Rect)
	return win
}

// type Window struct {
// 	//parent *xwindow.Window
// 	*xwindow.Window
// 	clkAdv OnClickFn
// 	clk    func()
// 	Rect
// 	title      string
// 	background colorful.Color
// 	view       *xwindow.Window
// 	isButton   bool
// 	isCheckBox bool
// 	checkState bool
// 	wg         sync.Mutex
// }
