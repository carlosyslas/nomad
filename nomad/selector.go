package nomad

func SelectLeftPane(state State) LeftPane {
	return state.LeftPane
}

func SelectRightPane(state State) RightPane {
	return state.RightPane
}
