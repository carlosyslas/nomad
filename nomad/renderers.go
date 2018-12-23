package nomad

import (
	"github.com/gdamore/tcell"
)

func puts(screen tcell.Screen, x int, y int, text string, style tcell.Style) {
	for i, ch := range text {
		screen.SetContent(x+i, y, ch, []rune{}, style)
	}
}

type Renderer interface {
	Render(screen tcell.Screen, state State)
}

type Position struct {
	left int
	top  int
}

type Size struct {
	width  int
	height int
}

type Box struct {
	position Position
	size     Size
}

type PathRenderer struct{}

func (r PathRenderer) getPathStyle() tcell.Style {
	return tcell.StyleDefault.Foreground(tcell.ColorBlue).Background(tcell.ColorBlack)
}

func (r PathRenderer) Render(screen tcell.Screen, state State) {
	style := r.getPathStyle()

	puts(screen, 0, 0, state.Path, style)
}

type ListRenderer struct {
	position Position
	size     Size
}

func (r ListRenderer) Render(screen tcell.Screen, state State) {
	for i, item := range state.LeftList.Items {
		style := tcell.StyleDefault
		if i == state.LeftList.Selected {
			style = style.Background(tcell.ColorWhite).Foreground(tcell.ColorBlack)
		}
		puts(screen, r.position.left+1, r.position.top+i, item, style)
	}
}

type DirectoryListRenderer struct {
	pathRenderer PathRenderer
	listRenderer ListRenderer
}

type LeftPaneRenderer struct {
	directoryListRenderer DirectoryListRenderer
}

// This could be a container instead of app
type appRenderer struct {
	children []Renderer
}

func (r appRenderer) Render(screen tcell.Screen, state State) {
	for _, child := range r.children {
		child.Render(screen, state)
	}
}

func NewAppRenderer() Renderer {
	return appRenderer{
		children: []Renderer{
			PathRenderer{},
			ListRenderer{position: Position{0, 1}},
		},
	}
}
