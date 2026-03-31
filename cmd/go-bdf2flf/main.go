package main

import (
	"fmt"
	"os"

	"github.com/mattn/go-figlet"
)

func main() {
	if len(os.Args) > 1 {
		fmt.Fprintln(os.Stderr, "Usage: go-bdf2flf < input.bdf > output.flf")
		fmt.Fprintln(os.Stderr, "BDF data is read from stdin and FIGlet font data is output to stdout.")
		os.Exit(0)
	}

	if err := figlet.BDF2FLF(os.Stdin, os.Stdout); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}
