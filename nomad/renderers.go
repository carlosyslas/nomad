package nomad

import (
	"github.com/gdamore/tcell"
)

func puts(screen tcell.Screen, x int, y int, text string, style tcell.Style) {
	for i, ch := range text {
		screen.SetContent(x+i, y, ch, []rune{}, style)
	}
}

type Container interface {
	GetParent() Container
	GetSize() Size
	GetPosition() Position
	GetChildren() []Renderer
}

type Renderer interface {
	Render(screen tcell.Screen, state State)
}

type Position struct {
	Left int
	Top  int
}

type Size struct {
	Width  int
	Height int
}

type Box struct {
	Position Position
	Size     Size
}

type PathRenderer struct {
	Parent Container
}

func (r PathRenderer) getPathStyle() tcell.Style {
	return tcell.StyleDefault.Foreground(tcell.ColorBlue).Background(tcell.ColorBlack)
}

func (r PathRenderer) Render(screen tcell.Screen, state State) {
	parentPosition := r.Parent.GetPosition()
	style := r.getPathStyle()

	puts(screen, parentPosition.Left, parentPosition.Top, state.Path, style)
}

type ListRenderer struct {
	Parent Container
}

func (r ListRenderer) Render(screen tcell.Screen, state State) {
	position := r.Parent.GetPosition()
	// TODO: Use parent selector
	for i, item := range state.LeftList.Items {
		style := tcell.StyleDefault
		if i == state.LeftList.Selected {
			style = style.Background(tcell.ColorWhite).Foreground(tcell.ColorBlack)
		}
		puts(screen, position.Left+1, position.Top+i+1, item, style)
	}
}

type DirectoryListRenderer struct {
	Parent       Container
	Selector     DirectoryListSelector
	pathRenderer PathRenderer
	ListRenderer ListRenderer
}

func (r DirectoryListRenderer) Render(screen tcell.Screen, state State) {
	for _, child := range r.GetChildren() {
		child.Render(screen, state)
	}
}

func (r DirectoryListRenderer) GetParent() Container {
	return r.Parent
}

func (r DirectoryListRenderer) GetSize() Size {
	return r.GetParent().GetSize()
}

func (r DirectoryListRenderer) GetPosition() Position {
	return r.GetParent().GetPosition()
}

func (r DirectoryListRenderer) GetChildren() []Renderer {
	return []Renderer{
		ListRenderer{
			Parent: r,
		},
		PathRenderer{
			Parent: r,
		},
	}
}

func NewDirectoryList(container Container, selector DirectoryListSelector) Renderer {
	return DirectoryListRenderer{
		Parent:   container,
		Selector: selector,
	}
}

type LeftPaneRenderer struct {
	Parent        Container
	DirectoryList Renderer
}

func (r LeftPaneRenderer) Render(screen tcell.Screen, state State) {
	for _, child := range r.GetChildren() {
		child.Render(screen, state)
	}
}

func (r LeftPaneRenderer) GetParent() Container {
	return r.Parent
}

func (r LeftPaneRenderer) GetSize() Size {
	parentSize := r.GetParent().GetSize()

	return Size{
		Width:  parentSize.Width / 2,
		Height: parentSize.Height,
	}
}

func (r LeftPaneRenderer) GetPosition() Position {
	return r.GetParent().GetPosition()
}

func (r LeftPaneRenderer) GetChildren() []Renderer {
	return []Renderer{
		NewDirectoryList(r, SelectLeftDirectoryList),
	}
}

type RightPaneRenderer struct {
	Parent        Container
	DirectoryList Renderer
}

func (r RightPaneRenderer) Render(screen tcell.Screen, state State) {
	for _, child := range r.GetChildren() {
		child.Render(screen, state)
	}
}

func (r RightPaneRenderer) GetParent() Container {
	return r.Parent
}

func (r RightPaneRenderer) GetSize() Size {
	parentSize := r.GetParent().GetSize()

	return Size{
		Width:  parentSize.Width / 2,
		Height: parentSize.Height,
	}
}

func (r RightPaneRenderer) GetChildren() []Renderer {
	return []Renderer{
		NewDirectoryList(r, SelectRightDirectoryList),
	}
}

func (r RightPaneRenderer) GetPosition() Position {
	parentPosition := r.GetParent().GetPosition()
	parentSize := r.GetParent().GetSize()

	return Position{
		Left: parentPosition.Left + parentSize.Width/2,
		Top:  parentPosition.Top,
	}
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

func (r appRenderer) GetParent() Container {
	return nil
}

func (r appRenderer) GetSize() Size {
	return Size{
		Width: 80,
	}
}
func (r appRenderer) GetPosition() Position {
	return Position{}
}

func (r appRenderer) GetChildren() []Renderer {
	return r.children
}

func NewAppRenderer() Renderer {
	renderer := appRenderer{}

	renderer.children = []Renderer{
		LeftPaneRenderer{
			Parent: renderer,
		},
		RightPaneRenderer{
			Parent: renderer,
		},
	}

	return renderer
}
