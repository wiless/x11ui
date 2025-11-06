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
	systemFG = colorful.LinearRgb(.8, .8, .8)
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
	prevRect         Rect
	onSizeChange     func(w, h int)
	onMove           func(x, y int)
	onMoveComplete   func(x, y int)
	onResizeComplete func(w, h int)
	title            string
	background       colorful.Color
	bgcolor          color.Color
	view             *xwindow.Window
	isButton         bool
	isCheckBox       bool
	checkState       bool
	wg               sync.Mutex
	margin           float64
	rawimage         *image.RGBA
	ximg             *xgraphics.Image
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

}

// OnClickAdv sets an advanced click handler that receives the window and click coordinates.
func (w *Window) OnClickAdv(fn OnClickFn) {
	w.clkAdv = fn

}

func (w *Window) Show() {
	w.Map()
}

func (w *Window) Hide() {
	w.Unmap()
}

func (w *Window) ToggleVisibility() {
	attrs, err := xproto.GetWindowAttributes(w.X().Conn(), w.Id).Reply()
	if err != nil {
		// Or handle the error more gracefully
		log.Println("Could not get window attributes:", err)
		return
	}
	if attrs.MapState == xproto.MapStateUnmapped {
		w.Map()
	} else {
		w.Unmap()
	}
}

func (w *Window) OnSizeChange(fn func(w, h int)) {
	w.onSizeChange = fn
}

func (w *Window) OnMove(fn func(x, y int)) {
	w.onMove = fn
}

func (w *Window) OnMoveComplete(fn func(x, y int)) {
	w.onMoveComplete = fn
}

func (w *Window) OnResizeComplete(fn func(w, h int)) {
	w.onResizeComplete = fn
}

