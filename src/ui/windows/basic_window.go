package windows

import (
	"github.com/ObserverUnit/arut/src/ui"
	"github.com/gdamore/tcell"
)

type BasicWindow struct {
	screen        tcell.Screen
	width, height int
	content       string
	x, y          int
}

func (w *BasicWindow) Screen() tcell.Screen {
	return w.screen
}

func (w *BasicWindow) Size() (int, int) {
	return w.width, w.height
}

func (w *BasicWindow) Position() (int, int) {
	return w.x, w.y
}

// sets the position of the window in precentages of the screen
// x and y must be between 0 and 100
func (w *BasicWindow) SetPositionRelative(x, y int) {
	max_width, max_height := w.Screen().Size()
	width, height := w.Size()

	if x < 0 || x > 100 {
		panic("window x must be between 0 and 100")
	}

	if y < 0 || y > 100 {
		panic("window y must be between 0 and 100")
	}

	w.SetPosition(x*(max_width-width)/100, y*(max_height-height)/100)

}

func (w *BasicWindow) SetPosition(x, y int) {
	max_width, max_height := w.screen.Size()
	max_x := max_width - w.width
	max_y := max_height - w.height

	if x >= 0 {
		w.x = min(x, max_x)
	}
	if y >= 0 {
		w.y = min(y, max_y)
	}
}

// sets the size of the window
func (w *BasicWindow) SetSizeRelative(width, height int) {
	max_width, max_height := w.Screen().Size()
	if width < 0 || width > 100 {
		panic("width must be between 0 and 100")
	}

	if height < 0 || height > 100 {
		panic("height must be between 0 and 100")
	}

	w.SetSize(width*max_width/100, height*max_height/100)
}

func (w *BasicWindow) SetSize(width, height int) {
	max_width, max_height := w.screen.Size()

	if width >= 0 {
		w.width = min(width, max_width)
	}

	if height >= 0 {
		w.height = min(height, max_height)
	}

	// fix the position after the size has been changed
	w.SetPosition(w.x, w.y)
}

func (w *BasicWindow) Render() {
	width, height := w.Size()
	// x and y are 2 to account for the borders
	x := 0
	y := 0
	for _, c := range w.content {
		if c == '\n' {
			x = 0
			y += 1
			continue
		}

		if x >= width-2*2 {
			x = 0
			y += 1
		}

		if y >= height-2*2 {
			break
		}

		ui.DrawRuneAt(w, x+2, y+2, c, nil, tcell.StyleDefault)
		x += 1
	}
}

func (w *BasicWindow) SetScreen(screen tcell.Screen) {
	w.screen = screen
}

func (w *BasicWindow) SetContent(content string) {
	w.content = content
}

func (w *BasicWindow) DrawRuneAtBody(x, y int, r rune, comb []rune, style tcell.Style) {
	draw_x, draw_y := x, y

	if draw_x >= w.width-4 {
		// how many w.width-2 x contains
		rows := draw_x / (w.width - 4)
		draw_y += rows
		draw_x -= rows * (w.width - 4)
	}

	if draw_y >= w.height-4 {
		return
	}

	draw_x += 2
	draw_y += 2

	ui.DrawRuneAt(w, draw_x, draw_y, r, comb, style)
}
func NewBasicWindow(width, height, x, y int, content string) *BasicWindow {
	if width < 0 || width > 100 {
		panic("width must be between 0 and 100")
	}

	if height < 0 || height > 100 {
		panic("height must be between 0 and 100")
	}

	window := &BasicWindow{
		screen:  ui.Screen(),
		content: content,
		width:   0,
		height:  0,
		x:       0,
		y:       0,
	}

	window.SetSizeRelative(width, height)
	window.SetPositionRelative(x, y)
	return window
}
