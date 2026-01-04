package main

import (
	"fmt"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "Usage: schego <file.scm>")
		os.Exit(1)
	}

	filename := os.Args[1]
	fmt.Printf("Schego Scheme Interpreter\nFile: %s\n", filename)

	// TODO: Implement file reading and interpretation
	fmt.Println("Interpreter not yet implemented")
}
