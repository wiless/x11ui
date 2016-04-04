package x11ui

import (
	"image"
	"image/color"

	"github.com/BurntSushi/xgbutil/xgraphics"
)

func StrokeBorder(img *xgraphics.Image, clr xgraphics.BGRA, margin, width int) {

	outset := img.Rect
	// outset.Max.Sub(image.Point{5, 5})
	size := outset.Size()
	inset := outset.Inset(width)
	for x := 0; x < size.X; x++ {
		for y := 0; y < size.Y; y++ {
			xcond := (outset.Min.X >= x && inset.Min.X > x) || (inset.Max.X < x)
			ycond := (outset.Min.Y >= y && inset.Min.Y > y) || (inset.Max.Y < y)
			if xcond || ycond {
				img.SetBGRA(x, y, clr)
			}
		}
	}
}
func StrokeBorderImg(img *image.RGBA, clr color.Color, margin, width int) {

	outset := img.Rect
	// outset.Max.Sub(image.Point{5, 5})
	size := outset.Size()
	inset := outset.Inset(width)
	for x := 0; x < size.X; x++ {
		for y := 0; y < size.Y; y++ {
			xcond := (outset.Min.X >= x && inset.Min.X > x) || (inset.Max.X < x)
			ycond := (outset.Min.Y >= y && inset.Min.Y > y) || (inset.Max.Y < y)
			if xcond || ycond {
				// xg.SetBGRA(x, y, clr)
				img.SetRGBA(x, y, clr.(color.RGBA))
			}
		}
	}
}
