package main

import (
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
)

const (
	BadArgumentErrorCode = 255
	CommandRunError      = 256
)

// RunCmd runs a command + arguments (cmd) with environment variables from env.
func RunCmd(cmd []string, env Environment, stdin io.Reader, stdout io.Writer, stderr io.Writer) (returnCode int) {
	if len(cmd) < 1 {
		return BadArgumentErrorCode
	}
	name := cmd[0]
	args := cmd[1:]

	command := exec.Command(name, args...)
	command.Stdin = stdin
	command.Stdout = stdout
	command.Stderr = stderr
	command.Env = prepareEnv(env)

	if err := command.Run(); err != nil {
		var exitErr *exec.ExitError
		if ok := errors.As(err, &exitErr); ok {
			return exitErr.ExitCode()
		}
		return CommandRunError
	}
	return 0
}

const envNameValueSeparator = "="

func prepareEnv(envModifier Environment) []string {
	originalEnv := os.Environ()
	resultEnvMap := make(Environment, len(originalEnv)) // make copy to avoid modifying function argument
	for name, value := range envModifier {
		resultEnvMap[name] = value
	}

	for _, originalEnvVar := range originalEnv {
		split := strings.SplitN(originalEnvVar, envNameValueSeparator, 2)
		name, value := split[0], split[1]
		_, ok := resultEnvMap[name]
		if !ok {
			resultEnvMap[name] = EnvValue{Value: value}
		}
	}

	resultEnv := make([]string, 0, len(resultEnvMap))
	for name, value := range resultEnvMap {
		if !value.NeedRemove {
			resultEnv = append(resultEnv, fmt.Sprintf("%s%s%s", name, envNameValueSeparator, value.Value))
		}
	}

	return resultEnv
}
