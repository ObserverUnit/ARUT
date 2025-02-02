package main

import (
	"iter"
	"maps"
	"os"
	"strings"

	"github.com/ObserverUnit/arut/src/ui"
	"github.com/gdamore/tcell"
)

var fileWindowCommandsTable = map[string]func(self *FileEditorWindow, args []string) string{
	"write": func(self *FileEditorWindow, args []string) string {
		return self.write()
	},
}

// an EditorWindow which is editing a file
type FileEditorWindow struct {
	EditorWindow
	file *os.File
}

func newFileEditorWindow(wm *ui.WindowManager, name string, width, height, x, y int) (*FileEditorWindow, error) {
	file, err := os.OpenFile(name, os.O_RDWR, 0)
	if err != nil {
		return nil, err
	}

	data, err := os.ReadFile(name)
	if err != nil {
		file.Close()
		return nil, err
	}

	inner := newEditorWindow(wm, width, height, x, y, string(data))
	return &FileEditorWindow{
		EditorWindow: *inner,
		file:         file,
	}, nil
}

func (w *FileEditorWindow) write() string {
	content := strings.Join(w.content, "\n")
	_, err := w.file.Seek(0, 0)
	if err != nil {
		return err.Error()
	}
	err = w.file.Truncate(0)
	if err != nil {
		return err.Error()
	}

	_, err = w.file.WriteString(content)

	if err != nil {
		return err.Error()
	}

	return ""
}

func (w *FileEditorWindow) Commands() iter.Seq[string] {
	return maps.Keys(fileWindowCommandsTable)
}

func (w *FileEditorWindow) ExecCommand(cmd string, args []string) string {
	return fileWindowCommandsTable[cmd](w, args)
}

func (w *FileEditorWindow) openCommandWindow() {
	w.wm.AddWindow(newCommandWindow(w, 50, 50))
}

func (w *FileEditorWindow) Close() {
	w.file.Close()
	w.wm.Close(w)
}
func (w *FileEditorWindow) OnEvent(event *tcell.Event) {
	// workaround for to make the command window work
	switch ev := (*event).(type) {
	case *tcell.EventKey:
		switch ev.Key() {
		case tcell.KeyRune:
			switch key := ev.Rune(); key {
			case ':':
				if w.mode == EditorMode(Normal) {
					w.openCommandWindow()
					return
				}
			}
		}
	}

	w.EditorWindow.OnEvent(event)
}
