# x11ui
Simple UI framework using github.com/BurntSushi/xgbutil package using Go. 

Look for documentation here https://godoc.org/github.com/wiless/x11ui 


![ScreenShot](x11ui.gif?raw=true "Title")



# Sample Usage

<a href="http://www.youtube.com/watch?feature=player_embedded&v=y7JasppN1FI" target="_blank"><img src="http://img.youtube.com/vi/y7JasppN1FI/0.jpg" alt="IMAGE ALT TEXT HERE" width="240" height="180" border="10" /></a>

# Theme Usage

Here's how you can use the new theme functions:

**1. Setting a Default Theme:**

The `init()` function in `theme.go` already sets `CurrentTheme` to `DarkTheme` by default. You can change this in `main.go` or any other initialization point.

```go
package main

import (
	"github.com/wiless/x11ui"
	"image/color"
)

func main() {
	app := x11ui.NewApplication("My App", 800, 600, false, false)

	// To use the default LightTheme:
	x11ui.UseLightTheme()

	// To use a DarkTheme with a custom accent color (e.g., red):
	accentColor := color.RGBA{R: 255, A: 255} // Red
	x11ui.SetTheme(x11ui.DarkThemeWithAccent(accentColor))

	// To use a LightTheme with a custom accent color (e.g., blue):
	// accentColor := color.RGBA{B: 255, A: 255} // Blue
	// x11ui.SetTheme(x11ui.LightThemeWithAccent(accentColor))

	// ... rest of your application setup
	app.Show()
}
```

**2. Switching Themes at Runtime:**

You can switch themes dynamically using `x11ui.SetTheme()`, `x11ui.UseDarkTheme()`, or `x11ui.UseLightTheme()`. After changing the theme, you'll typically need to trigger a redraw of your UI elements to apply the new colors.

```go
// Example: A button to toggle between dark and light themes
func createThemeToggleButton(app *x11ui.Application) *x11ui.Window {
	btn := app.AddButton("Toggle Theme")
	btn.OnClick(func() {
		if x11ui.CurrentTheme == x11ui.DarkTheme {
			x11ui.UseLightTheme()
		} else {
			x11ui.UseDarkTheme()
		}
		// You might need to explicitly redraw all widgets or the main application window
		// depending on how your UI is structured.
		app.AppWin().RePaint(x11ui.StateNormal) // Repaint the main window
		// For individual widgets, you might need to call their RePaint methods
		// or have a mechanism to propagate theme changes.
	})
	return btn
}
```

**3. Using Accent Themes:**

The `DarkThemeWithAccent(color)` and `LightThemeWithAccent(color)` functions allow you to customize the theme with a specific accent color. The accent color's hue will be used to derive complementary colors for elements like foreground, bar, and checked checkboxes.

```go
// Example: Create a dark theme with a vibrant purple accent
purpleAccent := color.RGBA{R: 128, G: 0, B: 128, A: 255}
customDarkTheme := x11ui.DarkThemeWithAccent(purpleAccent)
x11ui.SetTheme(customDarkTheme)

// Example: Create a light theme with a bright orange accent
orangeAccent := color.RGBA{R: 255, G: 165, B: 0, A: 255}
customLightTheme := x11ui.LightThemeWithAccent(orangeAccent)
x11ui.SetTheme(customLightTheme)
```

**Important Considerations:**

*   **Redrawing Widgets:** When you change `x11ui.CurrentTheme`, existing widgets *do not automatically update their colors*. You need to explicitly trigger a repaint or a theme application method for each widget that needs to reflect the new theme. The `RePaint` method on `Window` (and thus on `Widget` since `Window` is embedded) can be used for this.
*   **Theme Structure:** The `Theme` struct defines all the colors. When creating custom themes, ensure all necessary color fields are set.
*   **`go-colorful`:** The theme generation uses `go-colorful` for color manipulation. If you create your own custom themes, you can leverage `colorful.Hsl()`, `colorful.BlendLuvLCh()`, and `Clamped()` for sophisticated color generation.