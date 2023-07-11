package main

import (
	"bytes"
	"io"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRunCmd(t *testing.T) {
	cmd := []string{"printenv", "HELLO", "WORLD"}
	env := Environment{
		"HELLO": EnvValue{
			Value:      "hello",
			NeedRemove: false,
		},
		"WORLD": EnvValue{
			Value:      "world",
			NeedRemove: false,
		},
	}

	actualOutput := captureOutput(t, func() {
		RunCmd(cmd, env)
	})

	require.Equal(t, "hello\nworld\n", actualOutput)
}

func TestEmpty(t *testing.T) {
	cmd := []string{"printenv", "$PATH"}
	env := Environment{
		"PATH": EnvValue{
			Value:      "",
			NeedRemove: false,
		},
	}

	actualOutput := captureOutput(t, func() {
		RunCmd(cmd, env)
	})

	require.Equal(t, "", actualOutput)
}

func captureOutput(t *testing.T, f func()) string {
	t.Helper()
	reader, writer, _ := os.Pipe()
	savedStdout := os.Stdout
	os.Stdout = writer
	out := make(chan string)
	go func() {
		var buf bytes.Buffer
		io.Copy(&buf, reader)
		out <- buf.String()
	}()

	f()
	writer.Close()
	os.Stdout = savedStdout

	return <-out
}
