package main

import (
	"iter"
	"maps"

	. "github.com/ObserverUnit/arut/src/ui"
	. "github.com/ObserverUnit/arut/src/ui/windows"
	"github.com/gdamore/tcell"
)

var initWindowCommandTable = map[string]func(*InitWindow, []string) string{}

// A wrapper around BasicWindow that is basically a welcome screen
type InitWindow struct {
	BasicWindow
	wm *WindowManager
}

func newInitWindow(wm *WindowManager, width, height, x, y int) *InitWindow {
	inner := NewBasicWindow(wm, width, height, x, y, "Welcome to ARUT! Press enter to open the command window.")
	return &InitWindow{
		BasicWindow: *inner,
		wm:          wm,
	}
}

func (w *InitWindow) OnEvent(event *tcell.Event) {
	switch ev := (*event).(type) {
	case *tcell.EventKey:
		switch ev.Key() {
		case tcell.KeyEnter:
			w.wm.AddWindow(newCommandWindow(w, 50, 50))
		}
	}
}

func (w *InitWindow) Commands() iter.Seq[string] {
	return maps.Keys(initWindowCommandTable)
}

func (w *InitWindow) ExecCommand(cmd string, args []string) string {
	return initWindowCommandTable[cmd](w, args)
}

func (w *InitWindow) Close() {
	w.wm.Close(w)
}

func (w *InitWindow) WindowManager() *WindowManager {
	return w.wm
}
