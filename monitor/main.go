package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"runtime/debug"
	"time"

	"github.com/nsf/termbox-go"
)

var (
	fn = flag.String("fn", "access.log", "log file to monitor")
)

func main() {
	flag.Parse()

	defer func() {
		if e := recover(); e != nil {
			termbox.Close()
			trace := fmt.Sprintf("%s: %s", e, debug.Stack()) // line 20
			ioutil.WriteFile("trace.txt", []byte(trace), 0644)
		}
	}()

	line_chan := make(chan []byte, 16)
	data_chan := make(chan *LineData, 16)

	go startWatcher(*fn, line_chan)
	go startParser(line_chan, data_chan)
	go startStats(data_chan)
	go startCLI() // streams CLI commands to the main loop

	for {
		// select {}
		time.Sleep(time.Millisecond * 100)

		if quit == true {
			break
		}
	}

	termbox.Close()
}
