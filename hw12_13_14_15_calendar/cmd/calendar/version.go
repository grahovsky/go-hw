package main

import (
	"encoding/json"
	"fmt"
	"os"
)

var (
	release   = "1.0.0"
	buildDate = "20230802"
	gitHash   = "3230cde47aaa391ba9e96d02cc09faba6741c4be"
)

func printVersion() {
	if err := json.NewEncoder(os.Stdout).Encode(struct {
		Release   string
		BuildDate string
		GitHash   string
	}{
		Release:   release,
		BuildDate: buildDate,
		GitHash:   gitHash,
	}); err != nil {
		fmt.Printf("error while decode version info: %v\n", err)
	}
}
