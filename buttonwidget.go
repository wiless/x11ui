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
	state       WidgetState
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
	log.Printf("NewButtonWidget: Creating button '%s' with dims %v", title, dims)

	bw := &ButtonWidget{
		text:      title,
		isToggle:  false,
		checked:   false,
		state:     StateNormal,
		fsize:     12,
	}
	bw.Widget = WidgetFactory(p, dims...)
	bw.SetTitle(title) // Set the window title for debugging/WM
	bw.init()
	bw.SetLabel(title) // Initial rendering of the label

	return bw
}

// SetLabel sets the text displayed on the button.
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

func (bw *ButtonWidget) init() {
	bw.LoadTheme("") // Load default theme colors
	bw.gc.SetFontSize(bw.fsize)
	bw.updateButtonAppearance() // Initial draw
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
	bw.state = StateHovered
	bw.updateButtonAppearance()
}

func (bw *ButtonWidget) handleLeave() {
	log.Println("handleLeave called, setting state to StateNormal")
	bw.state = StateNormal
	bw.updateButtonAppearance()
}

func (bw *ButtonWidget) handleButtonClick() {
	// Briefly show pressed state
	bw.state = StatePressed
	bw.updateButtonAppearance()

	if bw.onClickFn != nil {
		bw.onClickFn()
	}


}

func (bw *ButtonWidget) handleButtonRelease(X *xgbutil.XUtil, e xevent.ButtonReleaseEvent) {
	if bw.isToggle && bw.checked {
		bw.state = StateChecked // Custom state for checked toggle button
	} else {
		bw.state = StateNormal
	}
	bw.updateButtonAppearance()
}

func (bw *ButtonWidget) updateButtonAppearance() {
	// Dynamically recalculate colors based on the current theme each time the button is drawn
	baseColor := toColorful(CurrentTheme.BarColor)
	bw.normalColor = CurrentTheme.BarColor
	bw.hoverColor = baseColor.BlendLuvLCh(colorful.Color{R: 1, G: 1, B: 1}, 0.2).Clamped()
	bw.pressColor = baseColor.BlendLuvLCh(colorful.Color{R: 0, G: 0, B: 0}, 0.4).Clamped()

	// Dynamically determine text color based on background luminance
	rVal, gVal, bVal, _ := baseColor.RGBA()
	avg := (float64(rVal>>8) + float64(gVal>>8) + float64(bVal>>8)) / 3.0 / 255.0

	if avg > 0.5 { // If background is light, use dark text
		bw.textColor = color.RGBA{R: 0, G: 0, B: 0, A: 255} // Black
	} else { // If background is dark, use light text
		bw.textColor = color.RGBA{R: 255, G: 255, B: 255, A: 255} // White
	}
	bw.txtColor = bw.textColor

	// Determine background color based on state
	var currentBgColor color.Color
	var currentLineColor color.Color

	switch bw.state {
	case StateNormal:
		currentBgColor = toRGBA(bw.normalColor)
		currentLineColor = CurrentTheme.LineColor
	case StateHovered:
		currentBgColor = toRGBA(bw.hoverColor)
		currentLineColor = CurrentTheme.LineColor
	case StatePressed:
		currentBgColor = toRGBA(bw.pressColor)
		currentLineColor = CurrentTheme.LineColor
	case StateChecked: // For toggle buttons when checked
		currentBgColor = CurrentTheme.CheckboxCheckedColor // Use theme color for checked state
		currentLineColor = CurrentTheme.LineColor
	default:
		if bw.isToggle && bw.checked {
			currentBgColor = CurrentTheme.CheckboxCheckedColor
			currentLineColor = CurrentTheme.LineColor
		} else {
			currentBgColor = bw.normalColor
			currentLineColor = CurrentTheme.LineColor
		}
	}
	log.Printf("ButtonWidget.updateButtonAppearance: state=%v, currentBgColor=%v, currentLineColor=%v", bw.state, currentBgColor, currentLineColor)

	// Set the embedded Widget's colors
	bw.Widget.bgColor = currentBgColor
	bw.Widget.lineColor = currentLineColor

	// Draw background
	bw.drawBackground()

	// Draw text
	tw, th := xgraphics.Extents(systemFont, bw.fsize, bw.text)
	xpos, ypos := (bw.Width()-tw)/2, (bw.Height()-th)/2
	log.Printf("ButtonWidget.updateButtonAppearance: Drawing text '%s' at (%d, %d) with color %v and font size %f", bw.text, xpos, ypos, bw.textColor, bw.fsize)
	bw.canvas.Text(xpos, ypos, bw.textColor, bw.fsize, systemFont, bw.text)

	// Draw border
	bw.drawBorder(bw.state) // This will use bw.Widget.bgColor or bw.Widget.lineColor

	bw.updateCanvas()
}

