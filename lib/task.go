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
	Env       Environment
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

func (t *Task) Execute() error {
	cmd := exec.Command(t.Name, t.Args...)

	cmd.Dir = t.Cwd
	cmd.Env = t.Env.Flatten()
	cmd.Stderr = t.IOStreams.errWriter
	cmd.Stdin = t.IOStreams.inReader
	cmd.Stdout = t.IOStreams.outWriter

	return cmd.Run()
}

func (t *Task) Run() (string, string, error) {
	stdoutBytes, stderrBytes, err := t.RunRaw()

	stderr := strings.TrimRight(string(stderrBytes), "\n")
	stdout := strings.TrimRight(string(stdoutBytes), "\n")

	return stdout, stderr, err
}

func (t *Task) RunRaw() ([]byte, []byte, error) {
	var (
		stderrBuffer bytes.Buffer
		stdoutBuffer bytes.Buffer
	)

	cmd := exec.Command(t.Name, t.Args...)

	cmd.Dir = t.Cwd
	cmd.Env = t.Env.Flatten()
	cmd.Stderr = &stderrBuffer
	cmd.Stdin = t.IOStreams.inReader
	cmd.Stdout = &stdoutBuffer

	err := cmd.Run()

	return stdoutBuffer.Bytes(), stderrBuffer.Bytes(), err
}