// onHoverEvent is the event handler for when the mouse enters the window's area.
func (w *Window) onHoverEvent(X *xgbutil.XUtil, e xevent.EnterNotifyEvent) {

	if w.isCheckBox && w.checkState {
		w.rePaint(StateHoveredChecked)
	} else {
		w.rePaint(StateHovered)
	}

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
		if w.isCheckBox {
			w.Toggle()
			if w.checkState {
				w.rePaint(StateSpecial)
			} else {
				w.rePaint(StateNormal)
			}
		} else {
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
	w.drawView(StateNormal)
	w.finishPaint(w.ximg)
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
func (w *Window) drawView(s WidgetState) {
	r := w.ImageRect()
	dest := image.NewRGBA(r)
	gc := draw2dimg.NewGraphicContext(dest)

	switch s {
	case StateNormal, StateReleased:
		gc.SetFillColor(toRGBA(w.bgcolor))
		gc.SetStrokeColor(toRGBA(CurrentTheme.LineColor))
	case StateHovered:
		gc.SetFillColor(toRGBA(toColorful(CurrentTheme.ForegroundColor).BlendLuvLCh(colorful.Color{R: 0, G: 0, B: 0}, 0.2).Clamped()))
		gc.SetStrokeColor(toRGBA(CurrentTheme.LineColor))
	case StatePressed:
		gc.SetFillColor(toRGBA(toColorful(CurrentTheme.ForegroundColor).BlendLuvLCh(colorful.Color{R: 0, G: 0, B: 0}, 0.4).Clamped()))
		gc.SetStrokeColor(toRGBA(CurrentTheme.LineColor))
	case StateSpecial:
		gc.SetFillColor(toRGBA(toColorful(w.bgcolor).BlendLuvLCh(toColorful(CurrentTheme.CheckboxCheckedColor), 0.5).Clamped()))
		gc.SetStrokeColor(toRGBA(CurrentTheme.LineColor))
	}

	gc.SetLineWidth(0)

	ww, hh := float64(w.Width), float64(w.Height)
	margin := w.margin
	ww, hh = ww-margin, hh-margin

	gc.BeginPath()
	gc.MoveTo(margin, margin)
	gc.LineTo(ww, margin)
	gc.LineTo(ww, hh)
	gc.LineTo(margin, hh)
	gc.LineTo(margin, margin)
	gc.FillStroke()
	gc.Close()

	w.rawimage = dest
	w.ximg = xgraphics.NewConvert(w.X(), dest)
	w.drawLabel(w.ximg, w.title, nil, margin, margin)
}

// drawBackground draws the background of the window based on its state.
func (w *Window) drawBackground(s WidgetState) {

	r := w.ImageRect()
	dest := image.NewRGBA(r)
	gc := draw2dimg.NewGraphicContext(dest)

	var fillColor, strokeColor color.Color

	switch s {
	case StateNormal, StateReleased:
		fillColor = toRGBA(CurrentTheme.BarColor)
		strokeColor = toRGBA(CurrentTheme.LineColor)
	case StateHovered:
		fillColor = toRGBA(toColorful(CurrentTheme.BarColor).BlendLuvLCh(colorful.Color{R: 1, G: 1, B: 1}, 0.2).Clamped())
		strokeColor = toRGBA(CurrentTheme.LineColor)
	case StatePressed:
		fillColor = toRGBA(toColorful(CurrentTheme.BarColor).BlendLuvLCh(colorful.Color{R: 0, G: 0, B: 0}, 0.4).Clamped())
		strokeColor = toRGBA(CurrentTheme.LineColor)
	case StateSpecial:
		fillColor = toRGBA(CurrentTheme.CheckboxCheckedColor)
		strokeColor = toRGBA(CurrentTheme.LineColor)
	case StateHoveredChecked:
		fillColor = toRGBA(toColorful(CurrentTheme.CheckboxCheckedColor).BlendLuvLCh(colorful.Color{R: 1, G: 1, B: 1}, 0.2).Clamped())
		strokeColor = toRGBA(CurrentTheme.LineColor)
	}

	gc.SetFillColor(fillColor)
	gc.SetStrokeColor(strokeColor)

	ww, hh := float64(w.Width), float64(w.Height)
	margin := w.margin
	ww, hh = ww-margin, hh-margin

	gc.BeginPath()
	gc.MoveTo(margin, margin)
	gc.LineTo(ww, margin)
	gc.LineTo(ww, hh)
	gc.LineTo(margin, hh)
	gc.LineTo(margin, margin)
	gc.FillStroke()
	gc.Close()

	w.rawimage = dest
	w.ximg = xgraphics.NewConvert(w.X(), dest)

	var buttonTextColor color.Color
	barColorColorful := toColorful(CurrentTheme.BarColor)
	rVal, gVal, bVal, _ := barColorColorful.RGBA()
	avg := (float64(rVal>>8) + float64(gVal>>8) + float64(bVal>>8)) / 3.0 / 255.0

	if avg > 0.5 { // If background is light, use dark text
		buttonTextColor = color.RGBA{R: 0, G: 0, B: 0, A: 255} // Black
	} else { // If background is dark, use light text
		buttonTextColor = color.RGBA{R: 255, G: 255, B: 255, A: 255} // White
	}
	w.drawLabel(w.ximg, w.title, buttonTextColor)
}

// drawLabel draws the window's title on the provided xgraphics.Image.
func (w *Window) drawLabel(g *xgraphics.Image, str string, textColor color.Color, pos ...float64) {
	if textColor == nil {
		textColor = CurrentTheme.TextColor
	}
	tw, th := xgraphics.Extents(systemFont, 13, w.title)
	x, y := (w.Width-tw)/2, (w.Height-th)/2
	if len(pos) == 2 {
		x, y = int(pos[0]), int(pos[1])
	}
	g.Text(x, y, textColor, 13, systemFont, w.title)
}

// In drawView, when calling w.drawLabel(w.ximg, w.title, margin, margin):
// w.drawLabel(w.ximg, w.title, nil, margin, margin)



type Container struct {
	Window // Embed the existing Window struct directly

	layoutDirection LayoutDirection // Horizontal or Vertical

	children []WindowProvider // List of child widgets/windows

	hspacing, vspacing int

	hpadding, vpadding int

	pvsChildRect Rect // To keep track of the previous child's position for cascading

}

func (c *Container) AddWidget(widget WindowProvider) {
	// Check if widget is already added
	for _, child := range c.children {
		if child == widget {
			return // Already there
		}
	}
	c.children = append(c.children, widget)
	c.RelayoutChildren() // Relayout and redraw all children
}
func (c *Container) RelayoutChildren() {
	log.Printf("RelayoutChildren: Container dimensions: Width=%d, Height=%d", c.Rect.Width, c.Rect.Height)
	c.ClearAll() // Clear the container's background

	// Previous child rect, in absolute coordinates. Start with a zero-width rect at the container's top-left.
	prevRect := Rect{c.Rect.X, c.Rect.Y, 0, 0}
	if c.layoutDirection == LayoutHor {
		prevRect.Width = -c.hspacing
	} else {
		prevRect.Height = -c.vspacing
	}

	for _, wp := range c.children {
		child := wp.Win() // Get the underlying *Window

		// Log original position
		log.Printf("RelayoutChildren: Child '%s' original absolute position: X=%d, Y=%d, W=%d, H=%d", child.Title(), child.Rect.X, child.Rect.Y, child.Rect.Width, child.Rect.Height)

		var newRect Rect // absolute coordinates
		newRect.Width = child.Rect.Width
		newRect.Height = child.Rect.Height

		switch c.layoutDirection {
		case LayoutHor:
			newRect.X = prevRect.X + prevRect.Width + c.hspacing
			newRect.Y = prevRect.Y
			newRect.Height = c.Rect.Height // Resize height to container's height
		case LayoutVer:
			newRect.X = prevRect.X
			newRect.Y = prevRect.Y + prevRect.Height + c.vspacing
			newRect.Width = c.Rect.Width // Resize width to container's width
		default:
			newRect.X = c.Rect.X + child.Rect.X
			newRect.Y = c.Rect.Y + child.Rect.Y
		}

		paddedRect := newRect.WithPadding(c.hpadding, c.vpadding)
		
		// Calculate relative new position for logging
		relativeNewX := paddedRect.X - c.Rect.X
		relativeNewY := paddedRect.Y - c.Rect.Y

		// Log new calculated relative position
		log.Printf("RelayoutChildren: Child '%s' calculated relative new position: X=%d, Y=%d, W=%d, H=%d", child.Title(), relativeNewX, relativeNewY, paddedRect.Width, paddedRect.Height)

		child.Move(relativeNewX, relativeNewY)
		if c.layoutDirection == LayoutVer {
			child.Resize(paddedRect.Width, child.Rect.Height)
		} else if c.layoutDirection == LayoutHor {
			child.Resize(child.Rect.Width, paddedRect.Height)
		} else {
			child.Resize(paddedRect.Width, paddedRect.Height)
		}
		prevRect = newRect
		child.PaintOnce()
	}
	c.pvsChildRect = prevRect
}

func (c *Container) SetSpacing(dx, dy int) {
	c.hspacing, c.vspacing = dx, dy
	c.RelayoutChildren() // Relayout after changing spacing
}

// WidgetState represents the current state of a UI widget

// rePaint redraws the window with a given state, if it's a button.
func (w *Window) rePaint(s WidgetState) {

	if w.isButton == true {
		w.drawBackground(s)
		w.finishPaint(w.ximg)
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

	// For now, simply repaint the entire window on expose events.
	// More optimized drawing could be implemented later (e.g., only redraw exposed region).
	w.PaintOnce()
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
	w.bgcolor = color.RGBA{0, 0, 0, 255} // Set default window background to black

	// mask := xproto.GcForeground | xproto.GcGraphicsExposures
	// values := []uint32{s.BlackPixel, 0}
		win.Create(parent, r.X, r.Y, r.Width, r.Height, 0)
		win.Map()
	
		// win.Create(parent, r.X, r.Y, r.Width, r.Height, mask, values...)
		if err != nil {
			log.Fatal(err)
		}
		// win.MoveResize(r.X, r.Y, r.Width, r.Height)
		// if p == nil {
		// 	win.Change(xproto.CwBackPixel, 0x000000) // Set root window background to black
		// } else {
		// win.Change(xproto.CwBackPixel, 0xFFAA00)
		// }

	//if p == nil {
	win.Listen(xproto.EventMaskKeyPress, xproto.EventMaskKeyRelease, xproto.EventMaskButtonPress, xproto.EventMaskButtonRelease, xproto.EventMaskExposure, xproto.EventMaskEnterWindow, xproto.EventMaskLeaveWindow, xproto.EventMaskKeyPress)
	//}

	w.Rect = r
	w.prevRect = r

	w.Window = win
	win.Map()

	// It's important that the map comes after setting WMGracefulClose, since
	// the WM isn't obliged to watch updates to the WM_PROTOCOLS property.

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
	if w.isButton {
		w.drawBackground(StateNormal)
	} else {
		w.drawView(StateNormal)
	}
	w.finishPaint(w.ximg)
}

// XWin returns the underlying *xwindow.Window.
func (w *Window) XWin() *xwindow.Window {
	return w.Window
}

// XProtoWin returns the X protocol window ID.
func (w *Window) XProtoWin() xproto.Window {
	return w.Window.Id
}

func (w *Window) Reparent(newParent xproto.Window, x, y int) error {
	err := xproto.ReparentWindow(w.X().Conn(), w.Id, newParent, int16(x), int16(y)).Check()
	if err != nil {
		return fmt.Errorf("failed to reparent window: %w", err)
	}
	// Update the internal Rect to reflect the new position relative to the new parent
	w.Rect.X = x
	w.Rect.Y = y
	return nil
}

func (w *Window) Move(x, y int) {
	w.Window.Move(x, y)
	w.Rect.X = x
	w.Rect.Y = y
}

func (w *Window) Resize(width, height int) {
	w.Window.Resize(width, height)
	w.Rect.Width = width
	w.Rect.Height = height
}

func (w *Window) Align(alignment Alignment) {
	parent, err := w.Parent()
	if err != nil {
		log.Printf("Error getting parent window for alignment: %v", err)
		return
	}
	parentGeom, err := parent.Geometry()
	if err != nil {
		log.Printf("Error getting parent geometry for alignment: %v", err)
		return
	}



	windowGeom := w.Rect


	newX, newY := windowGeom.X, windowGeom.Y

	switch alignment {
	case AlignTopLeft:
		newX = parentGeom.X()
		newY = parentGeom.Y()
	case AlignTopCenter:
		newX = parentGeom.X() + (parentGeom.Width()-windowGeom.Width)/2
		newY = parentGeom.Y()
	case AlignTopRight:
		newX = parentGeom.X() + parentGeom.Width() - windowGeom.Width
		newY = parentGeom.Y()
	case AlignMiddleLeft:
		newX = parentGeom.X()
		newY = parentGeom.Y() + (parentGeom.Height()-windowGeom.Height)/2
	case AlignCenter:
		newX = parentGeom.X() + (parentGeom.Width()-windowGeom.Width)/2
		newY = parentGeom.Y() + (parentGeom.Height()-windowGeom.Height)/2
	case AlignMiddleRight:
		newX = parentGeom.X() + parentGeom.Width() - windowGeom.Width
		newY = parentGeom.Y() + (parentGeom.Height()-windowGeom.Height)/2
	case AlignBottomLeft:
		newX = parentGeom.X()
		newY = parentGeom.Y() + parentGeom.Height() - windowGeom.Height
	case AlignBottomCenter:
		newX = parentGeom.X() + (parentGeom.Width()-windowGeom.Width)/2
		newY = parentGeom.Y() + parentGeom.Height() - windowGeom.Height
	case AlignBottomRight:
		newX = parentGeom.X() + parentGeom.Width() - windowGeom.Width
		newY = parentGeom.Y() + parentGeom.Height() - windowGeom.Height
	}


	w.Move(newX-parentGeom.X(), newY-parentGeom.Y())
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

	w.ximg.XPaintRects(w.Id, image.Rect(0, 0, rr.Dx(), rr.Dy()))
	w.ximg.XDraw()

}

// ReDrawImage redraws the window from its raw image buffer.
func (w *Window) ReDrawImage() {
	if w.rawimage == nil {
	
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
