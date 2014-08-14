package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/go-fsnotify/fsnotify"
)

var watcher *fsnotify.Watcher

// setup and defferd closing in main()

func startWatcher(filename string, out chan *LineRaw) {
	// file stuff
	file, err := os.Open(filename)
	if err != nil {
		panic(err)
	}
	fi, err := os.Stat(filename)
	if err != nil {
		panic(err)
	}

	// heres where backfill stuff would happen
	last_sz := fi.Size()

	// go into watcher loops
	watcher, err = fsnotify.NewWatcher()
	if err != nil {
		panic(err)
	}
	defer watcher.Close()
	err = watcher.Add(filename)
	if err != nil {
		panic(err)
	}

	// main event loop
	for {
		select {
		case event := <-watcher.Events:
			if event.Op&fsnotify.Write == fsnotify.Write {
				// get updated size
				fi, err := os.Stat(filename)
				if err != nil {
					panic(err)
				}
				curr_sz := fi.Size()
				sz_chg := curr_sz - last_sz

				// make buf for reading
				buf := make([]byte, sz_chg)
				// read and check for errors
				n, err := file.ReadAt(buf, last_sz)
				if err != nil {
					fmt.Println("readat error:", err)
				}
				if n != int(sz_chg) {
					fmt.Println("n:", n, "!=", "sz_chg:", sz_chg)
				}

				// send data out
				out <- &LineRaw{
					line:    buf,
					logfile: filename,
				}

				// update counters
				last_sz = curr_sz
			}
		case err := <-watcher.Errors:
			fmt.Println("watch error:", err)
		}
		// check for ending
		if quit == true {
			return
		}
	}
}

func startWatcherList(listfile string, out chan *LineRaw) {
	listbytes, err := ioutil.ReadFile(listfile)
	if err != nil {
		panic(err)
	}
	lines := bytes.Fields(listbytes)
	for _, logfile := range lines {
		filename := string(logfile)
		fmt.Println("  watching: ", filename)
		addSiteStats(filename)
		go startWatcher(filename, out)
	}
}
