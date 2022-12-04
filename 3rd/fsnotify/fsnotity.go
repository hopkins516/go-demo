package main

import (
	"log"

	"github.com/fsnotify/fsnotify"
)

func FSNotifyDemo() {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatalln("fsnotify new watcher failed, ", err)
	}

	defer watcher.Close()

	done := make(chan bool)
	go func() {
		defer close(done)

		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}
				log.Printf("%s %s\n", event.Name, event.Op)

			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				log.Println("error: ", err)
			}
		}
	}()

	err = watcher.Add("./")
	if err != nil {
		log.Fatal("Add current dir failed: ", err)
	}

	<-done
}

func main() {
	FSNotifyDemo()
}
