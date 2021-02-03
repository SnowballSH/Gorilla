package main

import (
	"Gorilla/repl"
	"os"
)

func main() {
	repl.Start(os.Stdin, os.Stdout)
}
