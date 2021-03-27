package repl

import (
	"bufio"
	"fmt"
	"github.com/SnowballSH/Gorilla/exports"
	"github.com/SnowballSH/Gorilla/runtime"
	"os"
	"strings"
)

func Start() {
	reader := bufio.NewReader(os.Stdin)
	fmt.Println("Welcome to Gorilla repl. Type :quit to quit.")

	env := runtime.NewEnvironment()

	for {
		fmt.Print("> ")
		text, e := reader.ReadString('\n')
		if e != nil {
			return
		}
		text = strings.TrimSpace(text)
		if text == ":quit" {
			return
		}

		res, err := exports.CompileGorilla(text)

		if err != nil {
			fmt.Println(err.Error())
			continue
		}

		vm, lastPopped, err := exports.ExecuteGorillaBytecodeFromSourceAndEnv(res, text, env)
		env = vm.Environment

		if lastPopped != nil {
			fmt.Println(fmt.Sprintf("#=> %s", vm.LastPopped.Inspect()))
		}
	}
}
