package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
)

// RunCmd runs a command + arguments (cmd) with environment variables from env.
func RunCmd(cmd []string, env Environment) (returnCode int) {
	command := exec.Command(cmd[0], cmd[1:]...)
	command.Env = os.Environ()

	for name, envVar := range env {
		strEnvVar := fmt.Sprintf("%v=%v", name, envVar.Value)
		switch envVar.NeedRemove {
		case true:
			command.Env = removeEnv(command.Env, name)
		case false:
			command.Env = append(command.Env, strEnvVar)
		}
	}
	command.Stdout = os.Stdout
	command.Stderr = os.Stderr
	var exitCode int

	if err := command.Run(); err != nil {
		var valErrs *exec.ExitError
		if errors.As(err, &valErrs) {
			exitCode = valErrs.ExitCode()
		} else {
			log.Println(err)
			exitCode = 1
		}
	}

	return exitCode
}

func removeEnv(env []string, name string) []string {
	result := make([]string, 0, len(env))
	for _, item := range env {
		if !strings.HasPrefix(item, name+"=") {
			result = append(result, item)
		}
	}
	return result
}
