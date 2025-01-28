package ui

import "github.com/gdamore/tcell"

type Window interface {
	Screen() tcell.Screen
	// the (width, height) of the window in actual coordinates
	Size() (int, int)
	Render()
	// the position of the window in the screen in actual coordinates (x, y)
	Position() (int, int)
	// change the screen of the window
	SetScreen(tcell.Screen)
	// changes the size of the window
	SetSize(int, int)
	// changes the position of the window
	SetPosition(int, int)
}

type InteractiveWindow interface {
	Window
	OnKeyPress(event *tcell.Event)
}

// based on height and width of the screen and the actual width and height of the terminal we can translate screen coordinates to the terminal coordinates
func translateCoordinates(w Window, x, y int) (int, int) {
	window_x, window_y := w.Position()
	// window_x and window_y are the center of the window
	// width and height are the size of the window
	// center - half of the size = starting point of the window

	x = window_x + x
	y = window_y + y
	return x, y
}

func DrawRuneAt(w Window, x, y int, c rune, comb []rune, style tcell.Style) {
	s := w.Screen()
	x, y = translateCoordinates(w, x, y)
	s.SetContent(x, y, c, comb, style)
}

func drawBorder(w Window) {
	width, height := w.Size()

	for y := 0; y < height; y++ {
		DrawRuneAt(w, 0, y, '│', nil, tcell.StyleDefault)
		DrawRuneAt(w, width-1, y, '│', nil, tcell.StyleDefault)
	}

	for x := 0; x < width; x++ {
		DrawRuneAt(w, x, 0, '─', nil, tcell.StyleDefault)
		DrawRuneAt(w, x, height-1, '─', nil, tcell.StyleDefault)
	}
	// connect the borders using the corners
	DrawRuneAt(w, 0, 0, '┌', nil, tcell.StyleDefault)
	DrawRuneAt(w, width-1, 0, '┐', nil, tcell.StyleDefault)
	DrawRuneAt(w, 0, height-1, '└', nil, tcell.StyleDefault)
	DrawRuneAt(w, width-1, height-1, '┘', nil, tcell.StyleDefault)
}

func drawBackground(w Window) {
	width, height := w.Size()

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			DrawRuneAt(w, x, y, ' ', nil, tcell.StyleDefault)
		}
	}
}

func RenderWindow(w Window) {
	drawBackground(w)
	drawBorder(w)
	w.Render()
}
