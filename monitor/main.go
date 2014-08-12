package main

import (
	"flag"
	"fmt"
	"log"
	"time"

	"github.com/go-fsnotify/fsnotify"
)

var (
	fn = flag.String("fn", "access.log", "log file to monitor")
)

func main() {
	flag.Parse()
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	defer watcher.Close()

	go func() {
		for {
			select {
			case event := <-watcher.Events:
				log.Println("event:", event)
				if event.Op&fsnotify.Write == fsnotify.Write {
					log.Println("modified file:", event.Name)
				}
			case err := <-watcher.Errors:
				log.Println("error:", err)
			}
		}
	}()

	err = watcher.Add(*fn)
	if err != nil {
		log.Fatal(err)
	}

	for {
		fmt.Println("---")
		time.Sleep(time.Second)
	}
}
