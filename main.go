package main

import (
	"nexus/repl"
	"os"
)

func main() {
	repl.Start(os.Stdin, os.Stdout)
}
