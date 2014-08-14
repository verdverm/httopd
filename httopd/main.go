package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"runtime/debug"
	"time"

	// "github.com/go-fsnotify/fsnotify"
	"github.com/nsf/termbox-go"
)

var (
	fn     = flag.String("fn", "", "log file to monitor")
	fnList = flag.String("fnList", "", "text file with list of log files to monitor")
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

	// var err error
	// watcher, err = fsnotify.NewWatcher()
	// if err != nil {
	// 	panic(err)
	// }
	// defer watcher.Close()
	line_chan := make(chan *LineRaw, 64)
	data_chan := make(chan *LineData, 64)

	if *fnList != "" {
		fmt.Println("Starting watchers")
		go startWatcherList(*fnList, line_chan)
	} else if *fn != "" {
		fmt.Println("Starting watcher")
		go startWatcher(*fn, line_chan)
	} else {
		fmt.Println("must specify log file(s) to watch")
		return
	}

	numParsers := 1
	for i := 0; i < numParsers; i++ {
		go startParser(line_chan, data_chan)
	}

	// can only have one of these right now
	go startStats(data_chan)

	// View & Cmd loops
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
