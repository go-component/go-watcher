package main

import (
	"flag"
	"github.com/go-component/go-watcher/internal/runner"
	"github.com/go-component/go-watcher/internal/watcher"
	"log"
	"os"
	"os/signal"
	"syscall"
)

var watchPath string

func init(){
	flag.StringVar(&watchPath,"w", "", "path of watch, default exec path")
}

func main() {
	flag.Parse()

	r, err := runner.NewRunner(flag.Args(), watchPath)
	if err != nil {
		log.Fatal(err)
	}
	defer r.Cleanup()

	w := watcher.NewWatcher(r)

	go func() {
		err := r.Exec()
		if err != nil {
			r.Cleanup()
			log.Fatal(err)
		}
	}()

	go func() {
		err := w.Start()
		if err != nil {
			r.Cleanup()
			log.Fatal(err)
		}
	}()

	sigChan := make(chan os.Signal, 1)

	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM, syscall.SIGHUP)

	<-sigChan
}
