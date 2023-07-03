package main

import (
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

	if err := command.Run(); err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			exitCode := exitErr.ExitCode()
			return exitCode
		} else {
			log.Println(err)
			return 1
		}
	}

	return 0
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
