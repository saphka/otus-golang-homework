package main

import (
	"log"
	"os"
)

const (
	dirNamePos = 1
	cmdPos     = 2
)

func main() {
	args := os.Args
	if len(args) < cmdPos {
		log.Fatalln("Npt enough arguments to run. Requires dir and cmd")
	}

	env, err := ReadDir(args[dirNamePos])
	if err != nil {
		log.Fatalf("Error loading env: %s\n", err)
	}

	code := RunCmd(args[cmdPos:], env, os.Stdin, os.Stdout, os.Stderr)
	os.Exit(code)
}
