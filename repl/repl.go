package repl

import (
	"bufio"
	"fmt"
	"github.com/SnowballSH/Gorilla/exports"
	"github.com/SnowballSH/Gorilla/runtime"
	"github.com/fatih/color"
	"os"
	"strings"
)

// Start starts the repl
func Start() {
	fmt.Println("Welcome to Gorilla repl. Type :quit to quit.")

	env := runtime.NewEnvironment()

	history := []string{"print('Hello, world!')"}
	r := bufio.NewReader(os.Stdin)

	for {
		fmt.Print("> ")
		t, _, _ := r.ReadLine()
		text := strings.TrimSpace(string(t))
		if text == ":quit" {
			return
		}

		history = append(history, text)

		res, err := exports.CompileGorilla(text)

		if err != nil {
			color.Red(err.Error())
			continue
		}

		vm, lastPopped, err := exports.ExecuteGorillaBytecodeFromSourceAndEnv(res, text, env)

		if err != nil {
			color.Red(err.Error())
			continue
		}

		env = vm.Environment

		if lastPopped != nil {
			color.Green(fmt.Sprintf("#=> %s", vm.LastPopped.Inspect()))
		}
	}
}
