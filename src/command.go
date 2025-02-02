package main

import (
	"errors"
	"strings"

	"github.com/ObserverUnit/arut/ui"
	"github.com/gdamore/tcell"
)

// an EditorWindow that can execute commands on a parent window
type CommandWindow struct {
	EditorWindow
	parent ui.InteractiveWindow
}

var commands = map[string]func(self *CommandWindow, args []string) error{
	"quit": func(self *CommandWindow, args []string) error {
		// TODO: rethink this design choice
		// if for example we have an unsaved buffer in another window...
		self.wm.Quit(0)
		return nil
	},
	"open": func(self *CommandWindow, args []string) error {
		if len(args) < 1 {
			return errors.New("Not enough arguments")
		}
		name := args[0]

		window, err := newFileEditorWindow(self.wm, name, 100, 100, 50, 50)
		if err != nil {
			return err
		}

		self.wm.AddWindow(window)
		self.Close()
		self.parent.Close()
		return nil
	},
}

func newCommandWindow(parent ui.InteractiveWindow, x, y int) *CommandWindow {
	wm := parent.WindowManager()
	inner := newEditorWindow(wm, 50, 15, x, y, "")
	return &CommandWindow{
		EditorWindow: *inner,
		parent:       parent,
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
		w.Close()
		return
	}

	cmd := command[0]
	args := command[1:]

	for c, f := range commands {
		if strings.Contains(strings.ToLower(c), strings.ToLower(cmd)) {
			err := f(w, args)
			if err != nil {
				w.reset(err.Error())
			}
			return
		}
	}

	for c := range w.parent.Commands() {
		if strings.Contains(strings.ToLower(c), strings.ToLower(cmd)) {
			response := w.parent.ExecCommand(c, args)
			if response != "" {
				w.reset(response)
			} else {
				w.Close()
			}
			return
		}
	}

	w.reset("Command not found")
}

func (w *CommandWindow) Close() {
	// TODO: rethink this design choice
	// close has to be overrided or it panics
	w.wm.Close(w)
}

func (w *CommandWindow) OnEvent(event *tcell.Event) {
	switch ev := (*event).(type) {
	case *tcell.EventKey:
		switch ev.Key() {
		case tcell.KeyCtrlC:
			w.Close()
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
