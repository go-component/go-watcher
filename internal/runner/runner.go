package runner

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"syscall"
	"time"
)

type Runner struct {
	goExecFilePath string

	execFilePath string
	execArgs     []string

	sourceFilePath string
	WorkPath       string

	buildArgs []string
	process   *os.Process
}

func NewRunner(args []string) (*Runner, error) {

	sourcePath, err := filepath.Abs(args[0])

	if err != nil {
		return nil, err
	}
	goExecFilePath, err := exec.LookPath("go")
	if err != nil {
		return nil, err
	}
	workPath := filepath.Dir(sourcePath)

	rand.Seed(time.Now().UTC().UnixNano())
	execFilePath := fmt.Sprintf("%s_%d", strings.TrimRight(filepath.Base(sourcePath), ".go"), rand.Int())

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
	}, nil

}

func (r *Runner) build() error {
	procAttr := &os.ProcAttr{
		Files: []*os.File{os.Stdin, os.Stdout, os.Stderr},
	}

	process, err := os.StartProcess(r.goExecFilePath, r.buildArgs, procAttr)

	if err != nil {
		return err
	}
	_, err = process.Wait()
	if err != nil {
		return err
	}

	return nil
}

func (r *Runner) Exec() error {

	err := r.build()

	if err != nil {
		return err
	}

	procAttr := &os.ProcAttr{
		Files: []*os.File{os.Stdin, os.Stdout, os.Stderr},
	}

	process, err := os.StartProcess(r.execFilePath, r.execArgs, procAttr)

	if err != nil {
		return err
	}
	r.process = process

	_, err = r.process.Wait()
	if err != nil {
		return err
	}

	r.process = nil

	return nil
}

func (r *Runner) Cleanup() {
	os.Remove(r.execFilePath)
}

func (r *Runner) Restart() error {
	if r.process != nil {
		err := r.process.Signal(syscall.SIGTERM)
		if err != nil {
			panic(err)
		}
		log.Println("Closed server success...")
	}
	log.Println("starting server...")
	return r.Exec()
}
