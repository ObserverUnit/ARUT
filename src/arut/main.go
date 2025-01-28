// width:182 height:46
package main

import (
	"fmt"
	"os"
	"strings"

	. "github.com/ObserverUnit/arut/src/ui"
	"github.com/ObserverUnit/arut/src/ui/windows"
	"github.com/gdamore/tcell"
)

type EditorMode int

const (
	Normal EditorMode = iota
	Insert
)

var cursorStyle = tcell.StyleDefault.Reverse(true)

type EditorWindow struct {
	windows.BasicWindow
	mode    EditorMode
	cursorX int
	cursorY int
	content []string
}

func newEditorWindow(width, height, x, y int) *EditorWindow {
	inner := windows.NewBasicWindow(width, height, x, y, "")
	return &EditorWindow{
		BasicWindow: *inner,
		mode:        EditorMode(Insert),
		cursorX:     0,
		cursorY:     0,
		content:     []string{""},
	}
}

func (w *EditorWindow) addRune(c rune) {
	if c == '\n' {
		w.cursorY += 1
		if w.cursorY >= len(w.content) {
			w.content = append(w.content, "")
		} else {
			w.content = append(w.content, "")
			copy(w.content[w.cursorY+1:], w.content[w.cursorY:])
			w.content[w.cursorY] = ""
		}
		// move the content from the line above to the current line (empty line)
		currentLine := w.content[w.cursorY-1]
		if w.cursorX+1 < len(currentLine) {
			w.content[w.cursorY] = currentLine[w.cursorX:]
			w.content[w.cursorY-1] = currentLine[:w.cursorX]
		}

		w.cursorX = 0
	} else {
		mut := &w.content[w.cursorY]

		cursor_left := (*mut)[:w.cursorX]
		cursor_right := (*mut)[w.cursorX:]
		*mut = cursor_left + string(c) + cursor_right
		w.cursorX += 1
	}
}

func (w *EditorWindow) removeRune() {
	if w.cursorY <= 0 && w.cursorX <= 0 {
		return
	}

	if w.cursorX <= 0 {
		w.cursorY -= 1
		w.cursorX = len(w.content[w.cursorY])
		if w.cursorY < len(w.content) {
			// copy the content of the line below to the current line and remove the last line
			w.content[w.cursorY] += w.content[w.cursorY+1]
			w.content = w.content[:len(w.content)-1]
		}
		return
	}

	mut := &w.content[w.cursorY]

	left := (*mut)[:w.cursorX-1]
	right := (*mut)[w.cursorX:]

	(*mut) = left + right

	w.cursorX -= 1

}

func (w *EditorWindow) moveCursorBy(x, y int) {
	if w.cursorY+y >= len(w.content) {
		return
	}

	if w.cursorX+x < 0 {
		return
	}

	if w.cursorY+y < 0 {
		return
	}

	w.cursorY += y

	if w.cursorX+x >= len(w.content[w.cursorY]) {
		w.cursorX = len(w.content[w.cursorY])
	} else {
		w.cursorX += x
	}
}

func (w *EditorWindow) OnNormalModeKeyPress(event *tcell.Event) {
	switch ev := (*event).(type) {
	case *tcell.EventKey:
		switch ev.Key() {
		case tcell.KeyCtrlC:
			w.Screen().Fini()
			os.Exit(0)
		case tcell.KeyRune:
			switch key := ev.Rune(); key {
			case 'i':
				w.mode = EditorMode(Insert)
			}
		}
	}
}

func (w *EditorWindow) OnInsertModeKeyPress(event *tcell.Event) {
	switch ev := (*event).(type) {
	case *tcell.EventKey:
		switch ev.Key() {
		case tcell.KeyEsc:
			w.mode = EditorMode(Normal)
		case tcell.KeyRune:
			w.addRune(ev.Rune())
		case tcell.KeyBackspace | tcell.KeyBackspace2:
			w.removeRune()
		case tcell.KeyEnter:
			w.addRune('\n')
		}
	}
}

func (w *EditorWindow) OnKeyPress(event *tcell.Event) {
	switch ev := (*event).(type) {
	case *tcell.EventKey:
		switch ev.Key() {
		case tcell.KeyLeft:
			w.moveCursorBy(-1, 0)
		case tcell.KeyRight:
			w.moveCursorBy(1, 0)
		case tcell.KeyUp:
			w.moveCursorBy(0, -1)
		case tcell.KeyDown:
			w.moveCursorBy(0, 1)
		default:
			switch w.mode {
			case EditorMode(Insert):
				w.OnInsertModeKeyPress(event)
			case EditorMode(Normal):
				w.OnNormalModeKeyPress(event)
			}
		}
	}
}

func (w *EditorWindow) Render() {
	w.SetContent(strings.Join(w.content, "\n"))
	w.BasicWindow.Render()

	line := w.content[w.cursorY]
	char := ' '

	if w.cursorX < len(line) {
		char = rune(line[w.cursorX])
	}

	w.DrawRuneAtBody(w.cursorX, w.cursorY, char, nil, cursorStyle)
}

func main() {
	fmt.Println("Hello, world!")
	Init()

	window := AddWindow(newEditorWindow(100, 100, 50, 50))

	SetActiveWindow(window)

	defer Screen().Fini()

	for {
		Screen().Clear()

		for _, w := range OpenWindows() {
			RenderWindow(w)
		}

		Screen().Show()

		ev := Screen().PollEvent()
		switch ev.(type) {
		default:
			ActiveWindow().OnKeyPress(&ev)
		}
	}
}
