package x11ui

import (
	"image/color"
	"time"

	"log"
)

type ScrollView struct {
	*Widget
	title         string
	autoresize    bool
	align         AlignMode
	margins       int
	basewidgets   []*Widget
	viewWidget    *Widget
	scrollvisible bool
	stepSize      int
	stepLimits    int
	stepcount     int
}

func NewScrollView(title string, p *Window, dims ...int) *ScrollView {
	if p == nil {
		log.Fatal("Cannot Create Scroll View")
	}
	scrl := new(ScrollView)
	scrl.Widget = WidgetFactory(p, dims...)

	scrl.init()

	// scrl.SetLabel(title)
	return scrl
}

func (s *ScrollView) init() {
	s.margins = 0
	s.scrollvisible = false
	s.SetStepSize(10) // default step size 10px
	s.SetMaxSteps(-1)
	w, h := s.Width(), s.Height()
	s.viewWidget = s.CreateChild(s.margins, s.margins, w-2*s.margins, h-2*s.margins)
	s.viewWidget.Win().Detach()
	s.Win().SetBGcolor(color.RGBA{255, 255, 0, 255})
	// s.updateCanvas()
	s.SetFontSize(12)
	// scrl.updateCanvas()

}

func (s *ScrollView) SetBackground(c color.Color) {
	s.viewWidget.Win().SetBGcolor(c)
	// s.Win().SetBGcolor(c)
	// s.viewWidget.drawBackground()
	// s.updateCanvas()
	// log.Println("I am here")
	// s.viewWidget.Win().SetBGcolor(color.RGBA{255, 0, 0, 20})

}

// returns if the scrollbars are visible
func (s *ScrollView) IsVisible() bool {
	return s.scrollvisible
}

func (s *ScrollView) ShowScrollBars(show bool) {

	if show {
		if s.scrollvisible {
			// Already visible
			return
		}
		s.margins = 30
		s.scrollvisible = true
		s.createBaseWidgets()
	} else {
		if !s.scrollvisible {
			// Already hidden
			return
		}
		s.scrollvisible = false
		for i, v := range s.basewidgets {
			v.Close()
			s.basewidgets[i] = nil
			// .SetBackground(color.RGBA{100 + i*20, 200, 200, 0})
		}
		// s.basewidgets = nil
		s.margins = 0

	}
	w, h := s.Width(), s.Height()
	s.viewWidget.Resize(s.margins, s.margins, w-2*s.margins, h-2*s.margins)
}

// createBaseWidgets creates basic widgets like the inset widget where all child widgets will be added
func (s *ScrollView) createBaseWidgets() {
	s.basewidgets = make([]*Widget, 4)
	s.margin = 2
	w, h := s.Width(), s.Height()
	// top & bottom
	s.basewidgets[0] = s.CreateChild(0, 0, w, s.margins)
	s.basewidgets[1] = s.CreateChild(0, h-s.margins, w, s.margins)
	// left & right
	s.basewidgets[2] = s.CreateChild(0, s.margins, s.margins, h-2*s.margins)
	s.basewidgets[3] = s.CreateChild(w-s.margins, s.margins, s.margins, h-2*s.margins)

	var i uint8
	for i = 0; i < 4; i++ {
		s.basewidgets[i].SetBackground(color.RGBA{100 + i*20, 200, 200, 0})
	}

	s.basewidgets[0].ClkFn = s.ScrollUp
	s.basewidgets[1].ClkFn = s.ScrollDown

}

func (s *ScrollView) SetStepSize(step int) {
	s.stepSize = step
	s.stepcount = 0
}
func (s *ScrollView) SetMaxSteps(maxsteps int) {
	s.stepLimits = maxsteps
	s.stepcount = 0
}

func (s *ScrollView) ScrollUp() {

	if s.stepcount < s.stepLimits || s.stepLimits == -1 {

		s.ScrollChilds(0, -s.stepSize)
		s.stepcount++
	}

}

func (s *ScrollView) ScrollDown() {
	if s.stepcount > 0 {
		// log.Println("I will ANimate")
		// // Animate Duration
		// // aDuration:=1000
		// aSteps:=10;
		// for i := 1; i <= aSteps; i++ {
		// s.ScrollChilds(0, s.stepSize*i/aSteps)
		// time.Sleep(100*time.Millisecond)		
		// }


		// Non-Animated Scroll
		s.ScrollChilds(0, s.stepSize)

		s.stepcount--
	}
}
func (s *ScrollView) ScrollChilds(dx, dy int) {


	fn:= func(ddx,ddy int) {
	
		for _, v := range s.viewWidget.childs {
		x, y := v.X(), v.Y()

		x += ddx
		y += ddy
		v.Win().Move(x, y)
		}
	}

		// log.Println("I will ANimate")
		// // Animate Duration
		// // aDuration:=1000
		aSteps:=10;
		tinydx,tinydy:=dx/aSteps, dy/aSteps
		for i := 1; i <= aSteps; i++ {
		fn(tinydx,tinydy)
		time.Sleep(50*time.Millisecond)		
		}

}

func (s *ScrollView) AddWidget(w *Widget) {

	s.viewWidget.appendChild(w)
}

func (s *ScrollView) View() *Window {
	return s.viewWidget.Win()
}
