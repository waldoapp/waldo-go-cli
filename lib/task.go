package lib

import (
	"bytes"
	"os"
	"os/exec"
	"strings"
)

type Task struct {
	Name string
	Args []string
	Env  map[string]string
	Cwd  string
}

//-----------------------------------------------------------------------------

func NewTask(name string, args []string) *Task {
	return &Task{
		Name: name,
		Args: args}
}

func NewTaskCwd(name string, args []string, cwd string) *Task {
	return &Task{
		Name: name,
		Args: args,
		Cwd:  cwd}
}

func NewTaskEnv(name string, args []string, env map[string]string) *Task {
	return &Task{
		Name: name,
		Args: args,
		Env:  env}
}

func NewTaskEnvCwd(name string, args []string, env map[string]string, cwd string) *Task {
	return &Task{
		Name: name,
		Args: args,
		Env:  env,
		Cwd:  cwd}
}

//-----------------------------------------------------------------------------

func (t *Task) Execute() {
	cmd := exec.Command(t.Name, t.Args...)

	cmd.Dir = t.Cwd
	cmd.Env = t.convertEnv()
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout

	err := cmd.Run()

	if ee, ok := err.(*exec.ExitError); ok {
		os.Exit(ee.ExitCode())
	}
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
