package main

import (
	"./nomad"
	"fmt"
	"github.com/gdamore/tcell"
	"io/ioutil"
	"os"
)

func render(screen tcell.Screen, state nomad.State) {
	// Display path
	renderer := nomad.NewAppRenderer()
	renderer.Render(screen, state)

	r := '#'
	if state.ActivePane == nomad.RIGHT_PANE {
		r = '$'
	}

	screen.SetCell(10, 10, tcell.StyleDefault.Foreground(tcell.ColorBlack).Background(tcell.ColorWhite), r)
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
	update := make(chan nomad.Action)
	path, _ := os.Getwd()
	files, _ := ioutil.ReadDir(path)
	var leftList nomad.ListState
	for _, f := range files {
		leftList.Items = append(leftList.Items, f.Name())
	}
	state := nomad.State{
		Path:     path,
		LeftList: leftList,
	}

	go func() {
		for {
			event := screen.PollEvent()

			switch event := event.(type) {
			case *tcell.EventKey:
				{
					switch event.Key() {
					case tcell.KeyEscape:
						{
							close(quit)
							return
						}
					case tcell.KeyTab:
						{
							update <- nomad.TOGGLE_ACTIVE_PANE
						}
					case tcell.KeyRune:
						{
							switch event.Rune() {
							case 'n':
								update <- nomad.GOTO_NEXT_LINE
							case 'p':
								update <- nomad.GOTO_PREV_LINE
							}
						}
					}

				}
			case *tcell.EventResize:
				{
					screen.Sync()
				}
			}
		}
	}()

	for {
		render(screen, state)
		select {
		case <-quit:
			{
				return
				break
			}
		case action := <-update:
			{
				nomad.Update(&state, action)
				render(screen, state)
			}
		}

		screen.SetCell(10, 10, tcell.StyleDefault.Foreground(tcell.ColorBlack).Background(tcell.ColorWhite), '#')
		screen.Show()
	}

	screen.Fini()
}
