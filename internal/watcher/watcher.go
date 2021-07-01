package watcher

import (
	"github.com/fsnotify/fsnotify"
	"github.com/go-component/go-watcher/internal/runner"
	"log"
	"strings"
)

type Watcher struct {
	runner *runner.Runner
}

func NewWatcher(runner *runner.Runner) *Watcher {
	return &Watcher{runner: runner}
}



func(w *Watcher) Start() error{

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return err
	}
	defer watcher.Close()

	done := make(chan interface{})
	go func() {


		defer func() {
			done <- recover()
		}()

		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}

				if !strings.Contains(event.Name, ".go"){
					continue
				}

				if event.Op&fsnotify.Write == fsnotify.Write {
					go func() {
						err = w.runner.Restart()
						if err != nil{
							log.Println(err)
						}
					}()
				}
			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				log.Println("watcher error:", err)
			}
		}
	}()
	err = watcher.Add(w.runner.WorkPath)
	if err != nil {
		log.Fatal(err)
	}
	state := <-done

	if err,ok := state.(error); ok{
		return err
	}

	return nil
}
