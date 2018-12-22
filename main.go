package main

import (
	"fmt"
	"os"
	"github.com/gdamore/tcell"
)

func puts(screen tcell.Screen, x int, y int, text string, style tcell.Style) {
	for i, ch := range text {
		screen.SetContent(x + i, y, ch, []rune{}, style)
	}
}

type State struct {
	path string
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
	state := State {
		path: path,
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
