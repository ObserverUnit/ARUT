package ui

import (
	"errors"

	"github.com/gdamore/tcell"
)

type WindowManager struct {
	tcell.Screen
	windows []InteractiveWindow
	active  InteractiveWindow
}

func (wm *WindowManager) AddWindow(w InteractiveWindow) {
	wm.windows = append(wm.windows, w)
	wm.active = w
}

func (wm *WindowManager) RemoveWindow(w InteractiveWindow) error {
	for i, window := range wm.windows {
		if window == w {
			wm.windows = append(wm.windows[:i], wm.windows[i+1:]...)

			if w == wm.active {
				if i >= 1 && i-1 < len(wm.windows) {
					wm.active = wm.windows[i-1]
				} else {
					wm.active = nil
				}
			}
			return nil
		}
	}

	return errors.New("window not found")
}

func (wm *WindowManager) render() {
	for _, w := range wm.windows {
		RenderWindow(w)
	}
}

func (wm *WindowManager) onEvent(event *tcell.Event) {
	wm.active.OnEvent(event)
}

func (wm *WindowManager) Run() {
	defer wm.Fini()

	for {
		wm.Clear()
		wm.render()
		wm.Show()

		ev := wm.PollEvent()
		switch ev.(type) {
		default:
			wm.onEvent(&ev)
		}
	}
}

func NewWindowManager() (error, *WindowManager) {
	screen, err := tcell.NewScreen()
	if err != nil {
		return err, nil
	}

	err = screen.Init()
	if err != nil {
		return err, nil
	}

	return nil, &WindowManager{
		Screen:  screen,
		windows: []InteractiveWindow{},
	}
}

// Closes a given window and removes it from the window manager, panics if the window is not found
func (wm *WindowManager) Close(w InteractiveWindow) {
	if err := wm.RemoveWindow(w); err != nil {
		panic(err)
	}
}
