package runner

import (
	"errors"
	"fmt"
	"log"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync/atomic"
	"syscall"
	"time"
)

var ErrorRepeat = errors.New("repeat exec")

type Runner struct {
	goExecFilePath string

	execFilePath string
	execArgs     []string

	sourceFilePath string
	WorkPath       string

	buildArgs []string

	running  int32
	procAttr *os.ProcAttr

	process  atomic.Value
}

func execFilePathFormat(sourcePath string) string {
	rand.Seed(time.Now().UTC().UnixNano())
	return fmt.Sprintf("%s_%d", strings.TrimRight(filepath.Base(sourcePath), ".go"), rand.Int())
}

func NewRunner(args []string) (*Runner, error) {

	if len(args) < 1 {
		return nil, errors.New("command args at least one")
	}

	sourcePath, err := filepath.Abs(args[0])

	if err != nil {
		return nil, err
	}
	goExecFilePath, err := exec.LookPath("go")
	if err != nil {
		return nil, err
	}
	workPath := filepath.Dir(sourcePath)

	execFilePath := execFilePathFormat(sourcePath)

	return &Runner{
		goExecFilePath: goExecFilePath,

		execFilePath: execFilePath,
		execArgs: append([]string{
			execFilePath,
		}, args[1:]...),

		sourceFilePath: sourcePath,
		WorkPath:       workPath,

		buildArgs: []string{
			"",
			"build",
			"-o",
			execFilePath,
			sourcePath,
		},

		procAttr: &os.ProcAttr{
			Files: []*os.File{os.Stdin, os.Stdout, os.Stderr},
		},

	}, nil

}

func (r *Runner) build() error {
	process, err := os.StartProcess(r.goExecFilePath, r.buildArgs, r.procAttr)

	if err != nil {
		return err
	}
	_, err = process.Wait()
	if err != nil {
		return err
	}

	return nil
}

func(r *Runner) start()(*os.Process, error){

	if !atomic.CompareAndSwapInt32(&r.running, 0, 1) {
		return nil, ErrorRepeat
	}
	defer atomic.CompareAndSwapInt32(&r.running, 1, 0)

	r.Shutdown()

	log.Println("starting server...")

	err := r.build()

	if err != nil {
		return nil,err
	}

	process, err := os.StartProcess(r.execFilePath, r.execArgs, r.procAttr)

	if err != nil {
		return nil,err
	}

	r.process.Store(process)

	return process,nil
}

func (r *Runner) Exec() error {

	process, err := r.start()

	if errors.Is(err, ErrorRepeat){
		return nil
	}

	if err != nil{
		return err
	}

	_, err = process.Wait()
	if err != nil {
		return err
	}

	return nil
}

func (r *Runner) Cleanup() {
	r.Shutdown()
	_ = os.Remove(r.execFilePath)
}

func (r *Runner) Shutdown() {

	process, ok := r.process.Load().(*os.Process)

	if ok{
		err := process.Signal(syscall.SIGTERM)
		if err == nil{
			log.Println("Closed server...")
		}
	}
}

func (r *Runner) Restart() error {
	return r.Exec()
}
