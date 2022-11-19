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
	resultEnv := make([]string, 0, len(originalEnv))
	skip := make(map[string]struct{}, len(envModifier))
	for _, originalEnvVar := range originalEnv {
		name := strings.SplitN(originalEnvVar, envNameValueSeparator, 2)[0]
		if modifier, ok := envModifier[name]; ok {
			if modifier.NeedRemove {
				continue
			}
			resultEnv = append(resultEnv, fmt.Sprintf("%s%s%s", name, envNameValueSeparator, modifier.Value))
			skip[name] = struct{}{}
		} else {
			resultEnv = append(resultEnv, originalEnvVar)
		}
	}
	for name, value := range envModifier {
		if _, ok := skip[name]; ok {
			continue
		}
		if value.NeedRemove {
			continue
		}
		resultEnv = append(resultEnv, fmt.Sprintf("%s%s%s", name, envNameValueSeparator, value.Value))
	}

	return resultEnv
}
