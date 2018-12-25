package nomad

import (
	"os"
)

type Action int8

const (
	TOGGLE_ACTIVE_PANE Action = iota
	GOTO_NEXT_LINE
	GOTO_PREV_LINE
)

func Update(state *State, action Action) {
	switch action {
	case TOGGLE_ACTIVE_PANE:
		ToggleActivePane(state)
		break
	case GOTO_NEXT_LINE:
		GoToNextLine(state)
		break
	case GOTO_PREV_LINE:
		GoToPrevLine(state)
		break
	}
}

func ToggleActivePane(state *State) {
	if state.ActivePane == LEFT_PANE {
		state.ActivePane = RIGHT_PANE
	} else {
		state.ActivePane = LEFT_PANE
	}
}

func GoToNextLine(state *State) {
	state.LeftList.Selected++
}

func GoToPrevLine(state *State) {
	state.LeftList.Selected--
}

func SetLeftPaneDirectory(state *State, path string, files []os.File) {
}
