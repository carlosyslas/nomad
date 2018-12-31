package nomad

type DirectoryListSelector func(state State) DirectoryList

func SelectLeftPane(state State) LeftPane {
	return state.LeftPane
}

func SelectRightPane(state State) RightPane {
	return state.RightPane
}

func SelectRightDirectoryList(state State) DirectoryList {
	return state.RightPane.List
}
