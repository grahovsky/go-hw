package main

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
)

type Environment map[string]EnvValue

// EnvValue helps to distinguish between empty files and files with the first empty line.
type EnvValue struct {
	Value      string
	NeedRemove bool
}

const (
	wrongEnvNameChars = "="
	trimChars         = " \n\t"
)

// ReadDir reads a specified directory and returns map of env variables.
// Variables represented as files where filename is name of variable, file first line is a value.
func ReadDir(dir string) (Environment, error) {
	files, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}
	env := make(Environment)

	for _, file := range files {
		if !file.Type().IsRegular() || file.IsDir() {
			continue
		}
		if strings.ContainsAny(file.Name(), wrongEnvNameChars) {
			continue
		}
		filename := filepath.Join(dir, file.Name())
		data, err := os.ReadFile(filename)
		if err != nil {
			return nil, err
		}
		fi, err := os.Stat(filename)
		if err != nil {
			return nil, err
		}

		env[strings.TrimRight(file.Name(), "=")] = EnvValue{
			Value:      clearValue(data),
			NeedRemove: fi.Size() == 0,
		}
	}
	return env, nil
}

func clearValue(data []byte) string {
	value := strings.Split(string(data), "\n")[0]
	value = strings.TrimRight(value, trimChars)
	value = string(bytes.ReplaceAll([]byte(value), []byte("\x00"), []byte("\n")))
	return value
}
