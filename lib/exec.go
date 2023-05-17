package lib

import (
	"bytes"
	"os"
	"os/exec"
	"strings"
)

func Exec(name string, args ...string) error {
	cmd := exec.Command(name, args...)

	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout

	exitOnError(cmd.Run())

	return nil
}

func ExecDir(dir, name string, args ...string) error {
	cmd := exec.Command(name, args...)

	cmd.Dir = dir
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout

	exitOnError(cmd.Run())

	return nil
}

func ExecEnv(env []string, name string, args ...string) error {
	cmd := exec.Command(name, args...)

	cmd.Env = env
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout

	exitOnError(cmd.Run())

	return nil
}

func Run(name string, args ...string) (string, string, error) {
	var (
		stderrBuffer bytes.Buffer
		stdoutBuffer bytes.Buffer
	)

	cmd := exec.Command(name, args...)

	cmd.Stderr = &stderrBuffer
	cmd.Stdout = &stdoutBuffer

	err := cmd.Run()

	stderr := strings.TrimRight(stderrBuffer.String(), "\n")
	stdout := strings.TrimRight(stdoutBuffer.String(), "\n")

	return stdout, stderr, err
}

func RunDir(dir, name string, args ...string) (string, string, error) {
	var (
		stderrBuffer bytes.Buffer
		stdoutBuffer bytes.Buffer
	)

	cmd := exec.Command(name, args...)

	cmd.Dir = dir
	cmd.Stderr = &stderrBuffer
	cmd.Stdout = &stdoutBuffer

	err := cmd.Run()

	stderr := strings.TrimRight(stderrBuffer.String(), "\n")
	stdout := strings.TrimRight(stdoutBuffer.String(), "\n")

	return stdout, stderr, err
}

func RunEnv(env []string, name string, args ...string) (string, string, error) {
	var (
		stderrBuffer bytes.Buffer
		stdoutBuffer bytes.Buffer
	)

	cmd := exec.Command(name, args...)

	cmd.Env = env
	cmd.Stderr = &stderrBuffer
	cmd.Stdout = &stdoutBuffer

	err := cmd.Run()

	stderr := strings.TrimRight(stderrBuffer.String(), "\n")
	stdout := strings.TrimRight(stdoutBuffer.String(), "\n")

	return stdout, stderr, err
}

//-----------------------------------------------------------------------------

func exitOnError(err error) {
	if ee, ok := err.(*exec.ExitError); ok {
		os.Exit(ee.ExitCode())
	}
}
