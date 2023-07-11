package main

import (
	"log"
	"os"
)

func main() {
	dirName := os.Args[1]

	env, err := ReadDir(dirName)
	if err != nil {
		log.Fatal(err)
	}
	cmd := os.Args[2:]

	returnCode := RunCmd(cmd, env)
	// fmt.Println("return code:", returnCode)

	os.Exit(returnCode)
}
