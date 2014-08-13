package main

import (
	"fmt"
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

var alertDetailsView = false
var alertPage = ""

func handleEvents(eventChan chan termbox.Event) {
	for {
		ev := <-eventChan
		switch ev.Type {
		case termbox.EventKey:
			switch ev.Key {

			case termbox.KeyEnter:
				alertDetailsView = true
				alertPage = fmt.Sprintf("page%d", selectedRow+1)

			case termbox.KeyArrowDown:
				// if alertDetailsView {
				// 	maxSelectedRow = len(siteStats.AlertHist[alertPage]) - 1
				// 	if siteStats.OpenAlerts[alertPage] != nil {
				// 		maxSelectedRow++
				// 	}
				// } else {
				// 	maxSelectedRow = len(knownPages) - 1
				// }
				if !alertDetailsView && selectedRow < maxSelectedRow {
					selectedRow++
				}
			case termbox.KeyArrowUp:
				if !alertDetailsView && selectedRow > 0 {
					selectedRow--
				}
			case termbox.KeyHome:
				selectedRow = 0
			case termbox.KeyEnd:
				selectedRow = maxSelectedRow

			case termbox.KeyEsc:
				if alertDetailsView {
					maxSelectedRow = len(knownPages)
					alertDetailsView = false
				} else {
					goto endfunc
				}
			case termbox.KeyCtrlQ:
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
