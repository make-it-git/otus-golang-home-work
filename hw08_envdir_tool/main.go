package main

import (
	"fmt"
	"os"
)

func main() {
	if len(os.Args) < 3 {
		fmt.Println("Directory required")
		os.Exit(1)
	}

	environment, err := ReadDir(os.Args[1])
	if err != nil {
		os.Exit(1)
	}

	cmd := os.Args[2:]
	RunCmd(cmd, environment, os.Stdout, os.Stderr, os.Stdin)
}
