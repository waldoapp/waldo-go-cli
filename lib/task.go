package lib

import (
	"bytes"
	"os"
	"os/exec"
	"strings"
)

type Task struct {
	Name      string
	Args      []string
	Env       map[string]string
	Cwd       string
	IOStreams *IOStreams
}

//-----------------------------------------------------------------------------

func NewTask(name string, args ...string) *Task {
	return &Task{
		Name:      name,
		Args:      args,
		IOStreams: NewIOStreams(os.Stdin, os.Stdout, os.Stderr)}
}

//-----------------------------------------------------------------------------

// if ee, ok := err.(*exec.ExitError); ok {
// 	os.Exit(ee.ExitCode())
// }

func (t *Task) Execute() error {
	cmd := exec.Command(t.Name, t.Args...)

	cmd.Dir = t.Cwd
	cmd.Env = t.convertEnv()
	cmd.Stderr = t.IOStreams.errWriter
	cmd.Stdin = t.IOStreams.inReader
	cmd.Stdout = t.IOStreams.outWriter

	return cmd.Run()
}

func (t *Task) Run() (string, string, error) {
	var (
		stderrBuffer bytes.Buffer
		stdoutBuffer bytes.Buffer
	)

	cmd := exec.Command(t.Name, t.Args...)

	cmd.Dir = t.Cwd
	cmd.Env = t.convertEnv()
	cmd.Stderr = &stderrBuffer
	cmd.Stdin = t.IOStreams.inReader
	cmd.Stdout = &stdoutBuffer

	err := cmd.Run()

	stderr := strings.TrimRight(stderrBuffer.String(), "\n")
	stdout := strings.TrimRight(stdoutBuffer.String(), "\n")

	return stdout, stderr, err
}

//-----------------------------------------------------------------------------

func (t *Task) convertEnv() []string {
	var flatEnv []string

	for envvar, value := range t.Env {
		flatEnv = append(flatEnv, envvar+"="+value)
	}

	return flatEnv
}
