package main

import (
	"os"
	"strings"

	"github.com/ObserverUnit/arut/src/ui"
	"github.com/gdamore/tcell"
)

type CommandWindow struct {
	EditorWindow
}

func newCommandWindow(wm *ui.WindowManager, x, y int) *CommandWindow {
	inner := newEditorWindow(wm, 50, 15, x, y)
	return &CommandWindow{
		EditorWindow: *inner,
	}
}

func (w *CommandWindow) getCommand() string {
	return w.content[0]
}

func (w *CommandWindow) executeCommand() {
	command := w.getCommand()
	if strings.TrimSpace(command) == "exit" {
		w.Screen().Fini()
		os.Exit(0)
	}
}

func (w *CommandWindow) close() {
	w.wm.Close(w)
}

func (w *CommandWindow) OnEvent(event *tcell.Event) {
	switch ev := (*event).(type) {
	case *tcell.EventKey:
		switch ev.Key() {
		case tcell.KeyCtrlC:
			w.close()
		case tcell.KeyLeft:
			w.moveCursorBy(-1, 0)
		case tcell.KeyRight:
			w.moveCursorBy(1, 0)
		case tcell.KeyEnter:
			w.executeCommand()
		case tcell.KeyBackspace | tcell.KeyBackspace2:
			w.removeRune()
		case tcell.KeyRune:
			key := ev.Rune()
			w.addRune(key)
		}
	}
}
