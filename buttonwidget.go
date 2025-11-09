package x11ui

import (
	"image/color"
	"log"

	"github.com/BurntSushi/xgbutil"
	"github.com/BurntSushi/xgbutil/mousebind"
	"github.com/BurntSushi/xgbutil/xevent"
	"github.com/BurntSushi/xgbutil/xgraphics"
	"github.com/lucasb-eyer/go-colorful"
)

// ButtonWidget represents a clickable button UI element.
type ButtonWidget struct {
	*Widget
	text        string
	isToggle    bool
	checked     bool
	onClickFn   func()
	releaseFn   func(X *xgbutil.XUtil, e xevent.ButtonReleaseEvent)
	normalColor color.Color
	hoverColor  color.Color
	pressColor  color.Color
	textColor   color.Color
	fsize       float64
}

// NewButtonWidget creates a new ButtonWidget.
func NewButtonWidget(title string, p *Window, dims ...int) *ButtonWidget {
	if p == nil {
		log.Fatal("Cannot Create ButtonWidget without Application")
	}


	bw := &ButtonWidget{
		text:      title,
		isToggle:  false,
		checked:   false,
		fsize:     12,
	}
	bw.Widget = WidgetFactory(p, dims...)
	bw.Widget.isButton = true
	bw.SetTitle(title) // Set the window title for debugging/WM
	bw.init()
	bw.SetLabel(title) // Initial rendering of the label

	return bw
}

// SetTitle sets the title of the button. This will update the visible label.
func (bw *ButtonWidget) SetTitle(title string) {
	bw.SetLabel(title)
}

// SetLabel sets the text displayed on the button.
func (bw *ButtonWidget) Paint() {
	bw.updateButtonAppearance()
}

func (bw *ButtonWidget) SetLabel(lbl string) {
	bw.text = lbl
	bw.updateButtonAppearance()
}

// SetOnClick sets the function to be called when the button is clicked.
func (bw *ButtonWidget) SetOnClick(fn func()) {
	bw.onClickFn = fn
	bw.ClkFn = bw.handleButtonClick // Override the base Widget's ClkFn
}

// SetToggle makes the button a toggle button.
func (bw *ButtonWidget) SetToggle(toggle bool) {
	bw.isToggle = toggle
}

// IsChecked returns true if the toggle button is currently checked.
func (bw *ButtonWidget) IsChecked() bool {
	return bw.checked
}

// SetChecked sets the checked state of a toggle button.
func (bw *ButtonWidget) SetChecked(checked bool) {
	if bw.isToggle {
		bw.checked = checked
		bw.updateButtonAppearance()
	}
}

// SetFontSize sets the font size of the button's text.
func (bw *ButtonWidget) SetFontSize(size float64) {
	bw.fsize = size
	bw.gc.SetFontSize(bw.fsize) // Update the graphics context as well
	bw.updateButtonAppearance()
}

// SetKeybFn sets the function to be called when a keyboard event is received.
func (bw *ButtonWidget) SetKeybFn(fn func(key string)) {
	bw.KeybFn = fn
}

func (bw *ButtonWidget) init() {
	bw.LoadTheme("") // Load default theme colors
	bw.gc.SetFontSize(float64(bw.fsize))
	bw.updateButtonAppearance() // Initial draw
	bw.LeaveFn = bw.handleLeave
	bw.HoverFn = bw.handleHover
	bw.ClkFn = bw.handleButtonClick
	bw.releaseFn = bw.handleButtonRelease
	mousebind.ButtonReleaseFun(bw.releaseFn).Connect(bw.xu, bw.xwin.Id, "1", false, true)
}

func (bw *ButtonWidget) LoadTheme(str string) {
	bw.Widget.LoadTheme(str) // Call base widget theme loading
	bw.bgColor = bw.normalColor // Set initial background
	bw.txtColor = bw.textColor
}

func (bw *ButtonWidget) handleHover() {
	bw.Widget.state = StateHovered
	bw.updateButtonAppearance()
}

func (bw *ButtonWidget) handleLeave() {
	bw.Widget.state = StateNormal
	bw.updateButtonAppearance()
}

func (bw *ButtonWidget) handleButtonClick() {
	if bw.isToggle {
		bw.checked = !bw.checked
	}
	// Briefly show pressed state
	bw.Widget.state = StatePressed
	bw.updateButtonAppearance()

	if bw.onClickFn != nil {
		bw.onClickFn()
	}
}

func (bw *ButtonWidget) handleButtonRelease(X *xgbutil.XUtil, e xevent.ButtonReleaseEvent) {
	if bw.isToggle && bw.checked {
		bw.Widget.state = StateChecked // Custom state for checked toggle button
	} else {
		bw.Widget.state = StateNormal
	}
	bw.updateButtonAppearance()
}

func (bw *ButtonWidget) updateButtonAppearance() {
	var baseColor colorful.Color
	var currentLineColor color.Color

	// A toggle button only looks different when it's checked.
	// Otherwise, it behaves like a normal button.
	if bw.isToggle && bw.checked {
		baseColor = toColorful(CurrentTheme.CheckboxCheckedColor)
		currentLineColor = CurrentTheme.CheckboxBorderColor
	} else {
		baseColor = toColorful(CurrentTheme.BarColor)
		currentLineColor = CurrentTheme.LineColor
	}

	var currentBgColor color.Color
	switch bw.Widget.state {
	case StateHovered:
		// Make the color 20% darker for hover
		h, s, l := baseColor.Hsl()
		currentBgColor = colorful.Hsl(h, s, l*0.8).Clamped()
	case StatePressed:
		// Make the color 50% darker for press
		h, s, l := baseColor.Hsl()
		currentBgColor = colorful.Hsl(h, s, l*0.5).Clamped()
	default:
		currentBgColor = baseColor
	}

	// Dynamically determine text color based on background luminance
	rVal, gVal, bVal, _ := toColorful(currentBgColor).RGBA()
	avg := (float64(rVal>>8) + float64(gVal>>8) + float64(bVal>>8)) / 3.0 / 255.0

	if avg > 0.5 { // If background is light, use dark text
		bw.textColor = color.RGBA{R: 0, G: 0, B: 0, A: 255} // Black
	} else { // If background is dark, use light text
		bw.textColor = color.RGBA{R: 255, G: 255, B: 255, A: 255} // White
	}
	bw.txtColor = bw.textColor

	// Set the embedded Widget's colors
	bw.Widget.bgColor = currentBgColor
	bw.Widget.lineColor = currentLineColor

	// Draw background
	bw.drawBackground()

	// Draw text
	tw, th := xgraphics.Extents(systemFont, float64(bw.fsize), bw.text)
	xpos, ypos := (bw.Width()-tw)/2, (bw.Height()-th)/2

	bw.canvas.Text(xpos, ypos, bw.textColor, float64(bw.fsize), systemFont, bw.text)

	// Draw border
	bw.drawBorder(bw.Widget.state) // This will use bw.Widget.bgColor or bw.Widget.lineColor

	bw.updateCanvas()
}

