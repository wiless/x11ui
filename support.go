package x11ui

import "image"

type Rect struct {
	X, Y          int
	Width, Height int
}

func (r *Rect) ShiftRight(dx int) *Rect {
	r.X += dx
	return r

}
func (r *Rect) ShiftDown(dy int) *Rect {
	r.Y += dy
	return r
}
func (r *Rect) ReSize(w, h int) *Rect {
	r.Width, r.Height = w, h
	return r
}
func (r *Rect) Grow(dw, dh int) *Rect {
	r.Width += dw
	r.Height += dh

	return r
}

func (r *Rect) array() []int {
	result := []int{r.X, r.Y, r.Width, r.Height}
	return result
}

func newRect(dims ...int) Rect {
	r := Rect{0, 0, 100, 50}
	for i, v := range dims {
		switch i {
		case 0:
			r.X = v
		case 1:
			r.Y = v
		case 2:
			r.Width = v
		case 3:
			r.Height = v
		}
	}
	return r
}

func (r *Rect) CenterX() int {
	return r.X + r.Width/2
}

func (r *Rect) CenterY() int {
	return r.Y + r.Height/2
}

func (r *Rect) Center() (x, y int) {
	return r.CenterX(), r.CenterY()
}

func (r *Rect) ImageRect() image.Rectangle {
	return image.Rect(0, 0, r.Width, r.Height)

}
