package watcher

import (
	"github.com/fsnotify/fsnotify"
	"github.com/go-component/go-watcher/internal/runner"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

type Watcher struct {
	runner *runner.Runner
}

func NewWatcher(runner *runner.Runner) *Watcher {

	return &Watcher{runner: runner}
}

func registerWatchPath(dirPath string, callback func(path string) error) (err error) {

	dir, err := ioutil.ReadDir(dirPath)
	if err != nil {
		return err
	}

	PathSep := string(os.PathSeparator)

	for _, fi := range dir {

		path := dirPath + PathSep + fi.Name()

		if fi.IsDir() && !strings.HasPrefix(fi.Name(), ".") {
			err = callback(path)
			if err != nil {
				return err
			}

			err = registerWatchPath(path, callback)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (w *Watcher) Start() error {

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return err
	}

	defer func(watcher *fsnotify.Watcher) {
		err := watcher.Close()
		if err != nil {
			log.Println(err.Error())
		}
	}(watcher)

	done := make(chan interface{})

	go func() {

		defer func() {
			done <- recover()
		}()

		for {

			select {

			case event, ok := <-watcher.Events:

				if !ok || !strings.Contains(event.Name, ".go") {
					continue
				}

				if event.Op&fsnotify.Write == fsnotify.Write {
					go func() {
						err = w.runner.Restart()
						if err != nil {
							log.Println(err)
						}
					}()
				}

			case err, ok := <-watcher.Errors:

				if !ok {
					return
				}
				panic(err)
			}

		}
	}()

	err = watcher.Add(w.runner.WatchPath)
	if err != nil {
		return err
	}

	err = registerWatchPath(w.runner.WatchPath, func(path string) error {
		err = watcher.Add(path)
		if err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		return err
	}

	state := <-done

	if err, ok := state.(error); ok {
		return err
	}

	return nil
}
