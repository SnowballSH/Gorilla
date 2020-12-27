package main

import (
	"fmt"
	"os"

	"./repl"
)

func main() {
	fmt.Println("Welcome to the Gorilla Repl.")
	repl.Start(os.Stdin, os.Stdout)
}
