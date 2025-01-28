package ui

import "github.com/gdamore/tcell"

var screen tcell.Screen = nil
var open_windows []InteractiveWindow = nil

func AddWindow[T InteractiveWindow](w T) T {
	w.SetScreen(screen)
	open_windows = append(open_windows, w)
	return w
}

func ActiveWindow() InteractiveWindow {
	return open_windows[len(open_windows)-1]
}

func SetActiveWindow(w InteractiveWindow) {
	found := false
	for i, window := range open_windows {
		if window == w {
			open_windows[i] = open_windows[len(open_windows)-1]
			found = true
		}
	}

	if !found {
		panic("window not found in current open windows")
	}

	open_windows[len(open_windows)-1] = w
}

func Init() {
	var err error = nil
	screen, err = tcell.NewScreen()
	if err != nil {
		panic(err)
	}

	err = screen.Init()
	if err != nil {
		panic(err)
	}

}

func Screen() tcell.Screen {
	return screen
}

func OpenWindows() []InteractiveWindow {
	return open_windows
}
