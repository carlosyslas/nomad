package nomad

type PaneID int8

const (
	LEFT_PANE PaneID = iota
	RIGHT_PANE
)

type State struct {
	Path       string
	LeftList   ListState
	ActivePane PaneID
	LeftPane   LeftPane
	RightPane  RightPane
}

type LeftPane struct {
	list DirectoryList
}

type RightPane struct {
	list DirectoryList
}

type ListState struct {
	Items    []string
	Selected int
}

type DirectoryList struct {
	path string
	list ListState
}
