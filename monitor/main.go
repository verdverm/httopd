package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/go-fsnotify/fsnotify"
)

var (
	fn = flag.String("fn", "access.log", "log file to monitor")
)

func main() {
	flag.Parse()

	startWatcher(*fn)

}

func startWatcher(filename string) {
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

	go func() {
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
					fmt.Println("file change", sz_chg, curr_sz)

					buf := make([]byte, sz_chg)
					n, err := file.ReadAt(buf, last_sz)
					if err != nil {
						fmt.Println("readat error:", err)
					}
					fmt.Println(" ", n, "  ", string(buf))

					last_sz = curr_sz
				}
			case err := <-watcher.Errors:
				fmt.Println("watch error:", err)
			}
		}
	}()

	err = watcher.Add(filename)
	if err != nil {
		panic(err)
	}

	for {
		time.Sleep(time.Second)
	}
}
