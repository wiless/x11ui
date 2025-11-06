package x11ui

import (
	"image/color"

	"github.com/lucasb-eyer/go-colorful"
)

// Theme holds all the color definitions for the UI.
type Theme struct {
	BackgroundColor color.Color
	ForegroundColor color.Color
	LineColor       color.Color
	TextColor       color.Color
	BarColor        color.Color
	CheckboxCheckedColor   color.Color
	CheckboxUncheckedColor color.Color
	CheckboxBorderColor    color.Color
	BaseHue         float64 // Base hue for generating color variations
}

// CurrentTheme is the global instance of the active theme.
var CurrentTheme *Theme

// Predefined themes
var DarkTheme *Theme
var LightTheme *Theme

func init() {
	// Initialize predefined themes
	DarkTheme = createDarkTheme()
	LightTheme = createLightTheme()

	// Set DarkTheme as the initial default theme
	CurrentTheme = DarkTheme
}

// createDarkTheme generates a minimalistic dark theme.
func createDarkTheme() *Theme {
	// Minimalistic dark gray theme
	darkGray := color.RGBA{R: 0x28, G: 0x28, B: 0x28, A: 0xFF} // A common dark gray
	lightGray := color.RGBA{R: 0xBB, G: 0xBB, B: 0xBB, A: 0xFF} // Light gray for foreground/text
	mediumGray := color.RGBA{R: 0x40, G: 0x40, B: 0x40, A: 0xFF} // Medium gray for lines/borders
	accentColor := color.RGBA{R: 0x60, G: 0x60, B: 0x60, A: 0xFF} // A subtle accent for bars/checked states

	return &Theme{
		BackgroundColor:        darkGray,
		ForegroundColor:        lightGray,
		LineColor:              mediumGray,
		TextColor:              color.RGBA{R: 0xFF, G: 0xFF, B: 0xFF, A: 0xFF}, // White text
		BarColor:               accentColor,
		CheckboxCheckedColor:   accentColor,
		CheckboxUncheckedColor: darkGray,
		CheckboxBorderColor:    mediumGray,
		BaseHue:                0.0, // Not using hue for this minimalistic theme
	}
}

// createLightTheme generates a light theme using colorful.
func createLightTheme() *Theme {
	hue := 210.0 // Blue-ish hue
	saturation := 0.3
	lightness := 0.9
	baseColor := colorful.Hsl(hue, saturation, lightness)

	return &Theme{
		BackgroundColor: baseColor.Clamped(),
		ForegroundColor: colorful.Hsl(hue, saturation*0.7, lightness*0.6).Clamped(),
		LineColor:       colorful.Hsl(hue+60, saturation, lightness*0.8).Clamped(),
		TextColor:       color.RGBA{0, 0, 0, 255}, // Pure black for readability
		BarColor:        colorful.Hsl(hue-60, saturation, lightness*0.7).Clamped(),

		CheckboxCheckedColor:   colorful.Hsl(hue, saturation*1.5, lightness*0.4).Clamped(), // More saturated and darker for light theme
		CheckboxUncheckedColor: colorful.Hsl(hue, saturation*0.2, lightness*0.95).Clamped(), // Very light, less saturated
		CheckboxBorderColor:    colorful.Hsl(hue, saturation, lightness*0.7).Clamped(), // Slightly darker border
		BaseHue:         hue,
	}
}

// DarkThemeWithAccent generates a dark theme with a specified accent color.
func DarkThemeWithAccent(accent color.Color) *Theme {
	accentC, _ := colorful.MakeColor(accent)
	hue, _, _ := accentC.Hsl()

	// Start with the base dark theme colors
	theme := createDarkTheme()

	// Adjust colors based on the accent hue
	theme.ForegroundColor = colorful.Hsl(hue, 0.7, 0.8).Clamped()
	theme.BarColor = colorful.Hsl(hue, 0.6, 0.7).Clamped()
	theme.CheckboxCheckedColor = colorful.Hsl(hue, 0.8, 0.6).Clamped()
	theme.BaseHue = hue

	return theme
}

// LightThemeWithAccent generates a light theme with a specified accent color.
func LightThemeWithAccent(accent color.Color) *Theme {
	accentC, _ := colorful.MakeColor(accent)
	hue, _, _ := accentC.Hsl()

	// Start with the base light theme colors
	theme := createLightTheme()

	// Adjust colors based on the accent hue
	theme.ForegroundColor = colorful.Hsl(hue, 0.7, 0.3).Clamped()
	theme.BarColor = colorful.Hsl(hue, 0.6, 0.4).Clamped()
	theme.CheckboxCheckedColor = colorful.Hsl(hue, 0.8, 0.5).Clamped()
	theme.BaseHue = hue

	return theme
}

// SetTheme allows updating the global theme at runtime.
// To apply the new theme to existing widgets, they would need to be explicitly redrawn
// or have their LoadTheme/ApplyTheme methods called.
func SetTheme(newTheme *Theme) {
	CurrentTheme = newTheme
}

// UseDarkTheme sets the global theme to the predefined DarkTheme.
func UseDarkTheme() {
	SetTheme(DarkTheme)
}

// UseLightTheme sets the global theme to the predefined LightTheme.
func UseLightTheme() {
	SetTheme(LightTheme)
}
