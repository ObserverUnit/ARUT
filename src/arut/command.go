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

var commands = map[string]func(wm *ui.WindowManager, args []string){
	"quit": func(wm *ui.WindowManager, args []string) {
		// TODO: rethink this design choice
		// if for example we have an unsaved buffer in another window...
		wm.Fini()
		os.Exit(0)
	},
}

func newCommandWindow(wm *ui.WindowManager, x, y int) *CommandWindow {
	inner := newEditorWindow(wm, 50, 15, x, y)
	return &CommandWindow{
		EditorWindow: *inner,
	}
}

func (w *CommandWindow) getCommand() []string {
	command := w.content[0]
	command = strings.TrimSpace(command)

	return strings.Split(w.content[0], " ")
}

func (w *CommandWindow) reset(message string) {
	w.content = []string{"", message}
	w.cursorX = 0
	w.cursorY = 0
}
func (w *CommandWindow) executeCommand() {
	command := w.getCommand()
	if len(command) == 0 || command[0] == "" {
		w.close()
		return
	}

	cmd := command[0]
	args := command[1:]

	for c, f := range commands {
		if strings.Contains(strings.ToLower(c), strings.ToLower(cmd)) {
			f(w.wm, args)
			return
		}
	}

	w.reset("Command not found")
}

func (w *CommandWindow) close() {
	// TODO: rethink this design choice
	// close has to be overrided or it panics
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
