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

// createDarkTheme generates a dark theme using colorful.
func createDarkTheme() *Theme {
	hue := 150.0 // Green-cyan hue
	saturation := 0.6
	lightness := 0.4
	baseColor := colorful.Hsl(hue, saturation, lightness)

	return &Theme{
		BackgroundColor: baseColor.Clamped(),
		ForegroundColor: colorful.Hsl(hue, saturation*0.7, lightness*1.2).Clamped(),
		LineColor:       colorful.Hsl(hue+60, saturation, lightness*0.9).Clamped(),
		TextColor:       color.RGBA{255, 255, 255, 255}, // Pure white for readability
		BarColor:        colorful.Hsl(hue-60, saturation, lightness*1.1).Clamped(),

		CheckboxCheckedColor:   colorful.Hsl(hue, saturation*1.2, lightness*1.8).Clamped(), // More vibrant and brighter for dark theme
		CheckboxUncheckedColor: colorful.Hsl(hue, saturation*0.5, lightness*0.2).Clamped(), // Darker, less saturated
		CheckboxBorderColor:    colorful.Hsl(hue, saturation, lightness*0.8).Clamped(), // Slightly lighter border
		BaseHue:         hue,
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
