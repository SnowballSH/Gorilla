package main

import (
	"os"

	"./repl"
)

func main() {
	repl.Start(os.Stdin, os.Stdout)
}
