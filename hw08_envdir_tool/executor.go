package main

import (
	"errors"
	"io"
	"os/exec"
)

// RunCmd runs a command + arguments (cmd) with environment variables from env.
func RunCmd(cmd []string, env Environment, stdout, stderr io.Writer, stdin io.Reader) (returnCode int) {
	c := exec.Command(cmd[0], cmd[1:]...) //nolint:gosec
	c.Stdout = stdout
	c.Stderr = stderr
	c.Stdin = stdin

	c.Env = prepareEnvironment(env)
	err := c.Run()

	var exitError *exec.ExitError
	if errors.As(err, &exitError) {
		return exitError.ExitCode()
	}
	if err != nil {
		return -1
	}

	return 0
}
