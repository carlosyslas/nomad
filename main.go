package main

import (
	"fmt"
	"os"
	"io/ioutil"
	"github.com/gdamore/tcell"
)

func puts(screen tcell.Screen, x int, y int, text string, style tcell.Style) {
	for i, ch := range text {
		screen.SetContent(x + i, y, ch, []rune{}, style)
	}
}

type State struct {
	path string
	leftList ListState
}

type ListState struct {
	items []string
	selected int
}

type Renderer interface {
	Render(screen tcell.Screen, state State)
}

type pathRenderer struct {}

func (r pathRenderer) getPathStyle() tcell.Style {
	return tcell.StyleDefault.Foreground(tcell.ColorBlue).Background(tcell.ColorBlack)
}

func (r pathRenderer) Render(screen tcell.Screen, state State) {
	style := r.getPathStyle()

	puts(screen, 0, 0, state.path, style)
}

type Position struct {
	left int
	top int
}

type Size struct {
	width int
	height int
}

type listRenderer struct {
	position Position
	size Size
}

func (r listRenderer) Render(screen tcell.Screen, state State) {
	for i, item := range state.leftList.items {
		style := tcell.StyleDefault
		if i == state.leftList.selected {
			style = style.Background(tcell.ColorWhite).Foreground(tcell.ColorBlack)
		}
		puts(screen, r.position.left + 1, r.position.top + i, item, style)
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

func NewAppRenderer() Renderer {
	return appRenderer{
		children: []Renderer{
			pathRenderer{},
			listRenderer{ position: Position{0, 1} },
		},
	}
}

func render(screen tcell.Screen, state State) {
	// Display path
	renderer := NewAppRenderer()
	renderer.Render(screen, state)

	screen.SetCell(10, 10, tcell.StyleDefault.Foreground(tcell.ColorBlack).Background(tcell.ColorWhite), '#')
	screen.Show()
}

func main() {
	tcell.SetEncodingFallback(tcell.EncodingFallbackASCII)
	screen, err := tcell.NewScreen()
	if err != nil {
		fmt.Printf("Error creating screen: %v\n", err)
	}
	if err = screen.Init(); err != nil {
		fmt.Printf("Error initializing screen: %v\n", err)
	}

	screen.SetStyle(tcell.StyleDefault.Foreground(tcell.ColorWhite).Background(tcell.ColorBlack))
	screen.Clear()

	//run := true
	quit := make(chan int)
	path, _ := os.Getwd()
	files, _ := ioutil.ReadDir(path)
	var leftList ListState
	for _, f := range files {
		leftList.items = append(leftList.items, f.Name())
	}
	state := State {
		path: path,
		leftList: leftList,
	}

	go func() {
		for {
			event := screen.PollEvent()

			switch event := event.(type) {
			case *tcell.EventKey: {
				switch event.Key() {
				case tcell.KeyEscape:{
					close(quit)
					return
				}
				}
			}
			case *tcell.EventResize: {
				screen.Sync()
			}
			}
		}
	}()

	for {
		render(screen, state)
		select {
		case <-quit:
			return
			break
		}
		screen.SetCell(10, 10, tcell.StyleDefault.Foreground(tcell.ColorBlack).Background(tcell.ColorWhite), '#')
		screen.Show()
	}

	screen.Fini()
}
