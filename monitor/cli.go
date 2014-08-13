package main

import (
	"time"

	"github.com/nsf/termbox-go"
)

var quit = false

func startCLI() {
	err := termbox.Init()
	if err != nil {
		panic(err)
	}
	termbox.SetInputMode(termbox.InputEsc)
	redraw_all()

	// capture and process events from the CLI
	eventChan := make(chan termbox.Event, 16)
	go handleEvents(eventChan)
	go func() {
		for {
			ev := termbox.PollEvent()
			eventChan <- ev
		}
	}()

	// start update (redraw) ticker
	timer := time.Tick(time.Millisecond * 100)
	for {
		select {
		case <-timer:
			redraw_all()
		}
	}
}

const edit_box_width = 30

func handleEvents(eventChan chan termbox.Event) {
	for {
		ev := <-eventChan
		switch ev.Type {
		case termbox.EventKey:
			switch ev.Key {
			case termbox.KeyArrowLeft, termbox.KeyCtrlB:
				// edit_box.MoveCursorOneRuneBackward()
			case termbox.KeyArrowRight, termbox.KeyCtrlF:
				// edit_box.MoveCursorOneRuneForward()
			case termbox.KeyBackspace, termbox.KeyBackspace2:
				// edit_box.DeleteRuneBackward()
			case termbox.KeyDelete, termbox.KeyCtrlD:
				// edit_box.DeleteRuneForward()
			case termbox.KeyTab:
				// edit_box.InsertRune('\t')
			case termbox.KeySpace:
				// edit_box.InsertRune(' ')
			case termbox.KeyCtrlK:
				// edit_box.DeleteTheRestOfTheLine()
			case termbox.KeyHome, termbox.KeyCtrlA:
				// edit_box.MoveCursorToBeginningOfTheLine()
			case termbox.KeyEnd, termbox.KeyCtrlE:
				// edit_box.MoveCursorToEndOfTheLine()
			case termbox.KeyCtrlQ, termbox.KeyEsc:
				goto endfunc
			default:
				if ev.Ch != 0 {
					// edit_box.InsertRune(ev.Ch)
				}
			}
		case termbox.EventError:
			panic(ev.Err)
		}
	}
endfunc:
	quit = true
}
