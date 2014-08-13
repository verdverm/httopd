package main

import (
	"fmt"
	"os"

	"github.com/go-fsnotify/fsnotify"
)

func startWatcher(filename string, out chan []byte) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		panic(err)
	}
	defer watcher.Close()

	file, err := os.Open(filename)
	if err != nil {
		panic(err)
	}

	fi, err := os.Stat(filename)
	if err != nil {
		panic(err)
	}

	last_sz := fi.Size()

	err = watcher.Add(filename)
	if err != nil {
		panic(err)
	}

	for {
		select {
		case event := <-watcher.Events:
			if event.Op&fsnotify.Write == fsnotify.Write {
				fi, err := os.Stat(filename)
				if err != nil {
					panic(err)
				}

				curr_sz := fi.Size()
				sz_chg := curr_sz - last_sz

				buf := make([]byte, sz_chg)
				n, err := file.ReadAt(buf, last_sz)
				if err != nil {
					fmt.Println("readat error:", err)
				}
				if n != int(sz_chg) {
					fmt.Println("n:", n, "!=", "sz_chg:", sz_chg)
				}
				out <- buf

				last_sz = curr_sz
			}
		case err := <-watcher.Errors:
			fmt.Println("watch error:", err)
		}
		if quit == true {
			return
		}
	}
}
